// Package main demonstrates TelemetryFlow SDK v1.1.2 integration with a background worker.
//
// This example shows:
// - Background job processing with tracing
// - Job queue metrics (pending, processed, failed)
// - Worker health monitoring
// - TFO v2 API with collector identity
// - Graceful shutdown with job completion
//
// Compatible with TFO-Collector v1.1.2 (OCB-native)
//
// Run with:
//
//	export TELEMETRYFLOW_API_KEY_ID=tfk_your_key
//	export TELEMETRYFLOW_API_KEY_SECRET=tfs_your_secret
//	go run main.go
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow"
)

var client *telemetryflow.Client

// Job represents a background job to be processed
type Job struct {
	ID        string
	Type      string
	Payload   map[string]interface{}
	CreatedAt time.Time
	Attempts  int
}

// Worker processes jobs from a queue
type Worker struct {
	id     string
	jobs   chan Job
	quit   chan struct{}
	wg     *sync.WaitGroup
	client *telemetryflow.Client
	ctx    context.Context
}

func main() {
	// Initialize TelemetryFlow client with TFO v2 API support
	var err error
	client, err = telemetryflow.NewBuilder().
		WithAPIKeyFromEnv().
		WithEndpointFromEnv().
		WithService("worker-example", "1.1.2").
		WithServiceNamespace("telemetryflow").
		WithEnvironmentFromEnv().
		WithGRPC().
		// TFO v2 API settings (aligned with tfoexporter)
		WithV2API(true).
		// Collector Identity (aligned with tfoidentityextension)
		WithCollectorName("Worker Pool Example").
		WithCollectorDescription("TelemetryFlow Go SDK worker pool example").
		WithDatacenterFromEnv().
		WithCollectorTag("pool_size", "3").
		WithEnrichResources(true).
		// Enable exemplars for metrics-to-traces correlation
		WithExemplars(true).
		// Custom attributes
		WithCustomAttribute("example", "worker").
		WithCustomAttribute("worker_pool_size", "3").
		WithCustomAttribute("sdk_version", "1.1.2").
		Build()

	if err != nil {
		log.Fatalf("Failed to create TelemetryFlow client: %v", err)
	}

	ctx := context.Background()
	if err := client.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize TelemetryFlow: %v", err)
	}

	if err := client.LogInfo(ctx, "Worker pool starting", map[string]interface{}{
		"pool_size": 3,
		"version":   "1.0.0",
	}); err != nil {
		log.Printf("Failed to log info: %v", err)
	}

	// Create job queue
	jobs := make(chan Job, 100)
	quit := make(chan struct{})
	var wg sync.WaitGroup

	// Start workers
	workers := make([]*Worker, 3)
	for i := 0; i < 3; i++ {
		workers[i] = NewWorker(fmt.Sprintf("worker-%d", i+1), jobs, quit, &wg, client, ctx)
		go workers[i].Start()
	}

	// Start job producer
	go produceJobs(ctx, jobs, quit)

	// Start metrics reporter
	go reportMetrics(ctx, jobs, quit)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutdown signal received, stopping workers...")
	if err := client.LogInfo(ctx, "Worker pool shutdown initiated", nil); err != nil {
		log.Printf("Failed to log info: %v", err)
	}

	// Signal shutdown
	close(quit)

	// Wait for workers to finish current jobs
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Wait with timeout
	select {
	case <-done:
		log.Println("All workers finished gracefully")
	case <-time.After(30 * time.Second):
		log.Println("Timeout waiting for workers")
	}

	// Report final queue stats
	if err := client.RecordGauge(ctx, "worker.queue.pending", float64(len(jobs)), map[string]interface{}{
		"status": "shutdown",
	}); err != nil {
		log.Printf("Failed to record gauge: %v", err)
	}

	// Flush and shutdown telemetry
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Flush(shutdownCtx); err != nil {
		log.Printf("Failed to flush: %v", err)
	}
	if err := client.Shutdown(shutdownCtx); err != nil {
		log.Printf("Failed to shutdown: %v", err)
	}

	log.Println("Worker pool stopped")
}

