# Cloudflared GUI — Coding Agent Instructions  
## Go Backend + React Frontend + Turborepo Monorepo

## 1. Project Overview

Create a Turborepo monorepo containing:

- A **Go backend** that controls `cloudflared.service` using **systemd D-Bus**  
- A **React dashboard** for UI  
- Shared UI and type packages  

The system must:

- Start/stop/restart cloudflared  
- Query cloudflared status  
- Stream live logs over WebSockets  
- Read/write cloudflared config YAML  

Backend must run on Linux with systemd enabled.

## 2. Monorepo Folder Structure

```
cloudflared-gui/
├── turbo.json
├── package.json
├── apps/
│   ├── backend/        (Go API server)
│   └── dashboard/      (React dashboard using Vite or Next.js)
└── packages/
    ├── ui/             (shared UI components)
    └── types/          (shared TypeScript API types)
```

## 3. Install Dependencies

### Root-level (Turborepo)

```bash
npm init -y
npm install -D turbo
```

### Go backend

Inside `apps/backend/`:

```bash
go mod init backend
go get github.com/coreos/go-systemd/v22/dbus
go get github.com/coreos/go-systemd/v22/sdjournal
go get github.com/gorilla/mux
go get github.com/gorilla/websocket
go get gopkg.in/yaml.v3
```

### React frontend

Inside `apps/dashboard/`:

```bash
npm create vite@latest
# or:
npx create-next-app@latest
```

## 4. Turborepo Config

Create `turbo.json` in the root:

```json
{
  "pipeline": {
    "dev": {
      "dependsOn": ["^dev"],
      "cache": false
    },
    "build": {
      "dependsOn": ["^build"],
      "outputs": ["dist/**", "build/**"]
    }
  }
}
```

## 5. Backend Structure

```
apps/backend/
├── cmd/server/main.go
├── go.mod
├── go.sum
└── internal/
    ├── systemd/
    │   ├── control.go
    │   ├── status.go
    │   └── logs.go
    ├── config/
    │   └── config.go
    └── api/
        ├── handlers.go
        └── router.go
```

### systemd control module

```go
package systemd

import (
    "github.com/coreos/go-systemd/v22/dbus"
)

type SystemdService struct {
    conn *dbus.Conn
}

func New() (*SystemdService, error) {
    conn, err := dbus.NewSystemConnection()
    if err != nil {
        return nil, err
    }
    return &SystemdService{conn: conn}, nil
}

func (s *SystemdService) Start() error {
    _, err := s.conn.StartUnit("cloudflared.service", "replace", nil)
    return err
}

func (s *SystemdService) Stop() error {
    _, err := s.conn.StopUnit("cloudflared.service", "replace", nil)
    return err
}

func (s *SystemdService) Restart() error {
    _, err := s.conn.RestartUnit("cloudflared.service", "replace", nil)
    return err
}
```

### systemd status module

```go
package systemd

func (s *SystemdService) Status() (map[string]interface{}, error) {
    return s.conn.GetUnitProperties("cloudflared.service")
}
```

### Logs (D-Bus journal)

```go
package systemd

import (
    "github.com/coreos/go-systemd/v22/sdjournal"
)

func StreamLogs(ch chan<- string) error {
    j, err := sdjournal.NewJournal()
    if err != nil:
        return err
    }

    j.AddMatch("_SYSTEMD_UNIT=cloudflared.service")
    j.SeekTail()

    for {
        j.Next()
        entry, _ := j.GetEntry()
        ch <- entry.Fields["MESSAGE"]
    }
}
```

## 6. API Endpoints

Expose:

- `POST /api/service/start`
- `POST /api/service/stop`
- `POST /api/service/restart`
- `GET  /api/service/status`
- `WS   /api/service/logs`

## 7. React Dashboard Requirements

React app must include:

- Dashboard layout  
- Buttons for Start/Stop/Restart  
- Status indicator (running/stopped/failed)  
- Live log viewer using WebSockets  
- Config editor UI  

## 8. Security

- Backend binds to `127.0.0.1`
- Use Cloudflare Tunnel to expose UI securely
- Optional polkit rule for non-root control of systemd

## 9. Deliverables

Coding agent must deliver:

- Turborepo structure
- Go backend with systemd D-Bus modules
- REST API + WebSocket service
- React dashboard with components
- Dockerfile and systemd unit for backend
- README instructions
