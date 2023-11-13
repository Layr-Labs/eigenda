#!/bin/bash

for ((i=0;i<10;i++)); do
    #  Generate 1KB of random data and store it in a variable called "data"
    #  The data is stored in hex format
    data=$(printf '1%.0s' {1..1000} | base64 | tr -d '\n')
    grpcurl -plaintext -d "{\"data\": \"$data\", \"security_params\": [{\"quorum_id\": 0, \"adversary_threshold\": 50, \"quorum_threshold\": 100}]}" localhost:32003 disperser.Disperser/DisperseBlob
    sleep 0.5
done

