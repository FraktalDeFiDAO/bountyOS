# BountyOS API Specification

**Version:** 2.0.0  
**Status:** Draft  
**Last Updated:** March 13, 2026

---

## Table of Contents

1. [Overview](#overview)
2. [Authentication](#authentication)
3. [Endpoints](#endpoints)
   - [Health](#health)
   - [Bounties](#bounties)
   - [Platforms](#platforms)
   - [Submissions](#submissions)
   - [Applications](#applications)
   - [Users](#users)
   - [Analytics](#analytics)
4. [Data Models](#data-models)
5. [Error Handling](#error-handling)
6. [Rate Limiting](#rate-limiting)
7. [Webhooks](#webhooks)
8. [Platform Adapters](#platform-adapters)
9. [Testing & Validation](#testing--validation)

---

## Overview

### Base URL
```
Development: http://localhost:8000
Production: https://api.bountyos.com
```

### API Versioning
- URL prefix: `/api/v{version}`
- Current version: `v1`
- Deprecation policy: 6 months notice

### Response Format
All responses follow this structure:
```typescript
interface ApiResponse<T> {
  data: T;
  meta?: {
    total?: number;
    page?: number;
    limit?: number;
    totalPages?: number;
    timestamp?: string;
  };
  error?: {
    code: string;
    message: string;
    details?: Record<string, any>;
  };
}
```

---

## Authentication

### API Keys
```http
Authorization: Bearer {api_key}
```

### JWT Tokens
```http
Authorization: Bearer {jwt_token}
```

---

## Endpoints

### Health

#### GET /health
**Description:** API health check

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2026-03-13T21:00:00.000Z",
  "version": "2.0.0",
  "platforms": ["gitcoin", "superteam", "algora", "proxies-sx", "code4rena"]
}
```

---

### Bounties

#### GET /api/bounties
**Description:** Get all bounties with pagination and filters

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | number | 1 | Page number |
| `limit` | number | 20 | Items per page (max: 100) |
| `type` | string[] | - | Filter by type (DEV, GRANT, etc.) |
| `platform` | string[] | - | Filter by platform ID |
| `status` | string[] | - | Filter by status |
| `minReward` | number | - | Minimum reward amount |
| `maxReward` | number | - | Maximum reward amount |
| `tags` | string | - | Comma-separated tags |
| `sortBy` | string | createdAt | Sort field |
| `sortOrder` | string | desc | Sort order (asc/desc) |
| `search` | string | - | Full-text search query |

**Response:**
```json
{
  "data": [
    {
      "id": "proxies-sx-3944053546",
      "type": "DEV",
      "title": "[BOUNTY] X/Twitter Real-Time Search API",
      "description": "...",
      "rewardAmount": 100,
      "rewardCurrency": "USD",
      "rewardUsdEquivalent": 100,
      "status": "open",
      "platform": {
        "id": "proxies-sx",
        "name": "Proxies.sx",
        "url": "https://proxies.sx"
      },
      "tags": ["bounty", "api", "twitter"],
      "url": "https://github.com/...",
      "createdAt": "2026-02-15T13:03:38Z",
      "deadline": "2026-04-15T00:00:00Z",
      "metadata": {...},
      "requirements": {...}
    }
  ],
  "meta": {
    "total": 150,
    "page": 1,
    "limit": 20,
    "totalPages": 8,
    "platforms": {
      "gitcoin": 25,
      "superteam": 30,
      "algora": 20,
      "proxies-sx": 50,
      "code4rena": 25
    }
  }
}
```

**Status Codes:**
- `200` - Success
- `400` - Invalid parameters
- `429` - Rate limit exceeded

---

#### GET /api/bounties/:id
**Description:** Get single bounty by ID

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Bounty ID |

**Response:**
```json
{
  "data": {
    "id": "proxies-sx-3944053546",
    "type": "DEV",
    "title": "[BOUNTY] X/Twitter Real-Time Search API",
    "description": "...",
    "rewardAmount": 100,
    "rewardCurrency": "USD",
    "status": "open",
    "platform": {...},
    "submissions": [],
    "applications": []
  }
}
```

**Status Codes:**
- `200` - Success
- `404` - Bounty not found

---

#### GET /api/bounties/featured
**Description:** Get featured bounties (top paying, urgent)

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `limit` | number | 6 | Number of bounties |

**Response:**
```json
{
  "data": [
    {
      "id": "...",
      "isFeatured": true,
      "isUrgent": false,
      ...
    }
  ]
}
```

---

#### GET /api/bounties/urgent
**Description:** Get urgent bounties (deadline soon)

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `limit` | number | 6 | Number of bounties |
| `hours` | number | 72 | Hours until deadline |

---

#### GET /api/bounties/filters
**Description:** Get available filter options

**Response:**
```json
{
  "data": {
    "types": [
      { "id": "DEV", "name": "Development", "count": 45 },
      { "id": "GRANT", "name": "Grants", "count": 30 },
      { "id": "MICRO", "name": "Microtasks", "count": 25 },
      { "id": "SEC", "name": "Security", "count": 20 }
    ],
    "platforms": [
      { "id": "gitcoin", "name": "Gitcoin", "count": 25 },
      { "id": "superteam", "name": "Superteam", "count": 30 },
      { "id": "algora", "name": "Algora", "count": 20 },
      { "id": "proxies-sx", "name": "Proxies.sx", "count": 50 },
      { "id": "code4rena", "name": "Code4rena", "count": 25 }
    ],
    "statuses": [
      { "id": "open", "name": "Open", "count": 120 },
      { "id": "in_progress", "name": "In Progress", "count": 20 },
      { "id": "completed", "name": "Completed", "count": 10 }
    ],
    "tags": [
      { "id": "solidity", "name": "Solidity", "count": 15 },
      { "id": "react", "name": "React", "count": 20 },
      { "id": "rust", "name": "Rust", "count": 10 }
    ]
  }
}
```

---

#### GET /api/bounties/stats
**Description:** Get bounty statistics

**Response:**
```json
{
  "data": {
    "total": 150,
    "byStatus": {
      "open": 120,
      "in_progress": 20,
      "completed": 10
    },
    "byType": {
      "DEV": 45,
      "GRANT": 30,
      "MICRO": 25,
      "SEC": 20
    },
    "byPlatform": {
      "gitcoin": 25,
      "superteam": 30,
      "algora": 20,
      "proxies-sx": 50,
      "code4rena": 25
    },
    "averageReward": 250.50,
    "medianReward": 150.00,
    "totalRewardValue": 37575.00
  }
}
```

---

#### GET /api/bounties/types
**Description:** Get bounty type definitions

**Response:**
```json
[
  {
    "id": "DEV",
    "name": "Development",
    "description": "Feature development, OSS contributions",
    "icon": "code"
  },
  {
    "id": "GRANT",
    "name": "Grants",
    "description": "Protocol/ecosystem grants",
    "icon": "grant"
  },
  {
    "id": "MICRO",
    "name": "Microtasks",
    "description": "Small tasks, surveys, testing",
    "icon": "task"
  },
  {
    "id": "SEC",
    "name": "Security",
    "description": "Bug bounties, security audits",
    "icon": "security"
  }
]
```

---

#### GET /api/bounties/platforms
**Description:** Get platform definitions

**Response:**
```json
[
  {
    "id": "gitcoin",
    "name": "Gitcoin",
    "url": "https://gitcoin.co",
    "types": ["DEV", "GRANT", "HACK"],
    "description": "Open source funding platform",
    "logoUrl": "https://..."
  },
  {
    "id": "superteam",
    "name": "Superteam",
    "url": "https://superteam.fun",
    "types": ["DEV", "GRANT", "AMB"],
    "description": "Solana ecosystem bounties",
    "logoUrl": "https://..."
  }
]
```

---

#### POST /api/bounties
**Description:** Create a new bounty (Admin only)

**Headers:**
```http
Authorization: Bearer {admin_token}
```

**Body:**
```json
{
  "platformId": "proxies-sx",
  "type": "DEV",
  "externalId": "github-issue-123",
  "title": "New Bounty Title",
  "description": "Detailed description...",
  "rewardAmount": 100,
  "rewardCurrency": "USD",
  "rewardUsdEquivalent": 100,
  "status": "open",
  "deadline": "2026-04-15T00:00:00Z",
  "metadata": {...},
  "tags": ["api", "backend"],
  "requirements": {...}
}
```

**Response:**
```json
{
  "data": {
    "id": "...",
    ...
  }
}
```

**Status Codes:**
- `201` - Created
- `400` - Invalid data
- `401` - Unauthorized
- `403` - Forbidden

---

#### PUT /api/bounties/:id
**Description:** Update an existing bounty (Admin only)

**Body:**
```json
{
  "title": "Updated Title",
  "rewardAmount": 150,
  "status": "in_progress"
}
```

**Status Codes:**
- `200` - Updated
- `404` - Not found
- `401` - Unauthorized

---

#### DELETE /api/bounties/:id
**Description:** Delete a bounty (Admin only)

**Status Codes:**
- `204` - Deleted
- `404` - Not found
- `401` - Unauthorized

---

### Platforms

#### GET /api/platforms
**Description:** Get all platforms

**Response:** See `/api/bounties/platforms`

---

#### GET /api/platforms/:id
**Description:** Get platform by ID

**Response:**
```json
{
  "data": {
    "id": "proxies-sx",
    "name": "Proxies.sx",
    "url": "https://proxies.sx",
    "description": "...",
    "isActive": true,
    "bountyCount": 50,
    "totalRewards": 5000,
    "averageReward": 100,
    "lastSyncedAt": "2026-03-13T21:00:00Z"
  }
}
```

---

#### POST /api/platforms/:id/sync
**Description:** Trigger platform sync

**Response:**
```json
{
  "data": {
    "platformId": "proxies-sx",
    "status": "syncing",
    "startedAt": "2026-03-13T21:00:00Z",
    "estimatedCompletion": "2026-03-13T21:05:00Z"
  }
}
```

---

### Submissions

#### GET /api/bounties/:id/submissions
**Description:** Get submissions for a bounty

**Response:**
```json
{
  "data": [
    {
      "id": "...",
      "bountyId": "...",
      "userId": "...",
      "submissionUrl": "https://...",
      "description": "...",
      "status": "pending",
      "submittedAt": "2026-03-13T21:00:00Z"
    }
  ]
}
```

---

#### POST /api/bounties/:id/submissions
**Description:** Submit work for a bounty

**Body:**
```json
{
  "submissionUrl": "https://github.com/...",
  "description": "Description of work..."
}
```

---

### Applications

#### GET /api/bounties/:id/applications
**Description:** Get applications for a bounty

---

#### POST /api/bounties/:id/apply
**Description:** Apply for a bounty

**Body:**
```json
{
  "motivation": "Why I'm interested...",
  "experience": "My relevant experience..."
}
```

---

### Users

#### GET /api/user/stats
**Description:** Get current user's statistics

**Response:**
```json
{
  "data": {
    "totalApplied": 15,
    "totalCompleted": 5,
    "totalEarned": 2500,
    "activeSubmissions": 3,
    "successRate": 0.33
  }
}
```

---

#### GET /api/user/bounties
**Description:** Get user's bounty activity

---

### Analytics

#### GET /api/analytics/overview
**Description:** Get platform analytics overview

**Response:**
```json
{
  "data": {
    "totalBounties": 150,
    "totalValue": 37575,
    "activeUsers": 250,
    "completedThisMonth": 25,
    "averageTimeToComplete": 14.5
  }
}
```

---

#### GET /api/analytics/trends
**Description:** Get bounty trends over time

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `period` | string | day/week/month/year |
| `startDate` | string | Start date |
| `endDate` | string | End date |

---

## Data Models

### Bounty
```typescript
interface Bounty {
  id: string;
  type: BountyType;
  title: string;
  description: string;
  rewardAmount: number;
  rewardCurrency: string;
  rewardUsdEquivalent?: number;
  status: BountyStatus;
  platform: Platform;
  tags: string[];
  url: string;
  createdAt: string;
  updatedAt: string;
  deadline?: string;
  metadata?: BountyMetadata;
  requirements?: BountyRequirements;
  isFeatured?: boolean;
  isUrgent?: boolean;
  contributorCount?: number;
  submissionsCount?: number;
  difficulty?: DifficultyLevel;
  remote?: boolean;
  organization?: string;
  organizationLogo?: string;
}
```

### Platform
```typescript
interface Platform {
  id: string;
  name: string;
  url: string;
  description?: string;
  logoUrl?: string;
  isActive: boolean;
  supportedTypes: BountyType[];
  apiEndpoint?: string;
  scrapeUrl?: string;
  lastSyncedAt?: string;
}
```

### Platform Support (Granular)
```typescript
interface PlatformSupport {
  // Core capabilities
  supportsBounties: boolean;
  supportsGrants: boolean;
  supportsContests: boolean;
  
  // Data fields
  supportsTitle: boolean;
  supportsDescription: boolean;
  supportsRewardAmount: boolean;
  supportsDeadline: boolean;
  supportsTags: boolean;
  supportsOrganization: boolean;
  supportsLogo: boolean;
  
  // Actions
  supportsApplications: boolean;
  supportsSubmissions: boolean;
  supportsMessaging: boolean;
  
  // Filtering
  supportsTypeFilter: boolean;
  supportsStatusFilter: boolean;
  supportsRewardFilter: boolean;
  
  // Sorting
  supportsSortByReward: boolean;
  supportsSortByDate: boolean;
  supportsSortByDeadline: boolean;
  
  // Pagination
  supportsPagination: boolean;
  maxPageSize?: number;
  
  // Rate limiting
  rateLimitPerMinute?: number;
  requiresAuth: boolean;
  
  // Method implementations (nil = unsupported)
  methods: {
    fetchBounties?: () => Promise<Bounty[]>;
    fetchBounty?: (id: string) => Promise<Bounty>;
    submitApplication?: (data: ApplicationData) => Promise<Application>;
    submitWork?: (data: SubmissionData) => Promise<Submission>;
    getTypes?: () => Promise<BountyType[]>;
    getStatuses?: () => Promise<BountyStatus[]>;
  };
}
```

---

## Error Handling

### Error Response Format
```json
{
  "error": {
    "code": "BOUNTY_NOT_FOUND",
    "message": "Bounty with ID 'xyz' not found",
    "details": {
      "bountyId": "xyz",
      "timestamp": "2026-03-13T21:00:00Z"
    }
  }
}
```

### Error Codes
| Code | HTTP Status | Description |
|------|-------------|-------------|
| `BOUNTY_NOT_FOUND` | 404 | Bounty doesn't exist |
| `INVALID_PARAMETERS` | 400 | Invalid query/body params |
| `UNAUTHORIZED` | 401 | Missing/invalid auth |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `PLATFORM_ERROR` | 502 | Platform scrape/API failed |
| `INTERNAL_ERROR` | 500 | Server error |

---

## Rate Limiting

### Limits
| Tier | Requests/Minute | Requests/Day |
|------|-----------------|--------------|
| Anonymous | 60 | 1,000 |
| Authenticated | 300 | 10,000 |
| Admin | 1,000 | 100,000 |

### Headers
```http
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1647200000
```

---

## Webhooks

### Events
- `bounty.created`
- `bounty.updated`
- `bounty.deleted`
- `submission.created`
- `application.created`

### Payload
```json
{
  "event": "bounty.created",
  "timestamp": "2026-03-13T21:00:00Z",
  "data": {
    "id": "...",
    "type": "...",
    ...
  }
}
```

---

## Platform Adapters

### Adapter Interface
```typescript
interface IPlatformAdapter {
  // Platform info
  readonly platformId: string;
  readonly platformName: string;
  
  // Capabilities
  readonly support: PlatformSupport;
  
  // Core methods
  fetchBounties(params: FetchParams): Promise<Bounty[]>;
  fetchBounty(id: string): Promise<Bounty | null>;
  
  // Optional methods (return null if unsupported)
  submitApplication?(data: ApplicationData): Promise<Application | null>;
  submitWork?(data: SubmissionData): Promise<Submission | null>;
  
  // Metadata
  getTypes?(): Promise<BountyType[]>;
  getStatuses?(): Promise<BountyStatus[]>;
  
  // Health
  healthCheck(): Promise<PlatformHealth>;
}
```

### Adapter Registry
```typescript
interface AdapterRegistry {
  getAdapter(platformId: string): IPlatformAdapter | null;
  getAllAdapters(): IPlatformAdapter[];
  getSupportedPlatforms(): PlatformInfo[];
  registerAdapter(adapter: IPlatformAdapter): void;
  unregisterAdapter(platformId: string): void;
}
```

---

## Testing & Validation

### Test Categories

#### 1. Unit Tests
- Adapter implementations
- Data transformations
- Filter logic
- Validation functions

#### 2. Integration Tests
- API endpoint responses
- Database operations
- Platform scraping
- Cache operations

#### 3. E2E Tests
- Full user workflows
- Real-time updates
- Filter combinations
- Pagination

### Validation Rules

#### Bounty Validation
```typescript
const bountySchema = {
  id: { required: true, type: 'string', minLength: 1 },
  type: { required: true, type: 'enum', values: BountyType },
  title: { required: true, type: 'string', minLength: 5, maxLength: 200 },
  description: { required: true, type: 'string', minLength: 10, maxLength: 10000 },
  rewardAmount: { required: true, type: 'number', min: 0 },
  rewardCurrency: { required: true, type: 'string', minLength: 3 },
  status: { required: true, type: 'enum', values: BountyStatus },
  platform: { required: true, type: 'object' },
  url: { required: true, type: 'url' },
  createdAt: { required: true, type: 'iso8601' }
};
```

#### Platform Adapter Validation
```typescript
const adapterValidation = {
  platformId: { required: true, unique: true },
  platformName: { required: true, minLength: 2 },
  support: { required: true, type: 'object' },
  fetchBounties: { required: true, type: 'function' },
  healthCheck: { required: true, type: 'function' }
};
```

### CI/CD Checks

#### Pre-commit
- [ ] TypeScript compilation
- [ ] ESLint validation
- [ ] Unit tests pass
- [ ] Type definitions complete

#### Pre-merge
- [ ] All tests pass
- [ ] API spec compliance
- [ ] Adapter validation
- [ ] Integration tests pass
- [ ] E2E tests pass
- [ ] Performance benchmarks

#### Pre-deploy
- [ ] Security audit
- [ ] Load testing
- [ ] API documentation generated
- [ ] Changelog updated

---

## Appendix

### Changelog

#### v2.0.0 (2026-03-13)
- Added platform adapter architecture
- Added granular support indicators
- Added real-time bounty population
- Improved filter system
- Increased bounty limit to 100+
- Added comprehensive API spec

#### v1.0.0 (2026-02-01)
- Initial release

---

**Document Maintained By:** bountyOS Core Team  
**Contact:** api@bountyos.com
