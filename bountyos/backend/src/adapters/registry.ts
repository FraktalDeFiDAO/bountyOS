/**
 * Platform Adapter Registry
 * 
 * Central management system for all platform adapters
 * Provides unified interface for fetching bounties across platforms
 */

import {
  IPlatformAdapter,
  IPlatformInfo,
  IAdapterRegistry,
  IBounty,
  IFetchBountiesParams
} from './types';

/**
 * Registry implementation for managing platform adapters
 */
export class AdapterRegistry implements IAdapterRegistry {
  private adapters: Map<string, IPlatformAdapter> = new Map();
  private readonly supportedPlatforms: IPlatformInfo[] = [
    {
      id: 'proxies-sx',
      name: 'Proxies.sx',
      url: 'https://proxies.sx',
      description: 'Mobile proxy marketplace for web scraping',
      isActive: true,
      supportedTypes: ['DEV', 'MICRO']
    },
    {
      id: 'github',
      name: 'GitHub',
      url: 'https://github.com',
      description: 'GitHub Issues bounty tracking',
      isActive: true,
      supportedTypes: ['DEV', 'GRANT']
    },
    {
      id: 'superteam',
      name: 'Superteam',
      url: 'https://superteam.fun',
      description: 'Solana ecosystem bounties and grants',
      isActive: true,
      supportedTypes: ['DEV', 'GRANT', 'DESIGN', 'CONTENT']
    },
    {
      id: 'algora',
      name: 'Algora',
      url: 'https://algora.io',
      description: 'Open source bounties',
      isActive: true,
      supportedTypes: ['DEV']
    },
    {
      id: 'code4rena',
      name: 'Code4rena',
      url: 'https://code4rena.com',
      description: 'Smart contract audit contests',
      isActive: true,
      supportedTypes: ['AUDIT', 'SEC']
    },
    {
      id: 'gitcoin',
      name: 'Gitcoin',
      url: 'https://gitcoin.co',
      description: 'Open source funding platform',
      isActive: true,
      supportedTypes: ['DEV', 'GRANT', 'HACK']
    }
  ];

  /**
   * Get adapter for a specific platform
   */
  getAdapter(platformId: string): IPlatformAdapter | null {
    const adapter = this.adapters.get(platformId);
    if (!adapter) {
      console.warn(`No adapter found for platform: ${platformId}`);
      return null;
    }
    return adapter;
  }

  /**
   * Get all registered adapters
   */
  getAllAdapters(): IPlatformAdapter[] {
    return Array.from(this.adapters.values());
  }

  /**
   * Get information about all supported platforms
   */
  getSupportedPlatforms(): IPlatformInfo[] {
    return this.supportedPlatforms;
  }

  /**
   * Register a new platform adapter
   */
  registerAdapter(adapter: IPlatformAdapter): void {
    const platformId = adapter.platformInfo.id;
    
    if (this.adapters.has(platformId)) {
      console.warn(`Adapter for platform ${platformId} already exists, replacing...`);
    }

    this.adapters.set(platformId, adapter);
    console.log(`Registered adapter for platform: ${platformId}`);
  }

  /**
   * Unregister an adapter
   */
  unregisterAdapter(platformId: string): void {
    const deleted = this.adapters.delete(platformId);
    if (deleted) {
      console.log(`Unregistered adapter for platform: ${platformId}`);
    } else {
      console.warn(`No adapter found for platform: ${platformId}`);
    }
  }

  /**
   * Check if adapter exists for platform
   */
  hasAdapter(platformId: string): boolean {
    return this.adapters.has(platformId);
  }

  /**
   * Get adapter version
   */
  getAdapterVersion(platformId: string): string | null {
    const adapter = this.adapters.get(platformId);
    return adapter?.version || null;
  }

  /**
   * Fetch bounties from all registered adapters
   */
  async fetchAllBounties(params: IFetchBountiesParams = {}): Promise<IBounty[]> {
    const allBounties: IBounty[] = [];
    
    const promises = Array.from(this.adapters.values()).map(async (adapter) => {
      try {
        const bounties = await adapter.methods.fetchBounties(params);
        return bounties;
      } catch (error) {
        console.error(`Error fetching from ${adapter.platformInfo.id}:`, error);
        return [];
      }
    });

    const results = await Promise.all(promises);
    results.forEach(bounties => allBounties.push(...bounties));

    return allBounties;
  }

  /**
   * Fetch bounty by ID from any adapter
   */
  async fetchBountyById(id: string): Promise<IBounty | null> {
    // Extract platform ID from bounty ID format: "platform-id-12345"
    const platformId = this.extractPlatformId(id);
    
    if (!platformId) {
      console.error('Could not extract platform ID from bounty ID:', id);
      return null;
    }

    const adapter = this.getAdapter(platformId);
    if (!adapter) {
      console.error('No adapter found for platform:', platformId);
      return null;
    }

    try {
      const bounty = await adapter.methods.fetchBounty(id);
      return bounty;
    } catch (error) {
      console.error(`Error fetching bounty ${id} from ${platformId}:`, error);
      return null;
    }
  }

  /**
   * Extract platform ID from bounty ID
   * Format: "platform-id-bounty-id" (e.g., "proxies-sx-3944053546")
   */
  private extractPlatformId(bountyId: string): string | null {
    // Try to match known platform IDs
    for (const platform of this.supportedPlatforms) {
      if (bountyId.startsWith(`${platform.id}-`)) {
        return platform.id;
      }
    }

    // Fallback: try to extract first two parts as platform ID
    const parts = bountyId.split('-');
    if (parts.length >= 3) {
      const potentialPlatformId = `${parts[0]}-${parts[1]}`;
      const platform = this.supportedPlatforms.find(p => p.id === potentialPlatformId);
      if (platform) {
        return platform.id;
      }
    }

    return null;
  }

  /**
   * Get health status for all adapters
   */
  async getHealthStatus(): Promise<Map<string, { status: string; responseTime?: number; error?: string }>> {
    const healthStatus = new Map<string, { status: string; responseTime?: number; error?: string }>();

    const promises = Array.from(this.adapters.values()).map(async (adapter) => {
      try {
        const startTime = Date.now();
        const health = await adapter.methods.healthCheck();
        const responseTime = Date.now() - startTime;

        healthStatus.set(adapter.platformInfo.id, {
          status: health.status,
          responseTime,
          ...(health.message && { error: health.message })
        });
      } catch (error) {
        healthStatus.set(adapter.platformInfo.id, {
          status: 'error',
          error: error instanceof Error ? error.message : 'Unknown error'
        });
      }
    });

    await Promise.all(promises);
    return healthStatus;
  }

  /**
   * Get statistics about registered adapters
   */
  getStats() {
    return {
      totalAdapters: this.adapters.size,
      activePlatforms: this.supportedPlatforms.filter(p => p.isActive).length,
      platformIds: Array.from(this.adapters.keys()),
      supportedPlatformIds: this.supportedPlatforms.map(p => p.id)
    };
  }
}

// Singleton instance
export const adapterRegistry = new AdapterRegistry();
