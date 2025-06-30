package coverage

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"runtime/coverage"
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

func (d *Dumper) WatchForFileAndDumpCoverage(triggerFile string) {
	go func() {
		for {
			if _, err := os.Stat(triggerFile); err == nil {
				// File exists, trigger coverage dump
				if err := d.Dump(); err != nil {
					log.Printf("Failed to dump coverage data: %v", err)
				} else {
					log.Println("Successfully dumped coverage data")
				}
				// Remove the trigger file after dumping
				os.Remove(triggerFile)
			}
			time.Sleep(1 * time.Second)
		}
	}()
}
func (d *Dumper) PrintLog(err error) {
	log.Printf("outputDir: %v", d.outputDir)
	log.Printf("dumper: %v", d)
	log.Fatalf("Failed to create coverage dumper: %v", err)
}
