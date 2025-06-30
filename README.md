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