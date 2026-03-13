/**
 * GitHub Platform Adapter
 * 
 * Fetches bounties from GitHub Issues across multiple repositories
 * Supports bounty hunting platforms that use GitHub Issues
 * 
 * Capabilities:
 * - Multi-repository scanning
 * - Label-based filtering
 * - Pagination support
 * - Rate limit handling
 */

import {
  BasePlatformAdapter,
  IBounty,
  IFetchBountiesParams,
  IPlatformHealth,
  BountyType,
  BountyStatus
} from '../types';

/**
 * Repository configuration for bounty scanning
 */
interface RepositoryConfig {
  owner: string;
  repo: string;
  bountyLabels: string[];
  defaultReward: number;
  isActive: boolean;
}

/**
 * GitHub Adapter Configuration
 */
interface GitHubAdapterConfig {
  githubApiUrl: string;
  repositories: RepositoryConfig[];
  maxBountiesPerRepo: number;
  timeout: number;
  userAgent: string;
  rateLimitDelay: number; // ms between requests
}

/**
 * GitHub Issue API Response
 */
interface GitHubIssue {
  id: number;
  number: number;
  title: string;
  body: string | null;
  state: 'open' | 'closed';
  created_at: string;
  updated_at: string;
  closed_at: string | null;
  labels: Array<{ name: string; color: string }>;
  user: { login: string; avatar_url?: string };
  assignee: { login: string } | null;
  comments: number;
  html_url: string;
}

/**
 * GitHub Platform Adapter Implementation
 */
export class GitHubAdapter extends BasePlatformAdapter {
  readonly platformInfo = {
    id: 'github',
    name: 'GitHub',
    url: 'https://github.com',
    description: 'GitHub Issues bounty tracking across multiple repositories',
    logoUrl: 'https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png',
    isActive: true,
    supportedTypes: ['DEV', 'GRANT', 'BUG'] as BountyType[]
  };

  readonly capabilities = {
    // Data fields support
    supportsTitle: true,
    supportsDescription: true,
    supportsRewardAmount: true,
    supportsRewardCurrency: true,
    supportsDeadline: false,
    supportsTags: true,
    supportsOrganization: true,
    supportsOrganizationLogo: false,
    supportsDifficulty: false,
    supportsRemote: false,
    supportsContributorCount: true,
    supportsSubmissionsCount: false,

    // Action support
    supportsApplications: false,
    supportsSubmissions: false,
    supportsMessaging: false,
    supportsDirectApply: false,

    // Filtering support
    supportsTypeFilter: false,
    supportsStatusFilter: true,
    supportsRewardFilter: false,
    supportsTagFilter: true,
    supportsDateFilter: false,

    // Sorting support
    supportsSortByReward: false,
    supportsSortByDate: true,
    supportsSortByDeadline: false,
    supportsSortByPopularity: true,

    // Pagination support
    supportsPagination: true,
    maxPageSize: 100,
    defaultPageSize: 30,

    // Rate limiting
    rateLimitPerMinute: 30,
    requiresAuth: false,
    authMethod: 'none'
  } as const;

  readonly version = '1.0.0';
  readonly lastUpdated = new Date('2026-03-13');

  private config: GitHubAdapterConfig = {
    githubApiUrl: 'https://api.github.com',
    repositories: [
      {
        owner: 'bolivian-peru',
        repo: 'marketplace-service-template',
        bountyLabels: ['bounty'],
        defaultReward: 50,
        isActive: true
      },
      {
        owner: 'superteamwin',
        repo: 'bounties',
        bountyLabels: ['bounty', 'paid'],
        defaultReward: 100,
        isActive: true
      }
    ],
    maxBountiesPerRepo: 50,
    timeout: 10000,
    userAgent: 'BountyOS/1.0 (GitHub Bounty Aggregator)',
    rateLimitDelay: 2000
  };

  /**
   * Methods implementation
   */
  readonly methods = {
    fetchBounties: (params: IFetchBountiesParams) => this.fetchBounties(params),
    fetchBounty: (id: string) => this.fetchBounty(id),
    healthCheck: () => this.healthCheck()
  };

  /**
   * Fetch bounties from all configured repositories
   */
  async fetchBounties(params: IFetchBountiesParams = {}): Promise<IBounty[]> {
    const { limit = this.config.maxBountiesPerRepo * this.config.repositories.length } = params;
    
    const allBounties: IBounty[] = [];

    // Fetch from each repository
    for (const repo of this.config.repositories) {
      if (!repo.isActive) continue;

      try {
        const bounties = await this.fetchFromRepository(repo, limit / this.config.repositories.length);
        allBounties.push(...bounties);

        // Rate limiting
        if (this.config.repositories.indexOf(repo) < this.config.repositories.length - 1) {
          await this.delay(this.config.rateLimitDelay);
        }
      } catch (error) {
        console.error(`Error fetching from ${repo.owner}/${repo.repo}:`, error);
      }
    }

    // Sort by creation date (newest first) and limit
    const sorted = allBounties.sort((a, b) => 
      new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
    );

    return this.validateBounties(sorted.slice(0, limit));
  }

