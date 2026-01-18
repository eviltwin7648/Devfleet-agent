package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/eviltwin7648/devfleet-agent/internal/auth"
	"github.com/eviltwin7648/devfleet-agent/internal/config"
	"github.com/eviltwin7648/devfleet-agent/internal/heartbeat"
	"github.com/eviltwin7648/devfleet-agent/internal/jobs"
)

var startCmd = &cobra.Command{
    Use:   "start",
    Short: "Start the DevFleet agent",
    Run: func(cmd *cobra.Command, args []string) {

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
    rootCmd.AddCommand(startCmd)
}
