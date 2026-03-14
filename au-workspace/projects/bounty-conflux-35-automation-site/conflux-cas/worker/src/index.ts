import { createPublicClient, createWalletClient, http, getContractAddress, parseEventLogs } from 'viem';
import { privateKeyToAccount } from 'viem/accounts';
import { conflux } from 'viem/chains';

export interface Job {
  id: bigint;
  owner: `0x${string}`;
  jobType: 0 | 1;
  tokenIn: `0x${string}`;
  tokenOut: `0x${string}`;
  amountIn: bigint;
  amountOutMin: bigint;
  targetPrice: bigint;
  interval: bigint;
  lastExecutionTime: bigint;
  endTime: bigint;
  status: 0 | 1 | 2 | 3;
  dcaExecutedCount: bigint;
}

export interface JobExecution {
  jobId: bigint;
  success: boolean;
  amountOut?: bigint;
  error?: string;
  timestamp: Date;
}

export class PriceChecker {
  private rpcUrl: string;
  private priceAdapterAddress: `0x${string}`;

  constructor(rpcUrl: string, priceAdapterAddress: `0x${string}`) {
    this.rpcUrl = rpcUrl;
    this.priceAdapterAddress = priceAdapterAddress;
  }

  async getPrice(tokenIn: `0x${string}`, tokenOut: `0x${string}`): Promise<bigint> {
    const publicClient = createPublicClient({
      chain: conflux,
      transport: http(this.rpcUrl),
    });

    const result = await publicClient.readContract({
      address: this.priceAdapterAddress,
      abi: priceAdapterABI,
      functionName: 'getPrice',
      args: [tokenIn, tokenOut],
    });

    return result as bigint;
  }

  async checkLimitOrder(job: Job, currentPrice: bigint): Promise<boolean> {
    return currentPrice >= job.targetPrice;
  }

  async checkDCA(job: Job): Promise<boolean> {
    const now = BigInt(Math.floor(Date.now() / 1000));
    const timeSinceLastExecution = now - job.lastExecutionTime;
    const preExecutionBuffer = BigInt(15);
    
    return timeSinceLastExecution >= (job.interval - preExecutionBuffer);
  }
}

export class KeeperClient {
  private walletClient: ReturnType<typeof createWalletClient>;
  private publicClient: ReturnType<typeof createPublicClient>;
  private automationManagerAddress: `0x${string}`;
  private maxGasPriceGwei: bigint;
  private rpcTimeoutMs: number;

  constructor(
    executorPrivateKey: `0x${string}`,
    rpcUrl: string,
    automationManagerAddress: `0x${string}`,
    maxGasPriceGwei: bigint = 50n,
    rpcTimeoutMs: number = 120000
  ) {
    const account = privateKeyToAccount(executorPrivateKey);
    
    this.walletClient = createWalletClient({
      chain: conflux,
      transport: http(rpcUrl),
      account,
    });

    this.publicClient = createPublicClient({
      chain: conflux,
      transport: http(rpcUrl),
    });

    this.automationManagerAddress = automationManagerAddress;
    this.maxGasPriceGwei = maxGasPriceGwei;
    this.rpcTimeoutMs = rpcTimeoutMs;
  }

  async simulateExecution(jobId: bigint): Promise<{ success: boolean; revertReason?: string }> {
    try {
      const { request } = await this.publicClient.simulateContract({
        address: this.automationManagerAddress,
        abi: automationManagerABI,
        functionName: 'executeJob',
        args: [jobId],
      });

      return { success: true };
    } catch (error: any) {
      return { success: false, revertReason: error.message };
    }
  }

  async executeJob(jobId: bigint): Promise<{ txHash?: `0x${string}`; error?: string }> {
    try {
      const gasPrice = await this.publicClient.getGasPrice();
      
      if (gasPrice > this.maxGasPriceGwei * BigInt(1e9)) {
        return { error: `Gas price too high: ${gasPrice}` };
      }

      const simulation = await this.simulateExecution(jobId);
      if (!simulation.success) {
        return { error: `Simulation failed: ${simulation.revertReason}` };
      }

      const controller = new AbortController();
      const timeout = setTimeout(() => controller.abort(), this.rpcTimeoutMs);

      try {
        const txHash = await this.walletClient.writeContract({
          address: this.automationManagerAddress,
          abi: automationManagerABI,
          functionName: 'executeJob',
          args: [jobId],
          gas: 500000n,
        });

        await this.publicClient.waitForTransactionReceipt({
          hash: txHash,
          timeout: this.rpcTimeoutMs,
        });

        return { txHash };
      } finally {
        clearTimeout(timeout);
      }
    } catch (error: any) {
      if (error.code === 'ACTION_REJECTED') {
        return { error: 'User rejected transaction' };
      }
      return { error: error.message };
    }
  }

  async isPaused(): Promise<boolean> {
    const paused = await this.publicClient.readContract({
      address: this.automationManagerAddress,
      abi: automationManagerABI,
      functionName: 'paused',
    });
    return paused as boolean;
  }
}

export class RetryQueue {
  private maxRetries: number;
  private baseDelay: number;

  constructor(maxRetries: number = 3, baseDelay: number = 1000) {
    this.maxRetries = maxRetries;
    this.baseDelay = baseDelay;
  }

  async executeWithRetry<T>(fn: () => Promise<T>): Promise<T> {
    let lastError: Error | undefined;
    
    for (let attempt = 0; attempt < this.maxRetries; attempt++) {
      try {
        return await fn();
      } catch (error: any) {
        lastError = error;
        
        if (this.isTransientError(error)) {
          const delay = this.baseDelay * Math.pow(2, attempt);
          await this.sleep(delay);
          continue;
        }
        
        throw error;
      }
    }
    
    throw lastError;
  }

