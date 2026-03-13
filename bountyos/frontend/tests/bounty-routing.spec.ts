/**
 * Bounty Routing E2E Tests
 * 
 * Tests to ensure bounty card clicks properly redirect to bounty detail pages
 * and that the API integration works correctly.
 */

import { test, expect } from '@playwright/test';

// Test configuration
const FRONTEND_URL = process.env.FRONTEND_URL || 'http://localhost:3000';
const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:8000';

test.describe('Bounty Routing', () => {
  // Test 1: Verify backend API is accessible
  test('backend API should be healthy', async ({ request }) => {
    const response = await request.get(`${BACKEND_URL}/health`);
    expect(response.ok()).toBeTruthy();
    
    const data = await response.json();
    expect(data.status).toBe('healthy');
    expect(data.platforms).toBeDefined();
  });

  // Test 2: Verify bounties API returns data
  test('bounties API should return valid data', async ({ request }) => {
    const response = await request.get(`${BACKEND_URL}/api/bounties`);
    expect(response.ok()).toBeTruthy();
    
    const data = await response.json();
    expect(data.data).toBeDefined();
    expect(Array.isArray(data.data)).toBeTruthy();
    expect(data.data.length).toBeGreaterThan(0);
    
    // Check first bounty has required fields
    const firstBounty = data.data[0];
    expect(firstBounty.id).toBeDefined();
    expect(firstBounty.title).toBeDefined();
    expect(firstBounty.rewardAmount).toBeDefined();
  });

  // Test 3: Frontend loads successfully
  test('frontend should load successfully', async ({ page }) => {
    await page.goto(FRONTEND_URL);
    
    await expect(page).toHaveTitle(/BountyOS/);
    await expect(page.locator('text=Bounties')).toBeVisible();
  });

  // Test 4: Bounty cards are rendered
  test('bounty cards should be rendered on home page', async ({ page }) => {
    await page.goto(FRONTEND_URL);
    await page.waitForSelector('.bounty-card', { timeout: 10000 });
    
    const bountyCards = page.locator('.bounty-card');
    const count = await bountyCards.count();
    expect(count).toBeGreaterThan(0);
  });

  // Test 5: Clicking bounty card navigates to detail page
  test('clicking bounty card should navigate to detail page', async ({ page }) => {
    await page.goto(FRONTEND_URL);
    await page.waitForSelector('.bounty-card', { timeout: 10000 });
    
    // Get the first bounty card
    const firstCard = page.locator('.bounty-card').first();
    
    // Get the bounty ID from the card (we'll check the URL after click)
    const cardTitle = await firstCard.locator('h3').textContent();
    expect(cardTitle).toBeDefined();
    
    // Click the card
    await firstCard.click();
    
    // Wait for navigation
    await page.waitForURL(/\/bounties\/.+/, { timeout: 10000 });
    
    // Verify we're on the bounty detail page
    const currentUrl = page.url();
    expect(currentUrl).toMatch(/\/bounties\/[\w-]+/);
    
    // Verify the bounty title is displayed
    await expect(page.locator('h1')).toContainText(cardTitle || '');
  });

  // Test 6: View Details button works
  test('View Details button should navigate to bounty detail', async ({ page }) => {
    await page.goto(FRONTEND_URL);
    await page.waitForSelector('.bounty-card', { timeout: 10000 });
    
    // Find the first "View Details" button
    const viewDetailsBtn = page.locator('button:has-text("View Details")').first();
    await expect(viewDetailsBtn).toBeVisible();
    
    // Get bounty title before clicking
    const card = viewDetailsBtn.locator('..').locator('..').locator('..');
    const cardTitle = await card.locator('h3').textContent();
    
    // Click the button
    await viewDetailsBtn.click();
    
    // Wait for navigation
    await page.waitForURL(/\/bounties\/.+/, { timeout: 10000 });
    
    // Verify detail page loaded
    const currentUrl = page.url();
    expect(currentUrl).toMatch(/\/bounties\/[\w-]+/);
    
    // Verify title matches
    await expect(page.locator('h1')).toContainText(cardTitle || '');
  });

  // Test 7: Direct URL access works
  test('direct access to bounty URL should work', async ({ page, request }) => {
    // Get a bounty ID from API
    const response = await request.get(`${BACKEND_URL}/api/bounties`);
    const data = await response.json();
    const bountyId = data.data[0].id;
    
    // Navigate directly to the bounty URL
    await page.goto(`${FRONTEND_URL}/bounties/${bountyId}`);
    
    // Wait for page to load
    await page.waitForLoadState('networkidle');
    
    // Verify the bounty detail page loaded without errors
    await expect(page.locator('h1')).toBeVisible();
    
    // Verify no 404 error
    const errorElement = page.locator('text=Bounty not found');
    await expect(errorElement).not.toBeVisible();
  });

  // Test 8: Invalid bounty ID shows proper error
  test('invalid bounty ID should show error message', async ({ page }) => {
    await page.goto(`${FRONTEND_URL}/bounties/invalid-bounty-id-12345`);
    
    // Wait for page to load
    await page.waitForLoadState('networkidle');
    
    // Should show error state
    const errorElement = page.locator('text=/Bounty not found|not found/i');
    await expect(errorElement).toBeVisible();
    
    // Should have "Back to Bounties" link
    const backLink = page.locator('a:has-text("Back to Bounties")');
    await expect(backLink).toBeVisible();
  });

  // Test 9: Bounty detail has all required sections
  test('bounty detail page should have all required sections', async ({ page, request }) => {
    // Get a bounty ID from API
    const response = await request.get(`${BACKEND_URL}/api/bounties`);
    const data = await response.json();
    const bountyId = data.data[0].id;
    
    await page.goto(`${FRONTEND_URL}/bounties/${bountyId}`);
    await page.waitForLoadState('networkidle');
    
    // Check for required sections
    await expect(page.locator('h1')).toBeVisible(); // Title
    await expect(page.locator('text=Reward')).toBeVisible(); // Reward section
    await expect(page.locator('text=Description')).toBeVisible(); // Description
    await expect(page.locator('text=Requirements')).toBeVisible(); // Requirements
    await expect(page.locator('text=Apply')).toBeVisible(); // Apply button
  });

  // Test 10: Breadcrumb navigation works
  test('breadcrumb navigation should work', async ({ page, request }) => {
    // Get a bounty ID from API
    const response = await request.get(`${BACKEND_URL}/api/bounties`);
    const data = await response.json();
    const bountyId = data.data[0].id;
    
    await page.goto(`${FRONTEND_URL}/bounties/${bountyId}`);
    await page.waitForLoadState('networkidle');
    
    // Click breadcrumb "Bounties" link
    const breadcrumbLink = page.locator('nav a:has-text("Bounties")');
    await expect(breadcrumbLink).toBeVisible();
    await breadcrumbLink.click();
    
    // Should navigate back to bounties list
    await page.waitForURL(/\/bounties$/, { timeout: 10000 });
    await expect(page.locator('.bounty-card')).toHaveCount({ min: 1 });
  });
});

// Test 11: Test with specific bounty ID format (proxies-sx-*)
test.describe('Proxies.sx Bounty Routing', () => {
  test('proxies-sx bounty should load correctly', async ({ page }) => {
    const testBountyId = 'proxies-sx-3944053546';
    
    await page.goto(`${FRONTEND_URL}/bounties/${testBountyId}`);
    await page.waitForLoadState('networkidle');
    
    // Should not show 404
    const errorElement = page.locator('text=/Bounty not found/i');
    await expect(errorElement).not.toBeVisible();
    
    // Should show bounty details
    await expect(page.locator('h1')).toBeVisible();
  });
});
