/**
 * API Integration Tests
 * Live tests for the bounty aggregation API
 */

import { describe, it, expect } from 'vitest';

const API_BASE = 'http://localhost:8000/api';

describe('BountyOS API Integration Tests', () => {
  describe('GET /health', () => {
    it('should return healthy status', async () => {
      const response = await fetch(`${API_BASE}/health`);
      expect(response.status).toBe(200);
      
      const data = await response.json();
      expect(data.status).toBe('healthy');
      expect(data).toHaveProperty('timestamp');
      expect(data).toHaveProperty('version');
    });
  });

  describe('GET /api/bounties', () => {
    it('should return bounties from multiple platforms', async () => {
      const response = await fetch(`${API_BASE}/bounties?limit=20`);
      expect(response.status).toBe(200);
      
      const data = await response.json();
      expect(data).toHaveProperty('data');
      expect(data).toHaveProperty('meta');
      expect(data).toHaveProperty('filters');
      
      // Should have bounties
      expect(Array.isArray(data.data)).toBe(true);
      expect(data.data.length).toBeGreaterThan(0);
      
      // Each bounty should have required fields
      data.data.forEach((bounty: any) => {
        expect(bounty).toHaveProperty('id');
        expect(bounty).toHaveProperty('title');
        expect(bounty).toHaveProperty('url');
        expect(bounty).toHaveProperty('platform');
        expect(bounty.url).toMatch(/^https:\/\//);
      });
      
      // Should have platform breakdown
      expect(data.meta.platforms).toHaveProperty('proxies-sx');
      expect(data.meta.platforms).toHaveProperty('superteam');
      expect(data.meta.platforms).toHaveProperty('algora');
      expect(data.meta.platforms).toHaveProperty('code4rena');
    }, 30000);

    it('should filter by platform', async () => {
      const response = await fetch(`${API_BASE}/bounties?platform=proxies-sx&limit=5`);
      const data = await response.json();
      
      data.data.forEach((bounty: any) => {
        expect(bounty.platform.id).toBe('proxies-sx');
      });
    }, 15000);

    it('should filter by type', async () => {
      const response = await fetch(`${API_BASE}/bounties?type=DEV&limit=5`);
      const data = await response.json();
      
      data.data.forEach((bounty: any) => {
        expect(bounty.type).toBe('DEV');
      });
    }, 15000);

    it('should sort by reward (highest first)', async () => {
      const response = await fetch(`${API_BASE}/bounties?limit=10`);
      const data = await response.json();
      
      let lastReward = Infinity;
      data.data.forEach((bounty: any) => {
        expect(bounty.rewardAmount).toBeLessThanOrEqual(lastReward);
        lastReward = bounty.rewardAmount;
      });
    }, 15000);
  });

  describe('GET /api/bounties/featured', () => {
    it('should return featured bounties', async () => {
      const response = await fetch(`${API_BASE}/bounties/featured?limit=6`);
      expect(response.status).toBe(200);
      
      const data = await response.json();
      expect(Array.isArray(data)).toBe(true);
      expect(data.length).toBeLessThanOrEqual(6);
    }, 15000);
  });

  describe('GET /api/bounties/types', () => {
    it('should return bounty types', async () => {
      const response = await fetch(`${API_BASE}/bounties/types`);
      expect(response.status).toBe(200);
      
      const data = await response.json();
      expect(Array.isArray(data)).toBe(true);
      expect(data.length).toBeGreaterThan(0);
      
      data.forEach((type: any) => {
        expect(type).toHaveProperty('id');
        expect(type).toHaveProperty('name');
        expect(type).toHaveProperty('description');
      });
    });
  });

  describe('GET /api/bounties/platforms', () => {
    it('should return platform list', async () => {
      const response = await fetch(`${API_BASE}/bounties/platforms`);
      expect(response.status).toBe(200);
      
      const data = await response.json();
      expect(Array.isArray(data)).toBe(true);
      
      data.forEach((platform: any) => {
        expect(platform).toHaveProperty('id');
        expect(platform).toHaveProperty('name');
        expect(platform).toHaveProperty('url');
        expect(platform).toHaveProperty('types');
        expect(platform.url).toMatch(/^https:\/\//);
      });
    });
  });

  describe('URL Validation', () => {
    it('should have valid URLs for all bounties', async () => {
      const response = await fetch(`${API_BASE}/bounties?limit=50`);
      const data = await response.json();
      
      const invalidUrls = data.data.filter((bounty: any) => {
        try {
          new URL(bounty.url);
          return false;
        } catch {
          return true;
        }
      });
      
      expect(invalidUrls.length).toBe(0);
    }, 30000);
  });

  describe('Performance', () => {
    it('should respond within 30 seconds', async () => {
      const startTime = Date.now();
      const response = await fetch(`${API_BASE}/bounties?limit=50`);
      const endTime = Date.now();
      
      expect(response.status).toBe(200);
      expect(endTime - startTime).toBeLessThan(30000);
    }, 35000);
  });
});
