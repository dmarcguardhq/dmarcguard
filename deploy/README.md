# Deployment Templates

This directory contains deployment templates for various cloud providers and platforms.

## Quick Deploy Buttons

Use the deploy buttons in the main [README.md](../README.md) for one-click deployment.

## Templates

| Provider             | Template                                                   | Documentation                                                                |
| -------------------- | ---------------------------------------------------------- | ---------------------------------------------------------------------------- |
| Railway              | [railway.json](./railway.json)                             | [Railway Docs](https://docs.railway.app/)                                    |
| Render               | [render.yaml](./render.yaml)                               | [Render Docs](https://render.com/docs)                                       |
| Fly.io               | [fly.toml](./fly.toml)                                     | [Fly.io Docs](https://fly.io/docs/)                                          |
| Heroku               | [app.json](./app.json), [heroku.yml](./heroku.yml)         | [Heroku Docs](https://devcenter.heroku.com/)                                 |
| DigitalOcean         | [digitalocean-app.yaml](./digitalocean-app.yaml)           | [DO App Platform Docs](https://docs.digitalocean.com/products/app-platform/) |
| Koyeb                | [koyeb.yaml](./koyeb.yaml)                                 | [Koyeb Docs](https://www.koyeb.com/docs)                                     |
| Zeabur               | [zeabur.json](./zeabur.json)                               | [Zeabur Docs](https://zeabur.com/docs)                                       |
| Google Cloud Run     | [cloudbuild.yaml](./cloudbuild.yaml)                       | [Cloud Run Docs](https://cloud.google.com/run/docs)                          |
| Azure Container Apps | [azure-container-apps.bicep](./azure-container-apps.bicep) | [Azure Docs](https://learn.microsoft.com/en-us/azure/container-apps/)        |
| Northflank           | [northflank.json](./northflank.json)                       | [Northflank Docs](https://northflank.com/docs)                               |
| CapRover             | [captain-definition](./captain-definition)                 | [CapRover Docs](https://caprover.com/docs/)                                  |
| Coolify              | [coolify.yaml](./coolify.yaml)                             | [Coolify Docs](https://coolify.io/docs)                                      |

## Environment Variables

All deployments require the following environment variables:

### Required

| Variable                    | Description                   | Example                |
| --------------------------- | ----------------------------- | ---------------------- |
| `IMAP_HOST`     | IMAP server hostname          | `imap.gmail.com`       |
| `IMAP_USERNAME` | IMAP username/email           | `dmarc@yourdomain.com` |
| `IMAP_PASSWORD` | IMAP password or app password | `your-app-password`    |

### Optional (with defaults)

| Variable                    | Description          | Default           |
| --------------------------- | -------------------- | ----------------- |
| `IMAP_PORT`     | IMAP server port     | `993`             |
| `IMAP_MAILBOX`  | IMAP mailbox         | `INBOX`           |
| `IMAP_USE_TLS`  | Use TLS for IMAP     | `true`            |
| `DATABASE_PATH` | SQLite database path | `/data/db.sqlite` |
| `SERVER_PORT`   | HTTP server port     | `8080`            |
| `SERVER_HOST`   | HTTP server host     | `0.0.0.0`         |

## Manual Deployment

### Using Docker

```bash
docker run -d \
  --name parse-dmarc \
  -p 8080:8080 \
  -e IMAP_HOST=imap.gmail.com \
  -e IMAP_USERNAME=your-email@gmail.com \
  -e IMAP_PASSWORD=your-app-password \
  -v parse-dmarc-data:/data \
  ghcr.io/meysam81/parse-dmarc:latest
```

### Using Fly.io CLI

```bash
cd deploy
fly launch --copy-config
fly secrets set IMAP_HOST=imap.gmail.com
fly secrets set IMAP_USERNAME=your-email@gmail.com
fly secrets set IMAP_PASSWORD=your-app-password
fly deploy
```

### Using Railway CLI

```bash
railway init
railway up
railway variables set IMAP_HOST=imap.gmail.com
railway variables set IMAP_USERNAME=your-email@gmail.com
railway variables set IMAP_PASSWORD=your-app-password
```

### Using Google Cloud Run

```bash
# Build and push image
gcloud builds submit --config=deploy/cloudbuild.yaml

# Or deploy directly
gcloud run deploy parse-dmarc \
  --image ghcr.io/meysam81/parse-dmarc:latest \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars "IMAP_HOST=imap.gmail.com,IMAP_USERNAME=your-email@gmail.com" \
  --set-secrets "IMAP_PASSWORD=parse-dmarc-imap-password:latest"
```

### Using Azure Container Apps

```bash
az deployment group create \
  --resource-group your-resource-group \
  --template-file deploy/azure-container-apps.bicep \
  --parameters imapHost=imap.gmail.com imapUsername=your-email@gmail.com imapPassword=your-app-password
```
