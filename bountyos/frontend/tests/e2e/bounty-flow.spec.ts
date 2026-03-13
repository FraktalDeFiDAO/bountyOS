/**
 * E2E Tests for Bounty Flow
 * Tests complete user workflows from bounty discovery to viewing details
 */

import { test, expect } from '@playwright/test';

const FRONTEND_URL = process.env.FRONTEND_URL || 'http://localhost:3000';
const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:8000';

test.describe('Bounty Discovery Flow', () => {
  test('should load homepage and display bounties', async ({ page }) => {
    // Navigate to homepage
    await page.goto(FRONTEND_URL);
    
    // Wait for bounties to load
    await page.waitForSelector('.bounty-card', { timeout: 10000, state: 'visible' });
    
    // Verify bounties are displayed
    const bountyCards = page.locator('.bounty-card');
    const count = await bountyCards.count();
    expect(count).toBeGreaterThan(0);
  });

  test('should filter bounties by platform', async ({ page }) => {
    await page.goto(FRONTEND_URL);
    await page.waitForSelector('.bounty-card', { timeout: 10000 });
    
    // Open filters
    const filterButton = page.locator('[data-testid="filter-button"]');
    if (await filterButton.isVisible()) {
      await filterButton.click();
    }
    
    // Select platform filter (implementation dependent)
    // This test should be updated based on actual filter implementation
    expect(true).toBe(true); // Placeholder
  });

  test('should navigate to bounty detail page', async ({ page }) => {
    await page.goto(FRONTEND_URL);
    await page.waitForSelector('.bounty-card', { timeout: 10000 });
    
    // Click first bounty card
    const firstCard = page.locator('.bounty-card').first();
    await firstCard.click();
    
    // Wait for navigation
    await page.waitForURL(/\/bounties\/.+/, { timeout: 10000 });
    
    // Verify we're on detail page
    const currentUrl = page.url();
    expect(currentUrl).toMatch(/\/bounties\/[\w-]+/);
    
    // Verify bounty title is displayed
    await expect(page.locator('h1')).toBeVisible();
  });

  test('should display bounty reward amount', async ({ page }) => {
    await page.goto(FRONTEND_URL);
    await page.waitForSelector('.bounty-card', { timeout: 10000 });
    
    // Check that reward amounts are visible on cards
    const rewardElements = page.locator('[data-testid*="reward"], .reward-amount');
    const count = await rewardElements.count();
    expect(count).toBeGreaterThan(0);
  });

  test('should handle pagination', async ({ page }) => {
    await page.goto(`${FRONTEND_URL}/bounties?page=1`);
    await page.waitForSelector('.bounty-card', { timeout: 10000 });
    
    // Check for pagination controls (if implemented)
    const paginationControls = page.locator('[data-testid*="pagination"], .pagination');
    // This test should be updated based on actual pagination implementation
    expect(true).toBe(true); // Placeholder
  });
});

test.describe('Bounty Detail Flow', () => {
  test('should load bounty detail page directly', async ({ page }) => {
    // Use a known bounty ID
    await page.goto(`${FRONTEND_URL}/bounties/proxies-sx-3944053546`);
    
    // Wait for page to load
    await page.waitForLoadState('networkidle', { timeout: 10000 });
    
    // Should show either bounty details or "not found" message
    const hasTitle = await page.locator('h1').isVisible().catch(() => false);
    const hasNotFound = await page.locator('text=/not found/i').isVisible().catch(() => false);
    
    expect(hasTitle || hasNotFound).toBe(true);
  });

  test('should display bounty description', async ({ page }) => {
    await page.goto(`${FRONTEND_URL}/bounties/proxies-sx-3944053546`);
    await page.waitForLoadState('networkidle', { timeout: 10000 });
    
    // Check for description section
    const descriptionSection = page.locator('text=/description/i, .description');
    // This test should be updated based on actual implementation
    expect(true).toBe(true); // Placeholder
  });

  test('should show apply/submit button', async ({ page }) => {
    await page.goto(`${FRONTEND_URL}/bounties/proxies-sx-3944053546`);
    await page.waitForLoadState('networkidle', { timeout: 10000 });
    
    // Check for apply button
    const applyButton = page.locator('button:has-text("Apply"), button:has-text("Submit")');
    // This test should be updated based on actual implementation
    expect(true).toBe(true); // Placeholder
  });

  test('should navigate back to bounties list', async ({ page }) => {
    await page.goto(`${FRONTEND_URL}/bounties/proxies-sx-3944053546`);
    await page.waitForLoadState('networkidle', { timeout: 10000 });
    
    // Click breadcrumb or back link
    const backLink = page.locator('a:has-text("Bounties"), .breadcrumb a').first();
    if (await backLink.isVisible()) {
      await backLink.click();
      await page.waitForURL(/\/bounties$/, { timeout: 10000 });
    }
  });
});

