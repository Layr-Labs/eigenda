# Glossary

## Rollup Batcher

Sequencer rollup node component responsible for constructing and submitting to the settlement chain user transaction batches

## Rollup Nodes

Refers to any rollup node (e,g, validator, verifier) which syncs current chain state through an onchain sequencer inbox.

## EigenDA Proxy

Side car server as a part of rollup and used for secure and trustless communication with EigenDA.

## EigenDA Client

A collection of [clients](https://github.com/Layr-Labs/eigenda/tree/bb91b829995c28e813fce46412a77f9fa428b0af/api/clients/v2) used for securely dispersing and reading EigenDA blobs.

## Rollup Payload

Compressed batches of transactions or state diffs.

## DA Certificate (DACert)

An EigenDA Certificate (or DACert for short) contains all the information needed to retrieve a blob from the EigenDA network and validate it.

## EigenDA Blob Derivation

A sequence of procedures to convert a byte array representing a DA certificate to the final rollup payload.

## Preimage Oracle

An object with an interface for fetching additional data during EigenDA blob derivation by using some keys generated from the data. Multiple implementations of the preimage oracle show up in the EigenDA. In proxy, ETH rpc serves as the preimage oracle for DAcert validity; EigenDA network rpc
serves as the preimage oracle for EigenDA blob.

## Blob Field Element

EigenDA uses bn254 curve, a field element on the bn254 curve is an integer whose range is 0 <= x < 21888242871839275222246405745257275088548364400416034343698204186575808495617. 