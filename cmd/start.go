package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/eviltwin7648/devfleet-agent/internal/auth"
	"github.com/eviltwin7648/devfleet-agent/internal/config"
	"github.com/eviltwin7648/devfleet-agent/internal/heartbeat"
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
		fmt.Println("token", token)
        if !auth.VerifyAgent(token.APIKey) {
            fmt.Println("Authentication failed. Run `devfleet-agent login` again.")
            os.Exit(1)
        }

        fmt.Println("Authentication successful. Running agent...")

        // Start your loops later:
        go heartbeat.Start(token.APIKey, token.AgentID)
        go jobs.StartPolling(token.APIKey, token.AgentID)

        select {} // keep running
    },
}

func init() {
    rootCmd.AddCommand(startCmd)
}
