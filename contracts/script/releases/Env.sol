// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import "forge-std/Vm.sol";
import "forge-std/StdJson.sol";
import "zeus-templates/utils/ZEnvHelpers.sol";

import {TimelockController} from "@openzeppelin/contracts/governance/TimelockController.sol";
import "@openzeppelin/contracts/proxy/beacon/UpgradeableBeacon.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

/// Core contracts from eigenlayer-middleware
import {
    IPauserRegistry
} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IPauserRegistry.sol";
import {BLSApkRegistry} from "lib/eigenlayer-middleware/src/BLSApkRegistry.sol";
import {IndexRegistry} from "lib/eigenlayer-middleware/src/IndexRegistry.sol";
import {StakeRegistry} from "lib/eigenlayer-middleware/src/StakeRegistry.sol";
import {SocketRegistry} from "lib/eigenlayer-middleware/src/SocketRegistry.sol";
import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";

/// EigenDA core contracts
import {EigenDAServiceManager} from "src/core/EigenDAServiceManager.sol";
import {EigenDARegistryCoordinator} from "src/core/EigenDARegistryCoordinator.sol";
import {EigenDAThresholdRegistry} from "src/core/EigenDAThresholdRegistry.sol";
import {EigenDARelayRegistry} from "src/core/EigenDARelayRegistry.sol";
import {EigenDADisperserRegistry} from "src/core/EigenDADisperserRegistry.sol";
import {PaymentVault} from "src/core/PaymentVault.sol";
import {EigenDAAccessControl} from "src/core/EigenDAAccessControl.sol";
import {EigenDADirectory} from "src/core/EigenDADirectory.sol";

/// EigenDA periphery
import {EigenDAEjectionManager} from "src/periphery/ejection/EigenDAEjectionManager.sol";

/// Certificate verification
import {EigenDACertVerifier} from "src/integrations/cert/EigenDACertVerifier.sol";
import {EigenDACertVerifierRouter} from "src/integrations/cert/router/EigenDACertVerifierRouter.sol";

/// EigenLayer Contracts
import {
    IAVSDirectory
} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IAVSDirectory.sol";
import {
    IRewardsCoordinator
} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IRewardsCoordinator.sol";

