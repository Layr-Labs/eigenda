#!/bin/bash
set -o errexit -o nounset -o pipefail

# This script compiles Solidity contracts with Foundry and generates Go bindings using abigen.
# You can choose which contracts use abigen v1 vs v2. Outputs:
#   - v2 -> $binding_dir/v2/<Contract>/binding.go
#   - v1 -> $binding_dir/<Contract>/binding.go
#
# This allows us to migrate contracts over time to use abigen v2 without introducing a very large
# breaking change at once.
# Make sure that `forge build` has been run before executing this script.

binding_dir="./bindings"
artifacts_root="./out"
go_pkg_prefix="contract"
abi_gen_v1="v1"
abi_gen_v2="v2"

ABIGEN_V2_CONTRACTS=(
  "EigenDACertVerifier"
)

ABIGEN_V1_CONTRACTS=(
  "PaymentVault"
  "SocketRegistry"
  "AVSDirectory"
  "DelegationManager"
  "BitmapUtils"
  "OperatorStateRetriever"
  "EigenDARegistryCoordinator"
  "BLSApkRegistry"
  "IIndexRegistry"
  "StakeRegistry"
  "BN254"
  "EigenDAServiceManager"
  "IEigenDAServiceManager"
  "EjectionManager"
  "EigenDACertVerifierV1"
  "EigenDACertVerifierV2"
  "IEigenDACertTypeBindings"
  "EigenDACertVerifier"
  "EigenDACertVerifierRouter"
  "IEigenDACertVerifierLegacy"
  "EigenDAThresholdRegistry"
  "EigenDARelayRegistry"
  "IEigenDARelayRegistry"
  "EigenDADisperserRegistry"
  "IEigenDADirectory"
)

build_artifact_json_path() {
  # args: <contract>
  local contract="$1"
  echo "${artifacts_root}/${contract}.sol/${contract}.json"
}


create_golang_abi_binding() {
  # args: <contract> <abigen_version: v1|v2>
  local contract="$1"
  local abigen_version="$2"

  local contract_json
  contract_json="$(build_artifact_json_path "${contract}")"
  if [[ ! -f "${contract_json}" ]]; then
    echo "❌ Missing artifact JSON for ${contract} at ${contract_json}" >&2
    return 1
  fi

  # Extract contract's ABI from foundry build artifact JSON
  local solc_abi
  solc_abi="$(jq -r '.abi' < "${contract_json}")"
  if [[ -z "${solc_abi}" || "${solc_abi}" == "null" ]]; then
    echo "❌ No ABI found in ${contract_json}" >&2
    return 1
  fi

  # output ABI to temporary file referenced during go binding generation
  mkdir -p data
  echo "${solc_abi}" > data/tmp.abi

  local out_dir
  if [[ "${abigen_version}" == "v2" ]]; then
    out_dir="${binding_dir}/v2/${contract}"
  else
    out_dir="${binding_dir}/${contract}"
  fi
  mkdir -p "${out_dir}"

  # Remove any previous generated golang binding to avoid stale diffs
  rm -f "${out_dir}/binding.go"

  # Build abigen args
  local pkg="${go_pkg_prefix}${contract}"
  local args=( --abi=data/tmp.abi --pkg="${pkg}" --out="${out_dir}/binding.go" )
  if [[ "${abigen_version}" == "v2" ]]; then
    args=( --v2 "${args[@]}" )
  fi

  echo "🔧 abigen ${abigen_version} → ${out_dir}/binding.go (${contract})"
  abigen "${args[@]}"
}

main() {
  # abigen v1
  for contract in "${ABIGEN_V1_CONTRACTS[@]}"; do
    create_golang_abi_binding "${contract}" ${abi_gen_v1}
  done
  
  echo
  echo "======================================"
  echo

  # abigen v2
  for contract in "${ABIGEN_V2_CONTRACTS[@]}"; do
    create_golang_abi_binding "${contract}" ${abi_gen_v2}
  done

  echo 
  echo "✅ Done."
}

main "$@"