# Deployment Guide

Complete guide for deploying Cloudflared GUI in production.

## Prerequisites

- Linux server with systemd
- Cloudflared installed and configured
- Root or sudo access
- Node.js 20+ (for building frontend)
- Go 1.23+ (for building backend)

## Step-by-Step Deployment

### 1. Build the Application

```bash
# Clone repository
git clone <repository-url>
cd cloudflared-gui

# Install dependencies
npm install

# Build backend
cd apps/backend
go build -o server ./cmd/server
cd ../..

# Build frontend
cd apps/dashboard
npm run build
cd ../..
```

### 2. Deploy Backend

```bash
# Create application directory
sudo mkdir -p /opt/cloudflared-gui/backend

# Copy backend binary
sudo cp apps/backend/server /opt/cloudflared-gui/backend/

# Set permissions
sudo chmod +x /opt/cloudflared-gui/backend/server

# Install systemd service
sudo cp apps/backend/cloudflared-gui-backend.service /etc/systemd/system/

# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable cloudflared-gui-backend
sudo systemctl start cloudflared-gui-backend

# Check status
sudo systemctl status cloudflared-gui-backend
```

### 3. Deploy Frontend

#### Option A: Nginx

```bash
# Install nginx
sudo apt install nginx  # Debian/Ubuntu
sudo dnf install nginx  # Fedora/RHEL

# Create web directory
sudo mkdir -p /var/www/cloudflared-gui

# Copy frontend build
sudo cp -r apps/dashboard/dist/* /var/www/cloudflared-gui/

# Create nginx config
sudo nano /etc/nginx/sites-available/cloudflared-gui
```

Nginx configuration:

```nginx
server {
    listen 80;
    server_name localhost;  # Change to your domain

    root /var/www/cloudflared-gui;
    index index.html;

    # Frontend
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Backend API proxy
    location /api {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Health check
    location /health {
        proxy_pass http://127.0.0.1:8080;
    }
}
```

Enable and start:

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/cloudflared-gui /etc/nginx/sites-enabled/

# Test configuration
sudo nginx -t

# Restart nginx
sudo systemctl restart nginx
```

#### Option B: Caddy

```bash
# Install Caddy
sudo apt install caddy  # Debian/Ubuntu

