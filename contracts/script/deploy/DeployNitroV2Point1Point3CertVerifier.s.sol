// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import "forge-std/Script.sol";
import "src/interfaces/IEigenDAStructs.sol";
import "src/core/EigenDAThresholdRegistry.sol";
import "src/core/EigenDACertVerifier.sol";
import "src/core/EigenDAServiceManager.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

import "forge-std/console2.sol";

contract DeployNitroV2Point1Point3CertVerifier is Script {

    address constant INITIAL_OWNER = 0x85C2AE9B88baDf751228e307Ae9ab76B74d84f5c;

    EigenDAServiceManager constant EIGEN_DA_SERVICE_MANAGER = EigenDAServiceManager(0x870679E138bCdf293b7Ff14dD44b70FC97e12fc0);
    
    address proxyAdmin;

    address thresholdRegistryProxy;
    address thresholdRegistryImpl;
    
    address certVerifier;
    
    /// @dev This script is only to be run once, and the contracts deployed will be deprecated when EigenDA V2 is live on mainnet.
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
                        EIGEN_DA_SERVICE_MANAGER.quorumAdversaryThresholdPercentages(),
                        EIGEN_DA_SERVICE_MANAGER.quorumConfirmationThresholdPercentages(),
                        EIGEN_DA_SERVICE_MANAGER.quorumNumbersRequired(),
                        VERSIONED_BLOB_PARAMS()
                    )
                )
            )
        );

        // DEPLOY CERT VERIFIER
        certVerifier = address(
            new EigenDACertVerifier(
                IEigenDAThresholdRegistry(thresholdRegistryProxy),
                IEigenDABatchMetadataStorage(address(EIGEN_DA_SERVICE_MANAGER)),
                IEigenDASignatureVerifier(address(EIGEN_DA_SERVICE_MANAGER)),
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