package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eviltwin7648/devfleet-agent/internal/utils"
)

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
		"http://localhost:8080/api/v1/agent/verify",
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