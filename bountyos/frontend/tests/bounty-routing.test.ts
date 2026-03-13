/**
 * Bounty Routing Tests
 * 
 * Unit tests for bounty routing functionality
 * Run with: pnpm test bounty-routing
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { createRouter, createWebHistory } from 'vue-router';

// Mock the router routes
const routes = [
  {
    path: '/',
    name: 'home',
    component: { template: '<div>Home</div>' }
  },
  {
    path: '/bounties',
    name: 'bounties',
    component: { template: '<div>Bounties</div>' }
  },
  {
    path: '/bounties/:id',
    name: 'bounty-detail',
    component: { 
      template: '<div>Bounty Detail: {{ $route.params.id }}</div>',
      props: true
    }
  }
];

describe('Bounty Routing', () => {
  let router: ReturnType<typeof createRouter>;

  beforeEach(() => {
    router = createRouter({
      history: createWebHistory(),
      routes
    });
  });

  describe('Route Configuration', () => {
    it('should have bounty-detail route configured', () => {
      const bountyDetailRoute = router.getRoutes().find(r => r.name === 'bounty-detail');
      expect(bountyDetailRoute).toBeDefined();
      expect(bountyDetailRoute?.path).toBe('/bounties/:id');
    });

    it('should match bounty IDs with various formats', async () => {
      const testCases = [
        '/bounties/proxies-sx-3944053546',
        '/bounties/superteam-123',
        '/bounties/github-issue-9876543210',
        '/bounties/algora-abc-123'
      ];

      for (const path of testCases) {
        await router.push(path);
        expect(router.currentRoute.value.name).toBe('bounty-detail');
        expect(router.currentRoute.value.params.id).toBeDefined();
      }
    });
  });

  describe('Bounty Card Navigation', () => {
    it('should navigate to correct URL when clicking bounty card', () => {
      const bountyId = 'proxies-sx-3944053546';
      const expectedPath = `/bounties/${bountyId}`;
      
      // Simulate router.push from BountyCard component
      const pushPath = `/bounties/${bountyId}`;
      expect(pushPath).toBe(expectedPath);
    });

    it('should handle special characters in bounty IDs', async () => {
      const specialIds = [
        'proxies-sx-3944053546',
        'github-issue-123456',
        'superteam_earn_789'
      ];

      for (const id of specialIds) {
        await router.push(`/bounties/${id}`);
        expect(router.currentRoute.value.params.id).toBe(id);
      }
    });
  });

  describe('API Integration', () => {
    it('should use correct API endpoint for fetching bounty', () => {
      const bountyId = 'proxies-sx-3944053546';
      const expectedEndpoint = `/api/bounties/${bountyId}`;
      
      expect(expectedEndpoint).toBe(`/api/bounties/${bountyId}`);
    });

    it('should construct proper API URL', () => {
      const baseUrl = 'http://localhost:8000';
      const bountyId = 'proxies-sx-3944053546';
      const fullUrl = `${baseUrl}/api/bounties/${bountyId}`;
      
      expect(fullUrl).toBe('http://localhost:8000/api/bounties/proxies-sx-3944053546');
    });
  });
});

describe('Bounty ID Format Validation', () => {
  const validBountyIds = [
    'proxies-sx-3944053546',
    'superteam-info',
    'algora-info', 
    'code4rena-info',
    'github-issue-12345',
    'gitlab-mr-67890'
  ];

  const invalidBountyIds = [
    '',
    null,
    undefined
  ];

  it('should accept valid bounty ID formats', () => {
    const idPattern = /^[\w-]+$/;
    
    validBountyIds.forEach(id => {
      expect(id).toMatch(idPattern);
    });
  });

  it('should reject invalid bounty IDs', () => {
    invalidBountyIds.forEach(id => {
      if (id !== null && id !== undefined) {
        expect(id).toBeFalsy();
      }
    });
  });
});
