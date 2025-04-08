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
