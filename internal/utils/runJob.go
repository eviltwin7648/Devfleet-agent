package utils

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type JobStatus string

const (
	JobSuccess JobStatus = "Success"
	JobFailed  JobStatus = "Failure"
	JobRunning JobStatus = "Running"
	JobTimeout JobStatus = "Timeout"
	JobCancelled JobStatus = "Cancelled"
)

type JobResult struct {
	ExitCode   int
	Stdout     string
	Stderr     string
	Duration   time.Duration
	Status     JobStatus
	StartedAt  time.Time
	FinishedAt time.Time
	Error      string
}

// JobDefinition holds the script and config, nested inside a JobExecution from the backend.
type JobDefinition struct {
	Script     string            `json:"script"`
	Env        map[string]string `json:"env"`
	TimeoutSec int               `json:"timeoutSec"`
}

// Job represents the JobExecution returned by the backend /jobs/pull endpoint.
// The script lives inside the nested JobDefinition (field "job").
type Job struct {
	ID         string        `json:"id"`
	Status     string        `json:"status"`
	Definition JobDefinition `json:"job"` // nested JobDefinition
}

func RunJob(job Job) JobResult {
	start := time.Now()

	// Default timeout 10 minutes if not specified
	timeout := 10 * time.Minute
	if job.Definition.TimeoutSec > 0 {
		timeout = time.Duration(job.Definition.TimeoutSec) * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fmt.Printf("[RunJob] Executing script for job %s: %q\n", job.ID, job.Definition.Script)
	cmd := exec.CommandContext(ctx, "bash", "-c", job.Definition.Script)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Set environment variables
	cmd.Env = os.Environ()
	for k, v := range job.Definition.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	if err := cmd.Start(); err != nil {
		return JobResult{
			Status:     JobFailed,
			Error:      err.Error(),
			StartedAt:  start,
			FinishedAt: time.Now(),
		}
	}

	err := cmd.Wait()
	end := time.Now()

	result := JobResult{
		Stdout:     stdout.String(),
		Stderr:     stderr.String(),
		StartedAt:  start,
		FinishedAt: end,
		Duration:   end.Sub(start),
	}

	if ctx.Err() == context.DeadlineExceeded {
		result.Status = JobTimeout
		result.Error = "job timed out"
		return result
	}

	if err != nil {
		result.Status = JobFailed
		if cmd.ProcessState != nil {
			result.ExitCode = cmd.ProcessState.ExitCode()
		}
		result.Error = err.Error()
		return result
	}

	result.Status = JobSuccess
	result.ExitCode = cmd.ProcessState.ExitCode()
	return result
}
