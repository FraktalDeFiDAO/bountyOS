/**
 * Platform Aggregator Service
 * Fetches real bounties from multiple platforms
 */

import Fastify from 'fastify';
import cors from '@fastify/cors';
import { scrapeSuperteamBounties, scrapeAlgoraBounties, scrapeCode4renaBounties, scanGitHubBountyRepos } from './services/web-scraper.js';

const fastify = Fastify({ logger: false });

// Register CORS
fastify.register(cors, { origin: true });

// ============================================
// REAL PLATFORM INTEGRATIONS
// ============================================

/**
 * Fetch bounties from Gitcoin
 * API: https://api.gitcoin.co/api/bounties
 */
async function fetchGitcoinBounties() {
  try {
    const response = await fetch('https://api.gitcoin.co/api/v1/bounties/?status=open&limit=20', {
      headers: { 'Accept': 'application/json' },
      signal: AbortSignal.timeout(5000)
    });
    
    if (!response.ok) throw new Error(`Gitcoin API error: ${response.status}`);
    
    const data = await response.json();
    
    return (data.results || data || []).map((bounty: any) => ({
      id: `gitcoin-${bounty.id || Math.random()}`,
      type: 'DEV',
      title: bounty.title || 'Untitled Bounty',
      description: (bounty.description || 'No description').substring(0, 500),
      rewardAmount: bounty.price || bounty.value || bounty.bounty_amount || 0,
      rewardCurrency: 'USD',
      status: 'open',
      platform: { id: 'gitcoin', name: 'Gitcoin' },
      tags: bounty.tags || ['Gitcoin', 'Open Source'],
      url: bounty.url || bounty.permalink || `https://gitcoin.co/issue/${bounty.id}`,
      createdAt: bounty.created_at || new Date().toISOString()
    }));
  } catch (error: any) {
    console.error('Gitcoin fetch error:', error.message);
    return []; // Return empty array on error
  }
}

/**
 * Fetch bounties from Superteam Earn
 * Uses web scraping since no public API
 */
async function fetchSuperteamBounties() {
  try {
    const scraped = await scrapeSuperteamBounties();
    if (scraped && scraped.length > 0 && scraped[0].id !== 'superteam-info') {
      console.log(`Superteam: Scraped ${scraped.length} bounties`);
      return scraped;
    }
  } catch (error: any) {
    console.error('Superteam error:', error.message);
  }
  
  // Fallback
  return [{
    id: 'superteam-info',
    type: 'GRANT',
    title: 'Superteam Earn - Solana Bounties & Grants',
    description: 'Superteam offers bounties, grants, and rewards for building on Solana.',
    rewardAmount: 0,
    rewardCurrency: 'USD',
    status: 'open',
    platform: { id: 'superteam', name: 'Superteam' },
    tags: ['Superteam', 'Solana', 'Grants'],
    url: 'https://superteam.fun/earn',
    createdAt: new Date().toISOString()
  }];
}

/**
 * Fetch bounties from Algora
 * Uses web scraping since API requires auth
 */
async function fetchAlgoraBounties() {
  try {
    const scraped = await scrapeAlgoraBounties();
    if (scraped && scraped.length > 0 && scraped[0].id !== 'algora-info') {
      console.log(`Algora: Scraped ${scraped.length} bounties`);
      return scraped;
    }
  } catch (error: any) {
    console.error('Algora error:', error.message);
  }
  
  // Fallback
  return [{
    id: 'algora-info',
    type: 'DEV',
    title: 'Algora - Open Source Bounties',
    description: 'Browse and contribute to open source bounties on Algora.',
    rewardAmount: 0,
    rewardCurrency: 'USD',
    status: 'open',
    platform: { id: 'algora', name: 'Algora' },
    tags: ['Algora', 'Open Source'],
    url: 'https://algora.io/bounties',
    createdAt: new Date().toISOString()
  }];
}

/**
 * Fetch bounties from Code4rena
 * Uses web scraping for contest listings
 */
