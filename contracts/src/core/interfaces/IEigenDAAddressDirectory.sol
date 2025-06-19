// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

interface IEigenDAAddressDirectory {
    event AddressSet(bytes32 indexed key, address indexed value);

    function setAddress(bytes32 key, address value) external;

    function getAddress(bytes32 key) external view returns (address);
}
