# Cloudflared GUI Dashboard

Modern React dashboard for managing cloudflared service.

## Quick Start

```bash
# Install dependencies
npm install

# Run development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

## Features

- **Service Controls**: Start, stop, restart cloudflared
- **Real-time Status**: Live service status updates
- **Log Viewer**: Stream logs via WebSocket
- **Config Editor**: Edit cloudflared configuration
- **Dark Mode**: Automatic dark/light theme support
- **Responsive**: Mobile-friendly design

## Development

### Tech Stack

- React 18.3
- TypeScript 5.6
- Vite 6.0
- TanStack Query 5.59
- Lucide React (icons)

### Project Structure

```
src/
├── api/
│   └── client.ts          # API client & WebSocket
├── components/
│   ├── Dashboard.tsx      # Main dashboard
│   ├── ServiceControls.tsx
│   ├── StatusDisplay.tsx
│   ├── LogViewer.tsx
│   └── ConfigEditor.tsx
├── App.tsx
└── main.tsx
```

### API Client

The API client (`src/api/client.ts`) provides:

- RESTful API methods
- WebSocket log streaming
- TypeScript types
- Error handling

### Styling

- Modern CSS with CSS variables
- Dark mode via `prefers-color-scheme`
- Mobile-responsive design
- Gradient backgrounds and shadows

## Building for Production

```bash
npm run build
```

Output will be in `dist/` directory.

### Deployment Options

1. **Static hosting**: Deploy `dist/` to any static host
2. **Nginx**: Serve `dist/` with nginx reverse proxy to backend
3. **Caddy**: Use Caddy as reverse proxy
4. **Cloudflare Pages**: Deploy to Cloudflare Pages

### Example Nginx Config

```nginx
server {
    listen 80;
    server_name cloudflared-gui.local;

    root /var/www/cloudflared-gui;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
    }
}
```

## Environment Variables

Create `.env.local` for custom configuration:

```env
VITE_API_URL=http://127.0.0.1:8080
```

## Browser Support

- Chrome/Edge (latest)
- Firefox (latest)
- Safari (latest)

## Contributing

See main [README.md](../../README.md) for contribution guidelines.

