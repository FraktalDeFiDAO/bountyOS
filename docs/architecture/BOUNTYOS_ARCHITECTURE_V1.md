# 🏗️ BountyOS Architecture - Deep Separation of Concerns

**Date:** March 12, 2026  
**Version:** 1.0.0  
**Principle:** Deep Separation of Concerns by Bounty Type

---

## 📋 EXECUTIVE SUMMARY

This architecture implements **deep separation of concerns** across all layers (frontend, backend, data) to properly handle the **8 distinct bounty types** identified in our research.

**Key Principles:**
1. **Bounty-Type Isolation** - Each bounty type has dedicated models, services, and UI
2. **Shared Infrastructure** - Common services (auth, payments, notifications) are shared
3. **Plugin Architecture** - New bounty types can be added without modifying existing code
4. **Unified Discovery** - Single search/browse interface across all bounty types
5. **Crypto-First Payments** - All bounties prefer crypto, support fiat fallback

---

## 🎯 BOUNTY TYPE TAXONOMY

### 8 Distinct Bounty Types:

| Type | ID | Description | Typical Reward | Payment Methods |
|------|-----|-------------|----------------|-----------------|
| **Development Bounties** | `DEV` | Feature development, OSS contributions | $50-$500K | Crypto (USDC, ETH, SOL, token) |
| **Grants** | `GRANT` | Protocol/ecosystem grants | $1K-$750K | Crypto (protocol token, USDC) |
| **Microtasks** | `MICRO` | Small tasks, surveys, testing | $1-$100 | Crypto (BTC, ETH, platform token) |
| **Freelance/Gigs** | `GIG` | Contract work, project-based | $500-$50K | Crypto + Fiat (Stripe, PayPal) |
| **Hackathons** | `HACK` | Time-limited competitions | $1K-$100K+ | Crypto (prize pool) |
| **Ambassador Programs** | `AMB` | Community building, advocacy | $100-$10K/mo | Crypto (token, USDC) |
| **Retroactive Funding** | `RETRO` | Retroactive public goods | $1K-$500K+ | Crypto (protocol token) |
| **Security Bounties** | `SEC` | Bug bounties, audits (OPTIONAL) | $100-$10M | Crypto (USDC, ETH, protocol token) |

**Note:** Security bounties are excluded per user preference, but architecture supports them.

---

## 🏛️ ARCHITECTURE OVERVIEW

