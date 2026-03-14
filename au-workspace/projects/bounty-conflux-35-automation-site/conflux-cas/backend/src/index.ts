import express from 'express';
import { createServer } from 'http';
import { drizzle } from 'drizzle-orm/sqlite3';
import SQLite from 'better-sqlite3';
import { jobs, executions, settings } from './db/schema.js';
import { eq, desc } from 'drizzle-orm';
import { verifyMessage } from 'viem';
import jwt from 'jsonwebtoken';
import { cors } from 'cors';

const sqlite = new SQLite('./cas.db');
const db = drizzle(sqlite);

const app = express();
app.use(express.json());
app.use(cors());

const JWT_SECRET = process.env.JWT_SECRET || 'dev-secret';

function generateToken(address: string): string {
  return jwt.sign({ address }, JWT_SECRET, { expiresIn: '7d' });
}

function authenticate(req: express.Request, res: express.Response, next: express.NextFunction) {
  const authHeader = req.headers.authorization;
  if (!authHeader?.startsWith('Bearer ')) {
    return res.status(401).json({ error: 'No token provided' });
  }

  try {
    const token = authHeader.substring(7);
    const decoded = jwt.verify(token, JWT_SECRET) as { address: string };
    (req as any).user = decoded;
    next();
  } catch {
    return res.status(401).json({ error: 'Invalid token' });
  }
}

app.post('/api/auth/login', async (req, res) => {
  const { signature, message } = req.body;
  
  try {
    const recoveredAddress = verifyMessage({
      message,
      signature: signature as `0x${string}`,
    });
    
    const token = generateToken(recoveredAddress);
    res.json({ token, address: recoveredAddress });
  } catch (error: any) {
    res.status(401).json({ error: 'Invalid signature' });
  }
});

app.get('/api/jobs', authenticate, async (req, res) => {
  try {
    const userAddress = (req as any).user.address.toLowerCase();
    const userJobs = await db.select().from(jobs)
      .where(eq(jobs.owner, userAddress))
      .orderBy(desc(jobs.createdAt));
    res.json(userJobs);
  } catch (error: any) {
    res.status(500).json({ error: error.message });
  }
});

app.post('/api/jobs', authenticate, async (req, res) => {
  const { jobType, tokenIn, tokenOut, amountIn, amountOutMin, targetPrice, interval, endTime, chainJobId } = req.body;
  const owner = (req as any).user.address.toLowerCase();

  try {
    const newJob = await db.insert(jobs).values({
      owner,
      jobType,
      tokenIn: tokenIn.toLowerCase(),
      tokenOut: tokenOut.toLowerCase(),
      amountIn: BigInt(amountIn),
      amountOutMin: BigInt(amountOutMin),
      targetPrice: targetPrice ? BigInt(targetPrice) : null,
      interval: interval ? BigInt(interval) : null,
      endTime: endTime ? new Date(endTime) : null,
      chainJobId: BigInt(chainJobId),
      status: 'ACTIVE',
      createdAt: new Date(),
      updatedAt: new Date(),
    }).returning();
    
    res.json(newJob[0]);
  } catch (error: any) {
    res.status(500).json({ error: error.message });
  }
});

app.patch('/api/jobs/:id/pause', authenticate, async (req, res) => {
  const { id } = req.params;
  const owner = (req as any).user.address.toLowerCase();

  try {
    await db.update(jobs)
      .set({ status: 'PAUSED', updatedAt: new Date() })
      .where(eq(jobs.id, parseInt(id)))
      .where(eq(jobs.owner, owner));
    
    res.json({ success: true });
  } catch (error: any) {
    res.status(500).json({ error: error.message });
  }
});

app.patch('/api/jobs/:id/resume', authenticate, async (req, res) => {
  const { id } = req.params;
  const owner = (req as any).user.address.toLowerCase();

  try {
    await db.update(jobs)
      .set({ status: 'ACTIVE', updatedAt: new Date() })
      .where(eq(jobs.id, parseInt(id)))
      .where(eq(jobs.owner, owner));
    
    res.json({ success: true });
  } catch (error: any) {
    res.status(500).json({ error: error.message });
  }
});

app.patch('/api/jobs/:id/cancel', authenticate, async (req, res) => {
  const { id } = req.params;
  const owner = (req as any).user.address.toLowerCase();

  try {
    await db.update(jobs)
      .set({ status: 'CANCELLED', updatedAt: new Date() })
      .where(eq(jobs.id, parseInt(id)))
      .where(eq(jobs.owner, owner));
    
    res.json({ success: true });
  } catch (error: any) {
    res.status(500).json({ error: error.message });
  }
});

app.get('/api/jobs/:id/executions', authenticate, async (req, res) => {
  const { id } = req.params;

  try {
    const jobExecutions = await db.select().from(executions)
      .where(eq(executions.jobId, parseInt(id)))
      .orderBy(desc(executions.timestamp));
    res.json(jobExecutions);
  } catch (error: any) {
    res.status(500).json({ error: error.message });
  }
});

app.get('/api/settings', async (req, res) => {
  try {
    const globalSettings = await db.select().from(settings).limit(1);
    res.json(globalSettings[0] || { paused: false });
  } catch (error: any) {
    res.status(500).json({ error: error.message });
  }
});

app.post('/api/admin/pause', authenticate, async (req, res) => {
  const { adminKey } = req.body;
  
  if (adminKey !== process.env.ADMIN_KEY) {
    return res.status(403).json({ error: 'Invalid admin key' });
  }

  try {
    await db.insert(settings).values({ paused: true, updatedAt: new Date() })
      .onConflictDoUpdate({ target: settings.id, set: { paused: true, updatedAt: new Date() } });
    res.json({ success: true });
  } catch (error: any) {
    res.status(500).json({ error: error.message });
  }
});

app.post('/api/admin/resume', authenticate, async (req, res) => {
  const { adminKey } = req.body;
  
  if (adminKey !== process.env.ADMIN_KEY) {
    return res.status(403).json({ error: 'Invalid admin key' });
  }

  try {
    await db.insert(settings).values({ paused: false, updatedAt: new Date() })
      .onConflictDoUpdate({ target: settings.id, set: { paused: false, updatedAt: new Date() } });
    res.json({ success: true });
  } catch (error: any) {
    res.status(500).json({ error: error.message });
  }
});

const server = createServer(app);

app.get('/api/sse/events', (req, res) => {
  res.setHeader('Content-Type', 'text/event-stream');
  res.setHeader('Cache-Control', 'no-cache');
  res.setHeader('Connection', 'keep-alive');

  const sendEvent = (data: any) => {
    res.write(`data: ${JSON.stringify(data)}\n\n`);
  };

  const executionHandler = (execution: any) => {
    sendEvent({ type: 'execution', data: execution });
  };

  res.on('close', () => {
  });

  sendEvent({ type: 'connected' });
});

const PORT = process.env.PORT || 3001;
app.listen(PORT, () => {
  console.log(`[Backend] Running on port ${PORT}`);
});

export { db };
