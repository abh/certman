# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Testing
- `go test ./...` - Run all tests
- `go test -v ./...` - Run tests with verbose output

### Code Formatting
- `gofumpt -w .` - Format all Go files (required before commits)

### Pre-commit Requirements
1. Run `gofumpt -w .` on all modified Go files
2. Run `go test ./...` and ensure all tests pass
3. Remove trailing whitespace from edited lines

## Architecture

### Core Components

**CertMan Struct (`certman.go:31-40`)**
- Manages TLS certificate and key pairs with live reloading
- Uses `fsnotify` to watch for file system changes
- Thread-safe with `sync.RWMutex` for certificate access
- Context-based cancellation for graceful shutdown

**Key Features:**
- **File Watching**: Monitors certificate and key file directories for changes
- **Kubernetes Support**: Handles `..data` symlink changes for Kubernetes secrets mounting
- **Atomic Reloading**: Waits 1 second after detecting changes to handle non-atomic cert/key updates
- **Error Resilience**: Continues serving old certificates if new ones fail to load

### API Design

**Constructor with Functional Options (`certman.go:70-99`)**
```go
certman.New(certFile, keyFile, opts...)
```

**Available Options:**
- `WithLogger(*slog.Logger)` - Custom logger (defaults to `setupLogs()`)
- `WithContext(context.Context)` - Custom context (defaults to `context.Background()`)

**Default Logger (`certman.go:42-47`)**
- Uses `log/slog` with text handler to stderr
- Info level by default
- Created by private `setupLogs()` function

### Threading Model

**Main Goroutine (`certman.go:152-196`)**
- Runs in `Watch()` method via `go cm.run()`
- Uses ticker for delayed reloading (handles atomic file updates)
- Monitors three channels: context cancellation, timer events, fs events, and errors
- Graceful shutdown via context cancellation

### File System Integration

**Path Handling (`certman.go:73-81`)**
- Converts relative paths to absolute paths
- Watches parent directories of cert/key files
- Handles cases where cert and key are in different directories

**Kubernetes Secrets Support (`certman.go:149-157`)**
- Detects `..data` symlink changes (Kubernetes secret mount pattern)
- Queues reload when Kubernetes updates secrets

### Logging

**Structured Logging with slog**
- Info: Certificate loaded, watching started/stopped, file modifications
- Error: Load failures, watch errors
- Debug: Internal state changes, reloading events

**Log Context**
- File modification events include file path and operation type
- Error messages include structured error context

### Certificate Management

**Thread-Safe Access (`certman.go:197-212`)**
- `GetCertificate()` for server TLS configs
- `GetClientCertificate()` for client TLS configs
- Read lock during certificate access
- Write lock during certificate updates

**Load Process (`certman.go:137-151`)**
- Uses `tls.LoadX509KeyPair()` for validation
- Atomic replacement under write lock
- Logs success/failure with structured context