```
┌─────────────────────────────────────────────────────────────────┐
│                         FRONTEND LAYER                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │   Discovery │  │   Bounty    │  │   User      │            │
│  │   & Search  │  │   Details   │  │   Dashboard │            │
│  └─────────────┘  └─────────────┘  └─────────────┘            │
│                                                                 │
│  ┌─────────────────────────────────────────────────┐          │
│  │         Bounty-Type Specific Components         │          │
│  ├─────────┬─────────┬─────────┬─────────┬───────┤          │
│  │   DEV   │  GRANT  │  MICRO  │   GIG   │ HACK  │          │
│  ├─────────┼─────────┼─────────┼─────────┼───────┤          │
│  │  AMB    │  RETRO  │  (SEC)  │         │       │          │
│  └─────────┴─────────┴─────────┴─────────┴───────┘          │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ API Gateway (REST + GraphQL)
                              │
┌─────────────────────────────────────────────────────────────────┐
│                         BACKEND LAYER                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────────────────────────────────────┐          │
│  │              Shared Infrastructure              │          │
│  ├──────────┬──────────┬──────────┬──────────────┤          │
│  │   Auth   │ Payments │   Notif. │   Analytics  │          │
│  │ Service  │ Service  │ Service  │   Service    │          │
│  └──────────┴──────────┴──────────┴──────────────┘          │
│                                                                 │
│  ┌─────────────────────────────────────────────────┐          │
│  │         Bounty-Type Specific Services           │          │
│  ├─────────┬─────────┬─────────┬─────────┬───────┤          │
│  │   DEV   │  GRANT  │  MICRO  │   GIG   │ HACK  │          │
│  │ Service │ Service │ Service │ Service │Service│          │
│  ├─────────┼─────────┼─────────┼─────────┼───────┤          │
│  │  AMB    │  RETRO  │  (SEC)  │         │       │          │
│  │ Service │ Service │ Service │         │       │          │
│  └─────────┴─────────┴─────────┴─────────┴───────┘          │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ ORM Layer
                              │
┌─────────────────────────────────────────────────────────────────┐
│                         DATA LAYER                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────────────────────────────────────┐          │
│  │              Shared Data Models                 │          │
│  ├──────────┬──────────┬──────────┬──────────────┤          │
│  │   User   │ Platform │ Payment  │   Activity   │          │
│  │  Model   │  Model   │  Model   │    Model     │          │
│  └──────────┴──────────┴──────────┴──────────────┘          │
│                                                                 │
│  ┌─────────────────────────────────────────────────┐          │
│  │        Bounty-Type Specific Data Models         │          │
│  ├─────────┬─────────┬─────────┬─────────┬───────┤          │
│  │   DEV   │  GRANT  │  MICRO  │   GIG   │ HACK  │          │
│  │  Model  │  Model  │  Model  │  Model  │Model  │          │
│  ├─────────┼─────────┼─────────┼─────────┼───────┤          │
│  │  AMB    │  RETRO  │  (SEC)  │         │       │          │
│  │  Model  │  Model  │  Model  │         │       │          │
│  └─────────┴─────────┴─────────┴─────────┴───────┘          │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 📁 DIRECTORY STRUCTURE

```
bountyos/
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── shared/              # Shared components
│   │   │   │   ├── Button.tsx
│   │   │   │   ├── Modal.tsx
│   │   │   │   └── ...
│   │   │   │
│   │   │   ├── discovery/           # Discovery & Search
│   │   │   │   ├── BountyCard.tsx
│   │   │   │   ├── SearchBar.tsx
│   │   │   │   ├── Filters.tsx
│   │   │   │   └── ...
│   │   │   │
│   │   │   ├── bounties/            # Bounty-Type Specific
│   │   │   │   ├── dev/             # Development Bounties
│   │   │   │   │   ├── DevBountyCard.tsx
│   │   │   │   │   ├── DevBountyDetails.tsx
│   │   │   │   │   ├── SubmissionForm.tsx
│   │   │   │   │   └── ...
│   │   │   │   ├── grant/           # Grants
│   │   │   │   │   ├── GrantCard.tsx
│   │   │   │   │   ├── GrantApplication.tsx
│   │   │   │   │   └── ...
│   │   │   │   ├── micro/           # Microtasks
│   │   │   │   ├── gig/             # Freelance/Gigs
│   │   │   │   ├── hackathon/       # Hackathons
│   │   │   │   ├── ambassador/      # Ambassador Programs
│   │   │   │   ├── retro/           # Retroactive Funding
│   │   │   │   └── security/        # Security Bounties (optional)
│   │   │   │
│   │   │   └── dashboard/           # User Dashboard
│   │   │       ├── Earnings.tsx
│   │   │       ├── Submissions.tsx
│   │   │       └── ...
│   │   │
│   │   ├── services/
│   │   │   ├── api.ts               # API client
│   │   │   ├── auth.ts              # Auth service
│   │   │   └── bounty-types.ts      # Bounty type definitions
│   │   │
│   │   └── types/
│   │       ├── bounty.ts            # Shared bounty types
│   │       ├── dev-bounty.ts        # Dev-specific types
│   │       ├── grant.ts             # Grant-specific types
│   │       └── ...                  # One file per bounty type
│   │
│   └── package.json
│
├── backend/
│   ├── src/
│   │   ├── shared/
│   │   │   ├── auth/                # Authentication service
│   │   │   ├── payments/            # Payment service (crypto + fiat)
│   │   │   ├── notifications/       # Notification service
│   │   │   ├── analytics/           # Analytics service
│   │   │   └── utils/               # Shared utilities
│   │   │
│   │   ├── bounties/                # Bounty-Type Specific
│   │   │   ├── dev/                 # Development Bounties
│   │   │   │   ├── dev-bounty.service.ts
│   │   │   │   ├── dev-bounty.controller.ts
│   │   │   │   ├── dev-bounty.schema.ts
│   │   │   │   └── dev-bounty.types.ts
│   │   │   │
│   │   │   ├── grant/               # Grants
│   │   │   │   ├── grant.service.ts
│   │   │   │   ├── grant.controller.ts
│   │   │   │   ├── grant.schema.ts
│   │   │   │   └── grant.types.ts
│   │   │   │
│   │   │   ├── micro/               # Microtasks
│   │   │   ├── gig/                 # Freelance/Gigs
│   │   │   ├── hackathon/           # Hackathons
│   │   │   ├── ambassador/          # Ambassador Programs
│   │   │   ├── retro/               # Retroactive Funding
│   │   │   └── security/            # Security Bounties (optional)
│   │   │
│   │   ├── platforms/               # Platform Integrations
│   │   │   ├── platform.service.ts
│   │   │   ├── platform.controller.ts
│   │   │   ├── platforms/
│   │   │   │   ├── gitcoin.ts
│   │   │   │   ├── superteam.ts
│   │   │   │   ├── algora.ts
│   │   │   │   ├── proxies-sx.ts
│   │   │   │   └── ...              # One file per platform
│   │   │   └── platform-registry.ts
│   │   │
│   │   ├── api/
│   │   │   ├── routes.ts            # Main router
│   │   │   └── middleware/          # API middleware
│   │   │
│   │   └── index.ts                 # Entry point
│   │
│   ├── tests/
│   ├── package.json
│   └── tsconfig.json
│
├── database/
│   ├── migrations/
│   ├── schemas/
│   │   ├── shared/                  # Shared schemas
│   │   │   ├── user.schema.ts
│   │   │   ├── platform.schema.ts
│   │   │   └── payment.schema.ts
│   │   │
│   │   ├── bounties/                # Bounty-Type Specific
│   │   │   ├── dev-bounty.schema.ts
│   │   │   ├── grant.schema.ts
│   │   │   ├── micro.schema.ts
│   │   │   ├── gig.schema.ts
│   │   │   ├── hackathon.schema.ts
│   │   │   ├── ambassador.schema.ts
│   │   │   ├── retro.schema.ts
│   │   │   └── security.schema.ts   # Optional
│   │   │
│   │   └── index.ts                 # Schema registry
│   │
│   └── seed/                        # Seed data
│
├── config/
│   ├── platforms.ts                 # Platform configurations
│   ├── payments.ts                  # Payment configurations
│   └── index.ts
│
└── docs/
    ├── architecture.md
    ├── bounty-types.md
    ├── platform-integrations.md
    └── api.md
