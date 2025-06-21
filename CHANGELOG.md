# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0] - 2025-06-21

### Changed
- **BREAKING**: Replaced custom logger interface with Go's standard `log/slog`
  - All logging now uses structured logging with `slog` (Info, Error, Debug levels)
  - Removed public `Logger()` method - use `WithLogger()` option instead
- **BREAKING**: Updated constructor to use functional options pattern
  - Constructor signature changed from `New(certFile, keyFile string)` to `New(certFile, keyFile string, opts ...Option)`
  - Added `WithLogger(*slog.Logger)` option for custom logger
  - Added `WithContext(context.Context)` option for custom context
- **BREAKING**: Minimum Go version bumped to 1.21 (required for `log/slog`)
- Replaced channel-based cancellation with `context.Context` for better lifecycle management
- Stop() method now uses context cancellation instead of channel communication

### Added
- Functional options pattern for cleaner API extensibility
- Context support for graceful shutdown and cancellation
- Private `setupLogs()` function for default logger configuration
- Comprehensive `CLAUDE.md` documentation for AI-assisted development

### Fixed
- Improved thread safety with context-based cancellation
- Better error handling with structured logging context

### Migration Guide

Before (v0.2.x):
```go
cm, err := certman.New("/path/to/cert", "/path/to/key")
if err != nil {
    log.Fatal(err)
}
cm.Logger(myLogger)
```

After (v0.3.0):
```go
cm, err := certman.New("/path/to/cert", "/path/to/key",
    certman.WithLogger(myLogger),
    certman.WithContext(ctx))
if err != nil {
    log.Fatal(err)
}
```