  private isTransientError(error: any): boolean {
    const transientMessages = [
      'PriceConditionNotMet',
      'DCAIntervalNotReached',
      'receipt not indexed',
    ];
    
    return transientMessages.some(msg => 
      error.message?.includes(msg)
    );
  }

  private sleep(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }
}

export class JobPoller {
  private publicClient: ReturnType<typeof createPublicClient>;
  private automationManagerAddress: `0x${string}`;
  private keeperClient: KeeperClient;
  private priceChecker: PriceChecker;
  private retryQueue: RetryQueue;
  private pollIntervalMs: number;
  private isRunning: boolean;
  private dbPaused: boolean;

  constructor(
    rpcUrl: string,
    automationManagerAddress: `0x${string}`,
    priceAdapterAddress: `0x${string}`,
    executorPrivateKey: `0x${string}`,
    pollIntervalMs: number = 30000
  ) {
    this.publicClient = createPublicClient({
      chain: conflux,
      transport: http(rpcUrl),
    });

    this.automationManagerAddress = automationManagerAddress;
    this.keeperClient = new KeeperClient(executorPrivateKey, rpcUrl, automationManagerAddress);
    this.priceChecker = new PriceChecker(rpcUrl, priceAdapterAddress);
    this.retryQueue = new RetryQueue();
    this.pollIntervalMs = pollIntervalMs;
    this.isRunning = false;
    this.dbPaused = false;
  }

  setDbPaused(paused: boolean): void {
    this.dbPaused = paused;
  }

  async start(): Promise<void> {
    this.isRunning = true;
    console.log('[JobPoller] Started');

    while (this.isRunning) {
      try {
        await this.poll();
      } catch (error: any) {
        console.error('[JobPoller] Poll error:', error.message);
      }
      
      await this.sleep(this.pollIntervalMs);
    }
  }

  stop(): void {
    this.isRunning = false;
    console.log('[JobPoller] Stopped');
  }

  private async poll(): Promise<void> {
    const onChainPaused = await this.keeperClient.isPaused();
    
    if (onChainPaused || this.dbPaused) {
      console.log('[JobPoller] Paused (onChain:', onChainPaused, 'db:', this.dbPaused, ')');
      return;
    }

    const jobCount = await this.publicClient.readContract({
      address: this.automationManagerAddress,
      abi: automationManagerABI,
      functionName: 'jobCount',
    });

    for (let i = 1n; i <= jobCount; i++) {
      try {
        await this.processJob(i);
      } catch (error: any) {
        console.error(`[JobPoller] Job ${i} error:`, error.message);
      }
    }
  }

  private async processJob(jobId: bigint): Promise<void> {
    const job = await this.publicClient.readContract({
      address: this.automationManagerAddress,
      abi: automationManagerABI,
      functionName: 'getJob',
      args: [jobId],
    }) as Job;

    if (job.status !== 0n) {
      return;
    }

    const jobType = Number(job.jobType);
    let shouldExecute = false;

    if (jobType === 0) {
      const currentPrice = await this.priceChecker.getPrice(job.tokenIn, job.tokenOut);
      shouldExecute = await this.priceChecker.checkLimitOrder(job, currentPrice);
    } else {
      shouldExecute = await this.priceChecker.checkDCA(job);
    }

    if (!shouldExecute) {
      return;
    }

    console.log(`[JobPoller] Executing job ${jobId}`);
    
    await this.retryQueue.executeWithRetry(async () => {
      const result = await this.keeperClient.executeJob(jobId);
      
      if (result.error) {
        throw new Error(result.error);
      }
      
      console.log(`[JobPoller] Job ${jobId} executed, tx: ${result.txHash}`);
    });
  }

  private sleep(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }
}

const automationManagerABI = [
  {
    type: 'function',
    name: 'jobCount',
    outputs: [{ type: 'uint256', name: '' }],
    stateMutability: 'view',
  },
  {
    type: 'function',
    name: 'getJob',
    inputs: [{ name: '_jobId', type: 'uint256' }],
    outputs: [{
      type: 'tuple',
      name: '',
      components: [
        { name: 'owner', type: 'address' },
        { name: 'jobType', type: 'uint8' },
        { name: 'tokenIn', type: 'address' },
        { name: 'tokenOut', type: 'address' },
        { name: 'amountIn', type: 'uint256' },
        { name: 'amountOutMin', type: 'uint256' },
        { name: 'targetPrice', type: 'uint256' },
        { name: 'interval', type: 'uint256' },
        { name: 'lastExecutionTime', type: 'uint256' },
        { name: 'endTime', type: 'uint256' },
        { name: 'status', type: 'uint8' },
        { name: 'dcaExecutedCount', type: 'uint256' },
      ],
    }],
    stateMutability: 'view',
  },
  {
    type: 'function',
    name: 'executeJob',
    inputs: [{ name: '_jobId', type: 'uint256' }],
    outputs: [],
    stateMutability: 'nonpayable',
  },
  {
    type: 'function',
    name: 'paused',
    outputs: [{ type: 'bool', name: '' }],
    stateMutability: 'view',
  },
] as const;

const priceAdapterABI = [
  {
    type: 'function',
    name: 'getPrice',
    inputs: [
      { name: '_tokenIn', type: 'address' },
      { name: '_tokenOut', type: 'address' },
    ],
    outputs: [{ type: 'uint256', name: '' }],
    stateMutability: 'view',
  },
  {
    type: 'function',
    name: 'getSwapRouter',
    outputs: [{ type: 'address', name: '' }],
    stateMutability: 'view',
  },
] as const;
