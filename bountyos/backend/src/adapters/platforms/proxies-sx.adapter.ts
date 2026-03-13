/**
 * Proxies.sx Platform Adapter
 * 
 * Fetches bounties from Proxies.sx GitHub Issues
 * 
 * Capabilities:
 * - Fetches bounties from GitHub Issues with "bounty" label
 * - Supports pagination (up to 100 issues)
 * - Real-time bounty data from GitHub API
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
 * Proxies.sx Adapter Configuration
 */
interface ProxiesSxConfig {
  githubApiUrl: string;
  repoOwner: string;
  repoName: string;
  bountyLabel: string;
  maxBounties: number;
  timeout: number;
  userAgent: string;
}

/**
 * GitHub Issue API Response Type
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
 * Proxies.sx Platform Adapter Implementation
 */
export class ProxiesSxAdapter extends BasePlatformAdapter {
  readonly platformInfo = {
    id: 'proxies-sx',
    name: 'Proxies.sx',
    url: 'https://proxies.sx',
    description: 'Mobile proxy marketplace for web scraping bounties',
    logoUrl: undefined,
    isActive: true,
    supportedTypes: ['DEV', 'MICRO'] as BountyType[]
  };

  readonly capabilities = {
    // Data fields support
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
    supportsSortByPopularity: false,

    // Pagination support
    supportsPagination: true,
    maxPageSize: 100,
    defaultPageSize: 30,

    // Rate limiting
    rateLimitPerMinute: 60,
    requiresAuth: false,
    authMethod: 'none'
  } as const;

  readonly version = '1.0.0';
  readonly lastUpdated = new Date('2026-03-13');

