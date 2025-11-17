# Cloudflared GUI Backend

Go backend server for managing cloudflared service via systemd D-Bus.

## Quick Start

```bash
# Install dependencies
go mod download

# Run in development
go run ./cmd/server

# Build for production
go build -o server ./cmd/server

# Run binary
./server
```

## API Documentation

See the main [README.md](../../README.md) for full API documentation.

## Development

### Prerequisites

- Go 1.23+
- Linux with systemd
- cloudflared installed and configured

### Running Tests

```bash
go test ./...
```

### Building

```bash
# Local build
go build -o server ./cmd/server

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# Docker build
docker build -t cloudflared-gui-backend .
```

## Configuration

Environment variables (all optional with defaults):

| Variable | Default | Description |
|----------|---------|-------------|
| `HOST` | `127.0.0.1` | Server host address |
| `PORT` | `8080` | Server port |
| `CLOUDFLARED_SERVICE_NAME` | `cloudflared.service` | Systemd service name to manage |
| `CLOUDFLARED_CONFIG_PATH` | `/etc/cloudflared/config.yml` | Path to cloudflared config file |

Example:

```bash
export HOST=127.0.0.1
export PORT=8080
export CLOUDFLARED_SERVICE_NAME=cloudflared.service
export CLOUDFLARED_CONFIG_PATH=/etc/cloudflared/config.yml
go run ./cmd/server
```

## Systemd Integration

The backend uses D-Bus to communicate with systemd:

- `internal/systemd/control.go` - Service start/stop/restart
- `internal/systemd/status.go` - Service status queries
- `internal/systemd/logs.go` - Log streaming via journald

## Security Considerations

1. **Run as root or configure polkit** - Required for systemd service control
2. **Bind to localhost** - Default configuration only allows local connections
3. **Use Cloudflare Tunnel** - Recommended for secure remote access

