package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/nccapo/kv-store/internal/api"
	"github.com/nccapo/kv-store/internal/raft"
	"github.com/nccapo/kv-store/store"
)

func main() {
	// Parse command line flags
	nodeID := flag.String("id", "", "Node ID")
	addr := flag.String("addr", ":8080", "HTTP server address")
	peers := flag.String("peers", "", "Comma-separated list of peer addresses")
	flag.Parse()

	if *nodeID == "" {
		log.Fatal("Node ID is required")
	}

	// Create Raft node
	node := raft.NewNode(*nodeID, parsePeers(*peers))

	// Create key-value store
	kvStore := store.NewClient(&store.ClientOptions{
		Addr: *addr,
	})

	// Create HTTP API server
	server := api.NewServer(kvStore, *addr)

	// Start Raft node
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := node.Start(ctx); err != nil {
		log.Fatalf("Failed to start Raft node: %v", err)
	}

	// Start HTTP server
	go func() {
		if err := server.Start(); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// Graceful shutdown
	log.Println("Shutting down...")
	if err := server.Stop(ctx); err != nil {
		log.Printf("Error stopping HTTP server: %v", err)
	}
	node.Stop()
}

func parsePeers(peers string) map[string]string {
	if peers == "" {
		return make(map[string]string)
	}

	result := make(map[string]string)
	for _, peer := range strings.Split(peers, ",") {
		parts := strings.Split(peer, "=")
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}
