# Bounty Monitoring Service Architecture

**Version:** 1.0.0  
**Date:** March 13, 2026  
**Status:** Design

---

## Overview

The Bounty Monitoring Service is a **background worker service** that continuously monitors all integrated platforms for:
- New bounty opportunities
- Status changes on existing bounties
- Deadline approaching alerts
- Reward amount changes
- Bounty assignments/claims

---

## Architecture

### High-Level Design

```
┌─────────────────────────────────────────────────────────────┐
│                   Bounty Monitoring Service                   │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐       │
│  │   Scheduler  │  │   Queue      │  │   Worker     │       │
│  │   (Cron)     │──│   (Bull)     │──│   Pool       │       │
│  └──────────────┘  └──────────────┘  └──────────────┘       │
│                           │                    │              │
│                           │                    ▼              │
│                           │          ┌──────────────┐       │
│                           │          │   Platform   │       │
│                           │          │   Adapters   │       │
│                           │          └──────────────┘       │
│                           │                    │              │
│                           ▼                    ▼              │
│                    ┌──────────────────────────────┐         │
│                    │      Change Detection        │         │
│                    │         Engine               │         │
│                    └──────────────────────────────┘         │
│                                      │                       │
│                                      ▼                       │
│                    ┌──────────────────────────────┐         │
│                    │       Event Emitter          │         │
│                    │    (WebSocket/Redis)         │         │
│                    └──────────────────────────────┘         │
│                                      │                       │
└──────────────────────────────────────┼───────────────────────┘
                                       │
                                       ▼
                    ┌──────────────────────────────┐
                    │       Database               │
                    │       (PostgreSQL)           │
                    └──────────────────────────────┘
                                       │
                                       ▼
                    ┌──────────────────────────────┐
                    │       Frontend               │
                    │    (Real-time Updates)       │
                    └──────────────────────────────┘
```

---

## Core Components

### 1. Scheduler Service

**Purpose:** Trigger monitoring jobs at configured intervals

**Implementation:**
```typescript
// src/monitoring/scheduler.ts
import { CronJob } from 'cron';
import { monitoringQueue } from './queue';

interface MonitoringSchedule {
  platformId: string;
  cronExpression: string; // e.g., '*/5 * * * *' = every 5 minutes
  priority: 'high' | 'medium' | 'low';
  enabled: boolean;
}

const DEFAULT_SCHEDULES: MonitoringSchedule[] = [
  { platformId: 'gitcoin', cronExpression: '*/10 * * * *', priority: 'high', enabled: true },
  { platformId: 'superteam', cronExpression: '*/5 * * * *', priority: 'high', enabled: true },
  { platformId: 'algora', cronExpression: '*/10 * * * *', priority: 'high', enabled: true },
  { platformId: 'proxies-sx', cronExpression: '*/3 * * * *', priority: 'high', enabled: true },
  { platformId: 'code4rena', cronExpression: '*/30 * * * *', priority: 'medium', enabled: true },
  { platformId: 'github', cronExpression: '*/15 * * * *', priority: 'medium', enabled: true },
];

export class MonitoringScheduler {
  private jobs: Map<string, CronJob> = new Map();

  async start() {
    for (const schedule of DEFAULT_SCHEDULES) {
      if (schedule.enabled) {
        this.schedulePlatform(schedule);
      }
    }
  }

  private schedulePlatform(schedule: MonitoringSchedule) {
    const job = new CronJob(schedule.cronExpression, async () => {
      await monitoringQueue.add('monitor-platform', {
        platformId: schedule.platformId,
        priority: schedule.priority,
        timestamp: new Date().toISOString()
      }, {
        priority: schedule.priority === 'high' ? 1 : 5,
        attempts: 3,
        backoff: {
          type: 'exponential',
          delay: 2000
        }
      });
    });

    this.jobs.set(schedule.platformId, job);
    job.start();
  }

  stop() {
    this.jobs.forEach(job => job.stop());
  }
}
```

---

### 2. Message Queue

**Purpose:** Manage monitoring tasks with priority and retry logic

