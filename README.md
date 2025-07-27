# cc-init

A Go command-line tool that quickly initializes Claude Code configuration for projects by copying embedded `.claude` template files to create the necessary directory structure.

## Features

- 🚀 **One-command setup** - Initialize Claude Code configuration instantly
- 📁 **Template embedding** - Templates are embedded in the binary using Go's `embed` package
- 🛡️ **Safe operations** - Skips existing files and directories to prevent overwrites
- 🔍 **Dry-run mode** - Preview changes before applying them
- 📊 **Progress tracking** - Colored output with detailed operation summaries
- 🎯 **Flexible targeting** - Initialize in current directory or specify target location
- 🖥️ **Cross-platform** - Works on Windows, macOS, and Linux

## Installation

## Download pre-built binaries

You can download pre-built binaries for your platform from the [releases page](https://github.com/ipfans/cc-init/releases).

## Build from source

```bash
go install github.com/ipfans/cc-init@latest
```

## Usage

### Basic usage

```bash
# Initialize .claude configuration in current directory
./cc-init

# Initialize in a specific directory
./cc-init --target /path/to/your/project

# Preview what would be created (dry-run)
./cc-init --dry-run

# Verbose output for debugging
./cc-init --verbose

# Show version
./cc-init --version

# Show help
./cc-init --help
```

### Command-line options

| Flag         | Short | Description                                   |
| ------------ | ----- | --------------------------------------------- |
| `--target`   | `-t`  | Target directory (default: current directory) |
| `--dry-run`  |       | Preview operations without making changes     |
| `--verbose`  | `-v`  | Enable verbose output for debugging           |
| `--no-color` |       | Disable colored output                        |
| `--version`  |       | Show version information                      |
| `--help`     | `-h`  | Show help message                             |

### Example output

```bash
$ ./cc-init --target myproject --verbose
✓ Created directory: .claude/commands
✓ Created file: .claude/commands/ask.md
✓ Created file: .claude/commands/spec-workflow.md
✓ Created file: .claude/settings.local.json

✓ Created 3 files and 1 directory
✓ Claude configuration initialized successfully
```

## What gets created

The tool creates a `.claude` directory structure with:

```
.claude/
├── commands/
│   ├── ask.md              # Custom ask command template
│   └── spec-workflow.md    # Specification workflow template
└── settings.local.json     # Local Claude Code settings
```

## Development

See [CLAUDE.md](CLAUDE.md) for detailed development instructions and architecture notes.

### Building

```bash
# Build the project
go build

# Build with specific output name
go build -o cc-init

# Run directly
go run .
```

### Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## Technical Details

- **Go Version**: 1.24.5
- **Module**: `github.com/ipfans/cc-init`
- **Architecture**: Component-based with clean separation of concerns
- **Security**: Path validation, input sanitization, and safe file operations
- **Testing**: Comprehensive test suite with 90%+ coverage target

## License

See [LICENSE](LICENSE) file for details.