// NewWorker creates a new worker
func NewWorker(id string, jobs chan Job, quit chan struct{}, wg *sync.WaitGroup, client *telemetryflow.Client, ctx context.Context) *Worker {
	return &Worker{
		id:     id,
		jobs:   jobs,
		quit:   quit,
		wg:     wg,
		client: client,
		ctx:    ctx,
	}
}

// Start begins processing jobs
func (w *Worker) Start() {
	log.Printf("[%s] Starting worker", w.id)
	if err := w.client.LogInfo(w.ctx, "Worker started", map[string]interface{}{
		"worker_id": w.id,
	}); err != nil {
		log.Printf("Failed to log info: %v", err)
	}

	for {
		select {
		case <-w.quit:
			log.Printf("[%s] Worker stopping", w.id)
			if err := w.client.LogInfo(w.ctx, "Worker stopped", map[string]interface{}{
				"worker_id": w.id,
			}); err != nil {
				log.Printf("Failed to log info: %v", err)
			}
			return

		case job := <-w.jobs:
			w.wg.Add(1)
			w.processJob(job)
			w.wg.Done()
		}
	}
}

// processJob handles a single job
func (w *Worker) processJob(job Job) {
	start := time.Now()
	ctx := w.ctx

	// Start job processing span
	spanID, err := w.client.StartSpan(ctx, fmt.Sprintf("job.process.%s", job.Type), "consumer", map[string]interface{}{
		"job.id":       job.ID,
		"job.type":     job.Type,
		"job.attempts": job.Attempts,
		"worker.id":    w.id,
	})
	if err != nil {
		log.Printf("[%s] Failed to start span: %v", w.id, err)
	}

	log.Printf("[%s] Processing job %s (type: %s)", w.id, job.ID, job.Type)

	// Record job picked up
	if err := w.client.IncrementCounter(ctx, "worker.jobs.picked_up", 1, map[string]interface{}{
		"job_type":  job.Type,
		"worker_id": w.id,
	}); err != nil {
		log.Printf("Failed to increment counter: %v", err)
	}

	// Record queue wait time
	waitTime := time.Since(job.CreatedAt)
	if err := w.client.RecordHistogram(ctx, "worker.job.queue_time", waitTime.Seconds(), "s", map[string]interface{}{
		"job_type": job.Type,
	}); err != nil {
		log.Printf("Failed to record histogram: %v", err)
	}

	// Process based on job type
	var processErr error
	switch job.Type {
	case "email":
		processErr = w.processEmailJob(ctx, job, spanID)
	case "notification":
		processErr = w.processNotificationJob(ctx, job, spanID)
	case "report":
		processErr = w.processReportJob(ctx, job, spanID)
	default:
		processErr = fmt.Errorf("unknown job type: %s", job.Type)
	}

	duration := time.Since(start)

	// Record job completion metrics
	status := "success"
	if processErr != nil {
		status = "failed"
		if err := w.client.IncrementCounter(ctx, "worker.jobs.failed", 1, map[string]interface{}{
			"job_type":  job.Type,
			"worker_id": w.id,
			"error":     processErr.Error(),
		}); err != nil {
			log.Printf("Failed to increment counter: %v", err)
		}
		if err := w.client.LogError(ctx, "Job processing failed", map[string]interface{}{
			"job_id":      job.ID,
			"job_type":    job.Type,
			"worker_id":   w.id,
			"error":       processErr.Error(),
			"duration_ms": duration.Milliseconds(),
		}); err != nil {
			log.Printf("Failed to log error: %v", err)
		}
	} else {
		if err := w.client.IncrementCounter(ctx, "worker.jobs.completed", 1, map[string]interface{}{
			"job_type":  job.Type,
			"worker_id": w.id,
		}); err != nil {
			log.Printf("Failed to increment counter: %v", err)
		}
	}

	if err := w.client.RecordHistogram(ctx, "worker.job.duration", duration.Seconds(), "s", map[string]interface{}{
		"job_type":  job.Type,
		"worker_id": w.id,
		"status":    status,
	}); err != nil {
		log.Printf("Failed to record histogram: %v", err)
	}

	// End span
	if spanID != "" {
		if err := w.client.EndSpan(ctx, spanID, processErr); err != nil {
			log.Printf("Failed to end span: %v", err)
		}
	}

	log.Printf("[%s] Completed job %s in %v (status: %s)", w.id, job.ID, duration, status)
}

