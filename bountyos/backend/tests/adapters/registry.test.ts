/**
 * Adapter Registry Unit Tests
 */

import { describe, it, expect, beforeEach, vi } from 'vitest';
import { AdapterRegistry } from '../registry';
import type { IPlatformAdapter, IBounty, IPlatformHealth } from '../types';

// Mock adapter for testing
class MockAdapter implements IPlatformAdapter {
  readonly platformInfo = {
    id: 'mock-platform',
    name: 'Mock Platform',
    url: 'https://mock.example.com',
    description: 'Mock platform for testing',
    isActive: true,
    supportedTypes: ['DEV' as const, 'GRANT' as const]
  };

  readonly capabilities = {
    supportsTitle: true,
    supportsDescription: true,
    supportsRewardAmount: true,
    supportsRewardCurrency: true,
    supportsDeadline: false,
    supportsTags: true,
    supportsOrganization: false,
    supportsOrganizationLogo: false,
    supportsDifficulty: false,
    supportsRemote: false,
    supportsContributorCount: false,
    supportsSubmissionsCount: false,
    supportsApplications: false,
    supportsSubmissions: false,
    supportsMessaging: false,
    supportsDirectApply: false,
    supportsTypeFilter: false,
    supportsStatusFilter: true,
    supportsRewardFilter: false,
    supportsTagFilter: false,
    supportsDateFilter: false,
    supportsSortByReward: false,
    supportsSortByDate: true,
    supportsSortByDeadline: false,
    supportsSortByPopularity: false,
    supportsPagination: true,
    maxPageSize: 100,
    defaultPageSize: 30,
    rateLimitPerMinute: 60,
    requiresAuth: false,
    authMethod: 'none' as const
  };

  readonly version = '1.0.0';
  readonly lastUpdated = new Date('2026-03-13');

  readonly methods = {
    fetchBounties: vi.fn(),
    fetchBounty: vi.fn(),
    healthCheck: vi.fn()
  };
}

