// Package instrumentation provides auto-instrumentation for TelemetryFlow SDK.
// This file provides database instrumentation for database/sql.
package instrumentation

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

// DBConfig holds database instrumentation configuration
type DBConfig struct {
	*Config
	// DatabaseSystem is the database system (e.g., "postgresql", "mysql")
	DatabaseSystem string
	// DatabaseName is the name of the database
	DatabaseName string
	// ServerAddress is the database server address
	ServerAddress string
	// ServerPort is the database server port
	ServerPort int
	// RecordStatement records the SQL statement in spans
	RecordStatement bool
	// RecordParameters records query parameters (use with caution for PII)
	RecordParameters bool
	// SanitizeStatement sanitizes the SQL statement
	SanitizeStatement bool
}

// DefaultDBConfig returns default database configuration
func DefaultDBConfig() *DBConfig {
	return &DBConfig{
		Config:            DefaultConfig(),
		DatabaseSystem:    "unknown",
		RecordStatement:   true,
		SanitizeStatement: true,
	}
}

// DBOption is a function that configures database instrumentation
type DBOption func(*DBConfig)

// WithDatabaseSystem sets the database system
func WithDatabaseSystem(system string) DBOption {
	return func(c *DBConfig) {
		c.DatabaseSystem = system
	}
}

// WithDatabaseName sets the database name
func WithDatabaseName(name string) DBOption {
	return func(c *DBConfig) {
		c.DatabaseName = name
	}
}

// WithServerAddress sets the server address and port
func WithServerAddress(address string, port int) DBOption {
	return func(c *DBConfig) {
		c.ServerAddress = address
		c.ServerPort = port
	}
}

// WithRecordStatement enables/disables SQL statement recording
func WithRecordStatement(enabled bool) DBOption {
	return func(c *DBConfig) {
		c.RecordStatement = enabled
	}
}

// WithSanitizeStatement enables/disables SQL statement sanitization
func WithSanitizeStatement(enabled bool) DBOption {
	return func(c *DBConfig) {
		c.SanitizeStatement = enabled
	}
}

// InstrumentedDB wraps a sql.DB with tracing and metrics
type InstrumentedDB struct {
	*sql.DB
	tracer  trace.Tracer
	config  *DBConfig
	metrics *DatabaseMetrics
}

// WrapDB wraps a sql.DB with tracing and metrics
func WrapDB(db *sql.DB, opts ...DBOption) *InstrumentedDB {
	cfg := DefaultDBConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	var tp trace.TracerProvider
	if cfg.TracerProvider != nil {
		tp = cfg.TracerProvider
	} else {
		tp = otel.GetTracerProvider()
	}

	tracer := tp.Tracer(InstrumentationName,
		trace.WithInstrumentationVersion(InstrumentationVersion),
	)

	var metrics *DatabaseMetrics
	if cfg.EnableMetrics && cfg.MeterProvider != nil {
		var err error
		metrics, err = NewDatabaseMetrics(cfg.MeterProvider)
		if err != nil {
			metrics = nil
		}
	}

	return &InstrumentedDB{
		DB:      db,
		tracer:  tracer,
		config:  cfg,
		metrics: metrics,
	}
}

// QueryContext executes a query with tracing
func (db *InstrumentedDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.queryContext(ctx, query, args...)
}

func (db *InstrumentedDB) queryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()

	operation := extractOperation(query)
	table := extractTable(query)
	spanName := fmt.Sprintf("%s %s", operation, table)

	attrs := db.createAttributes(query, operation, table)

	ctx, span := db.tracer.Start(ctx, spanName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	rows, err := db.DB.QueryContext(ctx, query, args...)

	duration := time.Since(start)
	hasError := err != nil

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	if db.metrics != nil {
		db.metrics.RecordQuery(ctx, operation, table, duration, hasError)
	}

	return rows, err
}

