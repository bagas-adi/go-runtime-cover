# Go Runtime Coverage Dumper

This project demonstrates how to use Go's runtime coverage functionality to dump coverage data at runtime, inspired by JaCoCo's dump feature in Java.

## Project Structure

- **Integration Example** (`examples/integration`):
  - Runs an HTTP server and demonstrates runtime coverage dumping.
  - Dumps coverage data when a trigger file is created.
- **Package Injector** (`cmd/pkg_injector`):
  - Utility to inject import statements and code snippets into Go files.
- **Coverage Dumper Library** (`pkg/coverage`):
  - Provides the `Dumper` type for handling coverage data dumps.

## Building

To build the integration and injector binaries:

```bash
# Build the integration example
GO111MODULE=on go build -cover -covermode=atomic -o bin/integration ./examples/integration

# Build the package injector
GO111MODULE=on go build -o bin/pkg_injector ./cmd/pkg_injector
```

## Usage

### Running the Integration Example

The integration binary starts an HTTP server and watches for a file named `generate_coverage` in the current directory. When this file is created, a coverage dump is triggered.

```bash
# Start the integration service (coverage data will be dumped to ./coverage-data/<timestamp> by default)
./bin/integration -coverage-dir ./coverage-data/<timestamp>
```

- The service listens on port 8080 by default.
- Coverage is dumped when the file `generate_coverage` is created in the working directory.
- Test endpoints:
  - `GET /test` - Test endpoint with coverage
  - `POST /test` - Test endpoint with different coverage

#### Using the Helper Script

You can use the provided script to build and start the integration service:

```bash
./scripts/start_services.sh
```

This script will:
- Build the integration binary
- Create a timestamped coverage data directory
- Start the integration service in the background
- Print the endpoints and coverage directory

To stop the service:
```bash
./scripts/stop_services.sh
```

### Using the Package Injector

The injector can add import statements and inject code snippets into Go files. This allows you to add coverage dumping logic on-the-fly to any of your Go code, making it easy to enable runtime coverage without manual editing.

Example usage:

```bash
./bin/pkg_injector -filename path/to/your.go -pkg go-runtime-cover/pkg/coverage
```

## Coverage Dumping Mechanism

- Coverage data is dumped when the file `generate_coverage` is created in the working directory.
- Output files are written to the specified coverage directory (e.g., `./coverage-data/<timestamp>`):
  - Meta data files: `covmeta.<hexid>`
  - Counter data files: `covcounters.<hexid>.<counter>.<timestamp>`

## Coverage Monitoring Script

You can monitor your runtime coverage data in real time using the provided script:

```bash
./scripts/get_coverage.sh
```

This script will:
- Continuously watch the latest coverage data directory (e.g., `./coverage-data/20250630_181441/`)
- Print the current coverage percentage for each package every second

**Sample Output:**
```
[2025-06-30 18:15:46]: (./coverage-data/20250630_181441/) Waiting for coverage data...
        github.com/bagas-adi/go-runtime-cover/examples/integration    coverage: 86.7% of statements
        github.com/bagas-adi/go-runtime-cover/examples/server         coverage: 89.5% of statements
        github.com/bagas-adi/go-runtime-cover/pkg/coverage           coverage: 65.0% of statements
[2025-06-30 18:15:47]: (./coverage-data/20250630_181441/) Waiting for coverage data...
        github.com/bagas-adi/go-runtime-cover/examples/integration    coverage: 86.7% of statements
        github.com/bagas-adi/go-runtime-cover/examples/server         coverage: 89.5% of statements
        github.com/bagas-adi/go-runtime-cover/pkg/coverage           coverage: 65.0% of statements
[2025-06-30 18:15:48]: (./coverage-data/20250630_181441/) Waiting for coverage data...
        github.com/bagas-adi/go-runtime-cover/examples/integration    coverage: 86.7% of statements
        github.com/bagas-adi/go-runtime-cover/examples/server         coverage: 89.5% of statements
        github.com/bagas-adi/go-runtime-cover/pkg/coverage           coverage: 65.0% of statements
```

- The script updates every second, showing the latest coverage for each package.
- This is useful for observing how your coverage changes as you interact with your running application.

## Stopping Coverage Monitoring (`Ctrl+C`)

When you stop the coverage monitoring script (`./scripts/get_coverage.sh`) by pressing `Ctrl+C`, the script will generate a detailed coverage report, showing coverage percentages for each function in your codebase. This helps you identify which parts of your code are well-tested and which are not.

**Sample Output:**
```
github.com/bagas-adi/go-runtime-cover/examples/integration/main.go:11:        main                            86.7%
github.com/bagas-adi/go-runtime-cover/examples/server/server.go:17:     NewExampleService               100.0%
github.com/bagas-adi/go-runtime-cover/examples/server/server.go:51:     processRequest                  100.0%
github.com/bagas-adi/go-runtime-cover/examples/server/server.go:66:     Start                           75.0%
github.com/bagas-adi/go-runtime-cover/examples/server/server.go:78:     Stop                            0.0%
github.com/bagas-adi/go-runtime-cover/pkg/coverage/dumper.go:16:        generateRandomHexID             75.0%
github.com/bagas-adi/go-runtime-cover/pkg/coverage/dumper.go:24:        generateRandomCounter           80.0%
github.com/bagas-adi/go-runtime-cover/pkg/coverage/dumper.go:40:        NewDumper                       71.4%
github.com/bagas-adi/go-runtime-cover/pkg/coverage/dumper.go:67:        Dump                            77.8%
github.com/bagas-adi/go-runtime-cover/pkg/coverage/dumper.go:87:        dumpMeta                        71.4%
github.com/bagas-adi/go-runtime-cover/pkg/coverage/dumper.go:102:       dumpCounters                    71.4%
github.com/bagas-adi/go-runtime-cover/pkg/coverage/dumper.go:117:       ClearCounters                   0.0%
github.com/bagas-adi/go-runtime-cover/pkg/coverage/dumper.go:124:       WatchForFileAndDumpCoverage     62.5%
github.com/bagas-adi/go-runtime-cover/pkg/coverage/dumper.go:141:       PrintLog                        0.0%
total:                                                                  (statements)                    73.4%
```

- Each line shows the file, line number, function name, and the percentage of statements covered by tests or runtime execution.
- The final line shows the total statement coverage across your codebase.

This summary is useful for quickly identifying coverage gaps and tracking improvements as you test your application.

## Graceful Shutdown

The application is designed to support graceful shutdown and can be extended to handle SIGINT/SIGTERM for a final coverage dump (see commented code in `examples/integration/main.go`).

## Important Notes

- Always use `-covermode=atomic` when building applications that use runtime coverage.
- Coverage data is only collected for code that is executed.
- Make sure to exercise your application's endpoints to collect meaningful coverage data.

## Snippets File

A code snippet is provided in `cmd/pkg_injector/snippets_sigterm.go`:

```go
	dumper, err := coverage.NewDumper()
	if err != nil {
		dumper.PrintLog(err)
	}
	dumper.WatchForFileAndDumpCoverage("./generate_coverage")
```

This snippet initializes the coverage dumper and sets up file-based coverage dumping. You can use the package injector to automatically insert this snippet after the `func main(){` line in your Go files. For example:

```bash
./bin/pkg_injector -filename path/to/your.go -injectfile cmd/pkg_injector/snippets_sigterm.go
```

- The `-injectfile` flag specifies the path to the snippet file to inject.

This will add the snippet after the `main` function declaration, enabling runtime coverage dumping in your application. 