# Bounty Routing Troubleshooting Guide

## Issue: 404 Error When Clicking "View Details"

**Symptom:**
```
Request failed with status code 404
Back to Bounties
```

**URL Example:**
```
http://localhost:3000/bounties/proxies-sx-3944053546
```

---

## Root Causes

### 1. **API Port Mismatch**

**Problem:** Frontend is configured to use wrong API port.

**Check:**
```bash
# Frontend should be using port 8000 (host network) or 39182 (port mapping)
cat .env
cat .env.local
```

**Expected:**
```env
VITE_API_BASE_URL=http://localhost:8000/api
# OR if using port mapping:
VITE_API_BASE_URL=http://localhost:39182/api
```

**Fix:**
```bash
# Create or update .env.local
echo "VITE_API_BASE_URL=http://localhost:8000/api" > .env.local
```

---

### 2. **Router Configuration Issue**

**Check:** `src/router/index.ts`

**Expected:**
```typescript
{
  path: '/bounties/:id',
  name: 'bounty-detail',
  component: () => import('@/views/BountyDetailView.vue')
}
```

**Verify:**
```bash
cat src/router/index.ts
```

---

### 3. **Bounty ID Format**

**Problem:** Bounty IDs contain special characters that need proper encoding.

**Example IDs:**
- ✅ `proxies-sx-3944053546` - Valid
- ✅ `superteam-info` - Valid
- ✅ `github-issue-12345` - Valid

**Check BountyCard.vue:**
```typescript
function handleClick() {
  router.push({ name: 'bounty-detail', params: { id: props.bounty.id } });
}
```

---

### 4. **Backend API Not Responding**

**Test:**
```bash
# Test health endpoint
curl http://localhost:8000/health

# Test bounties endpoint
curl http://localhost:8000/api/bounties

# Test specific bounty
curl http://localhost:8000/api/bounties/proxies-sx-3944053546
```

**Expected Response:**
```json
{
  "status": "healthy",
  "platforms": ["gitcoin", "superteam", ...]
}
```

---

### 5. **Frontend Build Issues**

**Symptom:** Old build cached, routes not working.

**Fix:**
```bash
cd bountyos/frontend

# Clear cache and rebuild
rm -rf node_modules/.vite
rm -rf dist
pnpm install
pnpm dev
```

---

## Testing Checklist

### Manual Testing

1. **Backend Health**
   ```bash
   curl http://localhost:8000/health
   ```
   ✅ Should return `{"status": "healthy"}`

2. **API Bounties**
   ```bash
   curl http://localhost:8000/api/bounties | jq '.data[0].id'
   ```
   ✅ Should return a bounty ID

3. **Specific Bounty**
   ```bash
   curl http://localhost:8000/api/bounties/proxies-sx-3944053546
   ```
   ✅ Should return bounty details

4. **Frontend Load**
   ```
   Open: http://localhost:3000
   ```
   ✅ Should show BountyOS homepage

5. **Bounty Cards**
   ```
   Check: .bounty-card elements exist
   ```
   ✅ Should see multiple bounty cards

6. **Navigation**
   ```
   Click: "View Details" button
   ```
   ✅ Should navigate to `/bounties/{id}`
   ✅ Should show bounty details (NOT 404)

---

### Automated Testing

**Run Unit Tests:**
```bash
cd bountyos/frontend
pnpm test bounty-routing
```

**Run E2E Tests:**
```bash
cd bountyos/frontend
pnpm test:e2e
```

---

## Common Fixes

### Fix 1: Update API URL

```bash
cd bountyos/frontend

# Create .env.local with correct API URL
cat > .env.local << EOF
VITE_API_BASE_URL=http://localhost:8000/api
EOF

# Restart dev server
pnpm dev
```

### Fix 2: Clear Router Cache

```bash
# In browser DevTools:
# Application > Storage > Clear Site Data

# Or hard refresh:
Ctrl+Shift+R (Windows/Linux)
Cmd+Shift+R (Mac)
```

### Fix 3: Rebuild Frontend

```bash
cd bountyos/frontend

# Stop dev server (Ctrl+C)

# Clear all caches
rm -rf node_modules/.vite
rm -rf .vite

# Restart
pnpm dev
```

---

## Port Configuration

### Development (Host Network)

```bash
# Frontend: localhost:3000
# Backend:  localhost:8000

# .env.local
VITE_API_BASE_URL=http://localhost:8000/api
```

### Production (Port Mapping)

```bash
# Frontend: localhost:28473 -> container:3000
# Backend:  localhost:39182 -> container:8000

# .env.production
VITE_API_BASE_URL=http://localhost:39182/api
```

---

## Debug Mode

**Enable Vue Router Debug:**
```typescript
// src/router/index.ts
const router = createRouter({
  history: createWebHistory(),
  routes,
  strict: true // Enable strict mode
});

// Add navigation guards for debugging
router.beforeEach((to, from, next) => {
  console.log('Navigation:', {
    from: from.path,
    to: to.path,
    params: to.params
  });
  next();
});
```

**Check Browser Console:**
```
F12 > Console
Look for:
- Navigation logs
- API errors
- 404 errors
```

---

## Quick Fix Script

```bash
#!/bin/bash
# fix-routing.sh

echo "🔧 Fixing Bounty Routing..."

cd bountyos/frontend

# 1. Update .env.local
echo "VITE_API_BASE_URL=http://localhost:8000/api" > .env.local

# 2. Clear caches
rm -rf node_modules/.vite
rm -rf dist

# 3. Restart
echo "🚀 Restarting frontend..."
pnpm dev
```

---

## Verification

After applying fixes, verify:

1. ✅ Backend health: `curl localhost:8000/health`
2. ✅ API returns bounties: `curl localhost:8000/api/bounties`
3. ✅ Frontend loads: `http://localhost:3000`
4. ✅ Bounty cards visible
5. ✅ Click "View Details" → No 404
6. ✅ Bounty details page loads correctly

---

**Last Updated:** March 13, 2026  
**Maintained By:** bountyOS Frontend Team