  /**
   * Fetch bounties from a specific repository
   */
  private async fetchFromRepository(repo: RepositoryConfig, limit: number): Promise<IBounty[]> {
    const url = new URL(`${this.config.githubApiUrl}/repos/${repo.owner}/${repo.repo}/issues`);
    url.searchParams.append('labels', repo.bountyLabels.join(','));
    url.searchParams.append('state', 'open');
    url.searchParams.append('per_page', Math.min(limit, 100).toString());
    url.searchParams.append('sort', 'created');
    url.searchParams.append('direction', 'desc');

    const response = await fetch(url.toString(), {
      headers: {
        'Accept': 'application/vnd.github.v3+json',
        'User-Agent': this.config.userAgent
      },
      signal: AbortSignal.timeout(this.config.timeout)
    });

    if (!response.ok) {
      throw new Error(`GitHub API error: ${response.status}`);
    }

    const issues: GitHubIssue[] = await response.json();
    
    return issues.map(issue => this.transformIssueToBounty(issue, repo));
  }

  /**
   * Fetch single bounty by ID
   */
  async fetchBounty(id: string): Promise<IBounty | null> {
    try {
      // Extract repository and issue number from ID
      const parts = this.extractBountyInfo(id);
      if (!parts) {
        return null;
      }

      const { owner, repo, issueNumber } = parts;

      const url = `${this.config.githubApiUrl}/repos/${owner}/${repo}/issues/${issueNumber}`;

      const response = await fetch(url, {
        headers: {
          'Accept': 'application/vnd.github.v3+json',
          'User-Agent': this.config.userAgent
        },
        signal: AbortSignal.timeout(this.config.timeout)
      });

      if (!response.ok) {
        if (response.status === 404) {
          return null;
        }
        throw new Error(`GitHub API error: ${response.status}`);
      }

      const issue: GitHubIssue = await response.json();
      const repoConfig = this.config.repositories.find(
        r => r.owner === owner && r.repo === repo
      );

      return this.transformIssueToBounty(issue, repoConfig || this.config.repositories[0]);
    } catch (error) {
      console.error('GitHub adapter error fetching single bounty:', error);
      return null;
    }
  }

  /**
   * Health check for GitHub adapter
   */
  async healthCheck(): Promise<IPlatformHealth> {
    try {
      const startTime = Date.now();
      
      // Check rate limit
      const rateLimitResponse = await fetch(`${this.config.githubApiUrl}/rate_limit`, {
        headers: {
          'Accept': 'application/vnd.github.v3+json',
          'User-Agent': this.config.userAgent
        },
        signal: AbortSignal.timeout(5000)
      });

      const responseTime = Date.now() - startTime;

      if (!rateLimitResponse.ok) {
        return {
          status: 'degraded',
          responseTime,
          message: 'Cannot access rate limit endpoint'
        };
      }

      const rateLimitData = await rateLimitResponse.json();
      const remaining = rateLimitData.resources?.core?.remaining || 0;

      if (remaining < 10) {
        return {
          status: 'degraded',
          responseTime,
          message: `Rate limit almost exceeded: ${remaining} remaining`
        };
      }

      // Check if we can fetch from first repo
      const firstRepo = this.config.repositories[0];
      const issuesResponse = await fetch(
        `${this.config.githubApiUrl}/repos/${firstRepo.owner}/${firstRepo.repo}/issues?per_page=1`,
        {
          headers: {
            'Accept': 'application/vnd.github.v3+json',
            'User-Agent': this.config.userAgent
          },
          signal: AbortSignal.timeout(5000)
        }
      );

      if (!issuesResponse.ok) {
        return {
          status: 'degraded',
          responseTime,
          message: 'Cannot fetch issues from primary repository'
        };
      }

      return {
        status: 'healthy',
        responseTime,
        bountyCount: this.config.repositories.length * this.config.maxBountiesPerRepo
      };
    } catch (error) {
      return {
        status: 'down',
        message: error instanceof Error ? error.message : 'Unknown error'
      };
    }
  }

