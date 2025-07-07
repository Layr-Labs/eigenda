#!/bin/bash
set -o errexit -o nounset -o pipefail

# This script compiles the Solidity contracts and generates Go bindings using abigen.

function create_binding_abi_only {
  contract_dir=$1
  contract=$2
  binding_dir=$3
  echo $contract
  mkdir -p $binding_dir/${contract}
  contract_json="$contract_dir/out/${contract}.sol/${contract}.json"
  solc_abi=$(cat ${contract_json} | jq -r '.abi')

  mkdir -p data
  echo ${solc_abi} >data/tmp.abi

  rm -f $binding_dir/${contract}/binding.go
  # We generate the Go bindings only with the ABI, without bytecode.
  # If you need bindings that include the bytecode to be able to deploy the contract from golang code,
  # you will need to pass --bin=data/tmp.bin and create the tmp.bin file similarly to how tmp.abi is created,
  # using `jq -r '.bytecode.object'` on the contract JSON file.
  abigen --abi=data/tmp.abi --pkg=contract${contract} --out=$binding_dir/${contract}/binding.go
}

forge clean
forge build

contracts="PaymentVault \
  SocketRegistry \
  AVSDirectory \
  DelegationManager \
  BitmapUtils \
  OperatorStateRetriever \
  RegistryCoordinator \
  BLSApkRegistry \
  IIndexRegistry \
  StakeRegistry \
  BN254 \
  EigenDAServiceManager \
  IEigenDAServiceManager \
  EjectionManager \
  EigenDACertVerifierV1 \
  EigenDACertVerifierV2 \
  IEigenDACertTypeBindings \
  EigenDACertVerifier \
  EigenDACertVerifierRouter \
  IEigenDACertVerifierLegacy \
  EigenDAThresholdRegistry \
  EigenDARelayRegistry \
  IEigenDARelayRegistry \
  EigenDADisperserRegistry \
  IRelayRegistry \
  IEigenDADirectory"

for contract in $contracts; do
  create_binding_abi_only ./ $contract ./bindings
done