**Implementation:**
```typescript
// src/monitoring/queue.ts
import { Queue, Worker, Job } from 'bullmq';
import { Redis } from 'ioredis';

const redis = new Redis({
  host: process.env.REDIS_HOST || 'localhost',
  port: parseInt(process.env.REDIS_PORT || '6379'),
  maxRetriesPerRequest: null
});

export const monitoringQueue = new Queue('bounty-monitoring', {
  connection: redis,
  defaultJobOptions: {
    removeOnComplete: 100,
    removeOnFail: 1000,
    attempts: 3,
    backoff: {
      type: 'exponential',
      delay: 2000
    }
  }
});

// Worker pool
const worker = new Worker('bounty-monitoring', async (job: Job) => {
  switch (job.name) {
    case 'monitor-platform':
      return await handlePlatformMonitoring(job.data);
    case 'check-bounty-status':
      return await handleBountyStatusCheck(job.data);
    case 'process-bounty-changes':
      return await handleBountyChanges(job.data);
    default:
      throw new Error(`Unknown job type: ${job.name}`);
  }
}, {
  connection: redis,
  concurrency: 5 // Process 5 jobs concurrently
});

// Event handlers
worker.on('completed', (job) => {
  console.log(`Job ${job.name} completed for ${job.data.platformId}`);
});

worker.on('failed', (job, err) => {
  console.error(`Job ${job?.name} failed:`, err);
  // Send alert for repeated failures
  if (job?.attemptsMade >= 3) {
    sendAlert('monitoring_failure', {
      platformId: job.data.platformId,
      error: err.message
    });
  }
});
```

---

### 3. Change Detection Engine

**Purpose:** Identify new bounties and status changes

