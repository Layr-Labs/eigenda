#!/bin/bash

function create_binding {
    local contract_dir="$1"
    local contract="$2"
    local binding_dir="$3"
    echo "$contract"
    mkdir -p "$binding_dir/${contract}"
    local contract_json="$contract_dir/out/${contract}.sol/${contract}.json"
    local solc_abi
    local solc_bin
    solc_abi=$(jq -r '.abi' "$contract_json")
    solc_bin=$(jq -r '.bytecode.object' "$contract_json")

    mkdir -p data
    echo "$solc_abi" > data/tmp.abi
    echo "$solc_bin" > data/tmp.bin

    rm -f "$binding_dir/${contract}/binding.go"
    abigen --bin <(echo "$solc_bin") --abi <(echo "$solc_abi") --pkg "contract${contract}" --out "$binding_dir/${contract}/binding.go"
}

# Clean and build contracts
forge clean
forge build

# List of contracts
contracts="AVSDirectory DelegationManager BitmapUtils OperatorStateRetriever RegistryCoordinator BLSApkRegistry IndexRegistry StakeRegistry BN254 EigenDAServiceManager IEigenDAServiceManager MockRollup"

# Create bindings for each contract
for contract in $contracts; do
    create_binding ./ "$contract" ./bindings
done