```

---

## 🗄️ DATA MODELS

### Shared Models

```typescript
// Shared across all bounty types
interface User {
  id: string;
  email: string;
  walletAddress?: string;  // Primary crypto wallet
  paymentPreferences: {
    preferredCrypto: string[];  // e.g., ['USDC', 'ETH', 'SOL']
    acceptFiat: boolean;
    fiatMethods: string[];  // e.g., ['stripe', 'paypal', 'venmo']
  };
  skills: string[];
  reputation: number;
}

interface Platform {
  id: string;
  name: string;
  type: BountyType;
  url: string;
  apiEndpoint?: string;
  authMethod: 'api_key' | 'oauth' | 'none';
  paymentMethods: PaymentMethod[];
  bountyTypes: BountyType[];
}

interface Payment {
  id: string;
  userId: string;
  bountyId: string;
  amount: number;
  currency: string;
  method: 'crypto' | 'fiat';
  cryptoNetwork?: string;
  status: 'pending' | 'completed' | 'failed';
  txHash?: string;
}
```

### Bounty-Type Specific Models

```typescript
// Development Bounty
interface DevBounty {
  id: string;
  platformId: string;
  type: 'DEV';
  title: string;
  description: string;
  reward: {
    amount: number;
    currency: string;
    usdEquivalent?: number;
  };
  requirements: {
    skills: string[];
    experience: 'beginner' | 'intermediate' | 'advanced';
    deliverables: string[];
  };
  deadline?: Date;
  status: 'open' | 'in_progress' | 'completed' | 'closed';
  submissions?: Submission[];
  metadata: {
    githubUrl?: string;
    repo?: string;
    issue?: string;
  };
}

