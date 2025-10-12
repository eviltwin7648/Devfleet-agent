package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/eviltwin7648/devfleet-agent/internal/config"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
)

type RegisterPayload struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Hostname string `json:"hostname"`
	TotalMem uint64 `json:"totalmem"`
	ApiKey   string `json:"apiKey"`
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the service",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter your Agent API Key: ")
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)

		if key == "" {
			return fmt.Errorf("API key cannot be empty")
		}

		hostname, _ := os.Hostname()
		memInfo, err := mem.VirtualMemory()
		if err != nil {
			return fmt.Errorf("could not get memory info: %w", err)
		}
		payload := RegisterPayload{
			OS:       runtime.GOOS,
			Arch:     runtime.GOARCH,
			Hostname: hostname,
			TotalMem: memInfo.Total,
			ApiKey:   key,
		}
		jsonBody, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("could not marshal request body: %w", err)
		}
		resp, err := http.Post("http://localhost:8080/api/v1/agent/register", "application/json", bytes.NewBuffer(jsonBody))

		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("registration failed with status %d: %s", resp.StatusCode, string(bodyBytes))
		}

		type ValidateResponse struct {
			Username string `json:"username"`
		}

		var data ValidateResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		if err := config.SaveKey(key); err != nil {
			return fmt.Errorf("failed to save key: %w", err)

		}
		fmt.Println("Welcome", data.Username)
		fmt.Println("API key saved successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
