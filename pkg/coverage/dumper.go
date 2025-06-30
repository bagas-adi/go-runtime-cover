package coverage

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/coverage"
	"syscall"
	"time"
)

func generateRandomHexID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func generateRandomCounter() (int64, error) {
	max := big.NewInt(10_000_000) // Example: random counter < 10 million
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return n.Int64(), nil
}

// Dumper handles coverage data dumping operations
type Dumper struct {
	outputDir string
	hexID     string
}

// NewDumper creates a new coverage dumper
func NewDumper() (*Dumper, error) {
	outputDir := flag.String("coverage-dir", "", "Directory to store code coverage")
	flag.Parse()

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	hexID, err := generateRandomHexID()
	if err != nil {
		panic("failed to generate hex ID: " + err.Error())
	}

	dumper := &Dumper{outputDir: *outputDir, hexID: hexID}
	if err := dumper.Dump(); err != nil {
		return nil, fmt.Errorf("failed to dump coverage data: %w", err)
	}
	// Dump meta data
	metaFile := filepath.Join(dumper.outputDir, fmt.Sprintf("covmeta.%s", dumper.hexID))
	if err := dumper.dumpMeta(metaFile); err != nil {
		return nil, fmt.Errorf("failed to dump meta data: %w", err)
	}

	return dumper, nil
}

// Dump writes both meta and counter data to the output directory
func (d *Dumper) Dump() error {
	counter, err := generateRandomCounter()
	if err != nil {
		panic("failed to generate counter: " + err.Error())
	}

	timestamp := time.Now().UnixNano()

	counterKey := fmt.Sprintf("covcounters.%s.%d.%d", d.hexID, counter, timestamp)

	// Dump counter data
	counterFile := filepath.Join(d.outputDir, counterKey)
	if err := d.dumpCounters(counterFile); err != nil {
		return fmt.Errorf("failed to dump counter data: %w", err)
	}

	return nil
}

// dumpMeta writes coverage meta data to a file
func (d *Dumper) dumpMeta(filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create meta file: %w", err)
	}
	defer f.Close()

	if err := coverage.WriteMeta(f); err != nil {
		return fmt.Errorf("failed to write meta data: %w", err)
	}

	return nil
}

// dumpCounters writes coverage counter data to a file
func (d *Dumper) dumpCounters(filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create counter file: %w", err)
	}
	defer f.Close()

	if err := coverage.WriteCounters(f); err != nil {
		return fmt.Errorf("failed to write counter data: %w", err)
	}

	return nil
}

// ClearCounters resets all coverage counters
func (d *Dumper) ClearCounters() error {
	if err := coverage.ClearCounters(); err != nil {
		return fmt.Errorf("failed to clear counters: %w", err)
	}
	return nil
}

func (d *Dumper) SetSIGTERM(terminationDuration time.Duration) {
	signalChan := make(chan os.Signal, 1)
	// change SIGTERM with the other SIGNAL, or any other event like a creation of a file
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	termDuration := terminationDuration * time.Second
	go func() {
		<-signalChan
		d.Dump() // Final dump before exit
		log.Printf("Terminate service in %v", termDuration.String())
		time.Sleep(termDuration)
		os.Exit(0)
	}()
}

func (d *Dumper) PrintLog(err error) {
	log.Printf("outputDir: %v", d.outputDir)
	log.Printf("dumper: %v", d)
	log.Fatalf("Failed to create coverage dumper: %v", err)
}
