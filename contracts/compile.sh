#!/bin/bash

function create_binding {
    contract_dir=$1
    contract=$2
    binding_dir=$3
    echo $contract
    mkdir -p $binding_dir/${contract}
    contract_json="$contract_dir/out/${contract}.sol/${contract}.json"
    solc_abi=$(cat ${contract_json} | jq -r '.abi')
    solc_bin=$(cat ${contract_json} | jq -r '.bytecode.object')

    mkdir -p data
    echo ${solc_abi} > data/tmp.abi
    echo ${solc_bin} > data/tmp.bin

    rm -f $binding_dir/${contract}/binding.go
    abigen --bin=data/tmp.bin --abi=data/tmp.abi --pkg=contract${contract} --out=$binding_dir/${contract}/binding.go
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
  IndexRegistry \
  StakeRegistry \
  BN254 \
  EigenDAServiceManager \
  IEigenDAServiceManager \
  MockRollup \
  EjectionManager \
  EigenDABlobVerifier \
  EigenDAThresholdRegistry \
  EigenDARelayRegistry \
  EigenDADisperserRegistry"

for contract in $contracts; do
    create_binding ./ $contract ./bindings
done

# ./compile.sh ./ BitmapUtils ./bindings 
# ./compile.sh ./ BLSOperatorStateRetriever ./bindings
# ./compile.sh ./ BN254 ./bindings
# ./compile.sh ./ BLSRegistryCoordinatorWithIndices ./bindings
# ./compile.sh ./ IBLSPubkeyRegistry ./bindings
# ./compile.sh ./ IIndexRegistry ./bindings
# ./compile.sh ./ IStakeRegistry ./bindings
