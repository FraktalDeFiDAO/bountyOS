'use client';

import { WagmiProvider } from 'wagmi';
import { conflux } from 'viem/chains';
import { http, createConfig } from 'wagmi';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { injected } from 'wagmi/connectors';

const config = createConfig({
  chains: [conflux],
  transports: {
    [conflux.id]: http(),
  },
  connectors: [injected()],
});

const queryClient = new QueryClient();

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <head>
        <title>Conflux Automation Site</title>
      </head>
      <body>
        <WagmiProvider config={config}>
          <QueryClientProvider client={queryClient}>
            {children}
          </QueryClientProvider>
        </WagmiProvider>
      </body>
    </html>
  );
}
