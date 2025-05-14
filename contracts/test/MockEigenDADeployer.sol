// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "../lib/eigenlayer-middleware/test/utils/BLSMockAVSDeployer.sol";
import {EigenDAServiceManager, IRewardsCoordinator} from "src/core/EigenDAServiceManager.sol";
import {EigenDAServiceManager} from "src/core/EigenDAServiceManager.sol";
import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";
import {EigenDATypesV2 as DATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";
import {EigenDACertVerificationV1Lib} from "src/periphery/cert/legacy/v1/EigenDACertVerificationV1Lib.sol";
import {IEigenDAServiceManager} from "src/core/interfaces/IEigenDAServiceManager.sol";
import {EigenDACertVerifierV2} from "src/periphery/cert/legacy/v2/EigenDACertVerifierV2.sol";
import {EigenDAThresholdRegistry, IEigenDAThresholdRegistry} from "src/core/EigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "src/core/interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDASignatureVerifier} from "src/core/interfaces/IEigenDASignatureVerifier.sol";
import {IRegistryCoordinator} from "../lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {IEigenDARelayRegistry} from "src/core/interfaces/IEigenDARelayRegistry.sol";
import {EigenDARelayRegistry} from "src/core/EigenDARelayRegistry.sol";
import {IPaymentVault} from "src/core/interfaces/IPaymentVault.sol";
import {PaymentVault} from "src/core/PaymentVault.sol";
import {IEigenDADisperserRegistry} from "src/core/interfaces/IEigenDADisperserRegistry.sol";
import {EigenDADisperserRegistry} from "src/core/EigenDADisperserRegistry.sol";
import "forge-std/StdStorage.sol";