**Implementation:**
```typescript
// src/monitoring/change-detector.ts
import { IBounty } from '../adapters/types';
import { bountyService } from '../services/bounty.service';

interface BountyChange {
  type: 'NEW' | 'UPDATED' | 'STATUS_CHANGED' | 'DEADLINE_APPROACHING' | 'REWARD_CHANGED';
  bountyId: string;
  platformId: string;
  timestamp: string;
  previousState?: Partial<IBounty>;
  newState: Partial<IBounty>;
  significance: 'high' | 'medium' | 'low';
}

export class ChangeDetector {
  /**
   * Compare fetched bounties with database to detect changes
   */
  async detectChanges(platformId: string, fetchedBounties: IBounty[]): Promise<BountyChange[]> {
    const changes: BountyChange[] = [];

    // Get existing bounties for this platform
    const existingBounties = await this.getExistingBounties(platformId);
    const existingMap = new Map(existingBounties.map(b => [b.id, b]));

    for (const fetched of fetchedBounties) {
      const existing = existingMap.get(fetched.id);

      if (!existing) {
        // NEW BOUNTY DETECTED
        changes.push({
          type: 'NEW',
          bountyId: fetched.id,
          platformId,
          timestamp: new Date().toISOString(),
          newState: fetched,
          significance: this.calculateSignificance(fetched)
        });
      } else {
        // Check for updates
        const updateChanges = this.detectUpdateChanges(existing, fetched);
        changes.push(...updateChanges);
      }
    }

    // Check for removed bounties (no longer available)
    const fetchedIds = new Set(fetchedBounties.map(b => b.id));
    for (const existing of existingBounties) {
      if (!fetchedIds.has(existing.id) && existing.status === 'OPEN') {
        changes.push({
          type: 'STATUS_CHANGED',
          bountyId: existing.id,
          platformId,
          timestamp: new Date().toISOString(),
          previousState: { status: 'OPEN' },
          newState: { status: 'CLOSED' },
          significance: 'high'
        });
      }
    }

    return changes;
  }

  /**
   * Detect specific changes between old and new bounty state
   */
  private detectUpdateChanges(oldBounty: IBounty, newBounty: IBounty): BountyChange[] {
    const changes: BountyChange[] = [];

    // Status change
    if (oldBounty.status !== newBounty.status) {
      changes.push({
        type: 'STATUS_CHANGED',
        bountyId: newBounty.id,
        platformId: newBounty.platform.id,
        timestamp: new Date().toISOString(),
        previousState: { status: oldBounty.status },
        newState: { status: newBounty.status },
        significance: oldBounty.status === 'OPEN' ? 'high' : 'medium'
      });
    }

    // Reward change
    if (oldBounty.rewardAmount !== newBounty.rewardAmount) {
      const change = {
        type: 'REWARD_CHANGED' as const,
        bountyId: newBounty.id,
        platformId: newBounty.platform.id,
        timestamp: new Date().toISOString(),
        previousState: { rewardAmount: oldBounty.rewardAmount },
        newState: { rewardAmount: newBounty.rewardAmount },
        significance: this.calculateRewardChangeSignificance(
          oldBounty.rewardAmount,
          newBounty.rewardAmount
        )
      };
      changes.push(change);
    }

    // Deadline approaching (within 48 hours)
    if (newBounty.deadline) {
      const hoursUntilDeadline = this.getHoursUntilDeadline(newBounty.deadline);
      if (hoursUntilDeadline <= 48 && hoursUntilDeadline > 24) {
        // First time crossing 48h threshold
        const oldHours = this.getHoursUntilDeadline(oldBounty.deadline);
        if (oldHours > 48) {
          changes.push({
            type: 'DEADLINE_APPROACHING',
            bountyId: newBounty.id,
            platformId: newBounty.platform.id,
            timestamp: new Date().toISOString(),
            newState: { deadline: newBounty.deadline, hoursRemaining: hoursUntilDeadline },
            significance: 'high'
          });
        }
      }
    }

    return changes;
  }

  /**
   * Calculate significance of a new bounty
   */
  private calculateSignificance(bounty: IBounty): 'high' | 'medium' | 'low' {
    // High: reward > $500 or urgent
    if (bounty.rewardAmount >= 500 || bounty.isUrgent) {
      return 'high';
    }
    // Medium: reward > $100
    if (bounty.rewardAmount >= 100) {
      return 'medium';
    }
    return 'low';
  }

  private calculateRewardChangeSignificance(oldAmount: number, newAmount: number): 'high' | 'medium' | 'low' {
    const percentChange = Math.abs((newAmount - oldAmount) / oldAmount);
    if (percentChange > 0.5 || Math.abs(newAmount - oldAmount) > 200) {
      return 'high';
    }
    if (percentChange > 0.2) {
      return 'medium';
    }
    return 'low';
  }

  private getHoursUntilDeadline(deadline: string): number {
    const now = new Date();
    const deadlineDate = new Date(deadline);
    const diffMs = deadlineDate.getTime() - now.getTime();
    return Math.floor(diffMs / (1000 * 60 * 60));
  }

  private async getExistingBounties(platformId: string): Promise<IBounty[]> {
    // Fetch from database/cache
    return await bountyService.findAll({ platform: platformId, limit: 1000 });
  }
}
```

---

### 4. Assignment Detection

**Purpose:** Identify if bounties are assigned/claimed and track availability

