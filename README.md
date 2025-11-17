# Cloudflared GUI

A modern web-based GUI for managing cloudflared service on Linux systems. Built with Go backend and React frontend in a Turborepo monorepo structure.

## Features

- ✅ **Service Control**: Start, stop, and restart cloudflared service
- ✅ **Real-time Status**: Monitor service status with live updates
- ✅ **Live Logs**: Stream cloudflared logs in real-time via WebSockets
- ✅ **Configuration Editor**: Edit cloudflared config files through the UI
- ✅ **System Integration**: Uses systemd D-Bus for native Linux integration
- ✅ **Modern UI**: Beautiful, responsive React dashboard with dark mode support
- ✅ **Type-Safe**: Shared TypeScript types across frontend and backend

## Architecture

```
cloudflared-gui/
├── apps/
│   ├── backend/         # Go API server with systemd integration
│   └── dashboard/       # React + Vite frontend
└── packages/
    ├── types/           # Shared TypeScript types
    └── ui/              # Shared UI components
```

## Requirements

### System Requirements

- **OS**: Linux with systemd
- **Go**: 1.23 or higher
- **Node.js**: 20.0.0 or higher
- **npm**: 10.9.0 or higher
- **Cloudflared**: Installed and configured

### Backend Dependencies

- `github.com/coreos/go-systemd/v22` - systemd D-Bus integration
- `github.com/gorilla/mux` - HTTP router
- `github.com/gorilla/websocket` - WebSocket support
- `gopkg.in/yaml.v3` - YAML configuration parsing

### Frontend Dependencies

- React 18.3+
- Vite 6.0+
- TanStack Query (React Query) 5.59+
- Lucide React (icons)

## Installation

### 1. Clone the Repository

```bash
git clone <repository-url>
cd cloudflared-gui
```

### 2. Install Dependencies

```bash
# Install root dependencies (Turborepo)
npm install

# Install backend dependencies
cd apps/backend
go mod download
cd ../..

# Install frontend dependencies
cd apps/dashboard
npm install
cd ../..
```

### 3. Build the Project

```bash
# Build all packages
npm run build

# Or build individually
cd apps/backend && go build -o server ./cmd/server
cd apps/dashboard && npm run build
```

## Development

### Run in Development Mode

```bash
# Run all services in development mode
npm run dev
```

This will start:
- Backend API server on `http://127.0.0.1:8080`
- Frontend dev server on `http://127.0.0.1:5173`

### Run Backend Only

```bash
cd apps/backend
go run ./cmd/server
```

### Run Frontend Only

```bash
cd apps/dashboard
npm run dev
```

## Production Deployment

### Option 1: Systemd Service (Recommended)

1. **Build the backend**:

```bash
cd apps/backend
go build -o server ./cmd/server
sudo mkdir -p /opt/cloudflared-gui/backend
sudo cp server /opt/cloudflared-gui/backend/
```

2. **Install the systemd service**:

```bash
sudo cp cloudflared-gui-backend.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable cloudflared-gui-backend
sudo systemctl start cloudflared-gui-backend
```

3. **Build and serve the frontend**:

```bash
cd apps/dashboard
npm run build
# Serve the dist/ folder with your preferred web server (nginx, caddy, etc.)
```

### Option 2: Docker

1. **Build backend Docker image**:

```bash
cd apps/backend
docker build -t cloudflared-gui-backend .
```

2. **Run with Docker**:

```bash
docker run -d \
  --name cloudflared-gui-backend \
  --network host \
  -v /var/run/dbus:/var/run/dbus:ro \
  cloudflared-gui-backend
```

Note: The container needs access to the host's D-Bus socket to control systemd services.

## Configuration

### Backend Configuration

The backend listens on `127.0.0.1:8080` by default. To change this, modify `apps/backend/cmd/server/main.go`:

```go
addr := "127.0.0.1:8080"  // Change this line
```

### Cloudflared Config Path

By default, the backend reads/writes cloudflared config at `/etc/cloudflared/config.yml`. This can be changed in `apps/backend/internal/config/config.go`.

### CORS Configuration

For production, update CORS settings in `apps/backend/internal/api/router.go` to match your deployment domain.

