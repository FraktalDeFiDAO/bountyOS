import { sqliteTable, text, integer, real } from 'drizzle-orm/sqlite-core';

export const jobs = sqliteTable('jobs', {
  id: integer('id').primaryKey({ autoIncrement: true }),
  owner: text('owner').notNull(),
  jobType: text('job_type').notNull(),
  tokenIn: text('token_in').notNull(),
  tokenOut: text('token_out').notNull(),
  amountIn: integer('amount_in').notNull(),
  amountOutMin: integer('amount_out_min').notNull(),
  targetPrice: integer('target_price'),
  interval: integer('interval'),
  endTime: integer('end_time'),
  chainJobId: integer('chain_job_id').notNull(),
  status: text('status').notNull().default('ACTIVE'),
  createdAt: integer('created_at', { mode: 'timestamp' }).notNull(),
  updatedAt: integer('updated_at', { mode: 'timestamp' }).notNull(),
});

export const executions = sqliteTable('executions', {
  id: integer('id').primaryKey({ autoIncrement: true }),
  jobId: integer('job_id').notNull(),
  txHash: text('tx_hash'),
  amountIn: integer('amount_in').notNull(),
  amountOut: integer('amount_out'),
  success: integer('success', { mode: 'boolean' }).notNull(),
  error: text('error'),
  timestamp: integer('timestamp', { mode: 'timestamp' }).notNull(),
});

export const settings = sqliteTable('settings', {
  id: integer('id').primaryKey({ defaultValue: 1 }),
  paused: integer('paused', { mode: 'boolean' }).notNull().default(false),
  maxSwapUsd: real('max_swap_usd'),
  maxSlippageBps: integer('max_slippage_bps'),
  updatedAt: integer('updated_at', { mode: 'timestamp' }).notNull(),
});
