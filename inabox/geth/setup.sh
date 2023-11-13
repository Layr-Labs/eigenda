#!/bin/bash

function call_geth {
    command=$1
    docker run -v $PWD:/geth ethereum/client-go:v1.10.17 $command
}

function load_keystore {
	if [ "$#" -ne 2 ]; then
		tput setaf 1
		echo "Usage. num-key<int> private_keys_path<str>"
		tput sgr0
		exit 1
	fi

	limit=$1
	private_key_path=$2


	if [ ! -f "${private_key_path}" ]; then
		echo "need private-keys.txt"
		echo "each line in the file needs to be a private key"
		echo "without prefix 0x"
		exit 1
	fi

	keys=$(cat ${private_key_path})
	num=0
	for k in $keys; do
		echo ${k:2} > secret/tmp
		echo "importing $(cat secret/tmp)"
		call_geth "account import --password /geth/secret/geth-account-password --datadir /geth/data /geth/secret/tmp"
		num=$((num + 1))
		if [ $num	-ge $limit ]; then
			break
		fi
	done
	rm secret/tmp
}

function setup {
	if [ "$#" -ne 2 ]; then
		tput setaf 1
		echo "Usage. num-key<int> private_keys_path<str>"
		echo "example: setup 10 ./secret/private-keys.txt "
		tput sgr0
		exit 1
	fi


	num_key=$1
	# load key to geth state
	load_keystore $1 $2
	# init geth states
	call_geth "init --datadir /geth/data /geth/genesis.json"
}

setup $@