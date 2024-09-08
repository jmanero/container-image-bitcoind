package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmanero/container-image-bitcoind/monitor/client"
	"github.com/jmanero/container-image-bitcoind/monitor/exporter"
	"github.com/spf13/cobra"
)

// CLI root
var CLI = cobra.Command{
	Use:  "bitcoind-monitor",
	RunE: Check,
}

// Result return code for health checks
var Result int

func init() {
	CLI.AddCommand(&cobra.Command{
		Use:   "startup",
		Short: "startup health-check",
		Long:  "startup waits for the node's getblockchaininfo.initializing flag to become false. This implies that the node's RPC service is available",
		RunE:  Startup,
	})

	CLI.AddCommand(&cobra.Command{
		Use:   "check",
		Short: "runtime health-check",
		Long:  "check that the node's RPC service is available",
		RunE:  Check,
	})

	CLI.AddCommand(&cobra.Command{
		Use:   "serve",
		Short: "Serve a Prometheus exporter for chain and mempool information",
		RunE:  exporter.Serve,
	})
}

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	if CLI.ExecuteContext(ctx) != nil {
		os.Exit(1)
	}

	os.Exit(Result)
}

// Startup health-check
func Startup(cmd *cobra.Command, args []string) error {
	info, err := client.GetBlockchainInfo(cmd.Context())
	if err != nil {
		return err
	}

	if info.Initializing {
		cmd.PrintErrf("Node is performing initial block synchronization: %d blocks/%d headers/%f%%\n", info.Blocks, info.Headers, info.Progress*100)
		Result = 100

		return nil
	}

	cmd.Printf("Node is running: %d blocks/%d headers/%f%%\n", info.Blocks, info.Headers, info.Progress*100)
	return nil
}

// Check that the node's RPC service is available
func Check(cmd *cobra.Command, args []string) error {
	info, err := client.GetBlockchainInfo(cmd.Context())
	if err != nil {
		cmd.PrintErrln("Unable to connect to node", client.Endpoint.String())
		Result = 101

		return nil
	}

	cmd.Printf("Node is running: %d blocks/%d headers/%f%%\n", info.Blocks, info.Headers, info.Progress*100)
	return nil
}