library Env {
    using stdJson for string;
    using ZEnvHelpers for *;

    /// -----------------------------------------------------------------------
    /// Constants
    /// -----------------------------------------------------------------------

    Vm internal constant vm = Vm(address(uint160(uint256(keccak256("hevm cheat code")))));

    /// @dev Storage slot with the address of the current implementation.
    /// This is the keccak-256 hash of "eip1967.proxy.implementation" subtracted by 1, and is
    /// validated in the constructor.
    bytes32 internal constant _IMPLEMENTATION_SLOT = 0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;

    /// @dev Storage slot with the admin of the contract.
    /// This is the keccak-256 hash of "eip1967.proxy.admin" subtracted by 1, and is
    /// validated in the constructor.
    bytes32 internal constant _ADMIN_SLOT = 0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103;

    /// @dev The storage slot of the UpgradeableBeacon contract which defines the implementation for this proxy.
    /// This is bytes32(uint256(keccak256('eip1967.proxy.beacon')) - 1)) and is validated in the constructor.
    bytes32 internal constant _BEACON_SLOT = 0xa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d50;

    /// -----------------------------------------------------------------------
    /// Enums
    /// -----------------------------------------------------------------------

    // TODO: Use user defined types instead of enums, like 4 lines instead of 20.

    /// Dummy types and variables to facilitate syntax, e.g: `Env.proxy.serviceManager()`
    enum DeployedProxy {
        A
    }
    enum DeployedBeacon {
        A
    }
    enum DeployedImpl {
        A
    }
    enum DeployedInstance {
        A
    }

    DeployedProxy internal constant proxy = DeployedProxy.A;
    DeployedBeacon internal constant beacon = DeployedBeacon.A;
    DeployedImpl internal constant impl = DeployedImpl.A;
    DeployedInstance internal constant instance = DeployedInstance.A;

    /// -----------------------------------------------------------------------
    /// Environment Variables
    /// -----------------------------------------------------------------------

    function env() internal view returns (string memory) {
        return _string("ZEUS_ENV");
    }

    function envVersion() internal view returns (string memory) {
        return _string("ZEUS_ENV_VERSION");
    }

    function deployVersion() internal view returns (string memory) {
        return _string("ZEUS_DEPLOY_TO_VERSION");
    }

    function proxyAdmin() internal view returns (address) {
        return _envAddress("proxyAdmin");
    }

    function timelockController() internal view returns (TimelockController) {
        return TimelockController(payable(_envAddress("timelockController")));
    }

    /// @dev Usage: `Env.impl.owner()`
    function owner(DeployedImpl) internal view returns (address) {
        return _deployedImpl("Owner");
    }

    /// -----------------------------------------------------------------------
    /// Core EigenDA Contracts
    /// -----------------------------------------------------------------------

    /// @dev Usage: `Env.proxy.directory()`
    function directory(DeployedProxy) internal view returns (EigenDADirectory) {
        return EigenDADirectory(_deployedProxy("Directory"));
    }

    /// @dev Usage: `Env.impl.directory()`
    function directory(DeployedImpl) internal view returns (EigenDADirectory) {
        return EigenDADirectory(_deployedImpl("Directory"));
    }

    /// @dev Usage: `Env.proxy.serviceManager()`
    function serviceManager(DeployedProxy) internal view returns (EigenDAServiceManager) {
        return EigenDAServiceManager(_deployedProxy("ServiceManager"));
    }

    /// @dev Usage: `Env.impl.serviceManager()`
    function serviceManager(DeployedImpl) internal view returns (EigenDAServiceManager) {
        return EigenDAServiceManager(_deployedImpl("ServiceManager"));
    }

    /// @dev Usage: `Env.proxy.registryCoordinator()`
    function registryCoordinator(DeployedProxy) internal view returns (EigenDARegistryCoordinator) {
        return EigenDARegistryCoordinator(_deployedProxy("RegistryCoordinator"));
    }

    /// @dev Usage: `Env.impl.registryCoordinator()`
    function registryCoordinator(DeployedImpl) internal view returns (EigenDARegistryCoordinator) {
        return EigenDARegistryCoordinator(_deployedImpl("RegistryCoordinator"));
    }

    /// @dev Usage: `Env.proxy.blsApkRegistry()`
    function blsApkRegistry(DeployedProxy) internal view returns (BLSApkRegistry) {
        return BLSApkRegistry(_deployedProxy("BlsApkRegistry"));
    }

    /// @dev Usage: `Env.impl.blsApkRegistry()`
    function blsApkRegistry(DeployedImpl) internal view returns (BLSApkRegistry) {
        return BLSApkRegistry(_deployedImpl("BlsApkRegistry"));
    }

    /// @dev Usage: `Env.proxy.indexRegistry()`
    function indexRegistry(DeployedProxy) internal view returns (IndexRegistry) {
        return IndexRegistry(_deployedProxy(type(IndexRegistry).name));
    }

    /// @dev Usage: `Env.impl.indexRegistry()`
    function indexRegistry(DeployedImpl) internal view returns (IndexRegistry) {
        return IndexRegistry(_deployedImpl(type(IndexRegistry).name));
    }

    /// @dev Usage: `Env.proxy.stakeRegistry()`
    function stakeRegistry(DeployedProxy) internal view returns (StakeRegistry) {
        return StakeRegistry(_deployedProxy(type(StakeRegistry).name));
    }

    /// @dev Usage: `Env.impl.stakeRegistry()`
    function stakeRegistry(DeployedImpl) internal view returns (StakeRegistry) {
        return StakeRegistry(_deployedImpl(type(StakeRegistry).name));
    }

    /// @dev Usage: `Env.proxy.socketRegistry()`
    function socketRegistry(DeployedProxy) internal view returns (SocketRegistry) {
        return SocketRegistry(_deployedProxy(type(SocketRegistry).name));
    }

    /// @dev Usage: `Env.impl.socketRegistry()`
    function socketRegistry(DeployedImpl) internal view returns (SocketRegistry) {
        return SocketRegistry(_deployedImpl(type(SocketRegistry).name));
    }

    /// @dev Usage: `Env.proxy.thresholdRegistry()`
    function thresholdRegistry(DeployedProxy) internal view returns (EigenDAThresholdRegistry) {
        return EigenDAThresholdRegistry(_deployedProxy("ThresholdRegistry"));
    }

    /// @dev Usage: `Env.impl.thresholdRegistry()`
    function thresholdRegistry(DeployedImpl) internal view returns (EigenDAThresholdRegistry) {
        return EigenDAThresholdRegistry(_deployedImpl("ThresholdRegistry"));
    }

    /// @dev Usage: `Env.proxy.relayRegistry()`
    function relayRegistry(DeployedProxy) internal view returns (EigenDARelayRegistry) {
        return EigenDARelayRegistry(_deployedProxy("RelayRegistry"));
    }

    /// @dev Usage: `Env.impl.relayRegistry()`
    function relayRegistry(DeployedImpl) internal view returns (EigenDARelayRegistry) {
        return EigenDARelayRegistry(_deployedImpl("RelayRegistry"));
    }

    /// @dev Usage: `Env.proxy.disperserRegistry()`
    function disperserRegistry(DeployedProxy) internal view returns (EigenDADisperserRegistry) {
        return EigenDADisperserRegistry(_deployedProxy("DisperserRegistry"));
    }

    /// @dev Usage: `Env.impl.disperserRegistry()`
    function disperserRegistry(DeployedImpl) internal view returns (EigenDADisperserRegistry) {
        return EigenDADisperserRegistry(_deployedImpl("DisperserRegistry"));
    }

    /// @dev Usage: `Env.proxy.paymentVault()`
    function paymentVault(DeployedProxy) internal view returns (PaymentVault) {
        return PaymentVault(payable(_deployedProxy(type(PaymentVault).name)));
    }

    /// @dev Usage: `Env.impl.paymentVault()`
    function paymentVault(DeployedImpl) internal view returns (PaymentVault) {
        return PaymentVault(payable(_deployedImpl(type(PaymentVault).name)));
    }

    /// @dev Usage: `Env.impl.accessControl()`
    function accessControl(DeployedImpl) internal view returns (EigenDAAccessControl) {
        return EigenDAAccessControl(_deployedImpl("AccessControl"));
    }

    /// @dev Usage: `Env.impl.operatorStateRetriever()`
    function operatorStateRetriever(DeployedImpl) internal view returns (OperatorStateRetriever) {
        return OperatorStateRetriever(_deployedImpl(type(OperatorStateRetriever).name));
    }

    /// @dev Usage: `Env.impl.pauserRegistry()`
    function pauserRegistry(DeployedImpl) internal view returns (IPauserRegistry) {
        return IPauserRegistry(_deployedImpl("PauserRegistry"));
    }

    /// -----------------------------------------------------------------------
    /// Periphery Contracts
    /// -----------------------------------------------------------------------

    /// @dev Usage: `Env.proxy.ejectionManager()`
    function ejectionManager(DeployedProxy) internal view returns (EigenDAEjectionManager) {
        return EigenDAEjectionManager(_deployedProxy("EjectionManager"));
    }

    /// @dev Usage: `Env.impl.ejectionManager()`
    function ejectionManager(DeployedImpl) internal view returns (EigenDAEjectionManager) {
        return EigenDAEjectionManager(_deployedImpl("EjectionManager"));
    }

    /// -----------------------------------------------------------------------
    /// Cert Verification Contracts
    /// -----------------------------------------------------------------------

    /// @dev Usage: `Env.impl.certVerifier()`
    function certVerifier(DeployedImpl) internal view returns (EigenDACertVerifier) {
        return EigenDACertVerifier(_deployedImpl("CertVerifier"));
    }

    /// @dev Usage: `Env.proxy.certVerifierRouter()`
    function certVerifierRouter(DeployedProxy) internal view returns (EigenDACertVerifierRouter) {
        return EigenDACertVerifierRouter(_deployedProxy("CertVerifierRouter"));
    }

    /// @dev Usage: `Env.impl.certVerifierRouter()`
    function certVerifierRouter(DeployedImpl) internal view returns (EigenDACertVerifierRouter) {
        return EigenDACertVerifierRouter(_deployedImpl("CertVerifierRouter"));
    }

    /// -----------------------------------------------------------------------
    /// EigenLayer Contracts
    /// -----------------------------------------------------------------------

    /// @dev Usage: `Env.proxy.avsDirectory()`
    function avsDirectory(DeployedProxy) internal view returns (IAVSDirectory) {
        return IAVSDirectory(_deployedProxy("AVSDirectory"));
    }

    /// @dev Usage: `Env.proxy.rewardsCoordinator()`
    function rewardsCoordinator(DeployedProxy) internal view returns (IRewardsCoordinator) {
        return IRewardsCoordinator(_deployedProxy("RewardsCoordinator"));
    }

    /// -----------------------------------------------------------------------
    /// Private Zeus Helpers
    /// -----------------------------------------------------------------------

    function _deployedInstance(string memory name, uint256 idx) private view returns (address) {
        return ZEnvHelpers.state().deployedInstance(name, idx);
    }

    function _deployedInstanceCount(string memory name) private view returns (uint256) {
        return ZEnvHelpers.state().deployedInstanceCount(name);
    }

    function _deployedProxy(string memory name) private view returns (address) {
        return ZEnvHelpers.state().deployedProxy(name);
    }

    function _deployedBeacon(string memory name) private view returns (address) {
        return ZEnvHelpers.state().deployedBeacon(name);
    }

    function _deployedImpl(string memory name) private view returns (address) {
        return ZEnvHelpers.state().deployedImpl(name);
    }

    function _envAddress(string memory key) private view returns (address) {
        return ZEnvHelpers.state().envAddress(key);
    }

    function _envU256(string memory key) private view returns (uint256) {
        return ZEnvHelpers.state().envU256(key);
    }

    function _envU64(string memory key) private view returns (uint64) {
        return ZEnvHelpers.state().envU64(key);
    }

    function _envU32(string memory key) private view returns (uint32) {
        return ZEnvHelpers.state().envU32(key);
    }

    function _envU16(string memory key) private view returns (uint16) {
        return ZEnvHelpers.state().envU16(key);
    }

    function _envBool(string memory key) private view returns (bool) {
        return ZEnvHelpers.state().envBool(key);
    }

    function _string(string memory key) private view returns (string memory) {
        return ZEnvHelpers.state().envString(key);
    }

    /// -----------------------------------------------------------------------
    /// ERC-1967 Storage Accessors
    /// -----------------------------------------------------------------------

    /// @dev Query and return the implementation address of the proxy.
    function _getProxyImpl(address _proxy) internal view returns (address) {
        return address(uint160(uint256(vm.load(_proxy, _IMPLEMENTATION_SLOT))));
    }

    /// @dev Query and return the admin address of the proxy.
    function _getProxyAdmin(address _proxy) internal view returns (address) {
        return address(uint160(uint256(vm.load(_proxy, _ADMIN_SLOT))));
    }

    /// @dev Query and return the beacon address of the proxy.
    function _getBeacon(address _proxy) internal view returns (address) {
        return address(uint160(uint256(vm.load(_proxy, _BEACON_SLOT))));
    }

}
