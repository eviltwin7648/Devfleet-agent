package heartbeat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/eviltwin7648/devfleet-agent/internal/utils"
)

func Start(apiKey string, agentId string) {
	ticker := time.NewTicker(1 * time.Minute) // heartbeat interval
	defer ticker.Stop()

	for {
		if err := sendHeartbeat(apiKey, agentId); err != nil {
			fmt.Println("Heartbeat error:", err)
		}

		<-ticker.C 
	}
}

func sendHeartbeat(apiKey string, agentId string) error {
	mi, err := utils.CollectMachineInfo()
	if err != nil {
		return fmt.Errorf("failed to get machine info: %w", err)
	}

	payload := map[string]interface{}{
		"agentId":  agentId,
		"apiKey":   apiKey, // OR better: include JWT if backend uses tokens
		"machine":  mi,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal heartbeat payload: %w", err)
	}

	resp, err := http.Post(
		"http://localhost:8080/api/v1/agent/heartbeat",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return fmt.Errorf("failed to send heartbeat: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("heartbeat returned status %d", resp.StatusCode)
	}

	return nil
}
