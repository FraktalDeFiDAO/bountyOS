const { test, expect } = require('@playwright/test');

test.describe('BountyOS Web UI', () => {
  test('should load the dashboard', async ({ page }) => {
    // Go to the local server
    await page.goto('http://localhost:12496');

    // Check header
    const header = await page.locator('h1');
    await expect(header).toContainText('BOUNTY OS v8: OBSIDIAN');

    // Check if stats containers exist
    const stats = await page.locator('#stats');
    await expect(stats).toBeVisible();

    // Check if table exists
    const table = await page.locator('table');
    await expect(table).toBeVisible();

    // Wait for data to load (it should at least show the table headers)
    const tableHeader = await page.locator('th').first();
    await expect(tableHeader).toContainText('Score');
  });

  test('should show stats', async ({ page }) => {
    await page.goto('http://localhost:12496');
    
    // Wait for the stats to be populated (the stat-card class is added by JS)
    await page.waitForSelector('.stat-card');
    
    const statCards = await page.locator('.stat-card');
    const count = await statCards.count();
    expect(count).toBeGreaterThan(0);
  });

  test('should have valid bounty links', async ({ page }) => {
    await page.goto('http://localhost:12496');
    
    // Wait for bounties to load
    await page.waitForSelector('.link');
    
    const links = await page.locator('.link');
    const count = await links.count();
    expect(count).toBeGreaterThan(0);
    
    // Check the first link's href
    const firstLink = links.first();
    const href = await firstLink.getAttribute('href');
    expect(href).not.toBe('undefined');
    expect(href).toMatch(/^https?:\/\//);
  });
});
