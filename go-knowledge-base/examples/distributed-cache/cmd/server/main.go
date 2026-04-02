package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"distributed-cache/internal/cache"
	"distributed-cache/internal/ring"
	"distributed-cache/internal/server"
)

func main() {
	var (
		nodeID   = flag.String("node-id", "", "Unique node ID")
		addr     = flag.String("addr", ":8080", "HTTP server address")
		maxSize  = flag.String("max-size", "1GB", "Maximum cache size")
		vnodes   = flag.Int("vnodes", 150, "Virtual nodes per physical node")
	)
	flag.Parse()
	
	if *nodeID == "" {
		*nodeID = generateNodeID()
	}
	
	log.Printf("Starting cache node %s on %s", *nodeID, *addr)
	
	// Parse max size
	maxBytes := parseSize(*maxSize)
	
	// Create cache
	cacheConfig := &cache.Config{
		MaxSize:        maxBytes,
		EvictionPolicy: cache.LRU,
		DefaultTTL:     time.Hour,
	}
	cacheInstance := cache.New(cacheConfig)
	
	// Create ring for cluster awareness
	ringInstance := ring.New(*vnodes)
	
	// Create and start server
	srv := server.New(cacheInstance, ringInstance, *nodeID)
	
	httpServer := &http.Server{
		Addr:         *addr,
		Handler:      srv.Router(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	
	// Start server in goroutine
	go func() {
		log.Printf("HTTP server listening on %s", *addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	
	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down server...")
	
	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Server exited")
}

func generateNodeID() string {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}
	return hostname + "-" + time.Now().Format("20060102150405")
}

func parseSize(s string) int64 {
	// Simple size parser - in production use a proper library
	var multiplier int64 = 1
	if len(s) > 2 {
		switch s[len(s)-2:] {
		case "KB":
			multiplier = 1024
			s = s[:len(s)-2]
		case "MB":
			multiplier = 1024 * 1024
			s = s[:len(s)-2]
		case "GB":
			multiplier = 1024 * 1024 * 1024
			s = s[:len(s)-2]
		case "TB":
			multiplier = 1024 * 1024 * 1024 * 1024
			s = s[:len(s)-2]
		}
	}
	
	var size int64
	if _, err := log.Writer().Write([]byte(s)); err == nil {
		// Try to parse
		var n int
		if _, err := log.Writer().Write([]byte{}); err == nil {
			_ = n
		}
	}
	
	// Default to 1GB if parsing fails
	size = 1024 * 1024 * 1024
	
	return size * multiplier
}
