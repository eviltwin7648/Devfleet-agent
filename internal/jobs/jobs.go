package jobs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/eviltwin7648/devfleet-agent/internal/utils"
)

func StartPolling(token string, agentId string) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		poll(token, agentId)
		<-ticker.C
	}
	fmt.Println("Polling for jobs...")
}

func poll(token string, agentId string) {
	req, err := http.NewRequest("GET", "http://localhost:8080/api/v1/agent/jobs/pull", nil)
	if err != nil {
		fmt.Println("Error creating poll request:", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error while polling for jobs", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Polling for jobs failed with status:", resp.Status)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Parse the response which is { "job": { ... } } or { "job": null }
	var respData struct {
		Job *utils.Job `json:"job"`
	}
	if err := json.Unmarshal(body, &respData); err != nil {
		fmt.Println("Error parsing job response:", err)
		return
	}

	if respData.Job != nil {
		fmt.Println("Found job:", respData.Job.ID, respData.Job.Script)
		result := utils.RunJob(*respData.Job)
		fmt.Printf("Job finished with status: %s, Exit Code: %d\n", result.Status, result.ExitCode)
		
		if err := reportJobResult(token, respData.Job.ID, result); err != nil {
			fmt.Println("Failed to report job result:", err)
		} else {
             fmt.Println("Job result reported successfully.")
        }
	} else {
		fmt.Println("No jobs pending.")
	}
}

func reportJobResult(token string, jobID string, result utils.JobResult) error {
	payload := map[string]interface{}{
		"status":    result.Status, // "SUCCESS", "FAILED"
		"exit_code": result.ExitCode,
        "stdout": result.Stdout, // Sending logs optionally
        "stderr": result.Stderr,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	url := fmt.Sprintf("http://localhost:8080/api/v1/agent/job/%s/result", jobID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("backend returned status %s: %s", resp.Status, string(body))
	}

	return nil
}

