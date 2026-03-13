# bountyOS - Podman Containers Status

**Date:** March 13, 2026  
**Status:** ✅ Running  
**Security:** 🔒 Non-standard random ports

---

## 🚀 Running Containers (Detached Mode)

### Web UI - Frontend
- **Container Name:** `bountyos-frontend`
- **Status:** ✅ Running
- **Internal Port:** 3000
- **External Port:** 28473 (random non-standard)
- **Access:** http://localhost:28473
- **Health Check:** ✅ 200 OK

### Backend API
- **Container Name:** `bountyos-backend`
- **Status:** ✅ Running
- **Internal Port:** 8000
- **External Port:** 39182 (random non-standard)
- **Access:** http://localhost:39182
- **Health Check:** ⚠️ 404 (Expected - API routes required)

---

## ⚠️ Obsidian - Config Issue

The Obsidian container has a configuration parsing error that needs to be fixed.

**Error:** `parsing error: no value can start with b`

**Issue:** The config file `/config/test-config.yaml` has a YAML parsing error.

### To Fix Obsidian:

1. **Check config file:**
   ```bash
   cat config/test-config.yaml
   ```

2. **Fix YAML syntax** (likely a value starting with 'b' that needs quotes)

3. **Restart Obsidian:**
   ```bash
   podman rm -f bountyos-obsidian-prod-ssl
   cd /home/administrator/projects/bountyOS
   podman-compose -f docker-compose.prod.ssl.yml up -d obsidian
   ```

### Alternative: Run without config

```bash
podman run -d \
  --name bountyos-obsidian \
  --restart=unless-stopped \
  -p 12496:12496 \
  -v ./au-workspace:/app/au-workspace:ro \
  -e WEB_PORT=12496 \
  -e HEADLESS=true \
  localhost/bountyos_obsidian:latest
```

---

## 📝 Container Management Commands

### Start All Containers
```bash
cd /home/administrator/projects/bountyOS
podman start bountyos-frontend bountyos-backend
```

### Stop All Containers
```bash
podman stop bountyos-frontend bountyos-backend
```

### View Logs
```bash
# Frontend logs
podman logs -f bountyos-frontend

# Backend logs
podman logs -f bountyos-backend

# All logs
podman logs -f bountyos-frontend bountyos-backend
```

### Restart Containers
```bash
podman restart bountyos-frontend bountyos-backend
```

### Check Status
```bash
podman ps --filter "name=bountyos"
```

---

## 🔗 Access Points

| Service | URL | Port | Status |
|---------|-----|------|--------|
| Frontend | http://localhost:28473 | 28473 | ✅ Running |
| Backend API | http://localhost:39182 | 39182 | ✅ Running |
| Obsidian (when fixed) | http://localhost:12496 | 12496 | ⚠️ Config issue |

**Note:** Using random non-standard ports for security and to avoid conflicts.

---

## 🛠️ Troubleshooting

### Frontend not accessible
```bash
podman logs bountyos-frontend
podman restart bountyos-frontend
```

### Backend not accessible
```bash
podman logs bountyos-backend
podman restart bountyos-backend
```

### Port conflicts
```bash
# Check what's using ports
ss -tlnp | grep -E '3000|8000|12496'

# Kill conflicting processes or change container ports
```

---

**Last Updated:** March 13, 2026  
**Maintained By:** bountyOS Core Team
