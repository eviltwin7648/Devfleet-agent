package cmd

import (
	"fmt"
	"os"

	"github.com/eviltwin7648/devfleet-agent/internal/auth"
	"github.com/eviltwin7648/devfleet-agent/internal/config"
	"github.com/eviltwin7648/devfleet-agent/internal/heartbeat"
	"github.com/eviltwin7648/devfleet-agent/internal/jobs"
	"github.com/spf13/cobra"
)

var bootstrapToken string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the DevFleet agent",
	Run: func(cmd *cobra.Command, args []string) {
		if bootstrapToken != "" {
			fmt.Println("Registering agent using provided API key...")
			data, err := auth.RegisterAgent(bootstrapToken)
			if err != nil {
				fmt.Println("Registration failed:", err)
				os.Exit(1)
			}
			if err := config.SaveKey(bootstrapToken, data.AgentID); err != nil {
				fmt.Println("Failed to save key:", err)
				os.Exit(1)
			}
			fmt.Println("Registration successful. Starting agent...")
		}

		token, err := config.LoadKey()
		if err != nil {
			fmt.Println("No auth token found. Run `devfleet-agent login` first.")
			os.Exit(1)
		}
		// Verify and get JWT
		jwtToken, err := auth.VerifyAgent(token.APIKey)
		if err != nil {
			fmt.Println("Authentication failed:", err)
			os.Exit(1)
		}

		fmt.Println("Authentication successful. Running agent...")

		// Start your loops:
		go heartbeat.Start(jwtToken, token.AgentID)
		go jobs.StartPolling(jwtToken, token.AgentID)

		select {} // keep running
	},
}

func init() {
	startCmd.Flags().StringVar(&bootstrapToken, "token", "", "agent API key used for one-time bootstrap registration")
	rootCmd.AddCommand(startCmd)
}