func (w *Worker) processEmailJob(ctx context.Context, job Job, parentSpanID string) error {
	// Start email sending span
	spanID, _ := w.client.StartSpan(ctx, "email.send", "client", map[string]interface{}{
		"email.to":      job.Payload["to"],
		"email.subject": job.Payload["subject"],
	})

	// Simulate email sending
	time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)

	// Simulate occasional failures
	if rand.Float32() < 0.1 {
		err := fmt.Errorf("SMTP connection timeout")
		if err := w.client.EndSpan(ctx, spanID, err); err != nil {
			log.Printf("Failed to end span: %v", err)
		}
		return err
	}

	if err := w.client.AddSpanEvent(ctx, spanID, "email.sent", map[string]interface{}{
		"provider": "sendgrid",
	}); err != nil {
		log.Printf("Failed to add span event: %v", err)
	}

	if err := w.client.IncrementCounter(ctx, "emails.sent", 1, map[string]interface{}{
		"provider": "sendgrid",
	}); err != nil {
		log.Printf("Failed to increment counter: %v", err)
	}

	if err := w.client.EndSpan(ctx, spanID, nil); err != nil {
		log.Printf("Failed to end span: %v", err)
	}
	return nil
}

func (w *Worker) processNotificationJob(ctx context.Context, job Job, parentSpanID string) error {
	// Start notification span
	spanID, _ := w.client.StartSpan(ctx, "notification.push", "client", map[string]interface{}{
		"notification.type":   job.Payload["type"],
		"notification.target": job.Payload["target"],
	})

	// Simulate push notification
	time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)

	if err := w.client.AddSpanEvent(ctx, spanID, "notification.delivered", nil); err != nil {
		log.Printf("Failed to add span event: %v", err)
	}

	if err := w.client.IncrementCounter(ctx, "notifications.sent", 1, map[string]interface{}{
		"type": job.Payload["type"],
	}); err != nil {
		log.Printf("Failed to increment counter: %v", err)
	}

	if err := w.client.EndSpan(ctx, spanID, nil); err != nil {
		log.Printf("Failed to end span: %v", err)
	}
	return nil
}

func (w *Worker) processReportJob(ctx context.Context, job Job, parentSpanID string) error {
	// Start report generation span
	spanID, _ := w.client.StartSpan(ctx, "report.generate", "internal", map[string]interface{}{
		"report.type":   job.Payload["report_type"],
		"report.format": job.Payload["format"],
	})

	// Simulate report generation (longer task)
	steps := []string{"fetch_data", "process_data", "generate_charts", "compile_pdf"}
	for _, step := range steps {
		if err := w.client.AddSpanEvent(ctx, spanID, fmt.Sprintf("report.%s", step), nil); err != nil {
			log.Printf("Failed to add span event: %v", err)
		}
		time.Sleep(time.Duration(100+rand.Intn(150)) * time.Millisecond)
	}

	// Record report size
	reportSize := float64(rand.Intn(5000) + 1000) // KB
	if err := w.client.RecordHistogram(ctx, "report.size", reportSize, "KB", map[string]interface{}{
		"report_type": job.Payload["report_type"],
	}); err != nil {
		log.Printf("Failed to record histogram: %v", err)
	}

	if err := w.client.IncrementCounter(ctx, "reports.generated", 1, map[string]interface{}{
		"report_type": job.Payload["report_type"],
		"format":      job.Payload["format"],
	}); err != nil {
		log.Printf("Failed to increment counter: %v", err)
	}

	if err := w.client.EndSpan(ctx, spanID, nil); err != nil {
		log.Printf("Failed to end span: %v", err)
	}
	return nil
}