async function fetchCode4renaBounties() {
  try {
    const scraped = await scrapeCode4renaBounties();
    if (scraped && scraped.length > 0 && scraped[0].id !== 'code4rena-info') {
      console.log(`Code4rena: Scraped ${scraped.length} contests`);
      return scraped;
    }
  } catch (error: any) {
    console.error('Code4rena error:', error.message);
  }
  
  // Fallback
  return [{
    id: 'code4rena-info',
    type: 'SEC',
    title: 'Code4rena - Security Audit Contests',
    description: 'Participate in smart contract security audit contests.',
    rewardAmount: 0,
    rewardCurrency: 'USD',
    status: 'open',
    platform: { id: 'code4rena', name: 'Code4rena' },
    tags: ['Code4rena', 'Security', 'Audit'],
    url: 'https://code4rena.com/contests',
    createdAt: new Date().toISOString()
  }];
}

/**
 * Scan GitHub for bounty issues across multiple repos
 */
async function fetchGitHubBountyRepos() {
  try {
    const bounties = await scanGitHubBountyRepos();
    console.log(`GitHub: Found ${bounties.length} bounty issues`);
    return bounties;
  } catch (error: any) {
    console.error('GitHub scan error:', error.message);
    return [];
  }
}

/**
 * Fetch bounties from Proxies.sx (GitHub Issues)
 * API: GitHub API
 */
async function fetchProxiesSXBounties() {
  try {
    const response = await fetch(
      'https://api.github.com/repos/bolivian-peru/marketplace-service-template/issues?labels=bounty&state=open&limit=20',
      { headers: { 'Accept': 'application/vnd.github.v3+json' } }
    );
    
    if (!response.ok) throw new Error(`GitHub API error: ${response.status}`);
    
    const data = await response.json();
    
    return (data || []).map((issue: any) => ({
      id: `proxies-sx-${issue.id}`,
      type: 'DEV',
      title: issue.title,
      description: issue.body?.substring(0, 500) || 'No description',
      rewardAmount: 50, // Default for Proxies.sx bounties
      rewardCurrency: 'USD',
      status: 'open',
      platform: { id: 'proxies-sx', name: 'Proxies.sx' },
      tags: issue.labels?.map((l: any) => l.name) || ['Proxies.sx', 'Web Scraping'],
      url: issue.html_url,
      createdAt: issue.created_at
    }));
  } catch (error: any) {
    console.error('Proxies.sx fetch error:', error.message);
    return [];
  }
}

// ============================================
// API ENDPOINTS
// ============================================

// Health check
fastify.get('/health', async () => {
  return { 
    status: 'healthy', 
    timestamp: new Date().toISOString(),
    version: '0.1.0',
    platforms: ['gitcoin', 'superteam', 'algora', 'proxies-sx', 'code4rena']
  };
});

// API root
fastify.get('/api', async () => {
  return {
    message: 'BountyOS API - Real Bounty Aggregator',
    endpoints: {
      health: '/health',
      bounties: '/api/bounties',
      featured: '/api/bounties/featured',
      types: '/api/bounties/types',
      platforms: '/api/bounties/platforms'
    }
  };
});

// Get all bounties from all platforms
fastify.get('/api/bounties', async (request: any, reply: any) => {
  const { type, platform, status, limit = 100 } = request.query as any;
  
  console.log(`Fetching bounties with filters: type=${type}, platform=${platform}`);
  
  // Fetch from all platforms in parallel
  const [gitcoin, superteam, algora, proxies, code4rena, github] = await Promise.all([
    fetchGitcoinBounties(),
    fetchSuperteamBounties(),
    fetchAlgoraBounties(),
    fetchProxiesSXBounties(),
    fetchCode4renaBounties(),
    fetchGitHubBountyRepos()
  ]);
  
  // Combine all bounties
  let allBounties = [...gitcoin, ...superteam, ...algora, ...proxies, ...code4rena, ...github];
  
  // Remove duplicates by URL
  const seen = new Set();
  allBounties = allBounties.filter((b: any) => {
    if (seen.has(b.url)) return false;
    seen.add(b.url);
    return true;
  });
  
  // Apply filters
  if (type) {
    allBounties = allBounties.filter((b: any) => b.type === type);
  }
  
  if (platform) {
    allBounties = allBounties.filter((b: any) => b.platform.id === platform);
  }
  
  // Sort by reward (highest first)
  allBounties.sort((a: any, b: any) => (b.rewardAmount || 0) - (a.rewardAmount || 0));
  
  // Apply limit
  allBounties = allBounties.slice(0, Number(limit));
  
  return {
    data: allBounties,
    meta: {
      total: allBounties.length,
      page: 1,
      limit: Number(limit),
      totalPages: 1,
      platforms: {
        gitcoin: gitcoin.length,
        superteam: superteam.length,
        algora: algora.length,
        'proxies-sx': proxies.length,
        code4rena: code4rena.length,
        github: github.length
      }
    },
    filters: {
      types: ['DEV', 'GRANT', 'MICRO', 'GIG', 'HACK', 'AMB', 'RETRO', 'SEC'],
      platforms: ['gitcoin', 'superteam', 'algora', 'proxies-sx', 'code4rena', 'github'],
      status: ['open', 'in_progress', 'completed', 'closed']
    }
  };
});

