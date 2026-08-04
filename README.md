# Hublet

Hublet is a self-hosted dashboard and homepage for services on a local network. It uses a Go backend with a Svelte frontend and is distributed as a single Docker image.

## Build

Clone the repository and build the image:

```bash
git clone https://github.com/pzdera/hublet.git
cd hublet

docker build -t hublet:latest .
```

## Run with Docker

```bash
docker run -d \
  --name hublet \
  --restart unless-stopped \
  -p 3001:3000 \
  -e HUBLET_ADDR=:3000 \
  -e HUBLET_DATA_DIR=/app/data \
  -e HUBLET_ICON_DIR=/app/icons \
  -e HUBLET_WALLPAPER_DIR=/app/wallpapers \
  -v /opt/hublet/data:/app/data \
  -v /opt/hublet/icons:/app/icons \
  -v /opt/hublet/wallpapers:/app/wallpapers \
  ghcr.io/pzdera/hublet:latest
```

Open:

```text
http://SERVER_IP:3001
```

## Run with Docker Compose

Create a `compose.yaml` file in the Hublet project directory:

```yaml
services:
  hublet:
    container_name: hublet
    image: ghcr.io/pzdera/hublet:latest
    restart: unless-stopped

    ports:
      - "3001:3000"

    environment:
      HUBLET_ADDR: ":3000"
      HUBLET_DATA_DIR: "/app/data"
      HUBLET_ICON_DIR: "/app/icons"
      HUBLET_WALLPAPER_DIR: "/app/wallpapers"

    volumes:
      - ./data:/app/data
      - ./icons:/app/icons
      - ./wallpapers:/app/wallpapers
```

Pull the image and start the container:

```bash
docker compose pull
docker compose up -d
```

View the container status:

```bash
docker compose ps
```

View application logs:

```bash
docker compose logs -f hublet
```

Stop the container:

```bash
docker compose down
```

Update to the latest image:

```bash
docker compose pull
docker compose up -d
```

## Persistent directories

Hublet stores its persistent files in the following container directories:

- `/app/data` — configuration and configuration backups
- `/app/icons` — locally stored service icons
- `/app/wallpapers` — locally stored wallpapers

The Docker Compose example maps them to these project directories:

```text
hublet/
├── data/
├── icons/
└── wallpapers/
```

Removing or recreating the container does not remove these directories.

## Environment variables

- `HUBLET_ADDR` — listen address; default `:3000`
- `HUBLET_DATA_DIR` — configuration directory; default `/app/data`
- `HUBLET_ICON_DIR` — icon directory; default `/app/icons`
- `HUBLET_WALLPAPER_DIR` — wallpaper directory; default `/app/wallpapers`

## Health check

Check the API health endpoint:

```bash
curl http://127.0.0.1:3001/api/v2/health
```
