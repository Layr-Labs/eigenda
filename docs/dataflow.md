```mermaid
flowchart TB

disperser -- []byte --> blobstore

blobstore -- []byte --> encserver
blobstore --> relay.Server.GetBlob

encserver -- []byte --> prover
prover -- []encoding.Frame --> encserver

prover -- []fr.Element --> prover1
prover1 -- []bn254.G1Affine --> prover
%% prover1: polyFr ->(GetSlicesCoeff) coeffStore -> sumVec -> fft_inv -> fft -> proof
%% fft: fast compute polynomial value using Toeplitz Matrix

encserver -- v2.FragmentInfo --> encmgr

relay -- []merkletree.Proof []encoding.Frame --> node

dispatcher -- BatchHeader.BatchRoot []v2.BlobHeader --> node
node -- Signature --> dispatcher

encserver -- []encoding.Frame --> chunkstore
chunkstore -- v2.FragmentInfo --> encserver
chunkstore -- []encoding.Frame --> relay

node -- BatchHeader.BatchRoot []v2.BlobHeader --> validate1 
node -- []BlobHeader.BlobCommitment []merkletree.Proof []encoding.Frame --> validate2
node -- []BlobHeader.BlobCommitment --> validate3
node --> batchstore --> ServerV2.GetChunks

encmgr -- []v2.BlobHeader --> certstore
certstore -- []v2.BlobHeader --> newbatch

newbatch -- BatchHeader.BatchRoot []v2.BlobHeader --> dispatcher
dispatcher --> Dispatcher.HandleSignatures

newbatch -- []merkletree.Proof --> verifstore 
verifstore -- []merkletree.Proof --> relay

newbatch[Dispatcher.NewBatch]
dispatcher[Dispatcher.HandleBatch]
prover[Prover.GetFrames]
prover1[ComputeMultiFrameProof]
encmgr[EncodingManager.HandleBatch]
encserver[EncoderServerV2.handleEncodingToChunkStore]

relay[relay.Server.GetChunks]
node[*node* ServerV2.StoreChunks]
batchstore[(node.StoreV2)]
validate1[ValidateBatchHeader]
validate2[ValidateBlobs Verifier.UniversalVerify]
validate3[VerifyBlobLength VerifyCommitEquivalenceBatch]

blobstore[(blobstore)]
chunkstore[(chunkstore)]
verifstore[(BlobMetadataStore)]
certstore[(BlobMetadataStore)]
```