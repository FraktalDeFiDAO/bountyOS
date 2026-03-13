/**
 * Platform Adapter Architecture
 * 
 * Granular interface system for platform integration with
 * nil method support for unsupported features
 */

// ============================================================================
// CORE INTERFACES
// ============================================================================

/**
 * Base platform information
 */
export interface IPlatformInfo {
  readonly id: string;
  readonly name: string;
  readonly url: string;
  readonly description?: string;
  readonly logoUrl?: string;
  readonly isActive: boolean;
  readonly supportedTypes: BountyType[];
}

/**
 * Granular feature support indicators
 * Each field indicates whether the platform supports that specific feature
 */
export interface IPlatformCapabilities {
  // Data fields support
  readonly supportsTitle: boolean;
  readonly supportsDescription: boolean;
  readonly supportsRewardAmount: boolean;
  readonly supportsRewardCurrency: boolean;
  readonly supportsDeadline: boolean;
  readonly supportsTags: boolean;
  readonly supportsOrganization: boolean;
  readonly supportsOrganizationLogo: boolean;
  readonly supportsDifficulty: boolean;
  readonly supportsRemote: boolean;
  readonly supportsContributorCount: boolean;
  readonly supportsSubmissionsCount: boolean;
  
  // Action support
  readonly supportsApplications: boolean;
  readonly supportsSubmissions: boolean;
  readonly supportsMessaging: boolean;
  readonly supportsDirectApply: boolean;
  
  // Filtering support
  readonly supportsTypeFilter: boolean;
  readonly supportsStatusFilter: boolean;
  readonly supportsRewardFilter: boolean;
  readonly supportsTagFilter: boolean;
  readonly supportsDateFilter: boolean;
  
  // Sorting support
  readonly supportsSortByReward: boolean;
  readonly supportsSortByDate: boolean;
  readonly supportsSortByDeadline: boolean;
  readonly supportsSortByPopularity: boolean;
  
  // Pagination support
  readonly supportsPagination: boolean;
  readonly maxPageSize?: number;
  readonly defaultPageSize?: number;
  
  // Rate limiting
  readonly rateLimitPerMinute?: number;
  readonly requiresAuth: boolean;
  readonly authMethod?: 'api_key' | 'oauth' | 'jwt' | 'none';
}

/**
 * Method implementations - nil/undefined means unsupported
 */
export interface IPlatformMethods {
  // Required methods
  fetchBounties: (params: IFetchBountiesParams) => Promise<IBounty[]>;
  fetchBounty: (id: string) => Promise<IBounty | null>;
  
  // Optional methods (undefined = unsupported)
  submitApplication?: (data: IApplicationData) => Promise<IApplication | null>;
  submitWork?: (data: ISubmissionData) => Promise<ISubmission | null>;
  getTypes?: () => Promise<IBountyType[]>;
  getStatuses?: () => Promise<IBountyStatus[]>;
  getTags?: () => Promise<ITag[]>;
  search?: (query: string, params: ISearchParams) => Promise<IBounty[]>;
  
  // Health and metadata
  healthCheck: () => Promise<IPlatformHealth>;
  getLastSyncTime?: () => Promise<Date | null>;
}

/**
 * Complete platform adapter interface
 */
export interface IPlatformAdapter {
  // Platform identification
  readonly platformInfo: IPlatformInfo;
  
  // Capabilities
  readonly capabilities: IPlatformCapabilities;
  
  // Methods
  readonly methods: IPlatformMethods;
  
  // Adapter metadata
  readonly version: string;
  readonly lastUpdated: Date;
}

// ============================================================================
// BOUNTY INTERFACES
// ============================================================================

/**
 * Core bounty interface - all platforms must provide these fields
 */
export interface IBountyCore {
  readonly id: string;
  readonly type: BountyType;
  readonly title: string;
  readonly description: string;
  readonly status: BountyStatus;
  readonly url: string;
  readonly createdAt: string;
  readonly updatedAt?: string;
}

/**
 * Optional bounty fields - platforms may not support these
 */
export interface IBountyOptional {
  readonly rewardAmount?: number;
  readonly rewardCurrency?: string;
  readonly rewardUsdEquivalent?: number;
  readonly deadline?: string;
  readonly tags?: string[];
  readonly organization?: string;
  readonly organizationLogo?: string;
  readonly difficulty?: DifficultyLevel;
  readonly remote?: boolean;
  readonly contributorCount?: number;
  readonly submissionsCount?: number;
  readonly isFeatured?: boolean;
  readonly isUrgent?: boolean;
  readonly metadata?: Record<string, any>;
  readonly requirements?: BountyRequirements;
}