## Security

### Binding to Localhost

By default, the backend binds to `127.0.0.1` (localhost only). This prevents external access without additional configuration.

### Exposing the UI Securely

**Recommended**: Use Cloudflare Tunnel to securely expose the UI:

```yaml
# Add to your cloudflared config.yml
ingress:
  - hostname: cloudflared-gui.yourdomain.com
    service: http://127.0.0.1:5173
  - service: http_status:404
```

### Non-Root Access (Optional)

To allow non-root users to control systemd services, install the polkit rule:

```bash
sudo cp apps/backend/polkit/10-cloudflared-gui.rules /etc/polkit-1/rules.d/
sudo systemctl restart polkit

# Create the group and add users
sudo groupadd cloudflared-admin
sudo usermod -a -G cloudflared-admin your-username
```

## API Endpoints

### Service Control

- `POST /api/service/start` - Start cloudflared service
- `POST /api/service/stop` - Stop cloudflared service
- `POST /api/service/restart` - Restart cloudflared service
- `GET /api/service/status` - Get service status

### Logs

- `WS /api/service/logs` - WebSocket for live log streaming
- `GET /api/service/logs/recent` - Get recent log entries

### Configuration

- `GET /api/config` - Get cloudflared configuration
- `POST /api/config` - Update cloudflared configuration

### Health Check

- `GET /health` - Backend health check

## API Response Format

All API responses follow this format:

```json
{
  "success": true,
  "message": "Optional message",
  "data": { /* Optional response data */ }
}
```

Error responses:

```json
{
  "success": false,
  "error": "Error message"
}
```

## Troubleshooting

### Backend won't start

1. **Check systemd D-Bus connection**:
```bash
busctl status
```

2. **Verify cloudflared service exists**:
```bash
systemctl status cloudflared
```

3. **Check permissions**:
```bash
# Backend needs permission to access systemd D-Bus
sudo usermod -a -G systemd-journal $USER
```

### Frontend can't connect to backend

1. **Verify backend is running**:
```bash
curl http://127.0.0.1:8080/health
```

2. **Check CORS settings** in `apps/backend/internal/api/router.go`

### WebSocket connection fails

1. **Check firewall rules**
2. **Verify WebSocket upgrade is allowed** by reverse proxy (if using one)
3. **Check browser console** for detailed error messages

### Permission denied errors

1. **Run backend as root** (required for systemd control), OR
2. **Configure polkit** as described in Security section

## Development Tips

### Hot Reload

Both backend and frontend support hot reload in development mode:

- Frontend: Automatic reload via Vite
- Backend: Use tools like `air` or `reflex`:

```bash
go install github.com/air-verse/air@latest
cd apps/backend
air
```

### Debugging

Enable Go backend debug logging:

```go
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

Frontend uses React DevTools and TanStack Query DevTools.

## Project Structure

```
cloudflared-gui/
├── apps/
│   ├── backend/
│   │   ├── cmd/server/              # Main entry point
│   │   ├── internal/
│   │   │   ├── api/                 # HTTP handlers & router
│   │   │   ├── config/              # Config management
│   │   │   └── systemd/             # Systemd D-Bus integration
│   │   ├── Dockerfile
│   │   ├── cloudflared-gui-backend.service
│   │   └── polkit/                  # PolicyKit rules
│   └── dashboard/
│       ├── src/
│       │   ├── api/                 # API client
│       │   ├── components/          # React components
│       │   ├── App.tsx
│       │   └── main.tsx
│       ├── index.html
│       └── vite.config.ts
└── packages/
    ├── types/                       # Shared TypeScript types
    └── ui/                          # Shared UI components
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

MIT License - See LICENSE file for details

## Support

For issues, questions, or contributions, please open an issue on GitHub.

## Acknowledgments

- Built with [Turborepo](https://turbo.build/)
- Backend powered by [Go](https://golang.org/)
- Frontend built with [React](https://react.dev/) and [Vite](https://vitejs.dev/)
- SystemD integration via [go-systemd](https://github.com/coreos/go-systemd)
- Icons by [Lucide](https://lucide.dev/)

