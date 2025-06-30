package main

import (
	"log"
	"time"

	"github.com/bagas-adi/go-runtime-cover/example/server"
)

func main() {
	port := 8080

	// Create and start the service
	service := server.NewExampleService(port)
	if err := service.Start(); err != nil {
		log.Fatalf("Failed to start service: %v", err)
	}

	log.Printf("Service started on port %d", port)
	log.Printf("Test endpoints:")
	log.Printf("  GET /test - Test endpoint with coverage")
	log.Printf("  POST /test - Test endpoint with different coverage")
	log.Printf("  GET /coverage - Force a coverage dump")

	// Main loop
	for {
		time.Sleep(1 * time.Second)
	}
}