/**
 * Complete bounty interface
 */
export interface IBounty extends IBountyCore, IBountyOptional {
  readonly platform: IPlatformInfo;
  readonly _raw?: any; // Original platform data for debugging
  readonly _adapterVersion?: string; // Adapter version that fetched this
  readonly _fetchedAt?: string; // When this was fetched
  readonly _supportedFields: string[]; // Which fields this platform supports
  readonly _unsupportedFields: string[]; // Which fields this platform doesn't support
}

// ============================================================================
// SUPPORT INTERFACES
// ============================================================================

/**
 * Bounty type definitions
 */
export interface IBountyType {
  readonly id: string;
  readonly name: string;
  readonly description: string;
  readonly icon?: string;
  readonly color?: string;
}

/**
 * Bounty status definitions
 */
export interface IBountyStatus {
  readonly id: string;
  readonly name: string;
  readonly color: string;
  readonly icon?: string;
}

/**
 * Tag interface
 */
export interface ITag {
  readonly id: string;
  readonly name: string;
  readonly count?: number;
}

// ============================================================================
// ACTION INTERFACES
// ============================================================================

/**
 * Application data
 */
export interface IApplicationData {
  readonly bountyId: string;
  readonly motivation: string;
  readonly experience: string;
  readonly portfolioUrl?: string;
  readonly estimatedCompletionTime?: string;
  readonly metadata?: Record<string, any>;
}

/**
 * Application result
 */
export interface IApplication {
  readonly id: string;
  readonly bountyId: string;
  readonly status: 'pending' | 'approved' | 'rejected';
  readonly submittedAt: string;
  readonly metadata?: Record<string, any>;
}

/**
 * Submission data
 */
export interface ISubmissionData {
  readonly bountyId: string;
  readonly submissionUrl: string;
  readonly description: string;
  readonly metadata?: Record<string, any>;
}

/**
 * Submission result
 */
export interface ISubmission {
  readonly id: string;
  readonly bountyId: string;
  readonly status: 'pending' | 'reviewed' | 'accepted' | 'rejected';
  readonly submittedAt: string;
  readonly reviewedAt?: string;
  readonly metadata?: Record<string, any>;
}

// ============================================================================
// QUERY INTERFACES
// ============================================================================

/**
 * Fetch bounties parameters
 */
export interface IFetchBountiesParams {
  readonly page?: number;
  readonly limit?: number;
  readonly type?: BountyType[];
  readonly status?: BountyStatus[];
  readonly minReward?: number;
  readonly maxReward?: number;
  readonly tags?: string[];
  readonly sortBy?: 'reward' | 'deadline' | 'createdAt' | 'updatedAt' | 'title';
  readonly sortOrder?: 'asc' | 'desc';
  readonly searchQuery?: string;
}

/**
 * Search parameters
 */
export interface ISearchParams {
  readonly query: string;
  readonly filters?: IFetchBountiesParams;
  readonly limit?: number;
}

// ============================================================================
// HEALTH INTERFACES
// ============================================================================

/**
 * Platform health status
 */
export interface IPlatformHealth {
  readonly status: 'healthy' | 'degraded' | 'down';
  readonly lastSyncTime?: Date;
  readonly responseTimeMs?: number;
  readonly errorRate?: number;
  readonly bountyCount?: number;
  readonly message?: string;
}

// ============================================================================
// ENUMS
// ============================================================================

export type BountyType = 
  | 'DEV'      // Development
  | 'GRANT'    // Grants
  | 'MICRO'    // Microtasks
  | 'GIG'      // Freelance
  | 'HACK'     // Hackathons
  | 'AMB'      // Ambassador
  | 'RETRO'    // Retroactive
  | 'SEC'      // Security
  | 'DESIGN'   // Design
  | 'CONTENT'  // Content
  | 'AUDIT';   // Audit

export type BountyStatus =
  | 'OPEN'
  | 'IN_PROGRESS'
  | 'SUBMITTED'
  | 'REVIEW'
  | 'COMPLETED'
  | 'CANCELLED';

export type DifficultyLevel =
  | 'BEGINNER'
  | 'INTERMEDIATE'
  | 'ADVANCED'
  | 'EXPERT';

export interface BountyRequirements {
  readonly skills?: string[];
  readonly experience?: string;
  readonly deliverables?: string[];
  readonly timeline?: string;
}

// ============================================================================
// ADAPTER REGISTRY
// ============================================================================

/**
 * Registry for managing platform adapters
 */
