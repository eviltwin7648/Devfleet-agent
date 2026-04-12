package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/eviltwin7648/devfleet-agent/internal/utils"
)

const BackendURL = "http://localhost:8080"

type registerPayload struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Hostname string `json:"hostname"`
	TotalMem uint64 `json:"totalmem"`
	ApiKey   string `json:"apiKey"`
}

type registerResponse struct {
	Username string `json:"username"`
	AgentID  string `json:"agent_id"`
}

func RegisterAgent(apiKey string) (*registerResponse, error) {
	mi, err := utils.CollectMachineInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get machine info: %w", err)
	}

	payload := registerPayload{
		OS:       mi.OS,
		Arch:     mi.Arch,
		Hostname: mi.Hostname,
		TotalMem: mi.TotalMem,
		ApiKey:   apiKey,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request body: %w", err)
	}

	resp, err := http.Post(
		BackendURL+"/api/v1/agent/register",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("registration failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var data registerResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &data, nil
}

func VerifyAgent(apiKey string) (string, error) {
	mi, err := utils.CollectMachineInfo()
	if err != nil {
		return "", fmt.Errorf("failed to collect machine info: %w", err)
	}

	payload := map[string]interface{}{
		"apiKey":   apiKey,
		"hostname": mi.Hostname,
		"os":       mi.OS,
		"arch":     mi.Arch,
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal verify payload: %w", err)
	}

	resp, err := http.Post(
		BackendURL+"/api/v1/agent/verify",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return "", fmt.Errorf("error while verifying agent: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("agent verification failed with status: %s", resp.Status)
	}

	var result struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode verify response: %w", err)
	}

	fmt.Println("Agent Verified Successfully. Token received.")
	return result.Token, nil
}
