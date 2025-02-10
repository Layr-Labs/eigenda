#!/bin/bash
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
  MockRollup \
  EjectionManager \
  EigenDACertVerifier \
  EigenDAThresholdRegistry \
  EigenDARelayRegistry \
  IEigenDARelayRegistry \
  EigenDADisperserRegistry"

output_dir="./test/storage"
for contract in $contracts 
do
    echo "Checking storage change of $contract"
    [ -f "$output_dir/$contract" ] && mv "$output_dir/$contract" "$output_dir/$contract-old"
    forge inspect "$contract" --pretty storage > "$output_dir/$contract"
    diff "$output_dir/$contract-old" "$output_dir/$contract"
    if [[ $? != "0" ]]
    then
        CHANGED=1
    fi
done
if [[ $CHANGED == 1 ]]
then
    exit 1
fi