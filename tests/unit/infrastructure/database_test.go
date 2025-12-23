// Package infrastructure_test provides unit tests for infrastructure components.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 DevOpsCorner Indonesia. All rights reserved.
package infrastructure_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test database configuration and connection string patterns

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

func (c DatabaseConfig) ConnectionString() string {
	switch c.Driver {
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			c.User, c.Password, c.Host, c.Port, c.Name)
	case "sqlite":
		return c.Name
	default:
		return ""
	}
}

func (c DatabaseConfig) DSN() string {
	switch c.Driver {
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
	case "mysql":
		return fmt.Sprintf("mysql://%s:%s@%s:%d/%s",
			c.User, c.Password, c.Host, c.Port, c.Name)
	case "sqlite":
		return fmt.Sprintf("sqlite://%s", c.Name)
	default:
		return ""
	}
}

func TestDatabaseConfig(t *testing.T) {
	t.Run("should generate postgres connection string", func(t *testing.T) {
		config := DatabaseConfig{
			Driver:   "postgres",
			Host:     "localhost",
			Port:     5432,
			Name:     "testdb",
			User:     "postgres",
			Password: "secret",
			SSLMode:  "disable",
		}

		connStr := config.ConnectionString()

		assert.Contains(t, connStr, "host=localhost")
		assert.Contains(t, connStr, "port=5432")
		assert.Contains(t, connStr, "user=postgres")
		assert.Contains(t, connStr, "dbname=testdb")
		assert.Contains(t, connStr, "sslmode=disable")
	})

	t.Run("should generate postgres DSN", func(t *testing.T) {
		config := DatabaseConfig{
			Driver:   "postgres",
			Host:     "localhost",
			Port:     5432,
			Name:     "testdb",
			User:     "postgres",
			Password: "secret",
			SSLMode:  "require",
		}

		dsn := config.DSN()

		assert.True(t, strings.HasPrefix(dsn, "postgres://"))
		assert.Contains(t, dsn, "postgres:secret@")
		assert.Contains(t, dsn, "localhost:5432")
		assert.Contains(t, dsn, "testdb")
		assert.Contains(t, dsn, "sslmode=require")
	})

	t.Run("should generate mysql connection string", func(t *testing.T) {
		config := DatabaseConfig{
			Driver:   "mysql",
			Host:     "localhost",
			Port:     3306,
			Name:     "testdb",
			User:     "root",
			Password: "secret",
		}

		connStr := config.ConnectionString()

		assert.Contains(t, connStr, "root:secret@tcp")
		assert.Contains(t, connStr, "localhost:3306")
		assert.Contains(t, connStr, "/testdb")
		assert.Contains(t, connStr, "parseTime=true")
	})

	t.Run("should generate sqlite connection string", func(t *testing.T) {
		config := DatabaseConfig{
			Driver: "sqlite",
			Name:   "./test.db",
		}

		connStr := config.ConnectionString()

		assert.Equal(t, "./test.db", connStr)
	})

	t.Run("should return empty for unknown driver", func(t *testing.T) {
		config := DatabaseConfig{
			Driver: "unknown",
		}

		assert.Empty(t, config.ConnectionString())
		assert.Empty(t, config.DSN())
	})
}

func TestDatabaseMigrationPatterns(t *testing.T) {
	t.Run("should generate migration filename", func(t *testing.T) {
		version := 1
		name := "create_users"

		upFilename := fmt.Sprintf("%06d_%s.up.sql", version, name)
		downFilename := fmt.Sprintf("%06d_%s.down.sql", version, name)

		assert.Equal(t, "000001_create_users.up.sql", upFilename)
		assert.Equal(t, "000001_create_users.down.sql", downFilename)
	})

	t.Run("should generate CREATE TABLE statement", func(t *testing.T) {
		tableName := "users"
		columns := []struct {
			Name       string
			Type       string
			Constraint string
		}{
			{"id", "UUID", "PRIMARY KEY DEFAULT gen_random_uuid()"},
			{"name", "VARCHAR(255)", "NOT NULL"},
			{"email", "VARCHAR(255)", "NOT NULL UNIQUE"},
			{"created_at", "TIMESTAMP", "NOT NULL DEFAULT NOW()"},
			{"updated_at", "TIMESTAMP", "NOT NULL DEFAULT NOW()"},
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", tableName))

		for i, col := range columns {
			sb.WriteString(fmt.Sprintf("    %s %s %s", col.Name, col.Type, col.Constraint))
			if i < len(columns)-1 {
				sb.WriteString(",\n")
			} else {
				sb.WriteString("\n")
			}
		}
		sb.WriteString(");")

		sql := sb.String()

		assert.Contains(t, sql, "CREATE TABLE users")
		assert.Contains(t, sql, "id UUID PRIMARY KEY")
		assert.Contains(t, sql, "email VARCHAR(255) NOT NULL UNIQUE")
	})

	t.Run("should generate DROP TABLE statement", func(t *testing.T) {
		tableName := "users"
		sql := fmt.Sprintf("DROP TABLE IF EXISTS %s;", tableName)

		assert.Equal(t, "DROP TABLE IF EXISTS users;", sql)
	})

	t.Run("should generate ALTER TABLE statement", func(t *testing.T) {
		tableName := "users"
		columnName := "phone"
		columnType := "VARCHAR(20)"

		sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;",
			tableName, columnName, columnType)

		assert.Contains(t, sql, "ALTER TABLE users")
		assert.Contains(t, sql, "ADD COLUMN phone VARCHAR(20)")
	})
}

