// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import "../MockEigenDADeployer.sol";
import {EigenDACertVerificationLib as CertLib} from "src/periphery/cert/libraries/EigenDACertVerificationLib.sol";
import {EigenDATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";
import {EigenDATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";
import {EigenDACertTypes} from "src/periphery/cert/EigenDACertTypes.sol";

contract EigenDACertVerifierV2Unit is MockEigenDADeployer {
    using stdStorage for StdStorage;
    using BN254 for BN254.G1Point;

    address relay0 = address(uint160(uint256(keccak256(abi.encodePacked("relay0")))));
    address relay1 = address(uint160(uint256(keccak256(abi.encodePacked("relay1")))));

    function setUp() public virtual {
        quorumNumbersRequired = hex"00";
        _deployDA();
    }

    function _getDACert(uint256 seed) internal returns (EigenDACertTypes.EigenDACertV3 memory) {
        (EigenDATypesV2.SignedBatch memory signedBatch, EigenDATypesV2.BlobInclusionInfo memory blobInclusionInfo,) =
            _getSignedBatchAndBlobVerificationProof(seed, 0);

        (DATypesV1.NonSignerStakesAndSignature memory nonSignerStakesAndSignature, bytes memory signedQuorumNumbers) =
            CertLib.getNonSignerStakesAndSignature(operatorStateRetriever, registryCoordinator, signedBatch);

        return EigenDACertTypes.EigenDACertV3(
            signedBatch.batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers
        );
    }

    function test_verifyDACertV2(uint256 pseudoRandomNumber) public {
        EigenDACertTypes.EigenDACertV3 memory cert = _getDACert(pseudoRandomNumber);
        uint8 res = eigenDACertVerifier.checkDACert(abi.encode(cert));
        assertEq(res, 1);
    }

    function test_verifyDACertV2_revert_InclusionProofInvalid(uint256 pseudoRandomNumber) public {
        EigenDACertTypes.EigenDACertV3 memory cert = _getDACert(pseudoRandomNumber);

        cert.blobInclusionInfo.inclusionProof =
            abi.encodePacked(keccak256(abi.encode(pseudoRandomNumber, "inclusion proof")));
        uint8 res = eigenDACertVerifier.checkDACert(abi.encode(cert));
        assertEq(res, uint8(CertLib.StatusCode.INVALID_INCLUSION_PROOF));
    }

    function _getSignedBatchAndBlobVerificationProof(uint256 pseudoRandomNumber, uint8 version)
        internal
        returns (
            EigenDATypesV2.SignedBatch memory,
            EigenDATypesV2.BlobInclusionInfo memory,
            BLSSignatureChecker.NonSignerStakesAndSignature memory
        )
    {
        EigenDATypesV2.BlobHeaderV2 memory blobHeader1 = _getRandomBlobHeaderV2(pseudoRandomNumber, version);
        EigenDATypesV2.BlobHeaderV2 memory blobHeader2 = _getRandomBlobHeaderV2(pseudoRandomNumber, version);

        uint32[] memory relayKeys = new uint32[](2);
        relayKeys[0] = 0;
        relayKeys[1] = 1;

        EigenDATypesV2.BlobCertificate memory blobCertificate1 =
            EigenDATypesV2.BlobCertificate({blobHeader: blobHeader1, signature: hex"00", relayKeys: relayKeys});

        EigenDATypesV2.BlobCertificate memory blobCertificate2 =
            EigenDATypesV2.BlobCertificate({blobHeader: blobHeader2, signature: hex"0001", relayKeys: relayKeys});

        bytes32 batchRoot = keccak256(
            abi.encode(
                keccak256(abi.encode(CertLib.hashBlobCertificate(blobCertificate1))),
                keccak256(abi.encode(CertLib.hashBlobCertificate(blobCertificate2)))
            )
        );

        EigenDATypesV2.BlobInclusionInfo memory blobInclusionInfo = EigenDATypesV2.BlobInclusionInfo({
            blobCertificate: blobCertificate1,
            blobIndex: 0,
            inclusionProof: abi.encodePacked(keccak256(abi.encode(CertLib.hashBlobCertificate(blobCertificate2))))
        });

        (
            uint32 referenceBlockNumber,
            BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature
        ) = _registerSignatoriesAndGetNonSignerStakeAndSignatureRandom(pseudoRandomNumber, 0, 1);

        EigenDATypesV2.BatchHeaderV2 memory batchHeader =
            EigenDATypesV2.BatchHeaderV2({batchRoot: batchRoot, referenceBlockNumber: referenceBlockNumber});

        nonSignerStakesAndSignature.sigma =
            BN254.hashToG1(keccak256(abi.encode(batchHeader))).scalar_mul(aggSignerPrivKey);

        uint32[] memory quorumNumbers = new uint32[](1);
        quorumNumbers[0] = 0;

        EigenDATypesV2.Attestation memory attestation = EigenDATypesV2.Attestation({
            nonSignerPubkeys: nonSignerStakesAndSignature.nonSignerPubkeys,
            quorumApks: nonSignerStakesAndSignature.quorumApks,
            sigma: nonSignerStakesAndSignature.sigma,
            apkG2: nonSignerStakesAndSignature.apkG2,
            quorumNumbers: quorumNumbers
        });

        EigenDATypesV2.SignedBatch memory signedBatch =
            EigenDATypesV2.SignedBatch({batchHeader: batchHeader, attestation: attestation});

        return (signedBatch, blobInclusionInfo, nonSignerStakesAndSignature);
    }

    function _getRandomBlobHeaderV2(uint256 psuedoRandomNumber, uint8 version)
        internal
        pure
        returns (EigenDATypesV2.BlobHeaderV2 memory)
    {
        uint256[2] memory lengthCommitmentX = [
            uint256(keccak256(abi.encode(psuedoRandomNumber, "blobHeader.commitment.lengthCommitment.X"))),
            uint256(keccak256(abi.encode(psuedoRandomNumber, "blobHeader.commitment.lengthCommitment.X")))
        ];
        uint256[2] memory lengthCommitmentY = [
            uint256(keccak256(abi.encode(psuedoRandomNumber, "blobHeader.commitment.lengthCommitment.Y"))),
            uint256(keccak256(abi.encode(psuedoRandomNumber, "blobHeader.commitment.lengthCommitment.Y")))
        ];
        uint256[2] memory lengthProofX = [
            uint256(keccak256(abi.encode(psuedoRandomNumber, "blobHeader.commitment.lengthProof.X"))),
            uint256(keccak256(abi.encode(psuedoRandomNumber, "blobHeader.commitment.lengthProof.X")))
        ];
        uint256[2] memory lengthProofY = [
            uint256(keccak256(abi.encode(psuedoRandomNumber, "blobHeader.commitment.lengthProof.Y"))),
            uint256(keccak256(abi.encode(psuedoRandomNumber, "blobHeader.commitment.lengthProof.Y")))
        ];

        EigenDATypesV2.BlobHeaderV2 memory blobHeader = EigenDATypesV2.BlobHeaderV2({
            version: version,
            quorumNumbers: hex"00",
            commitment: EigenDATypesV2.BlobCommitment({
                commitment: BN254.G1Point(
                    uint256(keccak256(abi.encode(psuedoRandomNumber, "blobHeader.commitment.X"))),
                    uint256(keccak256(abi.encode(psuedoRandomNumber, "blobHeader.commitment.Y")))
                ),
                lengthCommitment: BN254.G2Point(lengthCommitmentX, lengthCommitmentY),
                lengthProof: BN254.G2Point(lengthProofX, lengthProofY),
                length: uint32(uint256(keccak256(abi.encode(psuedoRandomNumber, "blobHeader.length"))))
            }),
            paymentHeaderHash: keccak256(abi.encode(psuedoRandomNumber, "blobHeader.paymentHeaderHash"))
        });

        return blobHeader;
    }
}
