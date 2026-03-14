// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/access/Ownable.sol";

interface ISwappiPair {
    function getReserves() external view returns (uint112 reserve0, uint112 reserve1, uint32 blockTimestampLast);
    function token0() external view returns (address);
    function token1() external view returns (address);
}

interface IERC20Metadata is IERC20 {
    function decimals() external view returns (uint8);
}

contract SwappiPriceAdapter is Ownable {
    address public swapRouter;
    mapping(address => mapping(address => address)) public tokenPairs;

    event PairSet(address indexed token0, address indexed token1, address indexed pair);
    event RouterSet(address indexed router);

    constructor() Ownable() {}

    function setSwapRouter(address _router) external onlyOwner {
        require(_router != address(0), "Invalid router");
        swapRouter = _router;
        emit RouterSet(_router);
    }

    function setTokenPair(address _tokenA, address _tokenB, address _pair) external onlyOwner {
        require(_tokenA != address(0) && _tokenB != address(0), "Invalid tokens");
        require(_pair != address(0), "Invalid pair");
        tokenPairs[_tokenA][_tokenB] = _pair;
        tokenPairs[_tokenB][_tokenA] = _pair;
        emit PairSet(_tokenA, _tokenB, _pair);
    }

    function getPrice(address _tokenIn, address _tokenOut) external view override returns (uint256) {
        address pair = tokenPairs[_tokenIn][_tokenOut];
        require(pair != address(0), "Pair not found");

        ISwappiPair swappiPair = ISwappiPair(pair);
        (uint112 reserve0, uint112 reserve1, ) = swappiPair.getReserves();

        address token0 = swappiPair.token0();
        uint8 decimalsIn = IERC20Metadata(_tokenIn).decimals();
        uint8 decimalsOut = IERC20Metadata(_tokenOut).decimals();

        uint256 amountIn = 10 ** decimalsIn;
        uint256 amountOut;

        if (_tokenIn == token0) {
            amountOut = (amountIn * reserve1) / reserve0;
            amountOut = amountOut * (10 ** decimalsOut) / (10 ** decimalsIn);
        } else {
            amountOut = (amountIn * reserve0) / reserve1;
            amountOut = amountOut * (10 ** decimalsOut) / (10 ** decimalsIn);
        }

        return amountOut;
    }

    function getSwapRouter() external view override returns (address) {
        return swapRouter;
    }
}