func TestDatabaseQueryPatterns(t *testing.T) {
	t.Run("should build SELECT query", func(t *testing.T) {
		table := "users"
		columns := []string{"id", "name", "email"}
		condition := "active = true"

		sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
			strings.Join(columns, ", "), table, condition)

		assert.Equal(t, "SELECT id, name, email FROM users WHERE active = true", sql)
	})

	t.Run("should build INSERT query with RETURNING", func(t *testing.T) {
		table := "users"
		columns := []string{"name", "email"}
		placeholders := []string{"$1", "$2"}

		sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id",
			table, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

		assert.Equal(t, "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", sql)
	})

	t.Run("should build UPDATE query", func(t *testing.T) {
		table := "users"
		sets := []string{"name = $1", "email = $2", "updated_at = NOW()"}
		condition := "id = $3"

		sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
			table, strings.Join(sets, ", "), condition)

		assert.Contains(t, sql, "UPDATE users SET")
		assert.Contains(t, sql, "name = $1")
		assert.Contains(t, sql, "WHERE id = $3")
	})

	t.Run("should build DELETE query", func(t *testing.T) {
		table := "users"
		condition := "id = $1"

		sql := fmt.Sprintf("DELETE FROM %s WHERE %s", table, condition)

		assert.Equal(t, "DELETE FROM users WHERE id = $1", sql)
	})

	t.Run("should build paginated query", func(t *testing.T) {
		table := "users"
		orderBy := "created_at DESC"
		limit := 10
		offset := 20

		sql := fmt.Sprintf("SELECT * FROM %s ORDER BY %s LIMIT %d OFFSET %d",
			table, orderBy, limit, offset)

		assert.Contains(t, sql, "ORDER BY created_at DESC")
		assert.Contains(t, sql, "LIMIT 10")
		assert.Contains(t, sql, "OFFSET 20")
	})

	t.Run("should build COUNT query", func(t *testing.T) {
		table := "users"
		condition := "active = true"

		sql := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", table, condition)

		assert.Equal(t, "SELECT COUNT(*) FROM users WHERE active = true", sql)
	})
}

func TestDatabaseTypeMapping(t *testing.T) {
	typeMap := map[string]struct {
		goType   string
		pgType   string
		mysqlType string
	}{
		"string":    {"string", "VARCHAR(255)", "VARCHAR(255)"},
		"text":      {"string", "TEXT", "TEXT"},
		"int":       {"int", "INTEGER", "INT"},
		"int64":     {"int64", "BIGINT", "BIGINT"},
		"float64":   {"float64", "DECIMAL(10,2)", "DECIMAL(10,2)"},
		"bool":      {"bool", "BOOLEAN", "BOOLEAN"},
		"time":      {"time.Time", "TIMESTAMP", "DATETIME"},
		"uuid":      {"uuid.UUID", "UUID", "CHAR(36)"},
	}

	t.Run("should map Go types to PostgreSQL types", func(t *testing.T) {
		assert.Equal(t, "VARCHAR(255)", typeMap["string"].pgType)
		assert.Equal(t, "INTEGER", typeMap["int"].pgType)
		assert.Equal(t, "BOOLEAN", typeMap["bool"].pgType)
		assert.Equal(t, "TIMESTAMP", typeMap["time"].pgType)
		assert.Equal(t, "UUID", typeMap["uuid"].pgType)
	})

	t.Run("should map Go types to MySQL types", func(t *testing.T) {
		assert.Equal(t, "VARCHAR(255)", typeMap["string"].mysqlType)
		assert.Equal(t, "INT", typeMap["int"].mysqlType)
		assert.Equal(t, "DATETIME", typeMap["time"].mysqlType)
		assert.Equal(t, "CHAR(36)", typeMap["uuid"].mysqlType)
	})
}

func TestDatabaseConnectionPooling(t *testing.T) {
	type PoolConfig struct {
		MaxOpenConns    int
		MaxIdleConns    int
		ConnMaxLifetime int // seconds
		ConnMaxIdleTime int // seconds
	}

	t.Run("should have reasonable pool defaults", func(t *testing.T) {
		config := PoolConfig{
			MaxOpenConns:    25,
			MaxIdleConns:    10,
			ConnMaxLifetime: 300,  // 5 minutes
			ConnMaxIdleTime: 60,   // 1 minute
		}

		assert.GreaterOrEqual(t, config.MaxOpenConns, config.MaxIdleConns)
		assert.Greater(t, config.ConnMaxLifetime, 0)
		assert.Greater(t, config.ConnMaxIdleTime, 0)
	})

	t.Run("should validate pool configuration", func(t *testing.T) {
		validatePoolConfig := func(config PoolConfig) error {
			if config.MaxOpenConns <= 0 {
				return fmt.Errorf("MaxOpenConns must be positive")
			}
			if config.MaxIdleConns < 0 {
				return fmt.Errorf("MaxIdleConns cannot be negative")
			}
			if config.MaxIdleConns > config.MaxOpenConns {
				return fmt.Errorf("MaxIdleConns cannot exceed MaxOpenConns")
			}
			return nil
		}

		t.Run("valid config", func(t *testing.T) {
			config := PoolConfig{MaxOpenConns: 25, MaxIdleConns: 10}
			assert.NoError(t, validatePoolConfig(config))
		})

		t.Run("invalid MaxOpenConns", func(t *testing.T) {
			config := PoolConfig{MaxOpenConns: 0, MaxIdleConns: 10}
			assert.Error(t, validatePoolConfig(config))
		})

		t.Run("MaxIdleConns exceeds MaxOpenConns", func(t *testing.T) {
			config := PoolConfig{MaxOpenConns: 10, MaxIdleConns: 25}
			assert.Error(t, validatePoolConfig(config))
		})
	})
}
