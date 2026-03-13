/**
 * Proxies.sx Adapter Unit Tests
 */

import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';
import { ProxiesSxAdapter } from '../../src/adapters/platforms/proxies-sx.adapter';
import type { IBounty } from '../../src/adapters/types';

// Mock fetch globally
const mockFetch = vi.fn();
global.fetch = mockFetch;

describe('ProxiesSxAdapter', () => {
  let adapter: ProxiesSxAdapter;

  beforeEach(() => {
    adapter = new ProxiesSxAdapter();
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('Platform Info', () => {
    it('should have correct platform info', () => {
      expect(adapter.platformInfo.id).toBe('proxies-sx');
      expect(adapter.platformInfo.name).toBe('Proxies.sx');
      expect(adapter.platformInfo.isActive).toBe(true);
    });

    it('should support correct bounty types', () => {
      expect(adapter.platformInfo.supportedTypes).toContain('DEV');
      expect(adapter.platformInfo.supportedTypes).toContain('MICRO');
    });
  });

  describe('Capabilities', () => {
    it('should support title and description', () => {
      expect(adapter.capabilities.supportsTitle).toBe(true);
      expect(adapter.capabilities.supportsDescription).toBe(true);
    });

    it('should support reward amount', () => {
      expect(adapter.capabilities.supportsRewardAmount).toBe(true);
      expect(adapter.capabilities.supportsRewardCurrency).toBe(true);
    });

    it('should not support deadline', () => {
      expect(adapter.capabilities.supportsDeadline).toBe(false);
    });

    it('should support pagination', () => {
      expect(adapter.capabilities.supportsPagination).toBe(true);
      expect(adapter.capabilities.maxPageSize).toBe(100);
    });

    it('should not require auth', () => {
      expect(adapter.capabilities.requiresAuth).toBe(false);
      expect(adapter.capabilities.authMethod).toBe('none');
    });
  });

  describe('Fetch Bounties', () => {
    const mockGitHubResponse = [
      {
        id: 3944053546,
        number: 73,
        title: '[BOUNTY] X/Twitter Real-Time Search API — $100 paid in $SX token',
        body: '## What to Build\n\nAn API that searches tweets...',
        state: 'open' as const,
        created_at: '2026-02-15T13:03:38Z',
        updated_at: '2026-02-15T13:03:38Z',
        closed_at: null,
        labels: [{ name: 'bounty', color: '000000' }],
        user: { login: 'testuser', avatar_url: 'https://example.com/avatar.png' },
        assignee: null,
        comments: 0,
        html_url: 'https://github.com/test/repo/issues/73'
      }
    ];

    it('should fetch bounties successfully', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockGitHubResponse
      });

      const bounties = await adapter.methods.fetchBounties({ limit: 10 });

      expect(bounties).toHaveLength(1);
      expect(bounties[0].id).toBe('proxies-sx-3944053546');
      expect(bounties[0].title).toContain('X/Twitter');
      expect(mockFetch).toHaveBeenCalled();
    });

    it('should parse reward from title', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockGitHubResponse
      });

      const bounties = await adapter.methods.fetchBounties();

      expect(bounties[0].rewardAmount).toBe(100);
      expect(bounties[0].rewardCurrency).toBe('USD');
    });

    it('should determine correct bounty type', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockGitHubResponse
      });

      const bounties = await adapter.methods.fetchBounties();

      expect(bounties[0].type).toBe('DEV');
    });

    it('should determine correct status', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockGitHubResponse
      });

      const bounties = await adapter.methods.fetchBounties();

      expect(bounties[0].status).toBe('OPEN');
    });

    it('should handle assigned bounties', async () => {
      const assignedIssue = {
        ...mockGitHubResponse[0],
        assignee: { login: 'developer123' }
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => [assignedIssue]
      });

      const bounties = await adapter.methods.fetchBounties();

      expect(bounties[0].status).toBe('IN_PROGRESS');
      expect(bounties[0].metadata?.isAssigned).toBe(true);
      expect(bounties[0].metadata?.assignee).toBe('developer123');
    });

    it('should handle closed bounties', async () => {
      const closedIssue = {
        ...mockGitHubResponse[0],
        state: 'closed' as const,
        closed_at: '2026-03-01T00:00:00Z'
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => [closedIssue]
      });

      const bounties = await adapter.methods.fetchBounties();

      expect(bounties[0].status).toBe('COMPLETED');
    });

    it('should extract tags from labels', async () => {
      const issueWithLabels = {
        ...mockGitHubResponse[0],
        labels: [
          { name: 'bounty', color: '000000' },
          { name: 'api', color: '123456' },
          { name: 'twitter', color: 'abcdef' }
        ]
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => [issueWithLabels]
      });

      const bounties = await adapter.methods.fetchBounties();

      expect(bounties[0].tags).toContain('bounty');
      expect(bounties[0].tags).toContain('api');
      expect(bounties[0].tags).toContain('twitter');
    });

    it('should handle API errors gracefully', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 403,
        statusText: 'Rate limit exceeded'
      });

      const bounties = await adapter.methods.fetchBounties();

      expect(bounties).toHaveLength(0);
    });

    it('should handle timeout errors', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Timeout'));

      const bounties = await adapter.methods.fetchBounties();

      expect(bounties).toHaveLength(0);
    });

    it('should support pagination', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockGitHubResponse
      });

      await adapter.methods.fetchBounties({ limit: 50, page: 2 });

      expect(mockFetch).toHaveBeenCalled();
      const callArgs = mockFetch.mock.calls[0][0];
      expect(callArgs).toContain('per_page=50');
    });
  });

  describe('Fetch Single Bounty', () => {
    const mockIssue = {
      id: 3944053546,
      number: 73,
      title: '[BOUNTY] Test Bounty — $50',
      body: 'Test description',
      state: 'open' as const,
      created_at: '2026-02-15T13:03:38Z',
      updated_at: '2026-02-15T13:03:38Z',
      closed_at: null,
      labels: [{ name: 'bounty', color: '000000' }],
      user: { login: 'testuser' },
      assignee: null,
      comments: 0,
      html_url: 'https://github.com/test/repo/issues/73'
    };

    it('should fetch bounty by ID', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockIssue
      });

      const bounty = await adapter.methods.fetchBounty('proxies-sx-73');

      expect(bounty).not.toBeNull();
      expect(bounty?.id).toBe('proxies-sx-3944053546');
    });

    it('should return null for non-existent bounty', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 404
      });

      const bounty = await adapter.methods.fetchBounty('proxies-sx-999999');

      expect(bounty).toBeNull();
    });

    it('should handle invalid bounty ID format', async () => {
      const bounty = await adapter.methods.fetchBounty('invalid-format');

      expect(bounty).toBeNull();
    });
  });

  describe('Health Check', () => {
    it('should return healthy status', async () => {
      mockFetch
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({})
        })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => [{ id: 1 }]
        });

      const health = await adapter.methods.healthCheck();

      expect(health.status).toBe('healthy');
      expect(health.responseTime).toBeDefined();
    });

    it('should return degraded status when repo accessible but issues not', async () => {
      mockFetch
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({})
        })
        .mockResolvedValueOnce({
          ok: false,
          status: 403
        });

      const health = await adapter.methods.healthCheck();

      expect(health.status).toBe('degraded');
    });

    it('should return down status on error', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Network error'));

      const health = await adapter.methods.healthCheck();

      expect(health.status).toBe('down');
      expect(health.message).toContain('Network error');
    });
  });

  describe('Bounty Validation', () => {
    it('should add adapter version to bounty', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => [{
          id: 123,
          number: 1,
          title: 'Test',
          body: 'Test',
          state: 'open' as const,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
          closed_at: null,
          labels: [],
          user: { login: 'test' },
          assignee: null,
          comments: 0,
          html_url: 'https://example.com'
        }]
      });

      const bounties = await adapter.methods.fetchBounties();

      expect(bounties[0]._adapterVersion).toBe('1.0.0');
      expect(bounties[0]._fetchedAt).toBeDefined();
    });

    it('should track supported fields', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => [{
          id: 123,
          number: 1,
          title: 'Test',
          body: 'Test',
          state: 'open' as const,
          created_at: new Date().toISOString(),
          labels: [],
          user: { login: 'test' },
          assignee: null,
          comments: 0,
          html_url: 'https://example.com'
        }]
      });

      const bounties = await adapter.methods.fetchBounties();

      expect(bounties[0]._supportedFields).toContain('title');
      expect(bounties[0]._supportedFields).toContain('rewardAmount');
      expect(bounties[0]._unsupportedFields).toContain('deadline');
    });
  });

  describe('Reward Parsing', () => {
    it('should parse reward from title with dollar sign', async () => {
      const issue = {
        id: 123,
        number: 1,
        title: '[BOUNTY] Test — $150 paid',
        body: null,
        state: 'open' as const,
        created_at: new Date().toISOString(),
        labels: [],
        user: { login: 'test' },
        assignee: null,
        comments: 0,
        html_url: 'https://example.com'
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => [issue]
      });

      const bounties = await adapter.methods.fetchBounties();

      expect(bounties[0].rewardAmount).toBe(150);
    });

    it('should use default reward when not specified', async () => {
      const issue = {
        id: 123,
        number: 1,
        title: '[BOUNTY] Test bounty',
        body: 'No reward mentioned',
        state: 'open' as const,
        created_at: new Date().toISOString(),
        labels: [],
        user: { login: 'test' },
        assignee: null,
        comments: 0,
        html_url: 'https://example.com'
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => [issue]
      });

      const bounties = await adapter.methods.fetchBounties();

      expect(bounties[0].rewardAmount).toBe(50); // Default
    });
  });
});
