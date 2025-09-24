# SRS Utilities

This project provides tools for working with EigenDA's Structured Reference String (SRS). It includes tools to:

1. Download pre-processed SRS files directly from the EigenDA S3 bucket
2. Download precomputed SRS tables for EigenDA V2 encoding operations
3. Extract G1 and G2 points used by EigenDA from the ptau challenge file, created from the Perpetual Powers of Tau MPC ceremony run by the Ethereum Foundation
4. Verify that the extracted points are correct based on approaches used in the Ethereum Foundation's KZG ceremony

## Structured Reference String (SRS) Files

The SRS files are required for KZG commitments and proofs in EigenDA.

### Core SRS Files

| File Name          | Size   | Number of Points | Point Size | SHA256 Hash                                                      |
|--------------------|--------|------------------|------------|------------------------------------------------------------------|
| g1.point           | 16 MB  | 524,288          | 32 bytes   | 8f18b9c04ed4bddcdb73001fb693703197328cecabdfa9025f647410b0c50d7f |
| g2.point           | 32 MB  | 524,288          | 64 bytes   | a6942684aa751b4ec7873e2edb4660ac5c4516adb3b310441802cc0d489f645a |
| g2.trailing.point  | 32 MB  | 524,288          | 64 bytes   | 78fad17d74d28cecdb7f826fdd72dee08bdbe1e8ad66f2b24fcf2fc140176788 |
| g2.point.powerOf2  | 1.8 KB | 28               | 64 bytes   | 4d5ed827f742e1270f22b4a39129bf1d25445821b15824e2eb3a709a16f64518 |

These files represent only a portion of the total SRS data that exists for EigenDA. They are sufficiently large
to support the largest permitted blob size of 16MB. This maximum blob size may increase in the future,
at which point larger SRS files will be needed.

Note that the G2 point files (`g2.point` and `g2.trailing.point`) are twice the size of the G1 point file because G2 
points require twice as many bytes to represent as G1 points in the BN254 curve. Each G1 point requires 32 bytes 
of storage, while each G2 point requires 64 bytes.

The `g2.point.powerOf2` file contains only G2 points at power-of-2 indices (indices 1, 2, 4, 8, 16, ..., 2^27). This
optimized file contains just 28 G2 points instead of the full set, significantly reducing memory usage for operator
nodes. Since operators only perform multi-reveal proofs on blobs with power-of-2 polynomial degrees, they don't need
the complete G2 SRS. This file is optional and primarily used by operator nodes for memory efficiency.

### SRS Tables for EigenDA V2

EigenDA V2 uses precomputed SRS tables for efficient polynomial operations with specific chunk counts. These tables
contain coset evaluations that accelerate encoding and decoding operations.

In EigenDA V2, **blob version 0** specifically sets `numChunks=8192`, which is why the dimE8192 tables are the
primary SRS tables used in production.

#### Available Table Files

The SRS tables are organized by dimension (numChunks) and coset size:

| Dimension | Coset Sizes | Total Size | Description |
|-----------|-------------|------------|-------------|
| dimE8192  | 4, 8, 16, 32, 64, 128, 256, 512, 1024 | ~1 GB | Tables for numChunks=8192 (blob version 0) |

Each table file is named following the pattern: `<dimension>.coset<size>` (e.g., `dimE8192.coset256`)

#### Blob Size Calculation

The supported blob size depends on the coset size (chunk length) used:

```
Blob Size = (numChunks × cosetSize × 32 bytes) / codingRatio
```

Where:
- `numChunks` = 8192 (for blob version 0)
- `cosetSize` = chunk length (varies based on blob size)
- `32 bytes` = size of each BN254 field element
- `codingRatio` = 8 (fixed erasure coding expansion factor)

Supported blob sizes for dimE8192:
- cosetSize=4: blob size = 128 KB (minimum)
- cosetSize=8: blob size = 256 KB
- cosetSize=16: blob size = 512 KB
- cosetSize=32: blob size = 1 MB
- cosetSize=64: blob size = 2 MB
- cosetSize=128: blob size = 4 MB
- cosetSize=256: blob size = 8 MB
- cosetSize=512: blob size = 16 MB (current production limit)
- cosetSize=1024: blob size = 32 MB (future support)

