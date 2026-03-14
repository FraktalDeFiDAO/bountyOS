// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/utils/cryptography/EIP712.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract PermitHandler is EIP712 {
    using ECDSA for bytes32;

    bytes32 public constant PERMIT_TYPEHASH = keccak256(
        "Permit(address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)"
    );

    mapping(address => uint256) public nonces;

    constructor() EIP712("Conflux Automation Site", "1") {}

    function permit(
        address _token,
        address _owner,
        address _spender,
        uint256 _value,
        uint256 _deadline,
        uint8 _v,
        bytes32 _r,
        bytes32 _s
    ) external {
        require(block.timestamp <= _deadline, "Deadline expired");

        bytes32 structHash = keccak256(
            abi.encode(
                PERMIT_TYPEHASH,
                _owner,
                _spender,
                _value,
                nonces[_owner]++,
                _deadline
            )
        );

        bytes32 hash = _hashTypedDataV4(structHash);
        address signer = hash.recover(_v, _r, _s);

        require(signer == _owner, "Invalid signature");

        IERC20(_token).permit(_owner, _spender, _value, _deadline, _v, _r, _s);
    }

    function getNonce(address _owner) external view returns (uint256) {
        return nonces[_owner];
    }

    function getDomainSeparator() external view returns (bytes32) {
        return _domainSeparatorV4();
    }
}