// produceJobs simulates job production
func produceJobs(ctx context.Context, jobs chan Job, quit chan struct{}) {
	jobTypes := []string{"email", "notification", "report"}
	jobID := 0

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			jobID++
			jobType := jobTypes[rand.Intn(len(jobTypes))]

			job := Job{
				ID:        fmt.Sprintf("job_%05d", jobID),
				Type:      jobType,
				CreatedAt: time.Now(),
				Attempts:  1,
				Payload:   generatePayload(jobType),
			}

			select {
			case jobs <- job:
				if err := client.IncrementCounter(ctx, "worker.jobs.queued", 1, map[string]interface{}{
					"job_type": jobType,
				}); err != nil {
					log.Printf("Failed to increment counter: %v", err)
				}
			default:
				// Queue full, log and drop
				if err := client.LogWarn(ctx, "Job queue full, job dropped", map[string]interface{}{
					"job_id":   job.ID,
					"job_type": jobType,
				}); err != nil {
					log.Printf("Failed to log warning: %v", err)
				}
				if err := client.IncrementCounter(ctx, "worker.jobs.dropped", 1, map[string]interface{}{
					"job_type": jobType,
					"reason":   "queue_full",
				}); err != nil {
					log.Printf("Failed to increment counter: %v", err)
				}
			}
		}
	}
}

// generatePayload creates sample job payload
func generatePayload(jobType string) map[string]interface{} {
	switch jobType {
	case "email":
		return map[string]interface{}{
			"to":      fmt.Sprintf("user%d@example.com", rand.Intn(1000)),
			"subject": "Important notification",
			"body":    "This is the email body.",
		}
	case "notification":
		return map[string]interface{}{
			"type":    []string{"push", "sms", "in-app"}[rand.Intn(3)],
			"target":  fmt.Sprintf("user_%d", rand.Intn(1000)),
			"message": "You have a new message",
		}
	case "report":
		return map[string]interface{}{
			"report_type": []string{"daily", "weekly", "monthly"}[rand.Intn(3)],
			"format":      []string{"pdf", "xlsx", "csv"}[rand.Intn(3)],
			"user_id":     rand.Intn(1000),
		}
	default:
		return map[string]interface{}{}
	}
}

// reportMetrics periodically reports queue metrics
func reportMetrics(ctx context.Context, jobs chan Job, quit chan struct{}) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			queueLen := len(jobs)
			queueCap := cap(jobs)

			if err := client.RecordGauge(ctx, "worker.queue.pending", float64(queueLen), map[string]interface{}{
				"capacity": queueCap,
			}); err != nil {
				log.Printf("Failed to record gauge: %v", err)
			}

			if err := client.RecordGauge(ctx, "worker.queue.utilization", float64(queueLen)/float64(queueCap)*100, map[string]interface{}{
				"unit": "percent",
			}); err != nil {
				log.Printf("Failed to record gauge: %v", err)
			}

			// Log if queue is getting full
			if float64(queueLen)/float64(queueCap) > 0.8 {
				if err := client.LogWarn(ctx, "Job queue utilization high", map[string]interface{}{
					"pending":     queueLen,
					"capacity":    queueCap,
					"utilization": fmt.Sprintf("%.1f%%", float64(queueLen)/float64(queueCap)*100),
				}); err != nil {
					log.Printf("Failed to log warning: %v", err)
				}
			}
		}
	}
}
