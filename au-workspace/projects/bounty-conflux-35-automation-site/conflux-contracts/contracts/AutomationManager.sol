// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/security/Pausable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/Address.sol";

contract AutomationManager is ReentrancyGuard, Pausable {
    using SafeERC20 for IERC20;
    using Address for address;

    enum JobType { LIMIT_ORDER, DCA }
    enum JobStatus { ACTIVE, PAUSED, COMPLETED, CANCELLED }

    struct Job {
        address owner;
        JobType jobType;
        address tokenIn;
        address tokenOut;
        uint256 amountIn;
        uint256 amountOutMin;
        uint256 targetPrice;
        uint256 interval;
        uint256 lastExecutionTime;
        uint256 endTime;
        JobStatus status;
        uint256 dcaExecutedCount;
    }

    uint256 public jobCount;
    mapping(uint256 => Job) public jobs;
    mapping(address => uint256[]) public ownerJobs;

    event JobCreated(uint256 indexed jobId, address indexed owner, JobType jobType);
    event JobExecuted(uint256 indexed jobId, uint256 amountOut, uint256 timestamp);
    event JobPaused(uint256 indexed jobId);
    event JobResumed(uint256 indexed jobId);
    event JobCancelled(uint256 indexed jobId);
    event JobCompleted(uint256 indexed jobId);

    address public priceAdapter;
    address public admin;

    modifier onlyJobOwner(uint256 _jobId) {
        require(jobs[_jobId].owner == msg.sender, "Not job owner");
        _;
    }

    modifier onlyAdmin() {
        require(msg.sender == admin, "Not admin");
        _;
    }

    constructor(address _priceAdapter) {
        priceAdapter = _priceAdapter;
        admin = msg.sender;
    }

    function setPriceAdapter(address _priceAdapter) external onlyAdmin {
        priceAdapter = _priceAdapter;
    }

    function createLimitOrder(
        address _tokenIn,
        address _tokenOut,
        uint256 _amountIn,
        uint256 _amountOutMin,
        uint256 _targetPrice
    ) external whenNotPaused returns (uint256) {
        require(_tokenIn != address(0) && _tokenOut != address(0), "Invalid tokens");
        require(_amountIn > 0 && _targetPrice > 0, "Invalid params");

        jobCount++;
        uint256 jobId = jobCount;

        jobs[jobId] = Job({
            owner: msg.sender,
            jobType: JobType.LIMIT_ORDER,
            tokenIn: _tokenIn,
            tokenOut: _tokenOut,
            amountIn: _amountIn,
            amountOutMin: _amountOutMin,
            targetPrice: _targetPrice,
            interval: 0,
            lastExecutionTime: block.timestamp,
            endTime: 0,
            status: JobStatus.ACTIVE,
            dcaExecutedCount: 0
        });

        ownerJobs[msg.sender].push(jobId);
        emit JobCreated(jobId, msg.sender, JobType.LIMIT_ORDER);

        return jobId;
    }

    function createDCAJob(
        address _tokenIn,
        address _tokenOut,
        uint256 _amountIn,
        uint256 _amountOutMin,
        uint256 _interval,
        uint256 _endTime
    ) external whenNotPaused returns (uint256) {
        require(_tokenIn != address(0) && _tokenOut != address(0), "Invalid tokens");
        require(_amountIn > 0 && _interval > 0, "Invalid params");
        require(_endTime > block.timestamp, "Invalid end time");

        jobCount++;
        uint256 jobId = jobCount;

        jobs[jobId] = Job({
            owner: msg.sender,
            jobType: JobType.DCA,
            tokenIn: _tokenIn,
            tokenOut: _tokenOut,
            amountIn: _amountIn,
            amountOutMin: _amountOutMin,
            targetPrice: 0,
            interval: _interval,
            lastExecutionTime: block.timestamp,
            endTime: _endTime,
            status: JobStatus.ACTIVE,
            dcaExecutedCount: 0
        });

        ownerJobs[msg.sender].push(jobId);
        emit JobCreated(jobId, msg.sender, JobType.DCA);

        return jobId;
    }

    function executeJob(uint256 _jobId) external nonReentrant whenNotPaused {
        Job storage job = jobs[_jobId];
        require(job.status == JobStatus.ACTIVE, "Job not active");
        require(block.timestamp >= job.lastExecutionTime + job.interval || job.jobType == JobType.LIMIT_ORDER, "Interval not reached");

        if (job.jobType == JobType.DCA) {
            require(block.timestamp < job.endTime, "Job expired");
        }

        uint256 currentPrice = IPriceAdapter(priceAdapter).getPrice(job.tokenIn, job.tokenOut);
        
        if (job.jobType == JobType.LIMIT_ORDER) {
            require(currentPrice >= job.targetPrice, "Price condition not met");
        }

        uint256 balance = IERC20(job.tokenIn).balanceOf(job.owner);
        require(balance >= job.amountIn, "Insufficient balance");

        IERC20(job.tokenIn).safeTransferFrom(job.owner, address(this), job.amountIn);
        
        uint256 amountOut = swap(job.tokenIn, job.tokenOut, job.amountIn, job.amountOutMin);
        
        require(amountOut >= job.amountOutMin, "Slippage exceeded");

        IERC20(job.tokenOut).safeTransfer(job.owner, amountOut);

        job.lastExecutionTime = block.timestamp;
        job.dcaExecutedCount++;

        if (job.jobType == JobType.DCA && block.timestamp >= job.endTime) {
            job.status = JobStatus.COMPLETED;
            emit JobCompleted(_jobId);
        }

        emit JobExecuted(_jobId, amountOut, block.timestamp);
    }

    function swap(
        address _tokenIn,
        address _tokenOut,
        uint256 _amountIn,
        uint256 _amountOutMin
    ) internal returns (uint256) {
        address swapRouter = IPriceAdapter(priceAdapter).getSwapRouter();
        require(swapRouter != address(0), "No swap router");

        IERC20(_tokenIn).safeApprove(swapRouter, _amountIn);

        bytes memory data = abi.encodeWithSelector(
            bytes4(keccak256("swapExactTokensForTokens(uint256,uint256,address[],address,uint256)")),
            _amountIn,
            _amountOutMin,
            _getPath(_tokenIn, _tokenOut),
            address(this),
            block.timestamp + 300
        );

        (bool success, bytes memory result) = swapRouter.call(data);
        require(success, "Swap failed");

        return abi.decode(result, (uint256));
    }

    function _getPath(address _tokenIn, address _tokenOut) internal pure returns (address[] memory) {
        address[] memory path = new address[](2);
        path[0] = _tokenIn;
        path[1] = _tokenOut;
        return path;
    }

    function pauseJob(uint256 _jobId) external onlyJobOwner(_jobId) {
        require(jobs[_jobId].status == JobStatus.ACTIVE, "Job not active");
        jobs[_jobId].status = JobStatus.PAUSED;
        emit JobPaused(_jobId);
    }

    function resumeJob(uint256 _jobId) external onlyJobOwner(_jobId) {
        require(jobs[_jobId].status == JobStatus.PAUSED, "Job not paused");
        jobs[_jobId].status = JobStatus.ACTIVE;
        emit JobResumed(_jobId);
    }

    function cancelJob(uint256 _jobId) external onlyJobOwner(_jobId) {
        require(jobs[_jobId].status != JobStatus.CANCELLED, "Already cancelled");
        jobs[_jobId].status = JobStatus.CANCELLED;
        emit JobCancelled(_jobId);
    }

    function getJob(uint256 _jobId) external view returns (Job memory) {
        return jobs[_jobId];
    }

    function getOwnerJobs(address _owner) external view returns (uint256[] memory) {
        return ownerJobs[_owner];
    }

    function pause() external onlyAdmin {
        _pause();
    }

    function unpause() external onlyAdmin {
        _unpause();
    }
}

interface IPriceAdapter {
    function getPrice(address _tokenIn, address _tokenOut) external view returns (uint256);
    function getSwapRouter() external view returns (address);
}
