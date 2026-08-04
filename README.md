# Hublet v2

Hublet v2 is a self-hosted dashboard and homepage for services on a local network. It uses a Go backend with a Svelte frontend and is distributed as a single Docker image.

## Build

Clone the repository and build the image:

```bash
git clone https://github.com/pzdera/hublet-v2.git
cd hublet-v2

docker build -t hublet-v2:latest .
```

## Run with Docker

```bash
docker run -d \
  --name hublet-v2 \
  --restart unless-stopped \
  -p 3001:3000 \
  -e HUBLET_V2_ADDR=:3000 \
  -e HUBLET_V2_DATA_DIR=/app/data \
  -e HUBLET_V2_ICON_DIR=/app/icons \
  -e HUBLET_V2_WALLPAPER_DIR=/app/wallpapers \
  -v /opt/hublet-v2/data:/app/data \
  -v /opt/hublet-v2/icons:/app/icons \
  -v /opt/hublet-v2/wallpapers:/app/wallpapers \
  hublet-v2:latest
```

Open:

```text
http://SERVER_IP:3001
```

## Run with Docker Compose

Create a `compose.yaml` file in the Hublet v2 project directory:

```yaml
services:
  hublet-v2:
    container_name: hublet-v2
    image: hublet-v2:latest
    build:
      context: .
      dockerfile: Dockerfile
    restart: unless-stopped

    ports:
      - "3001:3000"

    environment:
      HUBLET_V2_ADDR: ":3000"
      HUBLET_V2_DATA_DIR: "/app/data"
      HUBLET_V2_ICON_DIR: "/app/icons"
      HUBLET_V2_WALLPAPER_DIR: "/app/wallpapers"

    volumes:
      - ./data:/app/data
      - ./icons:/app/icons
      - ./wallpapers:/app/wallpapers
```

Build and start the container:

```bash
docker compose up -d --build
```

View the container status:

```bash
docker compose ps
```

View application logs:

```bash
docker compose logs -f hublet-v2
```

Stop the container:

```bash
docker compose down
```

Update after pulling new source code:

```bash
git pull

docker compose up -d --build
```

## Persistent directories

Hublet v2 stores its persistent files in the following container directories:

- `/app/data` — configuration and configuration backups
- `/app/icons` — locally stored service icons
- `/app/wallpapers` — locally stored wallpapers

The Docker Compose example maps them to these project directories:

```text
hublet-v2/
├── data/
├── icons/
└── wallpapers/
```

Removing or recreating the container does not remove these directories.

## Environment variables

- `HUBLET_V2_ADDR` — listen address; default `:3000`
- `HUBLET_V2_DATA_DIR` — configuration directory; default `/app/data`
- `HUBLET_V2_ICON_DIR` — icon directory; default `/app/icons`
- `HUBLET_V2_WALLPAPER_DIR` — wallpaper directory; default `/app/wallpapers`

## Health check

Check the API health endpoint:

```bash
curl http://127.0.0.1:3001/api/v2/health
```
