/**
 * Web Scraper Tests
 * Tests for bounty platform web scraping functionality
 */

import { describe, it, expect } from 'vitest';
import { 
  scrapeSuperteamBounties, 
  scrapeAlgoraBounties, 
  scrapeCode4renaBounties,
  scanGitHubBountyRepos,
  parseReward
} from '../src/services/web-scraper.js';

describe('Web Scraper Service', () => {
  describe('parseReward', () => {
    it('should parse USD amounts', () => {
      expect(parseReward('$1000')).toBe(1000);
      expect(parseReward('$1,000')).toBe(1000);
      expect(parseReward('$1,000.50')).toBe(1000.50);
    });

    it('should parse crypto amounts', () => {
      expect(parseReward('1000 USDC')).toBe(1000);
      expect(parseReward('5 ETH')).toBe(5);
      expect(parseReward('100 SOL')).toBe(100);
    });

    it('should return 0 for invalid input', () => {
      expect(parseReward('')).toBe(0);
      expect(parseReward('negotiable')).toBe(0);
      expect(parseReward(null as any)).toBe(0);
    });
  });

  describe('scrapeSuperteamBounties', () => {
    it('should return at least one bounty or fallback', async () => {
      const result = await scrapeSuperteamBounties();
      expect(Array.isArray(result)).toBe(true);
      expect(result.length).toBeGreaterThan(0);
      
      const bounty = result[0];
      expect(bounty).toHaveProperty('title');
      expect(bounty).toHaveProperty('url');
      expect(bounty.url).toMatch(/^https:\/\//);
      expect(bounty.platform).toEqual({ id: 'superteam', name: 'Superteam' });
    }, 15000);
  });

  describe('scrapeAlgoraBounties', () => {
    it('should return at least one bounty or fallback', async () => {
      const result = await scrapeAlgoraBounties();
      expect(Array.isArray(result)).toBe(true);
      expect(result.length).toBeGreaterThan(0);
      
      const bounty = result[0];
      expect(bounty).toHaveProperty('title');
      expect(bounty).toHaveProperty('url');
      expect(bounty.url).toMatch(/^https:\/\//);
      expect(bounty.platform).toEqual({ id: 'algora', name: 'Algora' });
    }, 15000);
  });

  describe('scrapeCode4renaBounties', () => {
    it('should return at least one bounty or fallback', async () => {
      const result = await scrapeCode4renaBounties();
      expect(Array.isArray(result)).toBe(true);
      expect(result.length).toBeGreaterThan(0);
      
      const bounty = result[0];
      expect(bounty).toHaveProperty('title');
      expect(bounty).toHaveProperty('url');
      expect(bounty.url).toMatch(/^https:\/\//);
      expect(bounty.platform).toEqual({ id: 'code4rena', name: 'Code4rena' });
    }, 15000);
  });

  describe('scanGitHubBountyRepos', () => {
    it('should scan GitHub repos for bounty issues', async () => {
      const result = await scanGitHubBountyRepos();
      expect(Array.isArray(result)).toBe(true);
      
      // May return empty array if repos don't have bounty labels
      // but should not throw errors
    }, 10000);
  });
});