# Create Caddyfile
sudo nano /etc/caddy/Caddyfile
```

Caddyfile configuration:

```
localhost {  # Change to your domain
    root * /var/www/cloudflared-gui
    file_server

    # Reverse proxy API requests
    reverse_proxy /api/* 127.0.0.1:8080
    reverse_proxy /health 127.0.0.1:8080
}
```

Deploy:

```bash
# Copy frontend files
sudo mkdir -p /var/www/cloudflared-gui
sudo cp -r apps/dashboard/dist/* /var/www/cloudflared-gui/

# Restart Caddy
sudo systemctl restart caddy
```

### 4. Secure with Cloudflare Tunnel (Recommended)

```bash
# Add to your cloudflared config.yml
sudo nano /etc/cloudflared/config.yml
```

Add ingress rule:

```yaml
tunnel: <your-tunnel-id>
credentials-file: /etc/cloudflared/<tunnel-id>.json

ingress:
  - hostname: cloudflared-gui.yourdomain.com
    service: http://127.0.0.1:80  # or direct to backend: http://127.0.0.1:8080
  - service: http_status:404
```

Restart cloudflared:

```bash
sudo systemctl restart cloudflared
```

### 5. Configure Non-Root Access (Optional)

```bash
# Install polkit rule
sudo cp apps/backend/polkit/10-cloudflared-gui.rules /etc/polkit-1/rules.d/

# Restart polkit
sudo systemctl restart polkit

# Create group and add users
sudo groupadd cloudflared-admin
sudo usermod -a -G cloudflared-admin your-username

# Update systemd service to run as that user
sudo nano /etc/systemd/system/cloudflared-gui-backend.service
# Change User=root to User=your-username
```

## Security Hardening

### 1. Firewall Configuration

```bash
# Allow only local connections to backend
sudo ufw deny 8080
sudo ufw allow 'Nginx Full'  # or 80/tcp, 443/tcp
sudo ufw enable
```

### 2. SSL/TLS with Let's Encrypt (Nginx)

```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d cloudflared-gui.yourdomain.com

# Auto-renewal is configured automatically
```

### 3. SSL/TLS with Caddy

Caddy handles SSL automatically - just use your domain in the Caddyfile:

```
cloudflared-gui.yourdomain.com {
    # ... rest of config
}
```

## Monitoring

### Check Backend Logs

```bash
# View logs
sudo journalctl -u cloudflared-gui-backend -f

# View last 100 lines
sudo journalctl -u cloudflared-gui-backend -n 100

# View logs with timestamp
sudo journalctl -u cloudflared-gui-backend --since "1 hour ago"
```

### Check Frontend Access Logs (Nginx)

```bash
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

### Health Check

```bash
# Check backend health
curl http://127.0.0.1:8080/health

# Check service status via API
curl http://127.0.0.1:8080/api/service/status
```

## Updating

### Update Backend

```bash
# Build new version
cd apps/backend
go build -o server ./cmd/server

# Stop service
sudo systemctl stop cloudflared-gui-backend

# Replace binary
sudo cp server /opt/cloudflared-gui/backend/

# Start service
sudo systemctl start cloudflared-gui-backend
```

### Update Frontend

```bash
# Build new version
cd apps/dashboard
npm run build

# Replace files
sudo rm -rf /var/www/cloudflared-gui/*
sudo cp -r dist/* /var/www/cloudflared-gui/

# No need to restart nginx/caddy
```

## Backup

### Backend Binary

```bash
sudo cp /opt/cloudflared-gui/backend/server /backup/server-$(date +%Y%m%d)
```

### Frontend Files

```bash
sudo tar -czf /backup/frontend-$(date +%Y%m%d).tar.gz /var/www/cloudflared-gui
```

### Cloudflared Config

```bash
sudo cp /etc/cloudflared/config.yml /backup/config-$(date +%Y%m%d).yml
```

## Troubleshooting

### Backend won't start

```bash
# Check service status
sudo systemctl status cloudflared-gui-backend

# Check for port conflicts
sudo lsof -i :8080

# Check permissions
ls -la /opt/cloudflared-gui/backend/server

# Test binary manually
sudo /opt/cloudflared-gui/backend/server
```

### Frontend not loading

```bash
# Check nginx status
sudo systemctl status nginx

# Test nginx config
sudo nginx -t

# Check file permissions
ls -la /var/www/cloudflared-gui

# Check nginx error logs
sudo tail -f /var/log/nginx/error.log
```

### API requests failing

```bash
# Check backend is running
curl http://127.0.0.1:8080/health

# Check proxy configuration
sudo nginx -t

# Check firewall
sudo ufw status

# Check SELinux (if applicable)
sudo getenforce
sudo setsebool -P httpd_can_network_connect 1
```

## Rollback

If something goes wrong:

```bash
# Rollback backend
sudo systemctl stop cloudflared-gui-backend
sudo cp /backup/server-YYYYMMDD /opt/cloudflared-gui/backend/server
sudo systemctl start cloudflared-gui-backend

# Rollback frontend
sudo rm -rf /var/www/cloudflared-gui/*
sudo tar -xzf /backup/frontend-YYYYMMDD.tar.gz -C /
```

## Performance Tuning

### Nginx

```nginx
# Add to nginx config for better performance
gzip on;
gzip_types text/plain text/css application/json application/javascript;
gzip_min_length 1000;

# Enable caching for static assets
location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

### Backend

Increase system limits if needed:

```bash
# Edit systemd service
sudo systemctl edit cloudflared-gui-backend

# Add:
[Service]
LimitNOFILE=65536
```

## Production Checklist

- [ ] Backend built and deployed
- [ ] Frontend built and deployed
- [ ] Systemd service enabled and running
- [ ] Web server (nginx/caddy) configured
- [ ] Firewall rules configured
- [ ] SSL/TLS configured (if not using Cloudflare Tunnel)
- [ ] Cloudflare Tunnel configured (recommended)
- [ ] Polkit rules installed (if non-root access needed)
- [ ] Logs monitoring setup
- [ ] Backup strategy in place
- [ ] Health checks working
- [ ] Documentation reviewed

## Support

For issues or questions, refer to the main [README.md](README.md) or open an issue on GitHub.

