import { HardhatUserConfig } from "hardhat/config";

const config: HardhatUserConfig = {
  solidity: {
    version: "0.8.19",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200,
      },
    },
  },
  networks: {
    conflux: {
      url: "https://evm.confluxrpc.com",
      chainId: 1030,
    },
    confluxTestnet: {
      url: "https://evmtestnet.confluxrpc.com",
      chainId: 71,
    },
  },
};

export default config;