## Installation

```bash
go install github.com/Layr-Labs/eigenda/tools/srs-utils@latest
```

## How to use

Once installed, you can run:

```bash
srs-utils help
```

### Downloading SRS Files

The simplest way to get the required SRS files is to download the pre-processed files directly from the EigenDA
S3 bucket:

```bash
srs-utils download --blob-size-bytes 16777216
```

This will download the SRS files needed for 16MB blob support (the default size). The files will be saved to a directory
named "srs-files". A hash file is also generated during download for verification purposes.

Options:
- `--blob-size-bytes`: Size of the blob in bytes (default: 16777216, which is 16MB)
- `--output-dir`: Directory where the files will be saved (default: "srs-files")
- `--base-url`: Base URL for downloading (default: "https://srs-mainnet.s3.amazonaws.com/kzg")
- `--include-g2-power-of-2`: Include the g2.point.powerOf2 file in the download (optional, for power-of-2 polynomial operations)

To download with the power-of-2 points file:

```bash
srs-utils download --blob-size-bytes 16777216 --include-g2-power-of-2
```

### Downloading SRS Tables for EigenDA V2

To download the precomputed SRS tables used by EigenDA V2 for encoding operations with numChunks=8192:

```bash
srs-utils download-tables
```

This will download all coset tables for the default dimension (dimE8192). The files will be saved to
`resources/srs/SRSTables` directory by default.

Options:
- `--dimension`: The dimension to download (default: "dimE8192")
- `--output-dir`: Directory where the tables will be saved (default: "resources/srs/SRSTables")
- `--base-url`: Base URL for downloading (default: "https://srs-mainnet.s3.amazonaws.com/kzg/SRSTables")
- `--coset-sizes`: Comma-separated list of coset sizes to download (default: "4,8,16,32,64,128,256,512,1024")

Example with custom parameters:

```bash
# Download only specific coset sizes
srs-utils download-tables --coset-sizes 256,512,1024

# Download to a custom directory
srs-utils download-tables --output-dir ./my-srs-tables
```

The download will show progress for each file and display the total size downloaded upon completion.

### Alternative: Generating SRS Files from the Original Challenge File

For users who prefer to generate SRS files directly from the original trusted setup, follow these steps:

#### 1. Download the ptau challenge file

```bash
wget https://pse-trusted-setup-ppot.s3.eu-central-1.amazonaws.com/challenge_0085
```

See more information from:
1. https://docs.axiom.xyz/docs/transparency-and-security/kzg-trusted-setup
2. https://github.com/privacy-scaling-explorations/perpetualpowersoftau/tree/master 

The challenge file has 103079215232 Bytes.

#### 2. Parse G1, G2 points from the challenge file

```bash
srs-utils parse --ptau-path <Path to challenge file>
```

It produces two files, g1.point and g2.point. g1.point contains 8589934592 Bytes and g2.point 17179869184 Bytes

This procedure takes roughly 10 minutes.

Note: The challenge file contains 2^29 G1 points and 2^28 G2 points with secret tau. We use only the first 2^28 G1 points for EigenDA.

#### 3. Verify the parsed G1, G2 points

```bash
srs-utils verify --g1-path <Path to g1.point> --g2-path <Path to g2.point>
```

The verification is based on the method listed here: https://github.com/ethereum/kzg-ceremony-specs/blob/master/docs/sequencer/sequencer.md#pairing-checks

This procedure takes approximately 27 hours on an 8-thread machine.

The program periodically prints out the time spent and its progress in validating 2^28 G1 and G2 points. If no error messages appear and the program terminates with "Done. Everything is correct", then the SRS is deemed correct.

## Security Considerations

Using the correct SRS files is essential for the proper functioning of any software interacting with EigenDA. If a
piece of software has incorrect or tampered SRS files, the following would occur:

1. **Verification failures**: The software would be unable to successfully verify KZG commitments and proofs, making it 
   impossible to validate blob data from the network.

2. **Submission failures**: The software would be unable to submit data to the EigenDA network, as it would
   consistently fail to generate commitments that can be verified by other participants.

It's important to understand that this isn't a security concern for the broader network. Rather, having incorrect SRS
files simply results in self-isolation from the network.
