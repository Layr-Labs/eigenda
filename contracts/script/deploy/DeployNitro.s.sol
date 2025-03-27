// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import "forge-std/Script.sol";
import "src/interfaces/IEigenDAStructs.sol";
import "src/core/EigenDAThresholdRegistry.sol";
import "src/core/EigenDACertVerifier.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

import "forge-std/console2.sol";

contract DeployNitroCertVerifier is Script {

    address constant INITIAL_OWNER = 0x85C2AE9B88baDf751228e307Ae9ab76B74d84f5c;
    address constant EIGEN_DA_SERVICE_MANAGER = 0x870679E138bCdf293b7Ff14dD44b70FC97e12fc0;
    bytes constant QUORUM_ADVERSARY_THRESHOLD_PERCENTAGE = hex"212121";
    bytes constant QUORUM_CONFIRMATION_THRESHOLD_PERCENTAGE = hex"373737";
    bytes constant QUORUM_NUMBERS_REQUIRED = hex"0001";

    address proxyAdmin;

    address thresholdRegistryProxy;
    address thresholdRegistryImpl;
    
    address certVerifier;
    
    function run() external {
        vm.startBroadcast(INITIAL_OWNER);
        proxyAdmin = address(new ProxyAdmin());
        // DEPLOY THRESHOLD REGISTRY
        thresholdRegistryImpl = address(new EigenDAThresholdRegistry());
        thresholdRegistryProxy = address(
            new TransparentUpgradeableProxy(
                thresholdRegistryImpl,
                proxyAdmin,
                abi.encodeCall(
                    EigenDAThresholdRegistry.initialize,
                    (
                        INITIAL_OWNER,
                        QUORUM_ADVERSARY_THRESHOLD_PERCENTAGE,
                        QUORUM_CONFIRMATION_THRESHOLD_PERCENTAGE,
                        QUORUM_NUMBERS_REQUIRED,
                        VERSIONED_BLOB_PARAMS()
                    )
                )
            )
        );

        // DEPLOY CERT VERIFIER
        certVerifier = address(
            new EigenDACertVerifier(
                IEigenDAThresholdRegistry(thresholdRegistryProxy),
                IEigenDABatchMetadataStorage(EIGEN_DA_SERVICE_MANAGER),
                IEigenDASignatureVerifier(EIGEN_DA_SERVICE_MANAGER),
                IEigenDARelayRegistry(address(0)), // UNUSED IN V1
                OperatorStateRetriever(address(0)), // UNUSED IN V1
                IRegistryCoordinator(address(0)), // UNUSED IN V1
                SecurityThresholds({
                    confirmationThreshold: 0,
                    adversaryThreshold: 0
                }), // UNUSED IN V1
                hex"" // UNUSED IN V1
            )
        );
        vm.stopBroadcast();

        console2.log("PROXY ADMIN: ", proxyAdmin);
        console2.log("THRESHOLD REGISTRY IMPL: ", proxyAdmin);
        console2.log("THRESHOLD REGISTRY: ", proxyAdmin);
        console2.log("CERT VERIFIER: ", proxyAdmin);
    }

    // V1 verification does not use this, so can be empty
    function VERSIONED_BLOB_PARAMS() internal pure returns (VersionedBlobParams[] memory params) {
        return params;

    }
}