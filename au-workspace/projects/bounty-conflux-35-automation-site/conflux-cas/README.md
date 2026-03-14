# Conflux Automation Site

Non-custodial limit-order and DCA automation on Conflux eSpace.

## Overview

Conflux Automation Site (CAS) enables users to create and manage automated trading strategies without giving up custody of their funds. Users connect their wallet, configure a strategy once, sign an ERC-20 approval, and a keeper worker executes trades automatically.

## Architecture

The system consists of four integrated layers:

1. **Smart Contracts** (`conflux-contracts/`) - Solidity contracts for job management
2. **Execution Worker** (`conflux-cas/worker/`) - Node.js service that polls jobs and executes trades
3. **Backend API** (`conflux-cas/backend/`) - Express API for job CRUD and authentication
4. **Frontend** (`conflux-cas/frontend/`) - Next.js app with Strategy Builder and Dashboard

## Quick Start

### Prerequisites

- Node.js 18+
- pnpm
- Conflux eSpace wallet (MetaMask, Fluent, etc.)

### Environment Variables

```bash
# Backend
DATABASE_URL=./cas.db
JWT_SECRET=your-secret-key
ADMIN_KEY=your-admin-key
PORT=3001

# Worker
CONFLUX_RPC_URL=https://evm.confluxrpc.com
EXECUTOR_PRIVATE_KEY=0x...
AUTOMATION_MANAGER_ADDRESS=0x...
PRICE_ADAPTER_ADDRESS=0x...
WORKER_RPC_TIMEOUT_MS=120000

# Frontend
VITE_AUTOMATION_MANAGER_ADDRESS=0x...
VITE_API_URL=http://localhost:3001
```

### Deployment

1. **Deploy Smart Contracts**
```bash
cd conflux-contracts
npm install
npx hardhat run scripts/deploy.js --network conflux
```

2. **Start Backend**
```bash
cd conflux-cas/backend
npm install
npm run dev
```

3. **Start Worker**
```bash
cd conflux-cas/worker
npm install
npm run start
```

4. **Start Frontend**
```bash
cd conflux-cas/frontend
npm install
npm run dev
```

### Docker Compose

```bash
docker-compose up -d
```

## Features

- **Limit Orders**: Execute trades when price reaches target
- **DCA**: Automated recurring purchases at set intervals
- **Non-custodial**: Users retain custody of funds at all times
- **Safety Controls**: Global pause, per-job cancel, slippage protection
- **Real-time Updates**: SSE events for live execution status

## API Endpoints

- `POST /api/auth/login` - SIWE authentication
- `GET /api/jobs` - List user jobs
- `POST /api/jobs` - Create new job
- `PATCH /api/jobs/:id/pause` - Pause job
- `PATCH /api/jobs/:id/resume` - Resume job
- `PATCH /api/jobs/:id/cancel` - Cancel job
- `GET /api/jobs/:id/executions` - Get job execution history
- `GET /api/sse/events` - Real-time execution events

## Safety Controls

- **Global Pause**: Admin can pause all executions
- **Per-job Control**: Users can pause/resume/cancel their jobs
- **Slippage Protection**: Configurable minimum output
- **Gas Price Circuit Breaker**: Abort if gas price too high
- **Transient Error Handling**: Skip without burning retry budget

## Testing

```bash
# Run all tests
npm test

# Run contract tests
cd conflux-contracts && npx hardhat test

# Run worker tests
cd conflux-cas/worker && npm test
```

## License

MIT
