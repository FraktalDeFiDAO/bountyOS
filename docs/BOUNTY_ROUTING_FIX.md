# Bounty Routing Fix Summary

**Date:** March 13, 2026  
**Issue:** 404 errors when clicking "View Details" on bounty cards  
**Status:** ⚠️ Partially Fixed

---

## Problem

When clicking "View Details" on bounty cards, users received:
```
Request failed with status code 404
Back to Bounties
```

**Example URL:**
```
http://localhost:3000/bounties/proxies-sx-3944053546
```

---

## Root Cause Analysis

### Issue 1: Missing Backend Route ✅ FIXED

**Problem:** The backend (`bountyos/backend/src/index.ts`) did not have a route for `/api/bounties/:id`.

**Files Involved:**
- `bountyos/backend/src/index.ts` - Main backend entry point
- `bountyos/backend/src/routes/bounties.ts` - Unused bounty routes file

**Fix Applied:**
Added single bounty endpoint to `bountyos/backend/src/index.ts`:

```typescript
// Get single bounty by ID
fastify.get<{ Params: { id: string } }>('/api/bounties/:id', async (request, reply) => {
  const { id } = request.params;
  
  // Fetch all bounties and find the matching one
  const [gitcoin, superteam, algora, proxies, code4rena, github] = await Promise.all([
    fetchGitcoinBounties(),
    fetchSuperteamBounties(),
    fetchAlgoraBounties(),
    fetchProxiesSXBounties(),
    fetchCode4renaBounties(),
    fetchGitHubBountyRepos()
  ]);

  const allBounties = [...gitcoin, ...superteam, ...algora, ...proxies, ...code4rena, ...github];
  const bounty = allBounties.find(b => b.id === id);
  
  if (!bounty) {
    return reply.status(404).send({
      message: `Bounty not found: ${id}`,
      error: 'Not Found',
      statusCode: 404
    });
  }
  
  return { data: bounty };
});
```

**Status:** Code added, needs container rebuild to take effect.

---

### Issue 2: Frontend Configuration ✅ VERIFIED

**Checked:**
- `.env` file has correct API URL: `VITE_API_BASE_URL=http://localhost:8000/api`
- Router configuration is correct: `/bounties/:id` route exists
- BountyCard component routing logic is correct

**Status:** ✅ All correct

---

## Testing

### Manual Test Steps

1. **Backend Health**
   ```bash
   curl http://localhost:8000/health
   ```
   Expected: `{"status": "healthy", ...}`

2. **All Bounties**
   ```bash
   curl http://localhost:8000/api/bounties
   ```
   Expected: List of bounties

3. **Single Bounty** (NEW)
   ```bash
   curl http://localhost:8000/api/bounties/proxies-sx-3944053546
   ```
   Expected: Single bounty object

4. **Frontend Navigation**
   ```
   1. Open http://localhost:3000
   2. Click "View Details" on any bounty card
   3. Should navigate to /bounties/{id} without 404
   ```

### Automated Tests Created

**Files Added:**
1. `bountyos/frontend/tests/bounty-routing.test.ts` - Vitest unit tests
2. `bountyos/frontend/tests/bounty-routing.spec.ts` - Playwright E2E tests
3. `bountyos/frontend/ROUTING_TROUBLESHOOTING.md` - Troubleshooting guide

**Run Tests:**
```bash
cd bountyos/frontend

# Unit tests
pnpm test bounty-routing

# E2E tests (requires Playwright)
pnpm test:e2e
```

---

## Container Rebuild Required

The backend fix requires rebuilding the container image:

```bash
cd bountyos/backend

# Rebuild image
podman build -t localhost/bountyos_backend:latest .

# Restart container
podman stop bountyos-backend
podman rm bountyos-backend
podman run -d --name bountyos-backend --restart=unless-stopped --network host localhost/bountyos_backend:latest

# Verify
curl http://localhost:8000/api/bounties/proxies-sx-3944053546
```

**Expected Response:**
```json
{
  "data": {
    "id": "proxies-sx-3944053546",
    "type": "DEV",
    "title": "[BOUNTY] X/Twitter Real-Time Search API — $100 paid in $SX token",
    ...
  }
}
```

---

## Current Status

| Component | Status | Port | Notes |
|-----------|--------|------|-------|
| Frontend | ✅ Running | 3000 | Working correctly |
| Backend | ⚠️ Needs Rebuild | 8000 | Route added, needs rebuild |
| API `/api/bounties` | ✅ Working | 8000 | Returns all bounties |
| API `/api/bounties/:id` | ⚠️ Pending | 8000 | Route added, needs rebuild |

---

## Next Steps

1. **Rebuild Backend Container**
   ```bash
   cd bountyos/backend
   podman build -t localhost/bountyos_backend:latest .
   ```

2. **Restart Containers**
   ```bash
   podman restart bountyos-backend
   podman restart bountyos-frontend
   ```

3. **Verify Fix**
   ```bash
   # Test API
   curl http://localhost:8000/api/bounties/proxies-sx-3944053546
   
   # Test Frontend
   # Open browser to http://localhost:3000
   # Click "View Details" - should work without 404
   ```

4. **Run Tests**
   ```bash
   cd bountyos/frontend
   pnpm test bounty-routing
   ```

---

## Files Modified

### Backend
- `bountyos/backend/src/index.ts` - Added `/api/bounties/:id` route

### Frontend
- `bountyos/frontend/tests/bounty-routing.test.ts` - NEW: Unit tests
- `bountyos/frontend/tests/bounty-routing.spec.ts` - NEW: E2E tests  
- `bountyos/frontend/ROUTING_TROUBLESHOOTING.md` - NEW: Troubleshooting guide
- `bountyos/frontend/.env.local` - NEW: Local environment config

---

## Port Configuration

**Development:**
- Frontend: `http://localhost:3000`
- Backend: `http://localhost:8000`
- API Base URL: `http://localhost:8000/api`

**Production (if using port mapping):**
- Frontend: `http://localhost:28473` → container:3000
- Backend: `http://localhost:39182` → container:8000
- API Base URL: `http://localhost:39182/api`

---

**Last Updated:** March 13, 2026  
**Maintained By:** bountyOS Development Team
