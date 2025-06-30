package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

func fileLines(filename string) ([]string, error) {

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	// Parse command line arguments
	importsToAdd := flag.String("pkg", "go-runtime-cover/pkg/coverage", "Packages to add")
	filename := flag.String("filename", "", "File to add imports to")
	injectFile := flag.String("injectfile", "cmd/pkg_injector/snippets_sigterm.go", "File to inject after func main(){")

	flag.Parse()
	if *filename == "" {
		fmt.Printf("-filename cannot be empty!\n\n")
		flag.CommandLine.Usage()
		return
	}
	lines, err := fileLines(*filename)
	if err != nil {
		fmt.Printf("File %s not found!\n", filename)
		return
	}

	// Check which imports are missing
	var missing []string
	var importsArray []string
	for _, imp := range strings.Split(*importsToAdd, ",") {
		importsArray = append(importsArray, imp)
	}

	// Find multi-line import block
	importStart, importEnd := -1, -1
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "import (") {
			importStart = i
			for j := i + 1; j < len(lines); j++ {
				if strings.TrimSpace(lines[j]) == ")" {
					importEnd = j
					break
				}
			}
			break
		}
	}

	if importStart != -1 && importEnd != -1 {
		// Checking if packages are missing in multi-line import block
		for _, imp := range importsArray {
			for i, line := range lines {
				if !strings.Contains(line, imp) && !slices.Contains(missing, imp) {
					missing = append(missing, imp)
				}
				if i == importEnd-1 {
					break
				}
			}
		}
	}

	if len(missing) == 0 {
		fmt.Println("All imports already present.")
		return
	}

	if importStart != -1 && importEnd != -1 {
		// Add missing imports to multi-line import block
		var buf bytes.Buffer
		for i, line := range lines {
			buf.WriteString(line + "\n")
			if i == importEnd-1 {
				for _, imp := range missing {
					buf.WriteString("    \"" + imp + "\"\n")
				}
			}
		}
		os.WriteFile(*filename, buf.Bytes(), 0644)
		fmt.Println("Added missing imports to multi-line import block.")
		time.Sleep(2 * time.Second)
	}

	// Re-read file after adding imports
	lines, err = fileLines(*filename)
	if err != nil {
		fmt.Printf("File %s not found!\n", filename)
		return
	}
	// Inject text after func main(){ if requested
	if *injectFile != "" {
		injectTextBytes, err := os.ReadFile(*injectFile)
		if err != nil {
			fmt.Printf("Failed to read inject file: %v\n", err)
			return
		}
		injectText := string(injectTextBytes)
		var buf bytes.Buffer
		injected := false
		for _, line := range lines {
			buf.WriteString(line + "\n")
			if !injected && strings.Contains(line, "func main()") && strings.HasSuffix(strings.TrimSpace(line), "{") {
				buf.WriteString(injectText + "\n")
				injected = true
			}
		}
		os.WriteFile(*filename, buf.Bytes(), 0644)
		fmt.Println("Injected file contents after func main(){.")
	}
}