**Implementation:**
```typescript
// src/monitoring/assignment-tracker.ts
import { IBounty } from '../adapters/types';

export interface AssignmentStatus {
  bountyId: string;
  platformId: string;
  isAssigned: boolean;
  assignedTo?: string;
  assignedAt?: string;
  assignmentType: 'open' | 'claimed' | 'in_progress' | 'completed';
  canApply: boolean;
  applicationDeadline?: string;
  maxApplicants?: number;
  currentApplicants?: number;
}

export class AssignmentTracker {
  /**
   * Determine assignment status from bounty data
   */
  async trackAssignment(bounty: IBounty): Promise<AssignmentStatus> {
    const status: AssignmentStatus = {
      bountyId: bounty.id,
      platformId: bounty.platform.id,
      isAssigned: this.detectAssignment(bounty),
      assignedTo: this.extractAssignedTo(bounty),
      assignedAt: this.extractAssignedAt(bounty),
      assignmentType: this.determineAssignmentType(bounty),
      canApply: this.canApplyForBounty(bounty),
      applicationDeadline: bounty.deadline,
      maxApplicants: bounty.metadata?.maxApplicants,
      currentApplicants: bounty.contributorCount
    };

    // Store in database
    await this.saveAssignmentStatus(status);

    return status;
  }

  /**
   * Detect if bounty is assigned based on platform-specific signals
   */
  private detectAssignment(bounty: IBounty): boolean {
    // Platform-specific assignment detection
    switch (bounty.platform.id) {
      case 'gitcoin':
        return bounty.status !== 'OPEN' || bounty.contributorCount > 0;
      case 'superteam':
        return bounty.status === 'IN_PROGRESS' || bounty.status === 'COMPLETED';
      case 'algora':
        return bounty.metadata?.assigned === true;
      case 'proxies-sx':
        return bounty.status === 'CLOSED' || bounty.metadata?.claimed === true;
      case 'code4rena':
        return bounty.status === 'COMPLETED';
      default:
        return bounty.status !== 'OPEN';
    }
  }

  /**
   * Extract who the bounty is assigned to
   */
  private extractAssignedTo(bounty: IBounty): string | undefined {
    return bounty.metadata?.assignee || 
           bounty.metadata?.contributor || 
           bounty.metadata?.winner;
  }

  /**
   * Extract when the bounty was assigned
   */
  private extractAssignedAt(bounty: IBounty): string | undefined {
    return bounty.metadata?.assignedAt || 
           bounty.metadata?.claimedAt ||
           bounty.updatedAt;
  }

  /**
   * Determine assignment type
   */
  private determineAssignmentType(bounty: IBounty): AssignmentStatus['assignmentType'] {
    if (bounty.status === 'COMPLETED') {
      return 'completed';
    }
    if (bounty.status === 'IN_PROGRESS' || bounty.contributorCount > 0) {
      return 'in_progress';
    }
    if (this.detectAssignment(bounty)) {
      return 'claimed';
    }
    return 'open';
  }

  /**
   * Check if user can still apply
   */
  private canApplyForBounty(bounty: IBounty): boolean {
    // Can't apply if closed/completed
    if (['COMPLETED', 'CANCELLED', 'CLOSED'].includes(bounty.status)) {
      return false;
    }

    // Check if max applicants reached
    if (bounty.maxApplicants && bounty.contributorCount >= bounty.maxApplicants) {
      return false;
    }

    // Check deadline
    if (bounty.deadline && new Date(bounty.deadline) < new Date()) {
      return false;
    }

    return true;
  }

  private async saveAssignmentStatus(status: AssignmentStatus) {
    // Save to database for tracking
    // This enables historical analysis of assignment patterns
  }
}
```

---

### 5. Event Emitter (Real-time Updates)

**Purpose:** Broadcast changes to connected clients

**Implementation:**
```typescript
// src/monitoring/event-emitter.ts
import { Server } from 'socket.io';
import { Redis } from 'ioredis';
import { BountyChange } from './change-detector';

export class MonitoringEventEmitter {
  private io: Server;
  private redis: Redis;

  constructor() {
    // WebSocket server
    this.io = new Server({
      cors: {
        origin: process.env.FRONTEND_URL || 'http://localhost:3000',
        methods: ['GET', 'POST']
      }
    });

    // Redis for pub/sub
    this.redis = new Redis();
  }

  /**
   * Emit bounty changes to connected clients
   */
  async emitChanges(changes: BountyChange[]) {
    for (const change of changes) {
      // Emit to WebSocket clients
      this.io.emit('bounty:update', change);

      // Emit to specific channels based on change type
      this.io.emit(`bounty:${change.type.toLowerCase()}`, change);

      // Emit to platform-specific room
      this.io.to(`platform:${change.platformId}`).emit('bounty:update', change);

      // Publish to Redis for other services
      await this.redis.publish(
        'bounty:changes',
        JSON.stringify(change)
      );

      // High significance changes get special handling
      if (change.significance === 'high') {
        this.emitHighPriorityAlert(change);
      }
    }
  }

  /**
   * Emit high priority alerts
   */
  private emitHighPriorityAlert(change: BountyChange) {
    // Send to admin dashboard
    this.io.to('admins').emit('alert:high_priority', {
      type: change.type,
      bountyId: change.bountyId,
      platformId: change.platformId,
      timestamp: change.timestamp
    });

    // Could also send email/SMS/Push notification here
    if (change.type === 'NEW' && change.newState.rewardAmount >= 1000) {
      this.sendNotification({
        type: 'HIGH_VALUE_BOUNTY',
        title: `High Value Bounty: $${change.newState.rewardAmount}`,
        body: change.newState.title,
        url: `/bounties/${change.bountyId}`
      });
    }
  }

  /**
   * Subscribe client to platform room
   */
  subscribeToPlatform(socketId: string, platformId: string) {
    this.io.sockets.sockets.get(socketId)?.join(`platform:${platformId}`);
  }

  /**
   * Subscribe client to admin room
   */
  subscribeToAdmin(socketId: string) {
    this.io.sockets.sockets.get(socketId)?.join('admins');
  }

  private sendNotification(notification: any) {
    // Implement email/SMS/push notification logic
    console.log('Sending notification:', notification);
  }
}
```

