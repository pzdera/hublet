# Hublet v2

Hublet v2 is a self-hosted dashboard and homepage for services on a local network. It uses a Go backend and a Svelte frontend and is distributed as a single Docker image.

## Build

```bash
docker build -t hublet-v2:latest .
```

## Run

```bash
docker run -d \
  --name hublet-v2 \
  --restart unless-stopped \
  -p 3001:3000 \
  -e HUBLET_V2_ADDR=:3000 \
  -e HUBLET_V2_DATA_DIR=/app/data \
  -v /opt/hublet-v2/data:/app/data \
  -v /opt/hublet-v2/icons:/app/icons \
  -v /opt/hublet-v2/wallpapers:/app/wallpapers \
  hublet-v2:latest
```

Open `http://SERVER_IP:3001`.

## Persistent directories

- `/app/data` — configuration and backups
- `/app/icons` — locally stored service icons
- `/app/wallpapers` — locally stored wallpapers

## Environment variables

- `HUBLET_V2_ADDR` — listen address; default `:3000`
- `HUBLET_V2_DATA_DIR` — configuration directory; default `/app/data`
- `HUBLET_V2_ICON_DIR` — icon directory; default `/app/icons`
- `HUBLET_V2_WALLPAPER_DIR` — wallpaper directory; default `/app/wallpapers`
