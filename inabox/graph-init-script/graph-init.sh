#!/bin/bash -e

../graph-init

for d in $(find /subgraphs -maxdepth 1 -mindepth 1 -type d); do
	echo $d
	pushd $d
	cat package.json
	yarn create-docker
	yarn deploy-docker --version-label=v0.0.1
	popd
done

node ../server.js
