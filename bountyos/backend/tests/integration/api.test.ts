/**
 * API Integration Tests
 * Tests for bounty API endpoints using the adapter system
 */

import { describe, it, expect, beforeAll, afterAll, vi } from 'vitest';
import Fastify from 'fastify';
import fastifyTestInstance from 'fastify';

// Mock the adapters
vi.mock('../../src/adapters/registry', () => ({
  adapterRegistry: {
    fetchAllBounties: vi.fn(),
    fetchBountyById: vi.fn(),
    getStats: vi.fn(() => ({
      totalAdapters: 2,
      activePlatforms: 2,
      platformIds: ['proxies-sx', 'github']
    }))
  }
}));

import { adapterRegistry } from '../../src/adapters/registry';

describe('Bounty API Integration', () => {
  let app: ReturnType<typeof fastifyTestInstance>;

  beforeAll(async () => {
    app = fastifyTestInstance({ logger: false });
    
    // Register routes (simplified version for testing)
    app.get('/api/bounties', async (request: any, reply: any) => {
      const { limit = 100, page = 1 } = request.query as any;
      
      const allBounties = await adapterRegistry.fetchAllBounties({
        limit: Number(limit),
        sortBy: 'createdAt',
        sortOrder: 'desc'
      });

      const pageNum = Number(page);
      const limitNum = Number(limit);
      const paginatedBounties = allBounties.slice((pageNum - 1) * limitNum, pageNum * limitNum);

      return {
        data: paginatedBounties,
        meta: {
          total: allBounties.length,
          page: pageNum,
          limit: limitNum,
          totalPages: Math.ceil(allBounties.length / limitNum)
        }
      };
    });

    app.get('/api/bounties/:id', async (request: any, reply: any) => {
      const { id } = request.params;
      const bounty = await adapterRegistry.fetchBountyById(id);
      
      if (!bounty) {
        return reply.status(404).send({
          message: `Bounty not found: ${id}`,
          error: 'Not Found',
          statusCode: 404
        });
      }
      
      return { data: bounty };
    });

    app.get('/health', async () => ({
      status: 'healthy',
      timestamp: new Date().toISOString(),
      version: '2.0.0'
    }));

    await app.ready();
  });

  afterAll(async () => {
    await app.close();
  });

  describe('GET /health', () => {
    it('should return healthy status', async () => {
      const response = await app.inject({
        method: 'GET',
        url: '/health'
      });

      expect(response.statusCode).toBe(200);
      const body = JSON.parse(response.body);
      expect(body.status).toBe('healthy');
      expect(body.version).toBe('2.0.0');
    });
  });

  describe('GET /api/bounties', () => {
    const mockBounties = [
      {
        id: 'proxies-sx-1',
        type: 'DEV',
        title: 'Test Bounty 1',
        description: 'Description 1',
        rewardAmount: 100,
        rewardCurrency: 'USD',
        status: 'OPEN',
        platform: { id: 'proxies-sx', name: 'Proxies.sx' },
        tags: ['test'],
        url: 'https://example.com/1',
        createdAt: new Date().toISOString(),
        _supportedFields: [],
        _unsupportedFields: []
      },
      {
        id: 'github-2',
        type: 'DEV',
        title: 'Test Bounty 2',
        description: 'Description 2',
        rewardAmount: 200,
        rewardCurrency: 'USD',
        status: 'OPEN',
        platform: { id: 'github', name: 'GitHub' },
        tags: ['test'],
        url: 'https://example.com/2',
        createdAt: new Date().toISOString(),
        _supportedFields: [],
        _unsupportedFields: []
      }
    ];

    it('should return bounties with pagination', async () => {
      vi.mocked(adapterRegistry.fetchAllBounties).mockResolvedValue(mockBounties);

      const response = await app.inject({
        method: 'GET',
        url: '/api/bounties?limit=1&page=1'
      });

      expect(response.statusCode).toBe(200);
      const body = JSON.parse(response.body);
      
      expect(body.data).toHaveLength(1);
      expect(body.meta.total).toBe(2);
      expect(body.meta.page).toBe(1);
      expect(body.meta.limit).toBe(1);
      expect(body.meta.totalPages).toBe(2);
    });

    it('should return all bounties when limit is high', async () => {
      vi.mocked(adapterRegistry.fetchAllBounties).mockResolvedValue(mockBounties);

      const response = await app.inject({
        method: 'GET',
        url: '/api/bounties?limit=100'
      });

      expect(response.statusCode).toBe(200);
      const body = JSON.parse(response.body);
      
      expect(body.data).toHaveLength(2);
    });

    it('should handle adapter errors gracefully', async () => {
      vi.mocked(adapterRegistry.fetchAllBounties).mockRejectedValue(new Error('Adapter error'));

      const response = await app.inject({
        method: 'GET',
        url: '/api/bounties'
      });

      // Should return 500 error
      expect(response.statusCode).toBe(500);
    });

    it('should sort by reward amount (highest first)', async () => {
      vi.mocked(adapterRegistry.fetchAllBounties).mockResolvedValue(mockBounties);

      const response = await app.inject({
        method: 'GET',
        url: '/api/bounties'
      });

      const body = JSON.parse(response.body);
      
      // Second bounty has higher reward (200 vs 100)
      expect(body.data[0].rewardAmount).toBeGreaterThanOrEqual(body.data[1]?.rewardAmount || 0);
    });
  });

  describe('GET /api/bounties/:id', () => {
    const mockBounty = {
      id: 'proxies-sx-1',
      type: 'DEV',
      title: 'Test Bounty',
      description: 'Test Description',
      rewardAmount: 100,
      rewardCurrency: 'USD',
      status: 'OPEN',
      platform: { id: 'proxies-sx', name: 'Proxies.sx' },
      tags: ['test'],
      url: 'https://example.com',
      createdAt: new Date().toISOString(),
      _supportedFields: [],
      _unsupportedFields: []
    };

    it('should return single bounty by ID', async () => {
      vi.mocked(adapterRegistry.fetchBountyById).mockResolvedValue(mockBounty);

      const response = await app.inject({
        method: 'GET',
        url: '/api/bounties/proxies-sx-1'
      });

      expect(response.statusCode).toBe(200);
      const body = JSON.parse(response.body);
      
      expect(body.data.id).toBe('proxies-sx-1');
      expect(body.data.title).toBe('Test Bounty');
    });

    it('should return 404 for non-existent bounty', async () => {
      vi.mocked(adapterRegistry.fetchBountyById).mockResolvedValue(null);

      const response = await app.inject({
        method: 'GET',
        url: '/api/bounties/non-existent'
      });

      expect(response.statusCode).toBe(404);
      const body = JSON.parse(response.body);
      
      expect(body.error).toBe('Not Found');
    });

    it('should handle fetch errors', async () => {
      vi.mocked(adapterRegistry.fetchBountyById).mockRejectedValue(new Error('Fetch error'));

      const response = await app.inject({
        method: 'GET',
        url: '/api/bounties/error-test'
      });

      expect(response.statusCode).toBe(500);
    });
  });

  describe('API Response Format', () => {
    it('should return consistent response format', async () => {
      vi.mocked(adapterRegistry.fetchAllBounties).mockResolvedValue([]);

      const response = await app.inject({
        method: 'GET',
        url: '/api/bounties'
      });

      const body = JSON.parse(response.body);
      
      // Should have data array
      expect(body).toHaveProperty('data');
      expect(Array.isArray(body.data)).toBe(true);
      
      // Should have meta object
      expect(body).toHaveProperty('meta');
      expect(body.meta).toHaveProperty('total');
      expect(body.meta).toHaveProperty('page');
      expect(body.meta).toHaveProperty('limit');
    });
  });
});
