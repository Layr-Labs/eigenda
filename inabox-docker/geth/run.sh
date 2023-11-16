#!/bin/bash

# generate keys using hardhat
# go to data-layr hardhat-node dir
# copy keys here


function load_keystore {
	if [ "$#" -ne 3 ]; then
		tput setaf 1
		echo "Usage. private_keys_path<str> geth_password_path<str> num-key<int>"
		tput sgr0
		exit 1
	fi

	limit=$1
	private_key_path=$2
	geth_password_path=$3


	if [ ! -f "${private_key_path}" ]; then
		echo "need private-keys.txt"
		echo "each line in the file needs to be a private key"
		echo "without prefix 0x"
		exit 1
	fi
	#
	if [ ! -f "${geth_password_path}" ]; then
		echo "add geth node password in file geth-account-password"
		exit 1
	fi

	keys=$(cat ${private_key_path})
	num=0
	for k in $keys; do
		echo ${k:2} > secret/tmp
		echo "importing $(cat secret/tmp)"
		geth account import --password ${geth_password_path} --datadir data secret/tmp
		num=$((num + 1))
		if [ $num	-ge $limit ]; then
			break
		fi
	done
	rm secret/tmp
}

function setup {
	if [ "$#" -ne 3 ]; then
		tput setaf 1
		echo "Usage. num-key<int> private_keys_path<str> geth_password_path<str>"
		echo "exmaple: setup 10 ../secrets/ecdsa_keys/private_keys_hex.txt ./secret/geth-account-password"
		tput sgr0
		exit 1
	fi
	num_key=$1
	# load key to geth state
	load_keystore $1 $2 $3
	# generate genesis.jon
	rm genesis.json
	./scripts/gen-gensis-config.py ${num_key}
	# init geth states
	rm -rf data/geth
	geth init --datadir data genesis.json
}

function run_geth {
	nid=$(cat genesis.json | jq -r '.config | .chainId')
	extradata=$(cat genesis.json | jq -r '.extradata')
	address=0x${extradata:66:40} # 66 = 2 ( 32 + 1) hex * (32zero + '0x')
	echo "geth --datadir data --networkid $nid --nodiscover --netrestrict 127.0.0.1/0 --unlock $address  --mine --miner.etherbase=$address   --password secret/geth-account-password --rpc.gascap 0 --miner.gasprice 1 --rpc.allow-unprotected-txs   --miner.gaslimit '2922337203600000' "
	nohup geth --datadir data --networkid $nid --ws --ws.addr 0.0.0.0 --http --http.addr 0.0.0.0 --nodiscover --netrestrict 127.0.0.1/0 --unlock $address  --mine --allow-insecure-unlock  --password secret/geth-account-password --rpc.gascap 0 --miner.gasprice 1 --rpc.allow-unprotected-txs  --miner.gaslimit "2922337203600000" &


}

function console {
	geth attach data/geth.ipc
}

function openWS {
	geth attach data/geth.ipc --exec 'admin.startWS("127.0.0.1", 8546)'
}

function step_mining {
	#period=$1
	geth attach data/geth.ipc --exec 'miner.start()'
	# leave sufficient time to mine a block
	sleep 0.1
	geth attach data/geth.ipc --exec 'miner.stop()'
}


case "$1" in
  help)
        cat <<-EOF
    Setup
      setup                  Setup genesis.json and load keys
    Load geth keystore
      load-keystore          Load geth keystore
    Run geth
      run-geth

EOF
        ;;
    load-keystore)
			load_keystore ${@:2} ;;
		setup)
			setup ${@:2} ;;
		run-geth)
			run_geth ;;
		console)
			console ;;
		openWS)
			openWS ;;
		step-mining)
			step_mining ${@:2} ;;
    *)
        tput setaf 1
        echo "Unknown subcommand" $1
        echo "./local.sh help"
        tput sgr0 ;;
esac
