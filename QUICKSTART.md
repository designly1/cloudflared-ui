# Quick Start Guide

Get Cloudflared GUI up and running in 5 minutes.

## Prerequisites

Make sure you have:
- Linux with systemd
- Node.js 20+
- Go 1.23+
- cloudflared installed

## Installation

### 1. Clone and Setup

```bash
git clone <repository-url>
cd cloudflared-gui
chmod +x setup.sh
./setup.sh
```

### 2. Start Development Servers

```bash
npm run dev
```

This starts:
- Backend API on `http://127.0.0.1:8080`
- Frontend on `http://127.0.0.1:5173`

### 3. Access the Dashboard

Open your browser to: `http://127.0.0.1:5173`

## Features Overview

### Service Controls
- ▶️ Start cloudflared service
- ⏹️ Stop cloudflared service
- 🔄 Restart cloudflared service

### Status Monitoring
- 🟢 Real-time service status
- 📊 Memory and CPU usage
- 🔢 Process ID (PID)

### Live Logs
- 📡 WebSocket-powered log streaming
- 🔄 Auto-scrolling
- 🗑️ Clear logs button

### Config Editor
- ✏️ Edit cloudflared config
- ✅ JSON validation
- 💾 Save changes directly

## Production Deployment

For production deployment, see [DEPLOYMENT.md](DEPLOYMENT.md).

Quick production setup:

```bash
# Build
cd apps/backend && go build -o server ./cmd/server
cd ../dashboard && npm run build

# Deploy backend
sudo cp apps/backend/server /opt/cloudflared-gui/backend/
sudo cp apps/backend/cloudflared-gui-backend.service /etc/systemd/system/
sudo systemctl enable --now cloudflared-gui-backend

# Deploy frontend (with nginx)
sudo cp -r apps/dashboard/dist/* /var/www/cloudflared-gui/
```

## Troubleshooting

### Permission denied
Run backend as root or configure polkit:
```bash
sudo apps/backend/server
```

### Port already in use
Change port in `apps/backend/cmd/server/main.go`:
```go
addr := "127.0.0.1:8080"  // Change 8080 to another port
```

### Frontend can't connect
Verify backend is running:
```bash
curl http://127.0.0.1:8080/health
```

## Architecture

```
┌─────────────┐      HTTP/WS      ┌──────────┐      D-Bus      ┌──────────────┐
│   Browser   │ ◄───────────────► │    Go    │ ◄──────────────► │   systemd    │
│  (React)    │                   │  Backend │                  │ (cloudflared)│
└─────────────┘                   └──────────┘                  └──────────────┘
```

## Next Steps

- Read the full [README.md](README.md)
- Check [DEPLOYMENT.md](DEPLOYMENT.md) for production
- See [CONTRIBUTING.md](CONTRIBUTING.md) to contribute

## Support

- 📖 Documentation: [README.md](README.md)
- 🐛 Issues: GitHub Issues
- 💬 Discussions: GitHub Discussions

---

**Enjoy managing cloudflared with a beautiful GUI! 🚀**