// Grant
interface Grant {
  id: string;
  platformId: string;
  type: 'GRANT';
  programName: string;
  organization: string;
  description: string;
  reward: {
    minAmount: number;
    maxAmount: number;
    currency: string;
    type: 'milestone' | 'one-time' | 'retroactive';
  };
  focusAreas: string[];
  requirements: {
    proposalRequired: boolean;
    milestones?: Milestone[];
    reportingRequired: boolean;
  };
  deadline?: Date;
  status: 'open' | 'closed' | 'reviewing';
  applications?: Application[];
}

// Microtask
interface Microtask {
  id: string;
  platformId: string;
  type: 'MICRO';
  title: string;
  description: string;
  reward: {
    amount: number;
    currency: string;
    perTask: boolean;
  };
  taskType: 'survey' | 'testing' | 'social' | 'content' | 'other';
  estimatedTime: number;  // minutes
  difficulty: 'easy' | 'medium' | 'hard';
  status: 'available' | 'in_progress' | 'completed';
  completions?: number;
  maxCompletions?: number;
}

// ... (similar for GIG, HACK, AMB, RETRO, SEC)
```

---

## 🔌 SERVICE LAYER

### Shared Services

```typescript
// Payment Service (shared)
class PaymentService {
  async processPayment(payment: Payment): Promise<TxHash> {
    if (payment.method === 'crypto') {
      return this.processCryptoPayment(payment);
    } else {
      return this.processFiatPayment(payment);
    }
  }
  
  private async processCryptoPayment(payment: Payment): Promise<TxHash> {
    // Support multiple chains: ETH, SOL, BTC, etc.
  }
  
  private async processFiatPayment(payment: Payment): Promise<string> {
    // Support Stripe, PayPal, Venmo
  }
}

// Platform Registry (shared)
class PlatformRegistry {
  private platforms: Map<string, PlatformAdapter> = new Map();
  
  register(platform: PlatformAdapter): void {
    this.platforms.set(platform.id, platform);
  }
  
  getPlatform(id: string): PlatformAdapter {
    return this.platforms.get(id)!;
  }
  
  getPlatformsByType(type: BountyType): PlatformAdapter[] {
    return Array.from(this.platforms.values())
      .filter(p => p.supportsType(type));
  }
}
```

### Bounty-Type Specific Services

```typescript
// Development Bounty Service
class DevBountyService {
  constructor(
    private db: Database,
    private paymentService: PaymentService,
    private githubService: GitHubService
  ) {}
  
  async findBounties(filters: DevBountyFilters): Promise<DevBounty[]> {
    return this.db.devBounties.findMany({
      where: filters,
      include: ['platform', 'submissions']
    });
  }
  
  async submitWork(submission: DevSubmission): Promise<void> {
    // Validate submission (e.g., GitHub PR)
    await this.githubService.validatePR(submission.prUrl);
    
    // Create submission record
    await this.db.devSubmission.create(submission);
    
    // Notify bounty creator
    await this.notificationService.notify(submission.bountyId);
  }
  
  async approveSubmission(submissionId: string): Promise<void> {
    const submission = await this.db.devSubmission.findUnique(submissionId);
    
    // Process payment
    await this.paymentService.processPayment({
      userId: submission.userId,
      bountyId: submission.bountyId,
      amount: submission.bounty.reward.amount,
      currency: submission.bounty.reward.currency,
      method: 'crypto'
    });
    
    // Update submission status
    await this.db.devSubmission.update(submissionId, {
      status: 'approved'
    });
  }
}

// Grant Service
class GrantService {
  async applyForGrant(application: GrantApplication): Promise<void> {
    // Validate application
    // Submit to platform
    // Track application status
  }
  
  async submitMilestone(milestone: MilestoneSubmission): Promise<void> {
    // Validate milestone deliverables
    // Request payment for milestone
  }
}