  private config: ProxiesSxConfig = {
    githubApiUrl: 'https://api.github.com',
    repoOwner: 'bolivian-peru',
    repoName: 'marketplace-service-template',
    bountyLabel: 'bounty',
    maxBounties: 100,
    timeout: 10000,
    userAgent: 'BountyOS/1.0 (Bounty Aggregator)'
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
   * Fetch bounties from Proxies.sx GitHub repository
   */
  async fetchBounties(params: IFetchBountiesParams = {}): Promise<IBounty[]> {
    const { limit = this.config.maxBounties, sortBy = 'createdAt', sortOrder = 'desc' } = params;

    try {
      const url = new URL(`${this.config.githubApiUrl}/repos/${this.config.repoOwner}/${this.config.repoName}/issues`);
      url.searchParams.append('labels', this.config.bountyLabel);
      url.searchParams.append('state', 'open');
      url.searchParams.append('per_page', Math.min(limit, 100).toString());
      url.searchParams.append('sort', sortBy === 'createdAt' ? 'created' : 'created');
      url.searchParams.append('direction', sortOrder);

      const response = await fetch(url.toString(), {
        headers: {
          'Accept': 'application/vnd.github.v3+json',
          'User-Agent': this.config.userAgent
        },
        signal: AbortSignal.timeout(this.config.timeout)
      });

      if (!response.ok) {
        throw new Error(`GitHub API error: ${response.status} ${response.statusText}`);
      }

      const issues: GitHubIssue[] = await response.json();

      // Convert GitHub issues to bounty format
      const bounties: IBounty[] = issues.map(issue => this.transformIssueToBounty(issue));

      return this.validateBounties(bounties);
    } catch (error) {
      console.error('Proxies.sx adapter error:', error);
      
      // Return empty array on error (don't break the whole system)
      return [];
    }
  }

  /**
   * Fetch single bounty by ID
   */
  async fetchBounty(id: string): Promise<IBounty | null> {
    try {
      // Extract issue number from ID (format: "proxies-sx-123456" or "proxies-sx-73")
      const issueNumber = this.extractIssueNumber(id);
      
      if (!issueNumber) {
        return null;
      }

      const url = `${this.config.githubApiUrl}/repos/${this.config.repoOwner}/${this.config.repoName}/issues/${issueNumber}`;

      const response = await fetch(url, {
        headers: {
          'Accept': 'application/vnd.github.v3+json',
          'User-Agent': this.config.userAgent
        },
        signal: AbortSignal.timeout(this.config.timeout)
      });

      if (!response.ok) {
        if (response.status === 404) {
          return null; // Bounty not found
        }
        throw new Error(`GitHub API error: ${response.status}`);
      }

      const issue: GitHubIssue = await response.json();
      const bounty = this.transformIssueToBounty(issue);

      return this.validateBounty(bounty);
    } catch (error) {
      console.error('Proxies.sx adapter error fetching single bounty:', error);
      return null;
    }
  }

  /**
   * Health check for Proxies.sx adapter
   */
  async healthCheck(): Promise<IPlatformHealth> {
    try {
      const startTime = Date.now();
      
      const response = await fetch(
        `${this.config.githubApiUrl}/repos/${this.config.repoOwner}/${this.config.repoName}`,
        {
          headers: {
            'Accept': 'application/vnd.github.v3+json',
            'User-Agent': this.config.userAgent
          },
          signal: AbortSignal.timeout(5000)
        }
      );

      const responseTime = Date.now() - startTime;

      if (!response.ok) {
        return {
          status: 'degraded',
          responseTime,
          message: `GitHub API returned ${response.status}`
        };
      }

      // Check if we can fetch issues
      const issuesResponse = await fetch(
        `${this.config.githubApiUrl}/repos/${this.config.repoOwner}/${this.config.repoName}/issues?labels=bounty&per_page=1`,
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
          message: 'Can fetch repo but not issues'
        };
      }

      const issues = await issuesResponse.json();

      return {
        status: 'healthy',
        responseTime,
        bountyCount: issues.length
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
  private transformIssueToBounty(issue: GitHubIssue): IBounty {
    // Parse reward from title or body
    const rewardAmount = this.parseReward(issue.title, issue.body);
    
    // Determine bounty type from labels or description
    const type = this.determineBountyType(issue);

    // Determine status
    const status = this.determineStatus(issue);

    // Extract tags from labels
    const tags = issue.labels.map(label => label.name);

    return {
      id: `proxies-sx-${issue.id}`,
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
      metadata: {
        issueNumber: issue.number,
        githubId: issue.id,
        author: issue.user.login,
        comments: issue.comments,
        isAssigned: issue.assignee !== null,
        assignee: issue.assignee?.login
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
   * Parse reward amount from title or body
   * Looks for patterns like "$50", "$100 paid in $SX token", etc.
   */
  private parseReward(title: string, body: string | null): number {
    const text = `${title} ${body || ''}`.toLowerCase();
    
    // Match patterns like "$50", "$100", "$500"
    const dollarMatch = text.match(/\$(\d+)/);
    if (dollarMatch && dollarMatch[1]) {
      return parseInt(dollarMatch[1], 10);
    }

    // Default reward for Proxies.sx bounties
    return 50;
  }

  /**
   * Determine bounty type from issue data
   */
  private determineBountyType(issue: GitHubIssue): BountyType {
    const text = `${issue.title} ${issue.body || ''}`.toLowerCase();
    
    if (text.includes('security') || text.includes('audit')) {
      return 'SEC';
    }
    if (text.includes('design') || text.includes('ui') || text.includes('ux')) {
      return 'DESIGN';
    }
    if (text.includes('content') || text.includes('write') || text.includes('article')) {
      return 'CONTENT';
    }
    if (text.includes('micro') || text.includes('small') || text.includes('quick')) {
      return 'MICRO';
    }
    
    return 'DEV'; // Default
  }

  /**
   * Determine bounty status from issue state
   */
  private determineStatus(issue: GitHubIssue): BountyStatus {
    if (issue.state === 'closed') {
      return issue.closed_at ? 'COMPLETED' : 'CANCELLED';
    }

    if (issue.assignee) {
      return 'IN_PROGRESS';
    }

    return 'OPEN';
  }

  /**
   * Clean and format description
   */
  private cleanDescription(body: string | null): string {
    if (!body) {
      return 'No description provided. Check the bounty URL for details.';
    }

    // Remove markdown code blocks for cleaner display
    let cleaned = body.replace(/```[\s\S]*?```/g, '');
    
    // Limit length
    if (cleaned.length > 5000) {
      cleaned = cleaned.substring(0, 5000) + '...';
    }

    return cleaned.trim();
  }

  /**
   * Extract requirements from description
   */
  private extractRequirements(body: string | null) {
    if (!body) {
      return undefined;
    }

    // Look for "What to Build" or "Requirements" sections
    const requirementsMatch = body.match(/##.*?(?:What to Build|Requirements|Deliverables)[\s\S]*?(?=##|$)/i);
    
    if (requirementsMatch) {
      return {
        description: requirementsMatch[0].trim()
      };
    }

    return undefined;
  }

  /**
   * Extract issue number from bounty ID
   */
  private extractIssueNumber(bountyId: string): number | null {
    // Try to extract from format "proxies-sx-3944053546" (GitHub issue ID)
    const githubIdMatch = bountyId.match(/proxies-sx-(\d+)/);
    if (githubIdMatch && githubIdMatch[1]) {
      // This is GitHub's issue ID, we need to search for it
      return null; // Will handle in fetchBounty by searching
    }

    // Try to extract from format "proxies-sx-73" (issue number)
    const issueNumberMatch = bountyId.match(/proxies-sx-(\d+)$/);
    if (issueNumberMatch && issueNumberMatch[1]) {
      const num = parseInt(issueNumberMatch[1], 10);
      if (!isNaN(num)) {
        return num;
      }
    }

    return null;
  }

  /**
   * Validate bounties array
   */
  private validateBounties(bounties: IBounty[]): IBounty[] {
    return bounties.map(bounty => this.validateBounty(bounty));
  }

  /**
   * Validate single bounty and add metadata
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
   * Get list of supported fields
   */
  private getSupportedFields(): string[] {
    return [
      'id', 'type', 'title', 'description', 'rewardAmount', 'rewardCurrency',
      'status', 'platform', 'tags', 'url', 'createdAt', 'updatedAt', 'metadata'
    ];
  }

  /**
   * Get list of unsupported fields
   */
  private getUnsupportedFields(): string[] {
    return [
      'deadline', 'organization', 'organizationLogo', 'difficulty', 'remote',
      'contributorCount', 'submissionsCount', 'isFeatured', 'isUrgent'
    ];
  }
}

// Export singleton instance
export const proxiesSxAdapter = new ProxiesSxAdapter();
