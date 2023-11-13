#!/bin/bash

docker exec geth sh -c "cd /root/.ethereum/ && geth --exec 'tx = eth.getTransaction(\"$1\"); eth.call(tx,tx.blockNumber)' attach geth.ipc"

