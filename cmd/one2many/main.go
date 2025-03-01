package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/BeautyyuYanli/tcp_multi/pkg/config"
	"github.com/BeautyyuYanli/tcp_multi/pkg/logger"
	"github.com/BeautyyuYanli/tcp_multi/pkg/proxy"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.New(logger.Config{
		Level:    cfg.Logging.Level,
		Format:   cfg.Logging.Format,
		Output:   cfg.Logging.Output,
		FilePath: cfg.Logging.FilePath,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Create and start proxy
	p := proxy.New(cfg, log)
	
	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		if err := p.Start(); err != nil {
			log.Error("Failed to start proxy: %v", err)
			os.Exit(1)
		}
	}()

	log.Info("TCP proxy started. Press Ctrl+C to exit.")
	
	// Wait for termination signal
	<-sigCh
	log.Info("Shutting down...")
	
	// Stop the proxy
	if err := p.Stop(); err != nil {
		log.Error("Error during shutdown: %v", err)
	}
	
	log.Info("TCP proxy stopped")
}
