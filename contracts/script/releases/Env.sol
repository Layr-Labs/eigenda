// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import "forge-std/Vm.sol";
import "forge-std/StdJson.sol";

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

library Env {
    using stdJson for string;

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

    /// env
    function env() internal view returns (string memory) {
        return _string("ZEUS_ENV");
    }

    function envVersion() internal view returns (string memory) {
        return _string("ZEUS_ENV_VERSION");
    }

    function deployVersion() internal view returns (string memory) {
        return _string("ZEUS_DEPLOY_TO_VERSION");
    }

    function executorMultisig() internal view returns (address) {
        return _envAddress("executorMultisig");
    }

    function opsMultisig() internal view returns (address) {
        return _envAddress("operationsMultisig");
    }

    function communityMultisig() internal view returns (address) {
        return _envAddress("communityMultisig");
    }

    function pauserMultisig() internal view returns (address) {
        return _envAddress("pauserMultisig");
    }

    function proxyAdmin() internal view returns (address) {
        return _envAddress("proxyAdmin");
    }

    function timelockController() internal view returns (TimelockController) {
        return TimelockController(payable(_envAddress("timelockController")));
    }

    /// Core EigenDA Contracts

    /// Directory
    function directory(DeployedProxy) internal view returns (EigenDADirectory) {
        return EigenDADirectory(_deployedProxy(type(EigenDADirectory).name));
    }

    function directory(DeployedImpl) internal view returns (EigenDADirectory) {
        return EigenDADirectory(_deployedImpl(type(EigenDADirectory).name));
    }

    /// Service Manager
    function serviceManager(DeployedProxy) internal view returns (EigenDAServiceManager) {
        return EigenDAServiceManager(_deployedProxy(type(EigenDAServiceManager).name));
    }

    function serviceManager(DeployedImpl) internal view returns (EigenDAServiceManager) {
        return EigenDAServiceManager(_deployedImpl(type(EigenDAServiceManager).name));
    }

    /// Registry Coordinator
    function registryCoordinator(DeployedProxy) internal view returns (EigenDARegistryCoordinator) {
        return EigenDARegistryCoordinator(_deployedProxy(type(EigenDARegistryCoordinator).name));
    }

    function registryCoordinator(DeployedImpl) internal view returns (EigenDARegistryCoordinator) {
        return EigenDARegistryCoordinator(_deployedImpl(type(EigenDARegistryCoordinator).name));
    }

    /// BLS APK Registry
    function blsApkRegistry(DeployedProxy) internal view returns (BLSApkRegistry) {
        return BLSApkRegistry(_deployedProxy(type(BLSApkRegistry).name));
    }

    function blsApkRegistry(DeployedImpl) internal view returns (BLSApkRegistry) {
        return BLSApkRegistry(_deployedImpl(type(BLSApkRegistry).name));
    }

    /// Index Registry
    function indexRegistry(DeployedProxy) internal view returns (IndexRegistry) {
        return IndexRegistry(_deployedProxy(type(IndexRegistry).name));
    }

    function indexRegistry(DeployedImpl) internal view returns (IndexRegistry) {
        return IndexRegistry(_deployedImpl(type(IndexRegistry).name));
    }

    /// Stake Registry
    function stakeRegistry(DeployedProxy) internal view returns (StakeRegistry) {
        return StakeRegistry(_deployedProxy(type(StakeRegistry).name));
    }

    function stakeRegistry(DeployedImpl) internal view returns (StakeRegistry) {
        return StakeRegistry(_deployedImpl(type(StakeRegistry).name));
    }

    /// Socket Registry
    function socketRegistry(DeployedProxy) internal view returns (SocketRegistry) {
        return SocketRegistry(_deployedProxy(type(SocketRegistry).name));
    }

    function socketRegistry(DeployedImpl) internal view returns (SocketRegistry) {
        return SocketRegistry(_deployedImpl(type(SocketRegistry).name));
    }

    /// Threshold Registry
    function thresholdRegistry(DeployedProxy) internal view returns (EigenDAThresholdRegistry) {
        return EigenDAThresholdRegistry(_deployedProxy(type(EigenDAThresholdRegistry).name));
    }

    function thresholdRegistry(DeployedImpl) internal view returns (EigenDAThresholdRegistry) {
        return EigenDAThresholdRegistry(_deployedImpl(type(EigenDAThresholdRegistry).name));
    }

    /// Relay Registry
    function relayRegistry(DeployedProxy) internal view returns (EigenDARelayRegistry) {
        return EigenDARelayRegistry(_deployedProxy(type(EigenDARelayRegistry).name));
    }

    function relayRegistry(DeployedImpl) internal view returns (EigenDARelayRegistry) {
        return EigenDARelayRegistry(_deployedImpl(type(EigenDARelayRegistry).name));
    }

    /// Disperser Registry
    function disperserRegistry(DeployedProxy) internal view returns (EigenDADisperserRegistry) {
        return EigenDADisperserRegistry(_deployedProxy(type(EigenDADisperserRegistry).name));
    }

    function disperserRegistry(DeployedImpl) internal view returns (EigenDADisperserRegistry) {
        return EigenDADisperserRegistry(_deployedImpl(type(EigenDADisperserRegistry).name));
    }

    /// Payment Vault
    function paymentVault(DeployedProxy) internal view returns (PaymentVault) {
        return PaymentVault(payable(_deployedProxy(type(PaymentVault).name)));
    }

    function paymentVault(DeployedImpl) internal view returns (PaymentVault) {
        return PaymentVault(payable(_deployedImpl(type(PaymentVault).name)));
    }

    /// Access Control
    function accessControl(DeployedImpl) internal view returns (EigenDAAccessControl) {
        return EigenDAAccessControl(_deployedImpl(type(EigenDAAccessControl).name));
    }

    /// Operator State Retriever
    function operatorStateRetriever(DeployedImpl) internal view returns (OperatorStateRetriever) {
        return OperatorStateRetriever(_deployedImpl(type(OperatorStateRetriever).name));
    }

    /// Pauser Registry
    function pauserRegistry(DeployedImpl) internal view returns (IPauserRegistry) {
        return IPauserRegistry(_deployedImpl("PauserRegistry"));
    }

    /// Periphery Contracts

    /// Ejection Manager
    function ejectionManager(DeployedProxy) internal view returns (EigenDAEjectionManager) {
        return EigenDAEjectionManager(_deployedProxy(type(EigenDAEjectionManager).name));
    }

    function ejectionManager(DeployedImpl) internal view returns (EigenDAEjectionManager) {
        return EigenDAEjectionManager(_deployedImpl(type(EigenDAEjectionManager).name));
    }

    /// Certificate Verification Contracts

    /// Certificate Verifier
    function certVerifier(DeployedImpl) internal view returns (EigenDACertVerifier) {
        return EigenDACertVerifier(_deployedImpl(type(EigenDACertVerifier).name));
    }

    /// Certificate Verifier Router
    function certVerifierRouter(DeployedProxy) internal view returns (EigenDACertVerifierRouter) {
        return EigenDACertVerifierRouter(_deployedProxy(type(EigenDACertVerifierRouter).name));
    }

    function certVerifierRouter(DeployedImpl) internal view returns (EigenDACertVerifierRouter) {
        return EigenDACertVerifierRouter(_deployedImpl(type(EigenDACertVerifierRouter).name));
    }

    /// Legacy Certificate Verifiers
    function certVerifierLegacyV1(DeployedImpl) internal view returns (address) {
        return _deployedImpl("EigenDACertVerifierLegacyV1");
    }

    function certVerifierLegacyV2(DeployedImpl) internal view returns (address) {
        return _deployedImpl("EigenDACertVerifierLegacyV2");
    }

    /// Helpers

    address internal constant VM_ADDRESS = address(uint160(uint256(keccak256("hevm cheat code"))));
    Vm internal constant vm = Vm(VM_ADDRESS);

    /// @dev Returns the path to the deployment state file based on ZEUS_ENV and ZEUS_ENV_VERSION
    function _getStatePath() private view returns (string memory) {
        string memory envName = _string("ZEUS_ENV");
        string memory version = _string("ZEUS_ENV_VERSION");
        return string.concat("script/releases/state/", envName, "/", version, "/deployed.json");
    }

    /// @dev Load and parse the deployment state JSON
    function _getDeploymentState() private view returns (string memory) {
        try vm.readFile(_getStatePath()) returns (string memory json) {
            return json;
        } catch {
            return "";
        }
    }

    function _deployedInstance(string memory name, uint256 idx) private view returns (address) {
        string memory json = _getDeploymentState();
        if (bytes(json).length == 0) return address(0);

        string memory key = string.concat(".instances.", name, "[", vm.toString(idx), "]");
        try vm.parseJsonAddress(json, key) returns (address addr) {
            return addr;
        } catch {
            return address(0);
        }
    }

    function _deployedInstanceCount(
        string memory /*name*/
    )
        private
        view
        returns (uint256)
    {
        string memory json = _getDeploymentState();
        if (bytes(json).length == 0) return 0;

        // Try to parse the array and get its length
        // This is a simplified approach - you may need to adjust based on your JSON structure
        return 0;
    }

    function _deployedProxy(string memory name) private view returns (address) {
        string memory json = _getDeploymentState();
        if (bytes(json).length == 0) return address(0);

        string memory key = string.concat(".proxies.", name);
        try vm.parseJsonAddress(json, key) returns (address addr) {
            return addr;
        } catch {
            return address(0);
        }
    }

    function _deployedBeacon(string memory name) private view returns (address) {
        string memory json = _getDeploymentState();
        if (bytes(json).length == 0) return address(0);

        string memory key = string.concat(".beacons.", name);
        try vm.parseJsonAddress(json, key) returns (address addr) {
            return addr;
        } catch {
            return address(0);
        }
    }

    function _deployedImpl(string memory name) private view returns (address) {
        string memory json = _getDeploymentState();
        if (bytes(json).length == 0) return address(0);

        string memory key = string.concat(".implementations.", name);
        try vm.parseJsonAddress(json, key) returns (address addr) {
            return addr;
        } catch {
            return address(0);
        }
    }

    function _envAddress(string memory key) private view returns (address) {
        try vm.envAddress(key) returns (address addr) {
            return addr;
        } catch {
            return address(0);
        }
    }

    function _envU256(string memory key) private view returns (uint256) {
        try vm.envUint(key) returns (uint256 value) {
            return value;
        } catch {
            return 0;
        }
    }

    function _envU64(string memory key) private view returns (uint64) {
        try vm.envUint(key) returns (uint256 value) {
            return uint64(value);
        } catch {
            return 0;
        }
    }

    function _envU32(string memory key) private view returns (uint32) {
        try vm.envUint(key) returns (uint256 value) {
            return uint32(value);
        } catch {
            return 0;
        }
    }

    function _envU16(string memory key) private view returns (uint16) {
        try vm.envUint(key) returns (uint256 value) {
            return uint16(value);
        } catch {
            return 0;
        }
    }

    function _envBool(string memory key) private view returns (bool) {
        try vm.envBool(key) returns (bool value) {
            return value;
        } catch {
            return false;
        }
    }

    function _string(string memory key) private view returns (string memory) {
        try vm.envString(key) returns (string memory value) {
            return value;
        } catch {
            return "";
        }
    }

    /// Test Helpers

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

    function _strEq(string memory a, string memory b) internal pure returns (bool) {
        return keccak256(bytes(a)) == keccak256(bytes(b));
    }
}

