'use client';

import Link from 'next/link';
import { useAccount } from 'wagmi';

export default function Home() {
  const { isConnected } = useAccount();

  return (
    <div className="min-h-screen bg-gray-900">
      <div className="max-w-6xl mx-auto px-4 py-16">
        <div className="text-center mb-16">
          <h1 className="text-5xl font-bold text-white mb-4">
            Conflux Automation Site
          </h1>
          <p className="text-xl text-gray-400 mb-8">
            Non-custodial limit orders and DCA automation on Conflux eSpace
          </p>
          <div className="flex justify-center gap-4">
            {isConnected ? (
              <Link
                href="/strategy-builder"
                className="bg-blue-600 text-white px-8 py-3 rounded-lg text-lg hover:bg-blue-700"
              >
                Create Strategy
              </Link>
            ) : (
              <button className="bg-blue-600 text-white px-8 py-3 rounded-lg text-lg hover:bg-blue-700">
                Connect Wallet
              </button>
            )}
            <Link
              href="/dashboard"
              className="bg-gray-700 text-white px-8 py-3 rounded-lg text-lg hover:bg-gray-600"
            >
              Dashboard
            </Link>
          </div>
        </div>

        <div className="grid md:grid-cols-2 gap-8">
          <div className="bg-gray-800 rounded-lg p-6">
            <h2 className="text-2xl font-bold text-white mb-4">Limit Orders</h2>
            <p className="text-gray-400">
              Set your target price and automatically execute trades when the price hits your desired level.
              Never miss a trading opportunity again.
            </p>
          </div>
          <div className="bg-gray-800 rounded-lg p-6">
            <h2 className="text-2xl font-bold text-white mb-4">Dollar Cost Averaging</h2>
            <p className="text-gray-400">
              Automate your recurring purchases with DCA strategies. Set intervals and let the worker
              execute your trades automatically.
            </p>
          </div>
        </div>

        <div className="mt-16">
          <h2 className="text-3xl font-bold text-white mb-8 text-center">How It Works</h2>
          <div className="grid md:grid-cols-4 gap-4">
            <div className="text-center">
              <div className="w-12 h-12 bg-blue-600 rounded-full flex items-center justify-center mx-auto mb-4 text-white text-xl font-bold">
                1
              </div>
              <h3 className="text-white font-semibold mb-2">Connect Wallet</h3>
              <p className="text-gray-400 text-sm">Link your Conflux eSpace wallet</p>
            </div>
            <div className="text-center">
              <div className="w-12 h-12 bg-blue-600 rounded-full flex items-center justify-center mx-auto mb-4 text-white text-xl font-bold">
                2
              </div>
              <h3 className="text-white font-semibold mb-2">Configure Strategy</h3>
              <p className="text-gray-400 text-sm">Set your parameters</p>
            </div>
            <div className="text-center">
              <div className="w-12 h-12 bg-blue-600 rounded-full flex items-center justify-center mx-auto mb-4 text-white text-xl font-bold">
                3
              </div>
              <h3 className="text-white font-semibold mb-2">Approve Token</h3>
              <p className="text-gray-400 text-sm">Sign approval transaction</p>
            </div>
            <div className="text-center">
              <div className="w-12 h-12 bg-blue-600 rounded-full flex items-center justify-center mx-auto mb-4 text-white text-xl font-bold">
                4
              </div>
              <h3 className="text-white font-semibold mb-2">Auto Execute</h3>
              <p className="text-gray-400 text-sm">Worker handles the rest</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
