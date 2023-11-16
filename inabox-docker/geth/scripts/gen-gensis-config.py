#!/usr/bin/env python3

import sys
import json
import os
from os import listdir
from os.path import isfile, join
from pathlib import Path


template = '''
{
  "config": {
    "chainId": 40525,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "istanbulBlock": 0,
    "petersburgBlock": 0,
    "berlinBlock": 0,
    "londonBlock": 0,
    "clique": {
      "period": 0,
      "epoch": 30000
    }
  },
  "difficulty": "1",
  "gasLimit": "3000000000000",
  "extradata": {},
  "alloc": {}
}
'''

if len(sys.argv) < 2:
    print("need number of keys to load. max 1000")
    sys.exit(1)

num_key_limit = int(sys.argv[1])

genesis = json.loads(template)

key_dir = './data/keystore'
balance = "10000000000000000000000"
alloc = {}

keystoreFiles = [key_dir + '/' + f for f in listdir(key_dir) if isfile(join(key_dir, f))]
keystoreFiles.sort(key=os.path.getctime)

print(keystoreFiles)

addresses = ""
signer = ''
num = 0
for fpath in keystoreFiles:
    with open(fpath) as f:
        config = json.load(f)
        addresses += config['address']
        alloc[config['address']] = {"balance": balance}
        num += 1
        if signer == '':
            signer = config['address']

        if num >= num_key_limit:
            break


extra_data_template = "0x0000000000000000000000000000000000000000000000000000000000000000{}0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
extra_data = extra_data_template.format(signer)

genesis['alloc'] = alloc
genesis['extradata'] = extra_data 

with open('genesis.json', 'w') as f:
    json.dump(genesis, f, indent=4)
