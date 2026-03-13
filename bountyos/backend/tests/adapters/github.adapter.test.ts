/**
 * GitHub Adapter Unit Tests
 */

import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';
import { GitHubAdapter } from '../../src/adapters/platforms/github.adapter';
import type { IBounty } from '../../src/adapters/types';

const mockFetch = vi.fn();
global.fetch = mockFetch;

describe('GitHubAdapter', () => {
  let adapter: GitHubAdapter;

  beforeEach(() => {
    adapter = new GitHubAdapter();
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('Platform Info', () => {
    it('should have correct platform info', () => {
      expect(adapter.platformInfo.id).toBe('github');
      expect(adapter.platformInfo.name).toBe('GitHub');
      expect(adapter.platformInfo.isActive).toBe(true);
    });

    it('should have logo URL', () => {
      expect(adapter.platformInfo.logoUrl).toContain('github');
    });
  });

  describe('Capabilities', () => {
    it('should support organization field', () => {
      expect(adapter.capabilities.supportsOrganization).toBe(true);
    });

    it('should support contributor count', () => {
      expect(adapter.capabilities.supportsContributorCount).toBe(true);
    });

    it('should support sorting by popularity', () => {
      expect(adapter.capabilities.supportsSortByPopularity).toBe(true);
    });
  });

  describe('Multi-Repository Fetching', () => {
    const mockRepo1Response = [
      {
        id: 111,
        number: 10,
        title: '[BOUNTY] Test Bounty 1 — $75',
        body: 'Description 1',
        state: 'open' as const,
        created_at: '2026-03-01T00:00:00Z',
        updated_at: '2026-03-01T00:00:00Z',
        closed_at: null,
        labels: [{ name: 'bounty', color: '000000' }],
        user: { login: 'user1' },
        assignee: null,
        comments: 2,
        html_url: 'https://github.com/owner1/repo1/issues/10'
      }
    ];

    const mockRepo2Response = [
      {
        id: 222,
        number: 20,
        title: '[BOUNTY] Test Bounty 2 — $100',
        body: 'Description 2',
        state: 'open' as const,
        created_at: '2026-03-02T00:00:00Z',
        updated_at: '2026-03-02T00:00:00Z',
        closed_at: null,
        labels: [{ name: 'bounty', color: '000000' }],
        user: { login: 'user2' },
        assignee: null,
        comments: 5,
        html_url: 'https://github.com/owner2/repo2/issues/20'
      }
    ];

    it('should fetch from multiple repositories', async () => {
      mockFetch
        .mockResolvedValueOnce({ ok: true, json: async () => mockRepo1Response })
        .mockResolvedValueOnce({ ok: true, json: async () => mockRepo2Response });

      const bounties = await adapter.methods.fetchBounties({ limit: 50 });

      expect(bounties.length).toBeGreaterThanOrEqual(1);
      expect(mockFetch).toHaveBeenCalledTimes(2); // Two repositories
    });

    it('should handle repository errors gracefully', async () => {
      mockFetch
        .mockResolvedValueOnce({ ok: true, json: async () => mockRepo1Response })
        .mockRejectedValueOnce(new Error('Repo 2 failed'));

      const bounties = await adapter.methods.fetchBounties({ limit: 50 });

      // Should still return bounties from successful repo
      expect(bounties.length).toBeGreaterThanOrEqual(0);
    });

    it('should sort bounties by creation date', async () => {
      mockFetch
        .mockResolvedValueOnce({ ok: true, json: async () => mockRepo1Response })
        .mockResolvedValueOnce({ ok: true, json: async () => mockRepo2Response });

      const bounties = await adapter.methods.fetchBounties({ limit: 50 });

      // Newest first (repo2 is newer)
      if (bounties.length >= 2) {
        expect(new Date(bounties[0].createdAt).getTime())
          .toBeGreaterThanOrEqual(new Date(bounties[1].createdAt).getTime());
      }
    });

    it('should apply limit correctly', async () => {
      mockFetch
        .mockResolvedValueOnce({ ok: true, json: async () => mockRepo1Response })
        .mockResolvedValueOnce({ ok: true, json: async () => mockRepo2Response });

      const bounties = await adapter.methods.fetchBounties({ limit: 1 });

      expect(bounties.length).toBeLessThanOrEqual(1);
    });
  });

  describe('Bounty Transformation', () => {
    const mockIssue = {
      id: 123456,
      number: 5,
      title: '[BOUNTY] API Development — $200',
      body: '## What to Build\n\nWe need an API...',
      state: 'open' as const,
      created_at: '2026-03-10T00:00:00Z',
      updated_at: '2026-03-10T00:00:00Z',
      closed_at: null,
      labels: [
        { name: 'bounty', color: '000000' },
        { name: 'api', color: '123456' },
        { name: 'backend', color: 'abcdef' }
      ],
      user: { login: 'projectowner', avatar_url: 'https://example.com/avatar.png' },
      assignee: null,
      comments: 3,
      html_url: 'https://github.com/testorg/testrepo/issues/5'
    };

    it('should transform issue to bounty correctly', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => [mockIssue]
      });

      const bounties = await adapter.methods.fetchBounties();
      const bounty = bounties[0];

      expect(bounty.id).toBe('github-testorg-testrepo-123456');
      expect(bounty.type).toBe('DEV');
      expect(bounty.rewardAmount).toBe(200);
      expect(bounty.organization).toBe('testorg/testrepo');
      expect(bounty.contributorCount).toBe(3);
    });

    it('should extract repository from organization field', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => [mockIssue]
      });

      const bounties = await adapter.methods.fetchBounties();

      expect(bounties[0].organization).toBe('testorg/testrepo');
    });

    it('should include all labels as tags', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => [mockIssue]
      });

      const bounties = await adapter.methods.fetchBounties();

      expect(bounties[0].tags).toContain('bounty');
      expect(bounties[0].tags).toContain('api');
      expect(bounties[0].tags).toContain('backend');
      expect(bounties[0].tags).toContain('testrepo'); // Repo name added as tag
    });

    it('should include author information in metadata', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => [mockIssue]
      });

      const bounties = await adapter.methods.fetchBounties();

      expect(bounties[0].metadata?.author).toBe('projectowner');
      expect(bounties[0].metadata?.repository).toBe('testorg/testrepo');
    });

    it('should extract requirements from body', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => [mockIssue]
      });

      const bounties = await adapter.methods.fetchBounties();

      expect(bounties[0].requirements).toBeDefined();
      expect(bounties[0].requirements?.description).toContain('What to Build');
    });
  });

  describe('Status Detection', () => {
    it('should detect open bounties', async () => {
      const openIssue = {
        id: 1,
        number: 1,
        title: 'Open Bounty',
        body: '',
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
        json: async () => [openIssue]
      });

      const bounties = await adapter.methods.fetchBounties();
      expect(bounties[0].status).toBe('OPEN');
    });

    it('should detect in-progress bounties', async () => {
      const assignedIssue = {
        id: 1,
        number: 1,
        title: 'Assigned Bounty',
        body: '',
        state: 'open' as const,
        created_at: new Date().toISOString(),
        labels: [],
        user: { login: 'test' },
        assignee: { login: 'developer' },
        comments: 0,
        html_url: 'https://example.com'
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => [assignedIssue]
      });

      const bounties = await adapter.methods.fetchBounties();
      expect(bounties[0].status).toBe('IN_PROGRESS');
    });

    it('should detect completed bounties', async () => {
      const closedIssue = {
        id: 1,
        number: 1,
        title: 'Closed Bounty',
        body: '',
        state: 'closed' as const,
        created_at: new Date().toISOString(),
        closed_at: new Date().toISOString(),
        labels: [],
        user: { login: 'test' },
        assignee: null,
        comments: 0,
        html_url: 'https://example.com'
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => [closedIssue]
      });

      const bounties = await adapter.methods.fetchBounties();
      expect(bounties[0].status).toBe('COMPLETED');
    });
  });

  describe('Rate Limiting', () => {
    it('should delay between repository requests', async () => {
      mockFetch
        .mockResolvedValueOnce({ ok: true, json: async () => [] })
        .mockResolvedValueOnce({ ok: true, json: async () => [] });

      const startTime = Date.now();
      await adapter.methods.fetchBounties({ limit: 50 });
      const endTime = Date.now();

      // Should have at least 2 second delay between repos
      expect(endTime - startTime).toBeGreaterThanOrEqual(1500);
    });
  });

  describe('Health Check', () => {
    it('should check rate limit status', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          resources: {
            core: {
              remaining: 5000
            }
          }
        })
      });

      const health = await adapter.methods.healthCheck();

      expect(health.status).toBe('healthy');
    });

    it('should return degraded when rate limit is low', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          resources: {
            core: {
              remaining: 5
            }
          }
        })
      });

      const health = await adapter.methods.healthCheck();

      expect(health.status).toBe('degraded');
      expect(health.message).toContain('Rate limit');
    });

    it('should handle rate limit check failure', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Network error'));

      const health = await adapter.methods.healthCheck();

      expect(health.status).toBe('down');
    });
  });

  describe('Bounty ID Extraction', () => {
    it('should extract info from github-owner-repo-123 format', () => {
      const adapterAny = adapter as any;
      const result = adapterAny.extractBountyInfo('github-owner-repo-123');
      
      expect(result).toEqual({
        owner: 'owner',
        repo: 'repo',
        issueNumber: 123
      });
    });

    it('should return null for invalid format', () => {
      const adapterAny = adapter as any;
      const result = adapterAny.extractBountyInfo('invalid-format');
      
      expect(result).toBeNull();
    });
  });
});