---

### 6. Main Service Orchestrator

**Purpose:** Coordinate all monitoring components

**Implementation:**
```typescript
// src/monitoring/service.ts
import { MonitoringScheduler } from './scheduler';
import { ChangeDetector } from './change-detector';
import { AssignmentTracker } from './assignment-tracker';
import { MonitoringEventEmitter } from './event-emitter';
import { adapterRegistry } from '../adapters/registry';

export class BountyMonitoringService {
  private scheduler: MonitoringScheduler;
  private changeDetector: ChangeDetector;
  private assignmentTracker: AssignmentTracker;
  private eventEmitter: MonitoringEventEmitter;
  private isRunning: boolean = false;

  constructor() {
    this.scheduler = new MonitoringScheduler();
    this.changeDetector = new ChangeDetector();
    this.assignmentTracker = new AssignmentTracker();
    this.eventEmitter = new MonitoringEventEmitter();
  }

  /**
   * Start the monitoring service
   */
  async start() {
    if (this.isRunning) {
      console.warn('Monitoring service is already running');
      return;
    }

    console.log('🚀 Starting Bounty Monitoring Service...');

    // Start scheduler
    await this.scheduler.start();

    // Initialize WebSocket
    this.eventEmitter.io.listen(3001); // Separate port for monitoring WS

    this.isRunning = true;
    console.log('✅ Bounty Monitoring Service started');

    // Log status periodically
    setInterval(() => this.logStatus(), 60000); // Every minute
  }

  /**
   * Stop the monitoring service
   */
  async stop() {
    if (!this.isRunning) {
      return;
    }

    console.log('Stopping Bounty Monitoring Service...');

    this.scheduler.stop();
    this.eventEmitter.io.close();

    this.isRunning = false;
    console.log('✅ Bounty Monitoring Service stopped');
  }

  /**
   * Handle platform monitoring job
   */
  async monitorPlatform(platformId: string) {
    try {
      console.log(`Monitoring platform: ${platformId}`);

      // Get adapter for platform
      const adapter = adapterRegistry.getAdapter(platformId);
      if (!adapter) {
        throw new Error(`No adapter found for platform: ${platformId}`);
      }

      // Fetch bounties from platform
      const bounties = await adapter.methods.fetchBounties({
        limit: 100,
        sortBy: 'createdAt',
        sortOrder: 'desc'
      });

      console.log(`Fetched ${bounties.length} bounties from ${platformId}`);

      // Detect changes
      const changes = await this.changeDetector.detectChanges(platformId, bounties);
      console.log(`Detected ${changes.length} changes`);

      // Track assignments for each bounty
      for (const bounty of bounties) {
        await this.assignmentTracker.trackAssignment(bounty);
      }

      // Emit changes to clients
      if (changes.length > 0) {
        await this.eventEmitter.emitChanges(changes);
      }

      // Update database with latest bounties
      await this.updateDatabase(bounties);

      return {
        platformId,
        bountiesFound: bounties.length,
        changesDetected: changes.length,
        timestamp: new Date().toISOString()
      };
    } catch (error) {
      console.error(`Error monitoring ${platformId}:`, error);
      throw error;
    }
  }

  /**
   * Log current monitoring status
   */
  private logStatus() {
    const status = {
      isRunning: this.isRunning,
      timestamp: new Date().toISOString(),
      // Add metrics like:
      // - bounties monitored in last hour
      // - changes detected in last hour
      // - active WebSocket connections
      // - queue size
    };

    console.log('Monitoring Status:', JSON.stringify(status, null, 2));
  }

  private async updateDatabase(bounties: IBounty[]) {
    // Upsert bounties to database
    // This is where you'd implement the actual database logic
  }
}

// Singleton instance
export const monitoringService = new BountyMonitoringService();
```