// ... (similar for Micro, Gig, Hackathon, Ambassador, Retro, Security)
```

---

## 🎨 FRONTEND COMPONENTS

### Discovery Layer (Shared)

```typescript
// Unified search across all bounty types
function BountyDiscovery() {
  const [filters, setFilters] = useState<BountyFilters>({
    types: [],  // Can filter by multiple types
    minReward: 0,
    maxReward: Infinity,
    paymentMethods: [],
    platforms: []
  });
  
  const bounties = useBountySearch(filters);
  
  return (
    <div>
      <BountyFilters filters={filters} onChange={setFilters} />
      <BountyGrid bounties={bounties} />
    </div>
  );
}

// Bounty card renders differently based on type
function BountyCard({ bounty }: { bounty: Bounty }) {
  switch (bounty.type) {
    case 'DEV':
      return <DevBountyCard bounty={bounty} />;
    case 'GRANT':
      return <GrantCard bounty={bounty} />;
    case 'MICRO':
      return <MicrotaskCard bounty={bounty} />;
    // ... etc
  }
}
```

### Bounty-Type Specific Components

```typescript
// Development Bounty Details
function DevBountyDetails({ bounty }: { bounty: DevBounty }) {
  return (
    <div>
      <BountyHeader bounty={bounty} />
      <Requirements requirements={bounty.requirements} />
      <GitHubIntegration repo={bounty.metadata.repo} />
      <SubmissionForm onSubmit={handleSubmit} />
      <PaymentPreferences userId={currentUser.id} />
    </div>
  );
}

// Grant Application Form
function GrantApplicationForm({ grant }: { grant: Grant }) {
  return (
    <div>
      <GrantHeader grant={grant} />
      <ProposalForm fields={grant.requirements.proposalFields} />
      <MilestonePlanner milestones={grant.requirements.milestones} />
      <BudgetPlanner />
      <SubmitApplication />
    </div>
  );
}

// ... (similar for other bounty types)
```

---

## 💰 PAYMENT INTEGRATION

### Crypto-First, Fiat Fallback

```typescript
interface PaymentPreferences {
  preferredCrypto: string[];  // ['USDC', 'ETH', 'SOL']
  acceptFiat: boolean;
  fiatMethods: string[];  // ['stripe', 'paypal', 'venmo']
}

// Payment processing prioritizes crypto
async function processPayment(payment: Payment): Promise<void> {
  const user = await getUser(payment.userId);
  
  if (user.paymentPreferences.preferredCrypto.length > 0) {
    // Process crypto payment
    await cryptoPaymentService.process({
      ...payment,
      currency: user.paymentPreferences.preferredCrypto[0],
      network: getNetworkForCurrency(user.paymentPreferences.preferredCrypto[0])
    });
  } else if (user.paymentPreferences.acceptFiat) {
    // Fallback to fiat
    await fiatPaymentService.process({
      ...payment,
      method: user.paymentPreferences.fiatMethods[0]
    });
  } else {
    throw new Error('No payment method configured');
  }
}
```

### Supported Payment Methods

| Method | Type | Supported Currencies | Fee | Settlement |
|--------|------|---------------------|-----|------------|
| **USDC (Ethereum)** | Crypto | USDC | ~$1-10 | Instant |
| **USDC (Solana)** | Crypto | USDC | ~$0.01 | Instant |
| **ETH** | Crypto | ETH | ~$1-50 | ~15 sec |
| **SOL** | Crypto | SOL | ~$0.01 | ~5 sec |
| **BTC** | Crypto | BTC | ~$1-10 | ~10 min |
| **Stripe** | Fiat | USD, EUR, etc. | 2.9% + $0.30 | 2-7 days |
| **PayPal** | Fiat | USD, EUR, etc. | 2.9% + $0.30 | Instant |
| **Venmo** | Fiat | USD | 1.75% | Instant |

---

## 🔌 PLATFORM INTEGRATIONS

### Platform Adapter Pattern

```typescript
interface PlatformAdapter {
  id: string;
  name: string;
  bountyTypes: BountyType[];
  
  // Fetch bounties from platform
  fetchBounties(type: BountyType): Promise<Bounty[]>;
  
  // Submit work/application
  submitSubmission(submission: Submission): Promise<void>;
  
  // Check submission status
  checkStatus(submissionId: string): Promise<SubmissionStatus>;
  
  // Process payment (if platform handles payments)
  processPayment?(payment: Payment): Promise<TxHash>;
}