// Featured bounties (top paying)
fastify.get('/api/bounties/featured', async (request: any, reply: any) => {
  const { limit = 6 } = request.query as any;
  
  const [gitcoin, superteam, algora] = await Promise.all([
    fetchGitcoinBounties(),
    fetchSuperteamBounties(),
    fetchAlgoraBounties()
  ]);
  
  const allBounties = [...gitcoin, ...superteam, ...algora];
  allBounties.sort((a: any, b: any) => (b.rewardAmount || 0) - (a.rewardAmount || 0));
  
  return allBounties.slice(0, Number(limit));
});

// Bounty types
fastify.get('/api/bounties/types', async () => {
  return [
    { id: 'DEV', name: 'Development', description: 'Feature development, OSS contributions' },
    { id: 'GRANT', name: 'Grants', description: 'Protocol/ecosystem grants' },
    { id: 'MICRO', name: 'Microtasks', description: 'Small tasks, surveys, testing' },
    { id: 'GIG', name: 'Freelance', description: 'Contract work, project-based' },
    { id: 'HACK', name: 'Hackathons', description: 'Time-limited competitions' },
    { id: 'AMB', name: 'Ambassador', description: 'Community building' },
    { id: 'RETRO', name: 'Retroactive', description: 'Retroactive public goods' },
    { id: 'SEC', name: 'Security', description: 'Bug bounties, security audits' }
  ];
});

// Platforms
fastify.get('/api/bounties/platforms', async () => {
  return [
    { id: 'gitcoin', name: 'Gitcoin', url: 'https://gitcoin.co', types: ['DEV', 'GRANT', 'HACK'] },
    { id: 'superteam', name: 'Superteam', url: 'https://superteam.fun', types: ['DEV', 'GRANT', 'AMB'] },
    { id: 'algora', name: 'Algora', url: 'https://algora.io', types: ['DEV', 'GIG'] },
    { id: 'proxies-sx', name: 'Proxies.sx', url: 'https://proxies.sx', types: ['DEV', 'MICRO'] },
    { id: 'code4rena', name: 'Code4rena', url: 'https://code4rena.com', types: ['SEC'] }
  ];
});

// Get single bounty by ID
fastify.get<{ Params: { id: string } }>('/api/bounties/:id', async (request, reply) => {
  const { id } = request.params;
  
  console.log(`Fetching single bounty: ${id}`);
  
  // Fetch all bounties and find the matching one
  const [gitcoin, superteam, algora, proxies, code4rena, github] = await Promise.all([
    fetchGitcoinBounties(),
    fetchSuperteamBounties(),
    fetchAlgoraBounties(),
    fetchProxiesSXBounties(),
    fetchCode4renaBounties(),
    fetchGitHubBountyRepos()
  ]);

  const allBounties = [...gitcoin, ...superteam, ...algora, ...proxies, ...code4rena, ...github];
  const bounty = allBounties.find(b => b.id === id);
  
  if (!bounty) {
    return reply.status(404).send({
      message: `Bounty not found: ${id}`,
      error: 'Not Found',
      statusCode: 404
    });
  }
  
  return { data: bounty };
});

// Start server
const start = async () => {
  try {
    await fastify.listen({ port: 8000, host: '0.0.0.0' });
    console.log('🚀 BountyOS API running at http://0.0.0.0:8000');
    console.log('📍 Health: http://0.0.0.0:8000/health');
    console.log('📍 Bounties: http://0.0.0.0:8000/api/bounties');
    console.log('🎯 Fetching REAL bounties from Gitcoin, Superteam, Algora, Proxies.sx, Code4rena');
  } catch (err) {
    fastify.log.error(err);
    process.exit(1);
  }
};

start();
