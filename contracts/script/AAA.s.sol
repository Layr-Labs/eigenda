pragma solidity =0.8.12;

import {PauserRegistry} from "eigenlayer-core/contracts/permissions/PauserRegistry.sol";
import {EmptyContract} from "eigenlayer-core/test/mocks/EmptyContract.sol";
import {BLSApkRegistry} from "eigenlayer-middleware/BLSApkRegistry.sol";
import {RegistryCoordinator} from "eigenlayer-middleware/RegistryCoordinator.sol";
import {OperatorStateRetriever} from "eigenlayer-middleware/OperatorStateRetriever.sol";
import {IRegistryCoordinator} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {IndexRegistry} from "eigenlayer-middleware/IndexRegistry.sol";
import {IIndexRegistry} from "eigenlayer-middleware/interfaces/IIndexRegistry.sol";
import {StakeRegistry, IStrategy} from "eigenlayer-middleware/StakeRegistry.sol";
import {IStakeRegistry, IDelegationManager} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";
import {IServiceManager} from "eigenlayer-middleware/interfaces/IServiceManager.sol";
import {IBLSApkRegistry} from "eigenlayer-middleware/interfaces/IBLSApkRegistry.sol";
import {EigenDAServiceManager, IAVSDirectory, IRewardsCoordinator} from "../src/core/EigenDAServiceManager.sol";
import {EigenDAHasher} from "../src/libraries/EigenDAHasher.sol";
import {EigenDAThresholdRegistry} from "../src/core/EigenDAThresholdRegistry.sol";
import {EigenDABlobVerifier} from "../src/core/EigenDABlobVerifier.sol";
import {IEigenDAThresholdRegistry} from "../src/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "../src/interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDASignatureVerifier} from "../src/interfaces/IEigenDASignatureVerifier.sol";
import {IEigenDARelayRegistry} from "../src/interfaces/IEigenDARelayRegistry.sol";
import {EigenDARelayRegistry} from "../src/core/EigenDARelayRegistry.sol";
import {ISocketRegistry, SocketRegistry} from "eigenlayer-middleware/SocketRegistry.sol";
import {DeployOpenEigenLayer, ProxyAdmin, ERC20PresetFixedSupply, TransparentUpgradeableProxy, IPauserRegistry} from "./DeployOpenEigenLayer.s.sol";
import "forge-std/Test.sol";
import "forge-std/Script.sol";
import "forge-std/StdJson.sol";
import "../src/interfaces/IEigenDAStructs.sol";
import {PaymentVault} from "../src/payments/PaymentVault.sol";
import {IPaymentVault} from "../src/interfaces/IPaymentVault.sol";

//forge script script/AAA.s.sol:AAA --rpc-url $HOL --private-key $DATN -vvvv --etherscan-api-key D6ZFHU3MWZXE4Z17ICWBA1IR8A4JEPK1ZJ --verify
contract AAA is Script, Test {

    address proxyAdmin = 0x9Fd7E279f5bD692Dc04792151E14Ad814FC60eC1;

    address eigenDAServiceManager = 0x54A03db2784E3D0aCC08344D05385d0b62d4F432;
    address eigenDAServiceManagerImplementation;

    address avsDirectory = 0x141d6995556135D4997b2ff72EB443Be300353bC;
    address rewardsCoordinator = 0xb22Ef643e1E067c994019A4C19e403253C05c2B0;
    address registryCoordinator = 0x2c61EA360D6500b58E7f481541A36B443Bc858c6;
    address stakeRegistry = 0x53668EBf2e28180e38B122c641BC51Ca81088871;
    address eigenDAThresholdRegistry = 0x41AEE4A23770045e9977CC9f964d3380D6Ff9e4E;
    address eigenDARelayRegistry = 0xca1ca181fCb3c4192D320569c6eB4b5161B80328;
    address paymentVault = 0x46E024ca6e5E1172100930c28DCCcF49BE5462C9;


    function run() external {
        vm.startBroadcast();

        eigenDAServiceManagerImplementation = address(new EigenDAServiceManager(
            IAVSDirectory(avsDirectory),
            IRewardsCoordinator(rewardsCoordinator),
            IRegistryCoordinator(registryCoordinator),
            IStakeRegistry(stakeRegistry),
            IEigenDAThresholdRegistry(eigenDAThresholdRegistry),
            IEigenDARelayRegistry(eigenDARelayRegistry),
            IPaymentVault(paymentVault)
        ));

        ProxyAdmin(proxyAdmin).upgrade(
            TransparentUpgradeableProxy(payable(address(eigenDAServiceManager))),
            address(eigenDAServiceManagerImplementation)
        );

        vm.stopBroadcast();

        console.log("Deployed new EigenDAServiceManagerImplementation at address: ", eigenDAServiceManagerImplementation);
    }
}



