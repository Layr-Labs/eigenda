# Structured Reference String (SRS) Files

This directory contains the Structured Reference String (SRS) files required for KZG commitments and proofs in EigenDA.

## File Information

| File Name          | Size   | Number of Points | Point Size | SHA256 Hash                                                      |
|--------------------|--------|------------------|------------|------------------------------------------------------------------|
| g1.point           | 16 MB  | 524,288          | 32 bytes   | 8f18b9c04ed4bddcdb73001fb693703197328cecabdfa9025f647410b0c50d7f |
| g2.point           | 32 MB  | 524,288          | 64 bytes   | a6942684aa751b4ec7873e2edb4660ac5c4516adb3b310441802cc0d489f645a |
| g2.trailing.point  | 32 MB  | 524,288          | 64 bytes   | 78fad17d74d28cecdb7f826fdd72dee08bdbe1e8ad66f2b24fcf2fc140176788 |

These files are only a portion of the total SRS data that exists for eigenDA. The files here are sufficiently large
to support the largest permitted blob size of 16MB. This maximum blob size may be increased at some point in the future,
at which time larger SRS files will be needed.

Note that the G2 point files (`g2.point` and `g2.trailing.point`) are twice the size of the G1 point file because G2 
points require twice as many bytes to represent as G1 points in the BN254 curve. Each G1 point requires 32 bytes 
of storage, while each G2 point requires 64 bytes.

## Retrieving SRS Files

These SRS files can be fetched from the AWS S3 bucket at:
```
https://srs-mainnet.s3.amazonaws.com/kzg/
```

### Using the Download Script

The EigenDA repository provides a convenience script located at `tools/download-srs.sh` for downloading and 
generating SRS file hashes.

Usage:
```bash
# Download SRS files for a specified blob size (in bytes)
./tools/download-srs.sh <size_in_bytes>

# Example: Download SRS files for 16MB blobs
./tools/download-srs.sh 16777216
```

The script will:
1. Download the appropriate portions of g1.point and g2.point files
2. Calculate the correct byte range for g2.trailing.point
3. Save all files to a "srs-files" directory
4. Generate SHA256 hashes for verification
5. Create a hash file with labeled hashes for each downloaded file

## SRS Verification and Alternative Retrieval Method

For users who need to verify the integrity of SRS files, please refer to the
[SRS Utilities README](/tools/srs-utils/README.md) for detailed instructions. This tool provides:

1. Methods to extract and verify G1 and G2 points from the original Perpetual Powers of Tau challenge file
2. Verification procedures to ensure the correctness of the SRS points based on approaches used by the Ethereum 
   Foundation's KZG ceremony
3. Ability to parse the full 8GB SRS files from the original challenge file, which can then be manually truncated 
   to smaller sizes as needed

The SRS utilities provide an alternative approach to obtaining SRS files by downloading the original challenge file 
directly from the Ethereum Foundation's trusted setup, extracting the points, and verifying their integrity.

## Security Considerations

Using the correct SRS files is essential for the proper functioning of any software interacting with EigenDA. 
If a user has incorrect or tampered SRS files, the following would occur:

1. **Verification failures**: The user would be unable to successfully verify KZG commitments and proofs, making it 
   impossible to validate blob data from the network.

2. **Network incompatibility**: A node using incorrect SRS files would be unable to meaningfully interact with the 
   EigenDA network, as it would consistently fail to verify certificates from honest nodes.

3. **Self-isolation**: Rather than creating a security vulnerability, having incorrect SRS files simply results in 
   self-isolation from the network's consensus.

It's important to understand that this isn't a security concern for the broader network. If a disperser attempted 
to use incorrect SRS files, honest nodes would immediately detect this during certificate verification. The network's 
security properties rely on the fact that honest nodes using the correct SRS can detect and reject improperly generated 
certificates.
