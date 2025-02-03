// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import "../MockEigenDADeployer.sol";

contract EigenDAThresholdRegistryUnit is MockEigenDADeployer {

    event VersionedBlobParamsAdded(uint16 indexed version, VersionedBlobParams versionedBlobParams);
    event QuorumAdversaryThresholdPercentagesUpdated(bytes previousQuorumAdversaryThresholdPercentages, bytes newQuorumAdversaryThresholdPercentages);
    event QuorumConfirmationThresholdPercentagesUpdated(bytes previousQuorumConfirmationThresholdPercentages, bytes newQuorumConfirmationThresholdPercentages);
    event QuorumNumbersRequiredUpdated(bytes previousQuorumNumbersRequired, bytes newQuorumNumbersRequired);
    event DefaultSecurityThresholdsV2Updated(SecurityThresholds previousDefaultSecurityThresholdsV2, SecurityThresholds newDefaultSecurityThresholdsV2);

    function setUp() virtual public {
        _deployDA();
    }

    function test_initalize() public {
        VersionedBlobParams memory _versionedBlobParams = VersionedBlobParams({
            maxNumOperators: 3537,
            numChunks: 8192,
            codingRate: 8
        });

        assertEq(eigenDAThresholdRegistry.owner(), registryCoordinatorOwner);
        assertEq(keccak256(abi.encode(eigenDAThresholdRegistry.quorumAdversaryThresholdPercentages())), keccak256(abi.encode(quorumAdversaryThresholdPercentages)));
        assertEq(keccak256(abi.encode(eigenDAThresholdRegistry.quorumConfirmationThresholdPercentages())), keccak256(abi.encode(quorumConfirmationThresholdPercentages)));
        assertEq(keccak256(abi.encode(eigenDAThresholdRegistry.quorumNumbersRequired())), keccak256(abi.encode(quorumNumbersRequired)));
        (uint32 maxNumOperators, uint32 numChunks, uint8 codingRate) = eigenDAThresholdRegistry.versionedBlobParams(0); 
        assertEq(maxNumOperators, _versionedBlobParams.maxNumOperators);
        assertEq(numChunks, _versionedBlobParams.numChunks);
        assertEq(codingRate, _versionedBlobParams.codingRate);

        VersionedBlobParams[] memory versionedBlobParams = new VersionedBlobParams[](1);
        versionedBlobParams[0] = _versionedBlobParams;
        vm.expectRevert("Initializable: contract is already initialized");
        eigenDAThresholdRegistry.initialize(
            registryCoordinatorOwner,
            quorumAdversaryThresholdPercentages,
            quorumConfirmationThresholdPercentages,
            quorumNumbersRequired,
            versionedBlobParams
        );
    }

    function test_addVersionedBlobParams() public {
        VersionedBlobParams memory _versionedBlobParams = VersionedBlobParams({
            maxNumOperators: 999,
            numChunks: 999,
            codingRate: 9
        });
        vm.expectEmit(address(eigenDAThresholdRegistry));
        emit VersionedBlobParamsAdded(1, _versionedBlobParams);
        vm.prank(registryCoordinatorOwner);
        uint16 version = eigenDAThresholdRegistry.addVersionedBlobParams(_versionedBlobParams);
        assertEq(version, 1);
        (uint32 maxNumOperators, uint32 numChunks, uint8 codingRate) = eigenDAThresholdRegistry.versionedBlobParams(version); 
        assertEq(maxNumOperators, _versionedBlobParams.maxNumOperators);
        assertEq(numChunks, _versionedBlobParams.numChunks);
        assertEq(codingRate, _versionedBlobParams.codingRate);
    }

    function test_revert_onlyOwner() public {
        vm.expectRevert("Ownable: caller is not the owner");
        eigenDAThresholdRegistry.addVersionedBlobParams(VersionedBlobParams({
            maxNumOperators: 999,
            numChunks: 999,
            codingRate: 9
        }));
    }

    function test_getQuorumAdversaryThresholdPercentage() public {
        uint8 quorumNumber = 1;
        uint8 adversaryThresholdPercentage = eigenDAThresholdRegistry.getQuorumAdversaryThresholdPercentage(quorumNumber);
        assertEq(adversaryThresholdPercentage, uint8(quorumAdversaryThresholdPercentages[quorumNumber]));
    }

    function test_getQuorumConfirmationThresholdPercentage() public {
        uint8 quorumNumber = 1;
        uint8 confirmationThresholdPercentage = eigenDAThresholdRegistry.getQuorumConfirmationThresholdPercentage(quorumNumber);
        assertEq(confirmationThresholdPercentage, uint8(quorumConfirmationThresholdPercentages[quorumNumber]));
    }

    function test_getIsQuorumRequired() public {
        uint8 quorumNumber = 0;
        bool isQuorumRequired = eigenDAThresholdRegistry.getIsQuorumRequired(quorumNumber);
        assertEq(isQuorumRequired, true);
        quorumNumber = 1;
        isQuorumRequired = eigenDAThresholdRegistry.getIsQuorumRequired(quorumNumber);
        assertEq(isQuorumRequired, true);
        quorumNumber = 2;
        isQuorumRequired = eigenDAThresholdRegistry.getIsQuorumRequired(quorumNumber);
        assertEq(isQuorumRequired, false);
    }

    function test_getBlobParams() public {
        VersionedBlobParams memory _versionedBlobParams = VersionedBlobParams({
            maxNumOperators: 999,
            numChunks: 999,
            codingRate: 9
        });
        vm.prank(registryCoordinatorOwner);
        uint16 version = eigenDAThresholdRegistry.addVersionedBlobParams(_versionedBlobParams);
        VersionedBlobParams memory blobParams = eigenDAThresholdRegistry.getBlobParams(version);
        assertEq(blobParams.maxNumOperators, _versionedBlobParams.maxNumOperators);
        assertEq(blobParams.numChunks, _versionedBlobParams.numChunks);
        assertEq(blobParams.codingRate, _versionedBlobParams.codingRate);
    }
}