// Example: Gitcoin adapter
class GitcoinAdapter implements PlatformAdapter {
  id = 'gitcoin';
  name = 'Gitcoin';
  bountyTypes = ['DEV', 'GRANT', 'HACK'];
  
  async fetchBounties(type: BountyType): Promise<DevBounty[]> {
    // Fetch from Gitcoin API
  }
  
  async submitSubmission(submission: Submission): Promise<void> {
    // Submit work via Gitcoin
  }
}

// Example: Proxies.sx adapter
class ProxiesSXAdapter implements PlatformAdapter {
  id = 'proxies-sx';
  name = 'Proxies.sx';
  bountyTypes = ['DEV', 'MICRO'];
  
  async fetchBounties(type: BountyType): Promise<DevBounty[]> {
    // Fetch from GitHub API (Proxies.sx uses GitHub Issues)
  }
}
```

---

## 📊 API ENDPOINTS

### REST API Structure

```
GET  /api/bounties                    # Search all bounties
GET  /api/bounties/:id                # Get bounty details
POST /api/bounties/:id/submit         # Submit work

# Bounty-Type Specific
GET  /api/bounties/dev                # Development bounties
GET  /api/bounties/grant              # Grants
GET  /api/bounties/micro              # Microtasks
GET  /api/bounties/gig                # Freelance/Gigs
GET  /api/bounties/hackathon          # Hackathons
GET  /api/bounties/ambassador         # Ambassador programs
GET  /api/bounties/retro              # Retroactive funding
GET  /api/bounties/security           # Security bounties (optional)

# Platform-Specific
GET  /api/platforms                   # List all platforms
GET  /api/platforms/:id/bounties      # Bounties from specific platform
POST /api/platforms/:id/sync          # Sync bounties from platform

# User
GET  /api/user/profile                # User profile
GET  /api/user/earnings               # User earnings
GET  /api/user/submissions            # User submissions
PUT  /api/user/payment-preferences    # Update payment preferences

# Payments
POST /api/payments/process            # Process payment
GET  /api/payments/:id/status         # Payment status
```

---

## 🚀 IMPLEMENTATION ROADMAP

### Phase 1: Core Infrastructure (Week 1-2)
- [ ] Set up project structure
- [ ] Implement shared data models
- [ ] Build platform registry
- [ ] Create payment service (crypto only initially)
- [ ] Set up database schemas

### Phase 2: First Bounty Type - DEV (Week 3-4)
- [ ] Implement DevBounty service
- [ ] Create DevBounty UI components
- [ ] Integrate 2-3 platforms (Gitcoin, Proxies.sx, Algora)
- [ ] Build submission flow
- [ ] Test end-to-end

### Phase 3: Additional Bounty Types (Week 5-8)
- [ ] Grants (GRANT)
- [ ] Microtasks (MICRO)
- [ ] Freelance/Gigs (GIG)
- [ ] Hackathons (HACK)
- [ ] Ambassador (AMB)
- [ ] Retroactive (RETRO)

### Phase 4: Payment Expansion (Week 9-10)
- [ ] Add fiat payment support (Stripe, PayPal, Venmo)
- [ ] Multi-chain crypto support (ETH, SOL, BTC)
- [ ] Payment preferences UI
- [ ] Payment history & tracking

### Phase 5: Polish & Launch (Week 11-12)
- [ ] Unified search & discovery
- [ ] User dashboard
- [ ] Notifications
- [ ] Analytics
- [ ] Testing & bug fixes
- [ ] Launch!

---

## 📝 NEXT STEPS

1. **Review Architecture** - Confirm separation of concerns approach
2. **Prioritize Bounty Types** - Which types to implement first?
3. **Select Initial Platforms** - Which platforms to integrate first?
4. **Payment Method Priority** - Which crypto/fiat methods first?
5. **Start Implementation** - Begin Phase 1

---

**Questions for User:**
1. Does this architecture meet your separation of concerns requirements?
2. Which bounty types should we prioritize? (Recommendation: DEV → GRANT → MICRO)
3. Which platforms should we integrate first? (Recommendation: Proxies.sx, Gitcoin, Algora)
4. Payment preference: Crypto-only for MVP, or include fiat from start?
