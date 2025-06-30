package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// ExampleService represents a simple HTTP service
type ExampleService struct {
	server *http.Server
}

// NewExampleService creates a new service instance
func NewExampleService(port int) *ExampleService {
	mux := http.NewServeMux()

	// Add test endpoints
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		// This function will be covered when the /test endpoint is called
		result := processRequest(r)
		fmt.Fprintf(w, "Processed: %s", result)
	})

	// mux.HandleFunc("/coverage", func(w http.ResponseWriter, r *http.Request) {
	// 	// Force a coverage dump
	// 	if err := dumper.Dump(); err != nil {
	// 		http.Error(w, fmt.Sprintf("Failed to dump coverage: %v", err), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	fmt.Fprintf(w, "Coverage data dumped successfully")
	// })

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return &ExampleService{
		server: server,
	}
}

// processRequest is a function that will be covered when /test is called
func processRequest(r *http.Request) string {
	// Simulate some processing
	time.Sleep(100 * time.Millisecond)

	// Add some branching logic for coverage
	if r.Method == "POST" {
		return "POST request processed"
	} else if r.Method == "GET" {
		return "GET request processed"
	} else {
		return "Unknown method processed"
	}
}

// Start begins the service
func (s *ExampleService) Start() error {
	// Start the server in a goroutine
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	return nil
}

// Stop gracefully shuts down the service
func (s *ExampleService) Stop(ctx context.Context) error {
	// Perform final coverage dump
	// if err := s.dumper.Dump(); err != nil {
	// 	log.Printf("Error during final coverage dump: %v", err)
	// }

	// Shutdown the server
	return s.server.Shutdown(ctx)
}