export interface IAdapterRegistry {
  getAdapter(platformId: string): IPlatformAdapter | null;
  getAllAdapters(): IPlatformAdapter[];
  getSupportedPlatforms(): IPlatformInfo[];
  registerAdapter(adapter: IPlatformAdapter): void;
  unregisterAdapter(platformId: string): void;
  hasAdapter(platformId: string): boolean;
  getAdapterVersion(platformId: string): string | null;
}

// ============================================================================
// BASE ADAPTER CLASS
// ============================================================================

/**
 * Abstract base class for platform adapters
 * Provides default implementations for unsupported methods
 */
export abstract class BasePlatformAdapter implements IPlatformAdapter {
  abstract readonly platformInfo: IPlatformInfo;
  abstract readonly capabilities: IPlatformCapabilities;
  abstract readonly methods: IPlatformMethods;
  abstract readonly version: string;
  abstract readonly lastUpdated: Date;

  /**
   * Default implementation for unsupported methods
   * Returns null to indicate unsupported feature
   */
  protected unsupported<T>(methodName: string): T {
    console.warn(
      `Method ${methodName} is not supported by platform ${this.platformInfo.id}`
    );
    return null as T;
  }

  /**
   * Validate bounty data against platform capabilities
   */
  protected validateBounty(bounty: Partial<IBounty>): IBounty {
    const supportedFields: string[] = [];
    const unsupportedFields: string[] = [];

    // Check each optional field
    const optionalFields: (keyof IBountyOptional)[] = [
      'rewardAmount', 'rewardCurrency', 'rewardUsdEquivalent', 'deadline',
      'tags', 'organization', 'organizationLogo', 'difficulty', 'remote',
      'contributorCount', 'submissionsCount', 'isFeatured', 'isUrgent'
    ];

    optionalFields.forEach(field => {
      if (this.isFieldSupported(field)) {
        supportedFields.push(field);
      } else {
        unsupportedFields.push(field);
      }
    });

    return {
      ...bounty as IBounty,
      platform: this.platformInfo,
      _adapterVersion: this.version,
      _fetchedAt: new Date().toISOString(),
      _supportedFields: supportedFields,
      _unsupportedFields: unsupportedFields
    };
  }

  /**
   * Check if a specific field is supported by this platform
   */
  private isFieldSupported(field: keyof IBountyOptional): boolean {
    const fieldMap: Record<keyof IBountyOptional, keyof IPlatformCapabilities> = {
      rewardAmount: 'supportsRewardAmount',
      rewardCurrency: 'supportsRewardCurrency',
      rewardUsdEquivalent: 'supportsRewardAmount',
      deadline: 'supportsDeadline',
      tags: 'supportsTags',
      organization: 'supportsOrganization',
      organizationLogo: 'supportsOrganizationLogo',
      difficulty: 'supportsDifficulty',
      remote: 'supportsRemote',
      contributorCount: 'supportsContributorCount',
      submissionsCount: 'supportsSubmissionsCount',
      isFeatured: 'supportsTitle', // Assume title support = can mark featured
      isUrgent: 'supportsTitle',
      metadata: 'supportsDescription', // Assume description support = can have metadata
      requirements: 'supportsDescription'
    };

    const capabilityKey = fieldMap[field];
    return capabilityKey ? this.capabilities[capabilityKey] : false;
  }
}

// ============================================================================
// TYPE GUARDS
// ============================================================================

/**
 * Type guard to check if a method is supported
 */
export function isMethodSupported<T extends Function>(
  method: T | undefined,
  methodName: string
): method is T {
  if (method === undefined) {
    console.warn(`Method ${methodName} is not supported`);
    return false;
  }
  return true;
}

/**
 * Type guard to check if a bounty field is populated
 */
export function isFieldPopulated<T extends keyof IBountyOptional>(
  bounty: IBounty,
  field: T
): bounty is IBounty & Required<Pick<IBountyOptional, T>> {
  return bounty[field] !== undefined && bounty[field] !== null;
}

// ============================================================================
// EXPORTS
// ============================================================================

export {
  IPlatformInfo,
  IPlatformCapabilities,
  IPlatformMethods,
  IPlatformAdapter,
  IBountyCore,
  IBountyOptional,
  IBounty,
  IBountyType,
  IBountyStatus,
  ITag,
  IApplicationData,
  IApplication,
  ISubmissionData,
  ISubmission,
  IFetchBountiesParams,
  ISearchParams,
  IPlatformHealth,
  IAdapterRegistry,
  BasePlatformAdapter,
  isMethodSupported,
  isFieldPopulated
};
