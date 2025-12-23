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

	"{{.ModulePath}}/telemetry"
	"{{.ModulePath}}/telemetry/logs"
	"{{.ModulePath}}/telemetry/metrics"
	"{{.ModulePath}}/telemetry/traces"
)

const (
	numWorkers = {{.NumWorkers}}
	queueSize  = {{.QueueSize}}
)

// Job represents a unit of work
type Job struct {
	ID        string
	Type      string
	Payload   interface{}
	CreatedAt time.Time
}

func main() {
	// Initialize TelemetryFlow SDK
	if err := telemetry.Init(); err != nil {
		log.Fatal(err)
	}
	defer telemetry.Shutdown()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create job queue
	jobQueue := make(chan Job, queueSize)

	// Start worker pool
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, &wg, i, jobQueue)
	}

	logs.Info("Worker pool started", map[string]interface{}{
		"num_workers": numWorkers,
		"queue_size":  queueSize,
	})

	metrics.RecordGauge("worker.pool.size", float64(numWorkers), nil)

	// Start job producer (simulates incoming jobs)
	go producer(ctx, jobQueue)

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logs.Info("Shutdown signal received, stopping workers...", nil)

	// Cancel context and wait for workers to finish
	cancel()
	close(jobQueue)
	wg.Wait()

	logs.Info("All workers stopped gracefully", nil)
}

func worker(ctx context.Context, wg *sync.WaitGroup, id int, jobs <-chan Job) {
	defer wg.Done()

	workerID := fmt.Sprintf("worker-%d", id)
	logs.Info("Worker started", map[string]interface{}{
		"worker_id": workerID,
	})

	for {
		select {
		case <-ctx.Done():
			logs.Info("Worker stopping", map[string]interface{}{
				"worker_id": workerID,
			})
			return

		case job, ok := <-jobs:
			if !ok {
				return
			}
			processJob(ctx, workerID, job)
		}
	}
}

func processJob(ctx context.Context, workerID string, job Job) {
	start := time.Now()

	// Start a trace span for this job
	spanID, err := traces.StartSpan(ctx, fmt.Sprintf("job.process.%s", job.Type), map[string]interface{}{
		"job.id":        job.ID,
		"job.type":      job.Type,
		"worker.id":     workerID,
		"job.queue_age": time.Since(job.CreatedAt).Seconds(),
	})
	if err != nil {
		logs.Error("Failed to start span", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Track queue wait time
	queueWaitTime := time.Since(job.CreatedAt).Seconds()
	metrics.RecordHistogram("job.queue.wait_time", queueWaitTime, "s", map[string]interface{}{
		"job_type": job.Type,
	})

	// Simulate job processing with variable duration
	processingTime := time.Duration(50+rand.Intn(200)) * time.Millisecond
	time.Sleep(processingTime)

	// Simulate occasional errors (10% failure rate)
	var jobErr error
	if rand.Float32() < 0.1 {
		jobErr = fmt.Errorf("simulated processing error for job %s", job.ID)
	}

	// Record job completion
	duration := time.Since(start).Seconds()
	status := "success"
	if jobErr != nil {
		status = "failed"
		logs.Error("Job processing failed", map[string]interface{}{
			"job_id":    job.ID,
			"job_type":  job.Type,
			"worker_id": workerID,
			"duration":  duration,
			"error":     jobErr.Error(),
		})
	} else {
		logs.Info("Job processed successfully", map[string]interface{}{
			"job_id":    job.ID,
			"job_type":  job.Type,
			"worker_id": workerID,
			"duration":  duration,
		})
	}

	// Record metrics
	metrics.RecordHistogram("job.processing.duration", duration, "s", map[string]interface{}{
		"job_type": job.Type,
		"status":   status,
	})
	metrics.IncrementCounter("jobs.processed.total", 1, map[string]interface{}{
		"job_type": job.Type,
		"status":   status,
	})

	// End the trace span
	traces.EndSpan(ctx, spanID, jobErr)
}

func producer(ctx context.Context, jobs chan<- Job) {
	jobTypes := []string{"email", "notification", "report", "sync", "cleanup"}
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	jobCounter := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			jobCounter++
			job := Job{
				ID:        fmt.Sprintf("job-%d", jobCounter),
				Type:      jobTypes[rand.Intn(len(jobTypes))],
				Payload:   nil,
				CreatedAt: time.Now(),
			}

			select {
			case jobs <- job:
				metrics.IncrementCounter("jobs.queued.total", 1, map[string]interface{}{
					"job_type": job.Type,
				})
				metrics.RecordGauge("job.queue.length", float64(len(jobs)), nil)
			default:
				// Queue is full
				logs.Warn("Job queue full, dropping job", map[string]interface{}{
					"job_id":   job.ID,
					"job_type": job.Type,
				})
				metrics.IncrementCounter("jobs.dropped.total", 1, map[string]interface{}{
					"job_type": job.Type,
				})
			}
		}
	}
}
