# Structured Reference String (SRS) Files

This directory contains the Structured Reference String (SRS) files required for KZG commitments and proofs in EigenDA.

## File Information

| File Name          | Size   | Number of Points | Point Size |
|--------------------|--------|------------------|------------|
| g1.point           | 16 MB  | 524,288          | 32 bytes   |
| g2.point           | 32 MB  | 524,288          | 64 bytes   |
| g2.trailing.point  | 32 MB  | 524,288          | 64 bytes   |

These files are only a portion of the total SRS data that exists for eigenDA. The files here are sufficiently large
to support the largest permitted blob size of 16MB. This maximum blob size may be increased at some point in the future,
at which time larger SRS files will need to be committed.

## Retrieving SRS Files

These SRS files can be fetched from the AWS S3 bucket at:
```
https://srs-mainnet.s3.amazonaws.com/kzg/
```

### Important Notes:

1. The S3 bucket only directly contains `g1.point` and `g2.point` files.
2. To retrieve trailing G2 points, you must explicitly specify the byte range when using curl. For example:

   ```bash
   # To download the entire g1.point file
   curl -o g1.point https://srs-mainnet.s3.amazonaws.com/kzg/g1.point
   
   # To download the entire g2.point file
   curl -o g2.point https://srs-mainnet.s3.amazonaws.com/kzg/g2.point
   
   # To download g2.trailing.point (by specifying byte range)
   # The exact byte range will depend on your requirements
   curl -o g2.trailing.point -r <start>-<end> https://srs-mainnet.s3.amazonaws.com/kzg/g2.point
   ```

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
