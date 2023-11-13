#!/bin/bash

pids=""
function kill_processes {
    echo "STOP"
    for pid in $pids; do
        echo "killing process $pid"
        kill -9 $pid
    done
}

function start_trap {
    trap kill_processes SIGINT

    set -a
    source $testpath/envs/churner.env
    set +a
    ../churner/bin/server & 

    pid="$!"
    pids="$pids $pid"

    for FILE in $(ls $testpath/envs/dis*.env); do 
        set -a
        source $FILE
        set +a
        ../disperser/bin/server &

        pid="$!"
        pids="$pids $pid"
    done

    for FILE in $(ls $testpath/envs/enc*.env); do 
        set -a
        source $FILE
        set +a
        ../disperser/bin/encoder &

        pid="$!"
        pids="$pids $pid"
    done

    for FILE in $(ls $testpath/envs/batcher*.env); do 
        set -a
        source $FILE
        set +a
        ../disperser/bin/batcher &

        pid="$!"
        pids="$pids $pid"
    done
    
    files=($(ls $testpath/envs/opr*.env))
    last_index=$(( ${#files[@]} - 1 ))

    for i in "${!files[@]}"; do 
        if [ $i -eq $last_index ]; then
            sleep 5  # Sleep for 5 seconds before the last loop iteration
        fi
        FILE=${files[$i]}
        set -a
        source $FILE
        set +a
        ../node/bin/node & 

        pid="$!"
        pids="$pids $pid"
    done

    for pid in $pids; do 
        wait $pid
    done 
}

function start_detached {

    pids=""
    waiters=""
    pid_file="$testpath/pids"

    if [[ -f "$pid_file" ]]; then
        echo "Processes still running. Run ./bin.sh stop"
        return
    fi

    mkdir -p $testpath/logs

    set -a
    source $testpath/envs/churner.env
    set +a
    ../churner/bin/server > $testpath/logs/churner.log 2>&1 & 

    pid="$!"
    pids="$pids $pid"

    ./wait-for 0.0.0.0:${CHURNER_GRPC_PORT} -- echo "Churner up" &
    waiters="$waiters $!"

    for FILE in $(ls $testpath/envs/dis*.env); do 
        set -a
        source $FILE
        set +a
        id=$(basename $FILE | tr -d -c 0-9)
        ../disperser/bin/server > $testpath/logs/dis${id}.log 2>&1 & 

        pid="$!"
        pids="$pids $pid"

        ./wait-for 0.0.0.0:${DISPERSER_SERVER_GRPC_PORT} -- echo "Disperser up" &
        waiters="$waiters $!"
    done

    for FILE in $(ls $testpath/envs/enc*.env); do 
        set -a
        source $FILE
        set +a
        id=$(basename $FILE | tr -d -c 0-9)
        ../disperser/bin/encoder > $testpath/logs/enc${id}.log 2>&1 &

        pid="$!"
        pids="$pids $pid"

        ./wait-for 0.0.0.0:${DISPERSER_ENCODER_GRPC_PORT} -- echo "Encoder up" &
        waiters="$waiters $!"
    done

    for FILE in $(ls $testpath/envs/batcher*.env); do 
        set -a
        source $FILE
        set +a
        id=$(basename $FILE | tr -d -c 0-9)
        ../disperser/bin/batcher > $testpath/logs/batcher${id}.log 2>&1 & 

        pid="$!"
        pids="$pids $pid"
    done

    for FILE in $(ls $testpath/envs/retriever*.env); do 
        set -a
        source $FILE
        set +a
        ../retriever/bin/server > $testpath/logs/retriever.log 2>&1 & 

        pid="$!"
        pids="$pids $pid"
    done

    files=($(ls $testpath/envs/opr*.env))
    last_index=$(( ${#files[@]} - 1 ))

    for i in "${!files[@]}"; do 
        if [ $i -eq $last_index ]; then
            sleep 5  # Sleep for 5 seconds before the last loop iteration
        fi
        FILE=${files[$i]}
        set -a
        source $FILE
        set +a
        id=$(basename $FILE | tr -d -c 0-9)
        ../node/bin/node > $testpath/logs/opr${id}.log 2>&1 &

        pid="$!"
        pids="$pids $pid"

        ./wait-for 0.0.0.0:${NODE_DISPERSAL_PORT} -- echo "Node up" &
        waiters="$waiters $!"
    done

    echo $pids > $pid_file

    for waiter in $waiters; do 
        wait $waiter
    done
}


function stop_detached {
    
    pid_file="$testpath/pids"
    pids=$(cat $pid_file)

    kill_processes

    rm -f $pid_file
}

function start_anvil {

    echo "Starting anvil server ....."
    anvil --host 0.0.0.0 > /dev/null & 
    anvil_pid=$! 
    echo "Anvil server started ....."

    echo $anvil_pid > ./anvil.pid

}

function stop_anvil {

    pid_file="./anvil.pid"
    anvil_pid=$(cat $pid_file)

    kill $anvil_pid

    rm -f $pid_file
}

function start_graph {

    echo "Starting graph node ....."
    pushd ./thegraph
        docker compose up -d
    popd
     ./wait-for http://0.0.0.0:8000 -- echo "GraphQL up"

    echo "Graph node started ....."
}

function stop_graph {

    pushd ./thegraph
        docker compose down -v
    popd
}

testpath=$(ls -td ./testdata/*/ | head -1)

case "$1" in
    help)
        cat <<-EOF
        Binary experiment tool
EOF
        ;;
    start)
        start_trap ${@:2} ;;   
    start-detached)
        start_detached ${@:2} ;;    
    stop)
        stop_detached ${@:2} ;;
    start-anvil)
        start_anvil ${@:2} ;;   
    stop-anvil)
        stop_anvil ${@:2} ;;    
    start-graph)
        start_graph ${@:2} ;;   
    stop-graph)
        stop_graph ${@:2} ;;  
    *)
esac
