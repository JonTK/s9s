// Package main provides the entry point for the s9s command-line interface.
package main

import (
	"os"

	"github.com/jontk/s9s/internal/cli"
	"github.com/jontk/s9s/internal/logging"
)

func main() {
	// Initialize structured logging early
	logging.Init(logging.DefaultConfig())
	logger := logging.GetLogger()

	if err := cli.Execute(); err != nil {
		logger.Error().Err(err).Msg("Failed to execute command")
		os.Exit(1)
	}
}