describe('AdapterRegistry', () => {
  let registry: AdapterRegistry;
  let mockAdapter: MockAdapter;

  beforeEach(() => {
    registry = new AdapterRegistry();
    mockAdapter = new MockAdapter();
  });

  describe('Registration', () => {
    it('should register an adapter', () => {
      registry.registerAdapter(mockAdapter);
      
      expect(registry.hasAdapter('mock-platform')).toBe(true);
      expect(registry.getAdapter('mock-platform')).toBe(mockAdapter);
    });

    it('should return null for non-existent adapter', () => {
      expect(registry.getAdapter('non-existent')).toBeNull();
    });

    it('should replace existing adapter when re-registering', () => {
      const newAdapter = new MockAdapter();
      
      registry.registerAdapter(mockAdapter);
      registry.registerAdapter(newAdapter);
      
      expect(registry.getAdapter('mock-platform')).toBe(newAdapter);
    });

    it('should unregister an adapter', () => {
      registry.registerAdapter(mockAdapter);
      registry.unregisterAdapter('mock-platform');
      
      expect(registry.hasAdapter('mock-platform')).toBe(false);
      expect(registry.getAdapter('mock-platform')).toBeNull();
    });
  });

  describe('Platform Information', () => {
    it('should return all supported platforms', () => {
      const platforms = registry.getSupportedPlatforms();
      
      expect(platforms.length).toBeGreaterThan(0);
      expect(platforms.map(p => p.id)).toContain('proxies-sx');
      expect(platforms.map(p => p.id)).toContain('github');
    });

    it('should return all registered adapters', () => {
      registry.registerAdapter(mockAdapter);
      
      const adapters = registry.getAllAdapters();
      
      expect(adapters).toContain(mockAdapter);
    });

    it('should return adapter version', () => {
      registry.registerAdapter(mockAdapter);
      
      const version = registry.getAdapterVersion('mock-platform');
      
      expect(version).toBe('1.0.0');
    });

    it('should return null for non-existent adapter version', () => {
      const version = registry.getAdapterVersion('non-existent');
      
      expect(version).toBeNull();
    });
  });

  describe('Registry Statistics', () => {
    it('should return correct stats', () => {
      registry.registerAdapter(mockAdapter);
      
      const stats = registry.getStats();
      
      expect(stats.totalAdapters).toBe(1);
      expect(stats.platformIds).toContain('mock-platform');
      expect(stats.activePlatforms).toBeGreaterThan(0);
    });
  });

  describe('Bounty Fetching', () => {
    beforeEach(() => {
      registry.registerAdapter(mockAdapter);
    });

    it('should fetch bounties from all adapters', async () => {
      const mockBounties: IBounty[] = [
        {
          id: 'mock-1',
          type: 'DEV',
          title: 'Test Bounty 1',
          description: 'Test description',
          rewardAmount: 100,
          rewardCurrency: 'USD',
          status: 'OPEN',
          platform: mockAdapter.platformInfo,
          tags: ['test'],
          url: 'https://mock.example.com/1',
          createdAt: new Date().toISOString(),
          _supportedFields: [],
          _unsupportedFields: []
        }
      ];

      mockAdapter.methods.fetchBounties.mockResolvedValue(mockBounties);

      const bounties = await registry.fetchAllBounties({ limit: 10 });

      expect(bounties).toHaveLength(1);
      expect(bounties[0].id).toBe('mock-1');
      expect(mockAdapter.methods.fetchBounties).toHaveBeenCalled();
    });

    it('should handle adapter errors gracefully', async () => {
      mockAdapter.methods.fetchBounties.mockRejectedValue(new Error('Test error'));

      const bounties = await registry.fetchAllBounties({ limit: 10 });

      expect(bounties).toHaveLength(0);
    });

    it('should fetch bounty by ID', async () => {
      const mockBounty: IBounty = {
        id: 'mock-1',
        type: 'DEV',
        title: 'Test Bounty',
        description: 'Test',
        rewardAmount: 100,
        rewardCurrency: 'USD',
        status: 'OPEN',
        platform: mockAdapter.platformInfo,
        tags: [],
        url: 'https://mock.example.com',
        createdAt: new Date().toISOString(),
        _supportedFields: [],
        _unsupportedFields: []
      };

      mockAdapter.methods.fetchBounty.mockResolvedValue(mockBounty);

      const bounty = await registry.fetchBountyById('mock-platform-1');

      expect(bounty).toBe(mockBounty);
      expect(mockAdapter.methods.fetchBounty).toHaveBeenCalledWith('mock-platform-1');
    });

    it('should return null for non-existent bounty', async () => {
      mockAdapter.methods.fetchBounty.mockResolvedValue(null);

      const bounty = await registry.fetchBountyById('mock-platform-nonexistent');

      expect(bounty).toBeNull();
    });

    it('should extract platform ID from bounty ID', async () => {
      // Test with proxies-sx format
      const result = registry.fetchBountyById('proxies-sx-3944053546');
      
      // Should not throw error
      await expect(result).resolves.toBeDefined();
    });
  });

  describe('Health Monitoring', () => {
    it('should get health status for all adapters', async () => {
      const mockHealth: IPlatformHealth = {
        status: 'healthy',
        responseTime: 100,
        bountyCount: 50
      };

      mockAdapter.methods.healthCheck.mockResolvedValue(mockHealth);
      registry.registerAdapter(mockAdapter);

      const healthStatus = await registry.getHealthStatus();

      expect(healthStatus.has('mock-platform')).toBe(true);
      expect(healthStatus.get('mock-platform')?.status).toBe('healthy');
    });

    it('should handle health check errors', async () => {
      mockAdapter.methods.healthCheck.mockRejectedValue(new Error('Health check failed'));
      registry.registerAdapter(mockAdapter);

      const healthStatus = await registry.getHealthStatus();

      expect(healthStatus.get('mock-platform')?.status).toBe('error');
    });
  });

  describe('Platform ID Extraction', () => {
    it('should extract proxies-sx platform ID', () => {
      const registryAny = registry as any;
      const platformId = registryAny.extractPlatformId('proxies-sx-3944053546');
      expect(platformId).toBe('proxies-sx');
    });

    it('should extract github platform ID', () => {
      const registryAny = registry as any;
      const platformId = registryAny.extractPlatformId('github-owner-repo-123');
      expect(platformId).toBe('github');
    });

    it('should return null for unknown platform', () => {
      const registryAny = registry as any;
      const platformId = registryAny.extractPlatformId('unknown-platform-123');
      expect(platformId).toBeNull();
    });
  });
});
