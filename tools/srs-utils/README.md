# srs-verification

This project is written for EigenDA SRS related works. It includes
1. Extract G1 and G2 points used by EigenDA from ptau challenge file, created from the perpetual power of tau MPC instance run by EF
2. Verify the extracted points are indeed correct based on approaches used by EF KZG ceremony
3. Provide utility tools for miscellaneous tasks

## To download the ptau file

`wget https://pse-trusted-setup-ppot.s3.eu-central-1.amazonaws.com/challenge_0085`

See more information from 
1. https://docs.axiom.xyz/docs/transparency-and-security/kzg-trusted-setup
2. https://github.com/privacy-scaling-explorations/perpetualpowersoftau/tree/master 

The challenge file has 103079215232 Bytes.

## How to use

`go run main.go help`

###  How to parse G1, G2 points from the challenge file.

`go run main.go parse --ptau-path <Path to challenge file>`

It produces two files, g1.point and g2.point. g1.point contains 8589934592 Bytes and g2.point 17179869184 Bytes

This procedure takes roughly 10 minutes.

Note: The challenge files contains 2^29 G1 points and 2^28 G2 points with secret tau. We use only the first 2^28 G1 points for EigenDA.

### How to verify the parsed G1, G2 points

`go run main.go verify --g1-path <Path to g1.point> --g2-path <Path to g2.point>`

The verification is based on method listed here (https://github.com/ethereum/kzg-ceremony-specs/blob/master/docs/sequencer/sequencer.md#pairing-checks)

This procedures takes roughly 27 Hours on a 8 thread machines.

The program periodically prints out the time spent and its progress of validating 2^28 G1 and G2 points. If no error message is shown and program terminates with "Done. Everything is correct". Then SRS is deemed as correct. 

