#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

const MIN_E2E_TESTS = Number.parseInt(process.env.MIN_E2E_TESTS || '5', 10);

const rootDir = process.cwd();
const e2eDir = path.join(rootDir, 'e2e-tests');
const routerFile = path.join(rootDir, 'web', 'src', 'router', 'index.js');

const routeCoverageSource = process.env.REQUIRED_E2E_ROUTES
  ? 'REQUIRED_E2E_ROUTES'
  : 'web/src/router/index.js';
const REQUIRED_E2E_ROUTES = process.env.REQUIRED_E2E_ROUTES
  ? parseRoutesFromEnv(process.env.REQUIRED_E2E_ROUTES)
  : parseRoutesFromRouter(routerFile);

if (!fs.existsSync(e2eDir)) {
  fail(`e2e-tests directory not found at ${e2eDir}`);
}

const specFiles = collectSpecFiles(e2eDir);
if (!specFiles.length) {
  fail('No Playwright spec files found under e2e-tests.');
}
if (!REQUIRED_E2E_ROUTES.length) {
  fail(`No required frontend routes found from ${routeCoverageSource}.`);
}

let runnableTests = 0;
const discoveredRoutes = new Set();

for (const specFile of specFiles) {
  const content = fs.readFileSync(specFile, 'utf8');

  const runnableTestRegex = /\btest(?:\.only)?\s*\(\s*(['"`])[\s\S]*?\1\s*,/g;
  const fileTests = content.match(runnableTestRegex) || [];
  runnableTests += fileTests.length;

  const gotoRegex = /page\.goto\(\s*(['"`])([^'"`]+)\1/g;
  let match;
  while ((match = gotoRegex.exec(content)) !== null) {
    const normalized = normalizeRoute(match[2]);
    if (normalized) {
      discoveredRoutes.add(normalized);
    }
  }
}

const missingRoutes = REQUIRED_E2E_ROUTES.filter((route) => !discoveredRoutes.has(route));

if (runnableTests < MIN_E2E_TESTS) {
  fail(
    `E2E coverage gate failed: expected at least ${MIN_E2E_TESTS} runnable Playwright tests, found ${runnableTests}.`
  );
}

if (missingRoutes.length) {
  fail(
    `E2E route coverage gate failed: missing route coverage for ${missingRoutes.join(', ')}.`
  );
}

console.log(
  [
    'E2E coverage gates passed.',
    `Runnable tests: ${runnableTests}`,
    `Route source: ${routeCoverageSource}`,
    `Required routes covered: ${REQUIRED_E2E_ROUTES.join(', ')}`,
  ].join(' ')
);

function collectSpecFiles(dir) {
  const files = [];
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    const fullPath = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      files.push(...collectSpecFiles(fullPath));
      continue;
    }

    if (/\.(spec|test)\.(c|m)?(j|t)sx?$/.test(entry.name)) {
      files.push(fullPath);
    }
  }
  return files;
}

function normalizeRoute(raw) {
  if (!raw) return '';
  const value = raw.trim();

  if (value.startsWith('http://') || value.startsWith('https://')) {
    try {
      const parsed = new URL(value);
      return stripTrailingSlash(parsed.pathname || '/');
    } catch {
      return '';
    }
  }

  if (value.startsWith('/')) {
    return stripTrailingSlash(value);
  }

  return '';
}

function parseRoutesFromEnv(raw) {
  return raw
    .split(',')
    .map((route) => normalizeRoute(route.trim()))
    .filter(Boolean);
}

function parseRoutesFromRouter(filePath) {
  if (!fs.existsSync(filePath)) {
    fail(`Router file not found at ${filePath}`);
  }

  const content = fs.readFileSync(filePath, 'utf8');
  const discovered = new Set();

  const pathRegex = /\bpath\s*:\s*(['"`])([^'"`]+)\1/g;
  let pathMatch;
  while ((pathMatch = pathRegex.exec(content)) !== null) {
    const normalized = normalizeRoute(pathMatch[2]);
    if (isStaticRoute(normalized)) {
      discovered.add(normalized);
    }
  }

  const aliasArrayRegex = /\balias\s*:\s*\[([^\]]*)\]/g;
  let aliasArrayMatch;
  while ((aliasArrayMatch = aliasArrayRegex.exec(content)) !== null) {
    const aliases = aliasArrayMatch[1];
    const quoted = aliases.match(/(['"`])([^'"`]+)\1/g) || [];
    for (const token of quoted) {
      const normalized = normalizeRoute(token.slice(1, -1));
      if (isStaticRoute(normalized)) {
        discovered.add(normalized);
      }
    }
  }

  return [...discovered];
}

function isStaticRoute(route) {
  if (!route) return false;
  if (route.includes(':')) return false;
  if (route.includes('*')) return false;
  return true;
}

function stripTrailingSlash(route) {
  const stripped = route.replace(/\/+$/, '');
  return stripped === '' ? '/' : stripped;
}

function fail(message) {
  console.error(message);
  process.exit(1);
}
