package service

import (
	"context"
	"fmt"
	"jobqueue/entity"
	_interface "jobqueue/interface"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	MaxRetries     = 3
	WorkerPoolSize = 5 // Number of workers processing jobs
)

type jobService struct {
	jobRepo  _interface.JobRepository
	jobQueue chan *entity.Job // The job queue
	wg       sync.WaitGroup   // To wait for workers to finish (for graceful shutdown/tests)
	mu       sync.Mutex       // To protect any shared service-level data (if any)
}

// Initiator ...
type Initiator func(s *jobService) *jobService

// startWorkerPool initializes and starts job processing workers
func (q *jobService) startWorkerPool() {
	for i := 0; i < WorkerPoolSize; i++ {
		q.wg.Add(1)
		go q.worker(i)
	}
	logrus.Infof("Started %d job queue workers", WorkerPoolSize)
}

// worker function processes jobs from the jobQueue
func (q *jobService) worker(id int) {
	defer q.wg.Done()
	logrus.Debugf("Worker %d started", id)

	for job := range q.jobQueue {
		// Job processing logic
		q.processJob(job)
	}
	logrus.Debugf("Worker %d stopped", id)
}

// processJob handles the execution and retry logic for a job
func (q *jobService) processJob(job *entity.Job) {
	job.Status = entity.RUNNING
	ctx := context.Background()

	// Save status update to repo
	if err := q.jobRepo.Save(ctx, job); err != nil {
		logrus.Errorf("Failed to update job %s status to RUNNING: %v", job.ID, err)
		return
	}

	for job.Attempts < MaxRetries {
		job.Attempts++
		logrus.Infof("Processing job ID: %s, Task: %s, Attempt: %d", job.ID, job.Task, job.Attempts)

		// Simulate job execution
		success := q.executeTask(job.Task, job.Attempts)

		if success {
			job.Status = entity.COMPLETED
			logrus.Infof("Job ID: %s, Task: %s completed successfully", job.ID, job.Task)
			break
		}

		// If failed
		logrus.Warnf("Job ID: %s, Task: %s failed on attempt %d. Retrying...", job.ID, job.Task, job.Attempts)

		// Add delay before retrying
		if job.Attempts < MaxRetries {
			time.Sleep(time.Second * time.Duration(job.Attempts)) // Exponential backoff simulation
		}
	}

	if job.Status != entity.COMPLETED {
		job.Status = entity.FAILED
		logrus.Errorf("Job ID: %s, Task: %s ultimately failed after %d attempts", job.ID, job.Task, job.Attempts)
	}

	// Final save status
	if err := q.jobRepo.Save(ctx, job); err != nil {
		logrus.Errorf("Failed to save final status for job %s: %v", job.ID, err)
	}
}

// executeTask simulates the actual work being done
func (q *jobService) executeTask(taskName string, attempt int32) bool {
	if taskName == "unstable-job" {
		return attempt > 2 // Fails on attempt 1 and 2, passes on 3
	}
	// All other jobs pass immediately
	return true
}

// Enqueue adds a new job to the queue
func (q *jobService) Enqueue(ctx context.Context, taskName string) (string, error) {

	newJob := entity.NewJob(taskName)

	// Save job to repository with PENDING status
	if err := q.jobRepo.Save(ctx, newJob); err != nil {
		logrus.Errorf("Failed to save job to repository: %v", err)
		return "", fmt.Errorf("failed to save job: %w", err)
	}

	// Add job to the worker queue
	q.jobQueue <- newJob

	logrus.Infof("Job %s with task '%s' enqueued", newJob.ID, taskName)
	return newJob.ID, nil
}

// GetAllJobs retrieves all jobs from the repository
func (q *jobService) GetAllJobs(ctx context.Context) ([]*entity.Job, error) {
	return q.jobRepo.FindAll(ctx)
}

// GetJobByID retrieves a specific job from the repository
func (q *jobService) GetJobByID(ctx context.Context, id string) (*entity.Job, error) {
	return q.jobRepo.FindByID(ctx, id)
}

// GetJobStatus calculates and returns the counts for each job status
func (q *jobService) GetJobStatus(ctx context.Context) (entity.JobStatus, error) {
	jobs, err := q.jobRepo.FindAll(ctx)
	if err != nil {
		return entity.JobStatus{}, err
	}

	statusCounts := entity.JobStatus{}

	for _, job := range jobs {
		switch job.Status {
		case entity.PENDING:
			statusCounts.Pending++
		case entity.RUNNING:
			statusCounts.Running++
		case entity.FAILED:
			statusCounts.Failed++
		case entity.COMPLETED:
			statusCounts.Completed++
		}
	}

	return statusCounts, nil
}

// NewJobService ...
func NewJobService() Initiator {
	return func(s *jobService) *jobService {
		s.jobQueue = make(chan *entity.Job, 100) // Buffer size 100
		s.startWorkerPool()                      // Start the workers immediately
		return s
	}
}

// SetJobRepository ...
func (i Initiator) SetJobRepository(jobRepository _interface.JobRepository) Initiator {
	return func(s *jobService) *jobService {
		i(s).jobRepo = jobRepository
		return s
	}
}

// Build ...
func (i Initiator) Build() _interface.JobService {
	return i(&jobService{})
}