---

## Configuration

### Environment Variables

```bash
# Monitoring Service Configuration
MONITORING_ENABLED=true
MONITORING_PORT=3001

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Platform-specific intervals (in minutes)
MONITOR_GITCOIN_INTERVAL=10
MONITOR_SUPORTEAM_INTERVAL=5
MONITOR_ALGORA_INTERVAL=10
MONITOR_PROXIES_SX_INTERVAL=3
MONITOR_CODE4RENA_INTERVAL=30
MONITOR_GITHUB_INTERVAL=15

# Alert Configuration
ALERT_EMAIL_ENABLED=true
ALERT_EMAIL_TO=admin@bountyos.com
ALERT_SLACK_WEBHOOK=
ALERT_HIGH_VALUE_THRESHOLD=500

# WebSocket Configuration
WS_CORS_ORIGIN=http://localhost:3000
WS_AUTH_ENABLED=false
```

---

## Running the Service

### Development

```bash
# Start monitoring service
cd bountyos/backend
pnpm monitoring:start

# Or with Docker
docker-compose up monitoring
```

### Production

```bash
# Start as systemd service
sudo systemctl start bountyos-monitoring

# Or with Docker Swarm/Kubernetes
kubectl apply -f k8s/monitoring-deployment.yaml
```

---

## Monitoring Dashboard

### Metrics to Track

```typescript
interface MonitoringMetrics {
  // Platform metrics
  platformsMonitored: number;
  lastSyncTime: Map<string, Date>;
  syncDuration: Map<string, number>;
  
  // Bounty metrics
  totalBounties: number;
  newBountiesToday: number;
  changesDetectedToday: number;
  
  // Assignment metrics
  openBounties: number;
  claimedBounties: number;
  completedBounties: number;
  
  // System metrics
  queueSize: number;
  activeWorkers: number;
  websocketConnections: number;
  errorRate: number;
}
```

---

## Testing

### Unit Tests

```typescript
describe('ChangeDetector', () => {
  it('should detect new bounties', async () => {
    const detector = new ChangeDetector();
    const existingBounties: IBounty[] = [];
    const newBounties: IBounty[] = [
      { id: 'test-1', title: 'Test', ... }
    ];
    
    const changes = await detector.detectChanges('test', newBounties);
    
    expect(changes.length).toBe(1);
    expect(changes[0].type).toBe('NEW');
  });
  
  it('should detect status changes', async () => {
    const detector = new ChangeDetector();
    const existingBounties: IBounty[] = [
      { id: 'test-1', status: 'OPEN', ... }
    ];
    const updatedBounties: IBounty[] = [
      { id: 'test-1', status: 'IN_PROGRESS', ... }
    ];
    
    const changes = await detector.detectChanges('test', updatedBounties);
    
    expect(changes.length).toBe(1);
    expect(changes[0].type).toBe('STATUS_CHANGED');
  });
});
```

---

## Next Steps

1. **Implement Queue System** - Set up BullMQ with Redis
2. **Create Platform Adapters** - Implement for all 5 platforms
3. **Build WebSocket Server** - Real-time event broadcasting
4. **Database Schema Updates** - Add assignment tracking tables
5. **Frontend Integration** - Connect to WebSocket for live updates
6. **Monitoring Dashboard** - Build admin UI for service metrics

---

**Last Updated:** March 13, 2026  
**Maintained By:** bountyOS Core Team
