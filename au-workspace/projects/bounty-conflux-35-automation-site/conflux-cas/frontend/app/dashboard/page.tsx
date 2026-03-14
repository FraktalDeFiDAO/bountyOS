'use client';

import { useState, useEffect } from 'react';
import { useAccount } from 'wagmi';

interface Job {
  id: number;
  jobType: string;
  tokenIn: string;
  tokenOut: string;
  amountIn: string;
  amountOutMin: string;
  status: string;
  createdAt: string;
}

export default function Dashboard() {
  const { address, isConnected } = useAccount();
  const [jobs, setJobs] = useState<Job[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (address) {
      fetchJobs();
    }
  }, [address]);

  const fetchJobs = async () => {
    try {
      const token = localStorage.getItem('token');
      const res = await fetch('/api/jobs', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      if (res.ok) {
        const data = await res.json();
        setJobs(data);
      }
    } catch (error) {
      console.error('Failed to fetch jobs:', error);
    } finally {
      setLoading(false);
    }
  };

  const handlePause = async (jobId: number) => {
    try {
      const token = localStorage.getItem('token');
      await fetch(`/api/jobs/${jobId}/pause`, {
        method: 'PATCH',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      fetchJobs();
    } catch (error) {
      console.error('Failed to pause job:', error);
    }
  };

  const handleResume = async (jobId: number) => {
    try {
      const token = localStorage.getItem('token');
      await fetch(`/api/jobs/${jobId}/resume`, {
        method: 'PATCH',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      fetchJobs();
    } catch (error) {
      console.error('Failed to resume job:', error);
    }
  };

  const handleCancel = async (jobId: number) => {
    try {
      const token = localStorage.getItem('token');
      await fetch(`/api/jobs/${jobId}/cancel`, {
        method: 'PATCH',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      fetchJobs();
    } catch (error) {
      console.error('Failed to cancel job:', error);
    }
  };

  if (!isConnected) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-900">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-white mb-4">Connect Your Wallet</h1>
          <p className="text-gray-400">Please connect your wallet to view your strategies</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-900 p-8">
      <div className="max-w-6xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold text-white">Dashboard</h1>
          <a
            href="/strategy-builder"
            className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700"
          >
            Create Strategy
          </a>
        </div>

        {loading ? (
          <div className="text-white">Loading...</div>
        ) : jobs.length === 0 ? (
          <div className="text-center text-gray-400">
            <p>No strategies found</p>
            <a href="/strategy-builder" className="text-blue-400 hover:underline">
              Create your first strategy
            </a>
          </div>
        ) : (
          <div className="grid gap-4">
            {jobs.map((job) => (
              <div key={job.id} className="bg-gray-800 rounded-lg p-6">
                <div className="flex justify-between items-start">
                  <div>
                    <div className="flex items-center gap-3 mb-2">
                      <span className="text-white font-semibold">
                        {job.jobType === 'LIMIT_ORDER' ? 'Limit Order' : 'DCA'}
                      </span>
                      <span className={`px-2 py-1 rounded text-xs ${
                        job.status === 'ACTIVE' ? 'bg-green-600' :
                        job.status === 'PAUSED' ? 'bg-yellow-600' :
                        job.status === 'CANCELLED' ? 'bg-red-600' : 'bg-gray-600'
                      } text-white`}>
                        {job.status}
                      </span>
                    </div>
                    <p className="text-gray-400">
                      {job.amountIn} {job.tokenIn} → {job.amountOutMin} {job.tokenOut}
                    </p>
                    <p className="text-gray-500 text-sm mt-1">
                      Created: {new Date(job.createdAt).toLocaleDateString()}
                    </p>
                  </div>
                  <div className="flex gap-2">
                    {job.status === 'ACTIVE' && (
                      <button
                        onClick={() => handlePause(job.id)}
                        className="px-3 py-1 bg-yellow-600 text-white rounded text-sm hover:bg-yellow-700"
                      >
                        Pause
                      </button>
                    )}
                    {job.status === 'PAUSED' && (
                      <button
                        onClick={() => handleResume(job.id)}
                        className="px-3 py-1 bg-green-600 text-white rounded text-sm hover:bg-green-700"
                      >
                        Resume
                      </button>
                    )}
                    {job.status !== 'CANCELLED' && job.status !== 'COMPLETED' && (
                      <button
                        onClick={() => handleCancel(job.id)}
                        className="px-3 py-1 bg-red-600 text-white rounded text-sm hover:bg-red-700"
                      >
                        Cancel
                      </button>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
