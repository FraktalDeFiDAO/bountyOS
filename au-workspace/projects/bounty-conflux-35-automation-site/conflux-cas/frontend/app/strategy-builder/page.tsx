'use client';

import { useState } from 'react';
import { useAccount, useWriteContract, useWaitForTransactionReceipt } from 'wagmi';
import { parseEther } from 'viem';

type JobType = 'LIMIT_ORDER' | 'DCA';

interface StrategyForm {
  jobType: JobType;
  tokenIn: string;
  tokenOut: string;
  amountIn: string;
  amountOutMin: string;
  targetPrice: string;
  interval: string;
  endTime: string;
}

export default function StrategyBuilder() {
  const { address, isConnected } = useAccount();
  const [step, setStep] = useState(1);
  const [form, setForm] = useState<StrategyForm>({
    jobType: 'DCA',
    tokenIn: '',
    tokenOut: '',
    amountIn: '',
    amountOutMin: '',
    targetPrice: '',
    interval: '',
    endTime: '',
  });
  const [txHash, setTxHash] = useState<string | null>(null);

  const { writeContractAsync } = useWriteContract();
  const { isLoading: isConfirming, isSuccess } = useWaitForTransactionReceipt({
    hash: txHash as `0x${string}`,
  });

  const handleSubmit = async () => {
    if (!address) return;

    try {
      let tx;
      if (form.jobType === 'LIMIT_ORDER') {
        tx = await writeContractAsync({
          address: form.tokenIn as `0x${string}`,
          abi: ['function approve(address spender, uint256 amount) returns (bool)'],
          functionName: 'approve',
          args: [import.meta.env.VITE_AUTOMATION_MANAGER_ADDRESS, parseEther(form.amountIn)],
        });
      } else {
        tx = await writeContractAsync({
          address: form.tokenIn as `0x${string}`,
          abi: ['function approve(address spender, uint256 amount) returns (bool)'],
          functionName: 'approve',
          args: [import.meta.env.VITE_AUTOMATION_MANAGER_ADDRESS, parseEther(form.amountIn)],
        });
      }
      setTxHash(tx);
      setStep(3);
    } catch (error) {
      console.error('Transaction failed:', error);
    }
  };

  const handleCreateJob = async () => {
    if (!address || !txHash) return;

    try {
      const automationManagerAddress = import.meta.env.VITE_AUTOMATION_MANAGER_ADDRESS;
      
      let chainJobId;
      if (form.jobType === 'LIMIT_ORDER') {
        chainJobId = await writeContractAsync({
          address: automationManagerAddress as `0x${string}`,
          abi: ['function createLimitOrder(address _tokenIn, address _tokenOut, uint256 _amountIn, uint256 _amountOutMin, uint256 _targetPrice) returns (uint256)'],
          functionName: 'createLimitOrder',
          args: [
            form.tokenIn as `0x${string}`,
            form.tokenOut as `0x${string}`,
            parseEther(form.amountIn),
            parseEther(form.amountOutMin),
            parseEther(form.targetPrice),
          ],
        });
      } else {
        chainJobId = await writeContractAsync({
          address: automationManagerAddress as `0x${string}`,
          abi: ['function createDCAJob(address _tokenIn, address _tokenOut, uint256 _amountIn, uint256 _amountOutMin, uint256 _interval, uint256 _endTime) returns (uint256)'],
          functionName: 'createDCAJob',
          args: [
            form.tokenIn as `0x${string}`,
            form.tokenOut as `0x${string}`,
            parseEther(form.amountIn),
            parseEther(form.amountOutMin),
            BigInt(form.interval),
            BigInt(Math.floor(new Date(form.endTime).getTime() / 1000)),
          ],
        });
      }
      
      setStep(4);
    } catch (error) {
      console.error('Failed to create job:', error);
    }
  };

  if (!isConnected) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-900">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-white mb-4">Connect Your Wallet</h1>
          <p className="text-gray-400">Please connect your wallet to create automation strategies</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-900 p-8">
      <div className="max-w-2xl mx-auto">
        <h1 className="text-3xl font-bold text-white mb-8">Strategy Builder</h1>
        
        <div className="flex items-center mb-8">
          {[1, 2, 3, 4].map((s) => (
            <div key={s} className="flex items-center">
              <div className={`w-8 h-8 rounded-full flex items-center justify-center ${
                step >= s ? 'bg-blue-600' : 'bg-gray-700'
              } text-white`}>
                {s}
              </div>
              {s < 4 && (
                <div className={`w-16 h-1 ${step > s ? 'bg-blue-600' : 'bg-gray-700'}`} />
              )}
            </div>
          ))}
        </div>

        <div className="bg-gray-800 rounded-lg p-6">
          {step === 1 && (
            <div className="space-y-4">
              <h2 className="text-xl font-semibold text-white mb-4">Step 1: Configure Strategy</h2>
              
              <div>
                <label className="block text-gray-400 mb-2">Strategy Type</label>
                <div className="flex gap-4">
                  <button
                    className={`px-4 py-2 rounded ${
                      form.jobType === 'LIMIT_ORDER' ? 'bg-blue-600' : 'bg-gray-700'
                    } text-white`}
                    onClick={() => setForm({ ...form, jobType: 'LIMIT_ORDER' })}
                  >
                    Limit Order
                  </button>
                  <button
                    className={`px-4 py-2 rounded ${
                      form.jobType === 'DCA' ? 'bg-blue-600' : 'bg-gray-700'
                    } text-white`}
                    onClick={() => setForm({ ...form, jobType: 'DCA' })}
                  >
                    DCA
                  </button>
                </div>
              </div>

              <div>
                <label className="block text-gray-400 mb-2">From Token</label>
                <input
                  type="text"
                  className="w-full bg-gray-700 text-white rounded px-4 py-2"
                  placeholder="0x..."
                  value={form.tokenIn}
                  onChange={(e) => setForm({ ...form, tokenIn: e.target.value })}
                />
              </div>

              <div>
                <label className="block text-gray-400 mb-2">To Token</label>
                <input
                  type="text"
                  className="w-full bg-gray-700 text-white rounded px-4 py-2"
                  placeholder="0x..."
                  value={form.tokenOut}
                  onChange={(e) => setForm({ ...form, tokenOut: e.target.value })}
                />
              </div>

              <div>
                <label className="block text-gray-400 mb-2">Amount</label>
                <input
                  type="text"
                  className="w-full bg-gray-700 text-white rounded px-4 py-2"
                  placeholder="0.0"
                  value={form.amountIn}
                  onChange={(e) => setForm({ ...form, amountIn: e.target.value })}
                />
              </div>

              <div>
                <label className="block text-gray-400 mb-2">Min Output (Slippage)</label>
                <input
                  type="text"
                  className="w-full bg-gray-700 text-white rounded px-4 py-2"
                  placeholder="0.0"
                  value={form.amountOutMin}
                  onChange={(e) => setForm({ ...form, amountOutMin: e.target.value })}
                />
              </div>

              {form.jobType === 'LIMIT_ORDER' && (
                <div>
                  <label className="block text-gray-400 mb-2">Target Price</label>
                  <input
                    type="text"
                    className="w-full bg-gray-700 text-white rounded px-4 py-2"
                    placeholder="0.0"
                    value={form.targetPrice}
                    onChange={(e) => setForm({ ...form, targetPrice: e.target.value })}
                  />
                </div>
              )}

              {form.jobType === 'DCA' && (
                <>
                  <div>
                    <label className="block text-gray-400 mb-2">Interval (seconds)</label>
                    <input
                      type="text"
                      className="w-full bg-gray-700 text-white rounded px-4 py-2"
                      placeholder="86400"
                      value={form.interval}
                      onChange={(e) => setForm({ ...form, interval: e.target.value })}
                    />
                  </div>
                  <div>
                    <label className="block text-gray-400 mb-2">End Time</label>
                    <input
                      type="datetime-local"
                      className="w-full bg-gray-700 text-white rounded px-4 py-2"
                      value={form.endTime}
                      onChange={(e) => setForm({ ...form, endTime: e.target.value })}
                    />
                  </div>
                </>
              )}

              <button
                className="w-full bg-blue-600 text-white py-3 rounded mt-4"
                onClick={() => setStep(2)}
              >
                Continue
              </button>
            </div>
          )}

          {step === 2 && (
            <div className="space-y-4">
              <h2 className="text-xl font-semibold text-white mb-4">Step 2: Approve Token</h2>
              <p className="text-gray-400">
                You need to approve the Automation Manager to spend your {form.tokenIn}
              </p>
              <div className="bg-gray-700 p-4 rounded">
                <p className="text-white">Approve: {form.amountIn} {form.tokenIn}</p>
              </div>
              <button
                className="w-full bg-blue-600 text-white py-3 rounded"
                onClick={handleSubmit}
                disabled={isConfirming}
              >
                {isConfirming ? 'Confirming...' : 'Approve'}
              </button>
              <button
                className="w-full bg-gray-700 text-white py-3 rounded mt-2"
                onClick={() => setStep(1)}
              >
                Back
              </button>
            </div>
          )}

          {step === 3 && (
            <div className="space-y-4">
              <h2 className="text-xl font-semibold text-white mb-4">Step 3: Register on Chain</h2>
              <p className="text-gray-400">
                Now register your strategy on the blockchain
              </p>
              <button
                className="w-full bg-blue-600 text-white py-3 rounded"
                onClick={handleCreateJob}
                disabled={isConfirming}
              >
                {isConfirming ? 'Confirming...' : 'Create Job'}
              </button>
            </div>
          )}

          {step === 4 && (
            <div className="space-y-4">
              <h2 className="text-xl font-semibold text-white mb-4">Success!</h2>
              <p className="text-green-400">Your strategy has been created successfully!</p>
              <button
                className="w-full bg-blue-600 text-white py-3 rounded"
                onClick={() => window.location.href = '/dashboard'}
              >
                View Dashboard
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