  /**
   * Transform GitHub issue to bounty format
   */
  private transformIssueToBounty(issue: GitHubIssue, repo: RepositoryConfig): IBounty {
    const rewardAmount = this.parseReward(issue.title, issue.body, repo.defaultReward);
    const type = this.determineBountyType(issue);
    const status = this.determineStatus(issue);
    const tags = [...issue.labels.map(l => l.name), repo.repo];

    return {
      id: `github-${repo.owner}-${repo.repo}-${issue.id}`,
      type,
      title: issue.title,
      description: this.cleanDescription(issue.body),
      rewardAmount,
      rewardCurrency: 'USD',
      rewardUsdEquivalent: rewardAmount,
      status,
      platform: this.platformInfo,
      tags,
      url: issue.html_url,
      createdAt: issue.created_at,
      updatedAt: issue.updated_at,
      organization: `${repo.owner}/${repo.repo}`,
      contributorCount: issue.comments,
      metadata: {
        issueNumber: issue.number,
        githubId: issue.id,
        author: issue.user.login,
        authorAvatar: issue.user.avatar_url,
        comments: issue.comments,
        isAssigned: issue.assignee !== null,
        assignee: issue.assignee?.login,
        repository: `${repo.owner}/${repo.repo}`
      },
      requirements: this.extractRequirements(issue.body),
      _raw: issue,
      _adapterVersion: this.version,
      _fetchedAt: new Date().toISOString(),
      _supportedFields: this.getSupportedFields(),
      _unsupportedFields: this.getUnsupportedFields()
    };
  }

  /**
   * Parse reward from title/body
   */
  private parseReward(title: string, body: string | null, defaultReward: number): number {
    const text = `${title} ${body || ''}`.toLowerCase();
    
    // Match patterns like "$50", "$100", "$500"
    const dollarMatch = text.match(/\$(\d+)/);
    if (dollarMatch && dollarMatch[1]) {
      return parseInt(dollarMatch[1], 10);
    }

    return defaultReward;
  }

  /**
   * Determine bounty type
   */
  private determineBountyType(issue: GitHubIssue): BountyType {
    const text = `${issue.title} ${issue.body || ''}`.toLowerCase();
    const labels = issue.labels.map(l => l.name.toLowerCase());

    if (labels.includes('bug') || text.includes('bug')) {
      return 'DEV';
    }
    if (labels.includes('grant') || text.includes('grant')) {
      return 'GRANT';
    }
    if (labels.includes('security') || text.includes('security')) {
      return 'DEV'; // Could be 'SEC' if we add it
    }

    return 'DEV';
  }

  /**
   * Determine status
   */
  private determineStatus(issue: GitHubIssue): BountyStatus {
    if (issue.state === 'closed') {
      return 'COMPLETED';
    }
    if (issue.assignee) {
      return 'IN_PROGRESS';
    }
    return 'OPEN';
  }

  /**
   * Clean description
   */
  private cleanDescription(body: string | null): string {
    if (!body) {
      return 'No description provided.';
    }

    let cleaned = body.replace(/```[\s\S]*?```/g, '');
    
    if (cleaned.length > 5000) {
      cleaned = cleaned.substring(0, 5000) + '...';
    }

    return cleaned.trim();
  }

  /**
   * Extract requirements
   */
  private extractRequirements(body: string | null) {
    if (!body) return undefined;

    const requirementsMatch = body.match(/##.*?(?:What to Build|Requirements|Deliverables)[\s\S]*?(?=##|$)/i);
    
    if (requirementsMatch) {
      return { description: requirementsMatch[0].trim() };
    }

    return undefined;
  }

  /**
   * Extract bounty info from ID
   */
  private extractBountyInfo(bountyId: string): { owner: string; repo: string; issueNumber: number } | null {
    // Format: "github-owner-repo-123456"
    const match = bountyId.match(/github-([^-]+)-([^-]+)-(\d+)$/);
    if (match) {
      return {
        owner: match[1],
        repo: match[2],
        issueNumber: parseInt(match[3], 10)
      };
    }
    return null;
  }

  /**
   * Delay helper
   */
  private delay(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  /**
   * Validate bounties
   */
  private validateBounties(bounties: IBounty[]): IBounty[] {
    return bounties.map(bounty => this.validateBounty(bounty));
  }

  /**
   * Validate single bounty
   */
  private validateBounty(bounty: IBounty): IBounty {
    return {
      ...bounty,
      _adapterVersion: this.version,
      _fetchedAt: new Date().toISOString(),
      _supportedFields: this.getSupportedFields(),
      _unsupportedFields: this.getUnsupportedFields()
    };
  }

  /**
   * Get supported fields
   */
  private getSupportedFields(): string[] {
    return [
      'id', 'type', 'title', 'description', 'rewardAmount', 'rewardCurrency',
      'status', 'platform', 'tags', 'url', 'createdAt', 'updatedAt',
      'organization', 'contributorCount', 'metadata'
    ];
  }

  /**
   * Get unsupported fields
   */
  private getUnsupportedFields(): string[] {
    return [
      'deadline', 'organizationLogo', 'difficulty', 'remote',
      'submissionsCount', 'isFeatured', 'isUrgent'
    ];
  }
}

// Export singleton instance
export const githubAdapter = new GitHubAdapter();