// QueryRowContext executes a query that returns a single row with tracing
func (db *InstrumentedDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	start := time.Now()

	operation := extractOperation(query)
	table := extractTable(query)
	spanName := fmt.Sprintf("%s %s", operation, table)

	attrs := db.createAttributes(query, operation, table)

	ctx, span := db.tracer.Start(ctx, spanName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	row := db.DB.QueryRowContext(ctx, query, args...)

	duration := time.Since(start)
	if db.metrics != nil {
		db.metrics.RecordQuery(ctx, operation, table, duration, false)
	}

	span.SetStatus(codes.Ok, "")
	return row
}

// ExecContext executes a query that doesn't return rows with tracing
func (db *InstrumentedDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()

	operation := extractOperation(query)
	table := extractTable(query)
	spanName := fmt.Sprintf("%s %s", operation, table)

	attrs := db.createAttributes(query, operation, table)

	ctx, span := db.tracer.Start(ctx, spanName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	result, err := db.DB.ExecContext(ctx, query, args...)

	duration := time.Since(start)
	hasError := err != nil

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
		if rowsAffected, raErr := result.RowsAffected(); raErr == nil {
			span.SetAttributes(attribute.Int64("db.rows_affected", rowsAffected))
		}
	}

	if db.metrics != nil {
		db.metrics.RecordQuery(ctx, operation, table, duration, hasError)
	}

	return result, err
}

// PrepareContext prepares a statement with tracing
func (db *InstrumentedDB) PrepareContext(ctx context.Context, query string) (*InstrumentedStmt, error) {
	stmt, err := db.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return &InstrumentedStmt{
		Stmt:    stmt,
		query:   query,
		tracer:  db.tracer,
		config:  db.config,
		metrics: db.metrics,
	}, nil
}

// BeginTx starts a transaction with tracing
func (db *InstrumentedDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*InstrumentedTx, error) {
	_, span := db.tracer.Start(ctx, "BEGIN",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(db.createBaseAttributes()...),
	)
	defer span.End()

	tx, err := db.DB.BeginTx(ctx, opts)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "")

	return &InstrumentedTx{
		Tx:      tx,
		tracer:  db.tracer,
		config:  db.config,
		metrics: db.metrics,
	}, nil
}

func (db *InstrumentedDB) createBaseAttributes() []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		semconv.DBSystemKey.String(db.config.DatabaseSystem),
	}

	if db.config.DatabaseName != "" {
		attrs = append(attrs, attribute.String("db.name", db.config.DatabaseName))
	}
	if db.config.ServerAddress != "" {
		attrs = append(attrs, semconv.ServerAddress(db.config.ServerAddress))
	}
	if db.config.ServerPort > 0 {
		attrs = append(attrs, semconv.ServerPort(db.config.ServerPort))
	}

	return attrs
}

func (db *InstrumentedDB) createAttributes(query, operation, table string) []attribute.KeyValue {
	attrs := db.createBaseAttributes()

	attrs = append(attrs,
		attribute.String("db.operation", operation),
	)

	if table != "" {
		attrs = append(attrs, attribute.String("db.sql.table", table))
	}

	if db.config.RecordStatement {
		stmt := query
		if db.config.SanitizeStatement {
			stmt = sanitizeQuery(query)
		}
		attrs = append(attrs, attribute.String("db.statement", stmt))
	}

	return attrs
}

// InstrumentedStmt wraps a sql.Stmt with tracing
type InstrumentedStmt struct {
	*sql.Stmt
	query   string
	tracer  trace.Tracer
	config  *DBConfig
	metrics *DatabaseMetrics
}

// QueryContext executes the prepared query with tracing
func (s *InstrumentedStmt) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()

	operation := extractOperation(s.query)
	table := extractTable(s.query)
	spanName := fmt.Sprintf("%s %s", operation, table)

	attrs := []attribute.KeyValue{
		semconv.DBSystemKey.String(s.config.DatabaseSystem),
		attribute.String("db.operation", operation),
	}
	if table != "" {
		attrs = append(attrs, attribute.String("db.sql.table", table))
	}

	ctx, span := s.tracer.Start(ctx, spanName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	rows, err := s.Stmt.QueryContext(ctx, args...)

	duration := time.Since(start)
	hasError := err != nil

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	if s.metrics != nil {
		s.metrics.RecordQuery(ctx, operation, table, duration, hasError)
	}

	return rows, err
}

// ExecContext executes the prepared statement with tracing
func (s *InstrumentedStmt) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	start := time.Now()

	operation := extractOperation(s.query)
	table := extractTable(s.query)
	spanName := fmt.Sprintf("%s %s", operation, table)

	attrs := []attribute.KeyValue{
		semconv.DBSystemKey.String(s.config.DatabaseSystem),
		attribute.String("db.operation", operation),
	}
	if table != "" {
		attrs = append(attrs, attribute.String("db.sql.table", table))
	}

	ctx, span := s.tracer.Start(ctx, spanName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	result, err := s.Stmt.ExecContext(ctx, args...)

	duration := time.Since(start)
	hasError := err != nil

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	if s.metrics != nil {
		s.metrics.RecordQuery(ctx, operation, table, duration, hasError)
	}

	return result, err
}