contract MockEigenDADeployer is BLSMockAVSDeployer {
    using stdStorage for StdStorage;
    using BN254 for BN254.G1Point;

    address confirmer = address(uint160(uint256(keccak256(abi.encodePacked("confirmer")))));
    address notConfirmer = address(uint160(uint256(keccak256(abi.encodePacked("notConfirmer")))));
    address rewardsInitiator = address(uint160(uint256(keccak256(abi.encodePacked("rewardsInitiator")))));

    EigenDAServiceManager eigenDAServiceManager;
    EigenDAServiceManager eigenDAServiceManagerImplementation;
    EigenDARelayRegistry eigenDARelayRegistry;
    EigenDARelayRegistry eigenDARelayRegistryImplementation;
    EigenDAThresholdRegistry eigenDAThresholdRegistry;
    EigenDAThresholdRegistry eigenDAThresholdRegistryImplementation;
    EigenDADisperserRegistry eigenDADisperserRegistry;
    EigenDADisperserRegistry eigenDADisperserRegistryImplementation;
    PaymentVault paymentVault;
    PaymentVault paymentVaultImplementation;
    EigenDACertVerifier eigenDACertVerifier;

    ERC20 mockToken;

    bytes quorumAdversaryThresholdPercentages = hex"212121";
    bytes quorumConfirmationThresholdPercentages = hex"373737";
    bytes quorumNumbersRequired = hex"0001";
    DATypesV1.SecurityThresholds defaultSecurityThresholds = DATypesV1.SecurityThresholds(55, 33);

    uint32 defaultReferenceBlockNumber = 100;
    uint32 defaultConfirmationBlockNumber = 1000;
    uint32 defaultBatchId = 0;

    uint64 minNumSymbols = 1;
    uint64 pricePerSymbol = 3;
    uint64 priceUpdateCooldown = 6 days;
    uint64 globalSymbolsPerPeriod = 2;
    uint64 reservationPeriodInterval = 4;
    uint64 globalRatePeriodInterval = 5;

    mapping(uint8 => bool) public quorumNumbersUsed;

    function _deployDA() public {
        _setUpBLSMockAVSDeployer();

        eigenDAServiceManager = EigenDAServiceManager(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(proxyAdmin), ""))
        );

        eigenDAThresholdRegistry = EigenDAThresholdRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(proxyAdmin), ""))
        );

        eigenDARelayRegistry = EigenDARelayRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(proxyAdmin), ""))
        );

        paymentVault = PaymentVault(
            payable(address(new TransparentUpgradeableProxy(address(emptyContract), address(proxyAdmin), "")))
        );

        eigenDADisperserRegistry = EigenDADisperserRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(proxyAdmin), ""))
        );

        eigenDAServiceManagerImplementation = new EigenDAServiceManager(
            avsDirectory,
            rewardsCoordinator,
            registryCoordinator,
            stakeRegistry,
            eigenDAThresholdRegistry,
            eigenDARelayRegistry,
            paymentVault,
            eigenDADisperserRegistry
        );

        address[] memory confirmers = new address[](1);
        confirmers[0] = confirmer;

        cheats.prank(proxyAdminOwner);
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenDAServiceManager))),
            address(eigenDAServiceManagerImplementation),
            abi.encodeWithSelector(
                EigenDAServiceManager.initialize.selector,
                pauserRegistry,
                0,
                registryCoordinatorOwner,
                confirmers,
                registryCoordinatorOwner
            )
        );

        eigenDAThresholdRegistryImplementation = new EigenDAThresholdRegistry();

        DATypesV1.VersionedBlobParams[] memory versionedBlobParams = new DATypesV1.VersionedBlobParams[](1);
        versionedBlobParams[0] = DATypesV1.VersionedBlobParams({maxNumOperators: 3537, numChunks: 8192, codingRate: 8});

        cheats.prank(proxyAdminOwner);
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenDAThresholdRegistry))),
            address(eigenDAThresholdRegistryImplementation),
            abi.encodeWithSelector(
                EigenDAThresholdRegistry.initialize.selector,
                registryCoordinatorOwner,
                quorumAdversaryThresholdPercentages,
                quorumConfirmationThresholdPercentages,
                quorumNumbersRequired,
                versionedBlobParams
            )
        );

        eigenDARelayRegistryImplementation = new EigenDARelayRegistry();

        cheats.prank(proxyAdminOwner);
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenDARelayRegistry))),
            address(eigenDARelayRegistryImplementation),
            abi.encodeWithSelector(EigenDARelayRegistry.initialize.selector, registryCoordinatorOwner)
        );

        eigenDADisperserRegistryImplementation = new EigenDADisperserRegistry();

        cheats.prank(proxyAdminOwner);
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenDADisperserRegistry))),
            address(eigenDADisperserRegistryImplementation),
            abi.encodeWithSelector(EigenDADisperserRegistry.initialize.selector, registryCoordinatorOwner)
        );

        paymentVaultImplementation = PaymentVault(payable(address(new PaymentVault())));

        paymentVault = PaymentVault(
            payable(
                address(
                    new TransparentUpgradeableProxy(
                        address(paymentVaultImplementation),
                        address(proxyAdmin),
                        abi.encodeWithSelector(
                            PaymentVault.initialize.selector,
                            registryCoordinatorOwner,
                            minNumSymbols,
                            pricePerSymbol,
                            priceUpdateCooldown,
                            globalSymbolsPerPeriod,
                            reservationPeriodInterval,
                            globalRatePeriodInterval
                        )
                    )
                )
            )
        );

        mockToken = new ERC20("Mock Token", "MOCK");

        eigenDACertVerifier = new EigenDACertVerifier(
            IEigenDAThresholdRegistry(address(eigenDAThresholdRegistry)),
            IEigenDASignatureVerifier(address(eigenDAServiceManager)),
            defaultSecurityThresholds,
            quorumNumbersRequired
        );
    }

    function _getHeaderandNonSigners(uint256 _nonSigners, uint256 _pseudoRandomNumber, uint8 _threshold)
        internal
        returns (DATypesV1.BatchHeader memory, BLSSignatureChecker.NonSignerStakesAndSignature memory)
    {
        // register a bunch of operators
        uint256 quorumBitmap = 1;
        bytes memory quorumNumbers = BitmapUtils.bitmapToBytesArray(quorumBitmap);

        // 0 nonSigners
        (
            uint32 referenceBlockNumber,
            BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature
        ) = _registerSignatoriesAndGetNonSignerStakeAndSignatureRandom(_pseudoRandomNumber, _nonSigners, quorumBitmap);

        // get a random batch header
        DATypesV1.BatchHeader memory batchHeader =
            _getRandomBatchHeader(_pseudoRandomNumber, quorumNumbers, referenceBlockNumber, _threshold);

        // set batch specific signature
        bytes32 reducedBatchHeaderHash = keccak256(
            abi.encode(
                DATypesV1.ReducedBatchHeader({
                    blobHeadersRoot: batchHeader.blobHeadersRoot,
                    referenceBlockNumber: batchHeader.referenceBlockNumber
                })
            )
        );
        nonSignerStakesAndSignature.sigma = BN254.hashToG1(reducedBatchHeaderHash).scalar_mul(aggSignerPrivKey);

        return (batchHeader, nonSignerStakesAndSignature);
    }

    function _getRandomBatchHeader(
        uint256 pseudoRandomNumber,
        bytes memory quorumNumbers,
        uint32 referenceBlockNumber,
        uint8 threshold
    ) internal pure returns (DATypesV1.BatchHeader memory) {
        DATypesV1.BatchHeader memory batchHeader;
        batchHeader.blobHeadersRoot = keccak256(abi.encodePacked("blobHeadersRoot", pseudoRandomNumber));
        batchHeader.quorumNumbers = quorumNumbers;
        batchHeader.signedStakeForQuorums = new bytes(quorumNumbers.length);
        for (uint256 i = 0; i < quorumNumbers.length; i++) {
            batchHeader.signedStakeForQuorums[i] = bytes1(threshold);
        }
        batchHeader.referenceBlockNumber = referenceBlockNumber;
        return batchHeader;
    }

    function _generateRandomBlobHeader(uint256 pseudoRandomNumber, uint256 numQuorumsBlobParams)
        internal
        returns (DATypesV1.BlobHeader memory)
    {
        if (pseudoRandomNumber == 0) {
            pseudoRandomNumber = 1;
        }

        DATypesV1.BlobHeader memory blobHeader;
        blobHeader.commitment.X =
            uint256(keccak256(abi.encodePacked(pseudoRandomNumber, "blobHeader.commitment.X"))) % BN254.FP_MODULUS;
        blobHeader.commitment.Y =
            uint256(keccak256(abi.encodePacked(pseudoRandomNumber, "blobHeader.commitment.Y"))) % BN254.FP_MODULUS;

        blobHeader.dataLength =
            uint32(uint256(keccak256(abi.encodePacked(pseudoRandomNumber, "blobHeader.dataLength"))));

        blobHeader.quorumBlobParams = new DATypesV1.QuorumBlobParam[](numQuorumsBlobParams);
        blobHeader.dataLength =
            uint32(uint256(keccak256(abi.encodePacked(pseudoRandomNumber, "blobHeader.dataLength"))));
        for (uint256 i = 0; i < numQuorumsBlobParams; i++) {
            if (i < 2) {
                blobHeader.quorumBlobParams[i].quorumNumber = uint8(i);
            } else {
                blobHeader.quorumBlobParams[i].quorumNumber = uint8(
                    uint256(
                        keccak256(
                            abi.encodePacked(pseudoRandomNumber, "blobHeader.quorumBlobParams[i].quorumNumber", i)
                        )
                    )
                ) % 192;

                // make sure it isn't already used
                while (quorumNumbersUsed[blobHeader.quorumBlobParams[i].quorumNumber]) {
                    blobHeader.quorumBlobParams[i].quorumNumber =
                        uint8(uint256(blobHeader.quorumBlobParams[i].quorumNumber) + 1) % 192;
                }
                quorumNumbersUsed[blobHeader.quorumBlobParams[i].quorumNumber] = true;
            }

            blobHeader.quorumBlobParams[i].adversaryThresholdPercentage = eigenDAThresholdRegistry
                .getQuorumAdversaryThresholdPercentage(blobHeader.quorumBlobParams[i].quorumNumber);
            blobHeader.quorumBlobParams[i].chunkLength = uint32(
                uint256(
                    keccak256(abi.encodePacked(pseudoRandomNumber, "blobHeader.quorumBlobParams[i].chunkLength", i))
                )
            );
            blobHeader.quorumBlobParams[i].confirmationThresholdPercentage = eigenDAThresholdRegistry
                .getQuorumConfirmationThresholdPercentage(blobHeader.quorumBlobParams[i].quorumNumber);
        }
        // mark all quorum numbers as unused
        for (uint256 i = 0; i < numQuorumsBlobParams; i++) {
            quorumNumbersUsed[blobHeader.quorumBlobParams[i].quorumNumber] = false;
        }

        return blobHeader;
    }
}
