// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

interface IStakeRegistryTest {
    function registryCoordinator() external view returns (address);
    function delegation() external view returns (address);
}
