# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

cc-init is a Go command-line tool that quickly initializes Claude Code configuration for projects. It copies embedded `.claude` template files to a target directory, creating the necessary directory structure for Claude Code integration. The tool includes features like dry-run mode, verbose output, and intelligent file skipping to avoid overwriting existing configurations.

**Module**: `github.com/ipfans/cc-init`
**Go Version**: 1.24.5

## Development Commands

### Go Module Management
```bash
# Download dependencies
go mod download

# Tidy and verify dependencies
go mod tidy

# Verify module dependencies
go mod verify
```

### Building and Running
```bash
# Build the project
go build

# Run the project
go run .

# Build with specific output name
go build -o cc-init
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Code Quality
```bash
# Format code
go fmt ./...

# Run go vet for static analysis
go vet ./...

# Run golangci-lint (if installed)
golangci-lint run
```

## Architecture Notes

cc-init is a fully functional command-line tool with a clean, component-based architecture:

### Core Components

- **CLI Interface** (`cli.go`): Command-line argument parsing with flags for dry-run, verbose, target directory, version, and help
- **Engine** (`engine.go`): Core orchestration logic that coordinates all operations and tracks statistics
- **Template Manager** (`template.go`): Handles embedded template files using Go's `embed` package
- **File System Abstraction** (`fs.go`): Provides testable file operations with OS and dry-run implementations
- **Logger** (`logger.go`): Colored output with different verbosity levels and ANSI color support

### Key Features

- **Embedded Templates**: Uses `//go:embed` to include `.claude` template files in the binary
- **Directory Structure Preservation**: Maintains the complete `.claude` directory hierarchy in the target location
- **Intelligent File Handling**: Skips existing files and directories to prevent accidental overwrites
- **Dry Run Mode**: Preview operations without making actual changes
- **Cross-Platform Support**: Works on Windows, macOS, and Linux with proper path handling
- **User-Friendly Output**: Colored status messages with progress tracking and summary reports

### Usage Examples

```bash
# Initialize .claude in current directory
./cc-init

# Initialize in specific directory with dry-run
./cc-init --target /path/to/project --dry-run

# Verbose output for debugging
./cc-init --verbose

# Show help and available options
./cc-init --help

# Show version information
./cc-init --version

# Disable colored output
./cc-init --no-color
```

## Template Structure

The tool creates the following `.claude` directory structure:

```
.claude/
├── commands/
│   ├── ask.md              # Custom ask command template
│   └── spec-workflow.md    # Specification workflow template  
└── settings.local.json     # Local Claude Code settings
```

## Security Features

- **Path Validation**: Prevents directory traversal attacks with comprehensive path sanitization
- **Safe File Operations**: Skips existing files to prevent accidental overwrites
- **Input Validation**: Validates all command-line inputs and target directories
- **Embedded Templates**: Templates are embedded in the binary, eliminating external dependencies