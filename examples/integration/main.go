package main

import (
	"log"
	"time"

	"go-runtime-cover/pkg/coverage"
	"go-runtime-cover/server"
)

func main() {
	dumper, err := coverage.NewDumper()
	if err != nil {
		dumper.PrintLog(err)
	}
	// dumper.SetSIGTERM(15)
	dumper.WatchForFileAndDumpCoverage("./generate_coverage")
	// Parse command line arguments
	port := 8080
	// dumpInterval := flag.Duration("dump-interval", 5*time.Minute, "Interval between coverage dumps")
	// flag.Parse()

	// Create and start the service
	service := server.NewExampleService(port)
	if err := service.Start(); err != nil {
		log.Fatalf("Failed to start service: %v", err)
	}

	// rest := rest.NewRestMux()
	// rest.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Hello, World!")
	// })

	// Set up periodic coverage dump
	// ticker := time.NewTicker(*dumpInterval)
	// defer ticker.Stop()

	// Set up signal handling
	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Service started on port %d", port)
	// log.Printf("Coverage dumps will be saved to %s every %v", *outputDir, *dumpInterval)
	log.Printf("Test endpoints:")
	log.Printf("  GET /test - Test endpoint with coverage")
	log.Printf("  POST /test - Test endpoint with different coverage")
	log.Printf("  GET /coverage - Force a coverage dump")

	// Main loop
	for {
		time.Sleep(1 * time.Second)
		// select {
		// case <-ticker.C:
		// 	if err := dumper.Dump(); err != nil {
		// 		log.Printf("Failed to dump coverage data: %v", err)
		// 	} else {
		// 		log.Println("Successfully dumped coverage data")
		// 	}

		// case sig := <-sigChan:
		// 	log.Printf("Received signal: %v", sig)

		// 	// Create shutdown context with timeout
		// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// 	defer cancel()

		// 	// Stop the service
		// 	if err := service.Stop(ctx); err != nil {
		// 		log.Printf("Error during shutdown: %v", err)
		// 	}

		// 	return
		// default:
		// 	time.Sleep(1 * time.Second)
		// }
	}
}
