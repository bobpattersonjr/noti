# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Noti is a cross-platform CLI tool that monitors long-running processes and triggers notifications when they complete. It supports 19+ notification services including Slack, Telegram, Pushover, Discord, webhooks, and platform-native notifications (macOS banners, Linux freedesktop, Windows notifications).

## Build Commands

```bash
make build              # Build debug binary with race detector (output: out/noti)
make test               # Run unit tests with coverage (excludes integration tests)
make test-integration   # Run E2E tests (requires out/noti to exist)
make lint               # Run golangci-lint with custom configuration
make man                # Generate man pages from markdown (requires pandoc)
```

## Running a Single Test

```bash
go test -v -run TestName ./internal/command/...
go test -v -run TestName ./service/slack/...
```

## Architecture

### Service Interface Pattern

All notification services implement a simple interface with a `Send() error` method. Each service package (e.g., `service/slack/`, `service/telegram/`) contains its own `Notification` struct that implements this interface.

### Directory Structure

- `cmd/noti/main.go` - Entry point
- `internal/command/` - Core CLI logic, configuration, and service builders
  - `root.go` - Root command and notification orchestration
  - `config.go` - Configuration defaults and Viper setup
  - `cloud.go` - Cloud service builder functions (getSlack, getTelegram, etc.)
  - `local.go`, `local_darwin.go`, `local_windows.go` - OS-specific banner notifications
- `service/` - Individual notification service packages (23 total)
- `integration/` - E2E tests that require the built binary

### Configuration Precedence

1. `viper.Set` (runtime)
2. Command-line flags
3. Configuration file (noti.yaml)
4. Environment variables (NOTI_*)
5. Defaults (in `config.go`)

Note that the config file outranks environment variables, which is a
deliberate departure from both upstream and viper's own fixed ordering
(`Set` > flag > env > config > default). Viper's order can't be reconfigured,
so `bindNotiEnv` skips `BindEnv` for any key already present in the config
file. Two consequences:

- `setupConfigFile` **must** run before `bindNotiEnv` in `configureApp`.
- Precedence is per-key: an env var still applies to keys the file omits.

### Adding a New Notification Service

1. Create package in `service/<servicename>/` with `Notification` struct implementing `Send() error`
2. Add defaults to `baseDefaults` in `internal/command/config.go`
3. Add env bindings to `keyEnvBindings` in `config.go`
4. Add service name to the `services` map in `enabledFromSlice`, `hasServiceFlags`, and `enabledFromFlags`
5. Create builder function (e.g., `getServiceName`) in `cloud.go`
6. Add call to builder in `getNotifications` function
7. Update man pages in `docs/man/`

### OS-Specific Code

Platform-specific implementations use build tags:
- `*_darwin.go` - macOS
- `*_windows.go` - Windows
- `*_unix.go` - Linux/Unix

## Dependencies

Uses vendored dependencies (`vendor/` directory). Key libraries:
- `spf13/cobra` - CLI framework
- `spf13/viper` - Configuration management
- `godbus/dbus/v5` - Linux D-Bus notifications