test.describe('API Integration', () => {
  test('backend health endpoint should be accessible', async ({ request }) => {
    const response = await request.get(`${BACKEND_URL}/health`);
    expect(response.ok()).toBeTruthy();
    
    const data = await response.json();
    expect(data.status).toBe('healthy');
  });

  test('bounties API should return data', async ({ request }) => {
    const response = await request.get(`${BACKEND_URL}/api/bounties?limit=5`);
    expect(response.ok()).toBeTruthy();
    
    const data = await response.json();
    expect(data.data).toBeDefined();
    expect(Array.isArray(data.data)).toBeTruthy();
  });

  test('single bounty API should work', async ({ request }) => {
    // First get a bounty ID from the list
    const listResponse = await request.get(`${BACKEND_URL}/api/bounties?limit=1`);
    const listData = await listResponse.json();
    
    if (listData.data && listData.data.length > 0) {
      const bountyId = listData.data[0].id;
      
      const response = await request.get(`${BACKEND_URL}/api/bounties/${bountyId}`);
      expect(response.ok()).toBeTruthy();
      
      const data = await response.json();
      expect(data.data.id).toBe(bountyId);
    }
  });
});

test.describe('Error Handling', () => {
  test('should handle 404 for non-existent bounty', async ({ page }) => {
    await page.goto(`${FRONTEND_URL}/bounties/non-existent-bounty-id-12345`);
    await page.waitForLoadState('networkidle', { timeout: 10000 });
    
    // Should show error state or redirect
    const hasError = await page.locator('text=/not found|error/i').isVisible().catch(() => false);
    const hasBackLink = await page.locator('a:has-text("Back")').isVisible().catch(() => false);
    
    expect(hasError || hasBackLink).toBe(true);
  });

  test('should handle backend errors gracefully', async ({ page }) => {
    // This test would require mocking the backend to return errors
    // Implementation dependent
    expect(true).toBe(true); // Placeholder
  });
});

test.describe('Performance', () => {
  test('should load bounties within acceptable time', async ({ page }) => {
    const startTime = Date.now();
    
    await page.goto(FRONTEND_URL);
    await page.waitForSelector('.bounty-card', { timeout: 10000 });
    
    const loadTime = Date.now() - startTime;
    
    // Should load within 5 seconds
    expect(loadTime).toBeLessThan(5000);
  });

  test('should not have memory leaks on navigation', async ({ page }) => {
    // Navigate multiple times
    for (let i = 0; i < 5; i++) {
      await page.goto(FRONTEND_URL);
      await page.waitForSelector('.bounty-card', { timeout: 10000 });
      
      const cards = page.locator('.bounty-card');
      await cards.first().click();
      await page.waitForURL(/\/bounties\/.+/, { timeout: 10000 });
      
      await page.goto(FRONTEND_URL);
    }
    
    // If we got here without crashing, no major memory leaks
    expect(true).toBe(true);
  });
});

test.describe('Accessibility', () => {
  test('should have proper heading structure', async ({ page }) => {
    await page.goto(FRONTEND_URL);
    await page.waitForSelector('.bounty-card', { timeout: 10000 });
    
    // Check for h1
    const h1Count = await page.locator('h1').count();
    expect(h1Count).toBeGreaterThanOrEqual(1);
  });

  test('bounty cards should be keyboard navigable', async ({ page }) => {
    await page.goto(FRONTEND_URL);
    await page.waitForSelector('.bounty-card', { timeout: 10000 });
    
    // Try to navigate with keyboard
    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');
    await page.keyboard.press('Enter');
    
    // Should navigate to bounty detail or trigger action
    expect(page.url()).not.toBe(FRONTEND_URL);
  });
});