// InstrumentedTx wraps a sql.Tx with tracing
type InstrumentedTx struct {
	*sql.Tx
	tracer  trace.Tracer
	config  *DBConfig
	metrics *DatabaseMetrics
}

// QueryContext executes a query within the transaction with tracing
func (tx *InstrumentedTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()

	operation := extractOperation(query)
	table := extractTable(query)
	spanName := fmt.Sprintf("%s %s", operation, table)

	attrs := []attribute.KeyValue{
		semconv.DBSystemKey.String(tx.config.DatabaseSystem),
		attribute.String("db.operation", operation),
	}
	if table != "" {
		attrs = append(attrs, attribute.String("db.sql.table", table))
	}

	ctx, span := tx.tracer.Start(ctx, spanName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	rows, err := tx.Tx.QueryContext(ctx, query, args...)

	duration := time.Since(start)
	hasError := err != nil

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	if tx.metrics != nil {
		tx.metrics.RecordQuery(ctx, operation, table, duration, hasError)
	}

	return rows, err
}

// ExecContext executes a statement within the transaction with tracing
func (tx *InstrumentedTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()

	operation := extractOperation(query)
	table := extractTable(query)
	spanName := fmt.Sprintf("%s %s", operation, table)

	attrs := []attribute.KeyValue{
		semconv.DBSystemKey.String(tx.config.DatabaseSystem),
		attribute.String("db.operation", operation),
	}
	if table != "" {
		attrs = append(attrs, attribute.String("db.sql.table", table))
	}

	ctx, span := tx.tracer.Start(ctx, spanName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	result, err := tx.Tx.ExecContext(ctx, query, args...)

	duration := time.Since(start)
	hasError := err != nil

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	if tx.metrics != nil {
		tx.metrics.RecordQuery(ctx, operation, table, duration, hasError)
	}

	return result, err
}

// Commit commits the transaction with tracing
func (tx *InstrumentedTx) Commit() error {
	_, span := tx.tracer.Start(context.Background(), "COMMIT",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(semconv.DBSystemKey.String(tx.config.DatabaseSystem)),
	)
	defer span.End()

	err := tx.Tx.Commit()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return err
}

// Rollback rolls back the transaction with tracing
func (tx *InstrumentedTx) Rollback() error {
	_, span := tx.tracer.Start(context.Background(), "ROLLBACK",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(semconv.DBSystemKey.String(tx.config.DatabaseSystem)),
	)
	defer span.End()

	err := tx.Tx.Rollback()
	if err != nil && err != sql.ErrTxDone {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return err
}

// Helper functions

var operationRegex = regexp.MustCompile(`(?i)^\s*(\w+)`)
var tableRegex = regexp.MustCompile(`(?i)(?:FROM|INTO|UPDATE|TABLE)\s+["']?(\w+)["']?`)

func extractOperation(query string) string {
	matches := operationRegex.FindStringSubmatch(query)
	if len(matches) > 1 {
		return strings.ToUpper(matches[1])
	}
	return "QUERY"
}

func extractTable(query string) string {
	matches := tableRegex.FindStringSubmatch(query)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// sanitizeQuery removes sensitive values from SQL queries
func sanitizeQuery(query string) string {
	// Replace string literals with '?'
	sanitized := regexp.MustCompile(`'[^']*'`).ReplaceAllString(query, "'?'")
	// Replace numeric values with '?'
	sanitized = regexp.MustCompile(`\b\d+\b`).ReplaceAllString(sanitized, "?")
	return sanitized
}

// Ensure interfaces are implemented
var (
	_ driver.Connector = (*instrumentedConnector)(nil)
)

// instrumentedConnector wraps a driver.Connector
type instrumentedConnector struct {
	connector driver.Connector
	config    *DBConfig
}

// WrapConnector wraps a driver.Connector with instrumentation
func WrapConnector(connector driver.Connector, opts ...DBOption) driver.Connector {
	cfg := DefaultDBConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	return &instrumentedConnector{
		connector: connector,
		config:    cfg,
	}
}

func (c *instrumentedConnector) Connect(ctx context.Context) (driver.Conn, error) {
	return c.connector.Connect(ctx)
}

func (c *instrumentedConnector) Driver() driver.Driver {
	return c.connector.Driver()
}
