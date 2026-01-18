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
	JobTimeout JobStatus = "Timeout"
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

type Job struct {
	ID          string            `json:"id"`
	Script      string            `json:"script"`
	Env         map[string]string `json:"env"`
	TimeoutSec  int               `json:"timeoutSec"`
}

func RunJob(job Job) JobResult {
	start := time.Now()

	// Default timeout 10 minutes if not specified
	timeout := 10 * time.Minute
	if job.TimeoutSec > 0 {
		timeout = time.Duration(job.TimeoutSec) * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", job.Script)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Set environment variables
	cmd.Env = os.Environ()
	for k, v := range job.Env {
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
