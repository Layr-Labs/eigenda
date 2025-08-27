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

    # Check if churner is provided by testinfra
    if [ -z "$CHURNER_URL" ]; then
        set -a
        source $testpath/envs/churner.env
        set +a
        ../operators/churner/bin/server &

        pid="$!"
        pids="$pids $pid"
    else
        echo "Using churner from testinfra at $CHURNER_URL"
    fi

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

    for FILE in $(ls $testpath/envs/controller*.env); do
        set -a
        source $FILE
        set +a
        ../disperser/bin/controller > $testpath/logs/controller.log 2>&1 &

        pid="$!"
        pids="$pids $pid"
    done


    files=($(ls $testpath/envs/relay*.env))
    for i in "${!files[@]}"; do
        FILE=${files[$i]}
        set -a
        source $FILE
        set +a
        ../relay/bin/relay &

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

    echo "DEBUG: Starting detached mode with testpath: $testpath"
    echo "DEBUG: Process ID file will be: $pid_file"
    
    mkdir -p $testpath/logs

    # Check if churner is provided by testinfra
    if [ -n "$CHURNER_URL" ]; then
        echo "DEBUG: Using churner from testinfra at $CHURNER_URL"
    else
        echo "DEBUG: Starting churner service"
        if [ ! -f "../operators/churner/bin/server" ]; then
            echo "ERROR: Churner binary not found at ../operators/churner/bin/server"
            echo "DEBUG: Current directory: $(pwd)"
            echo "DEBUG: Available files in ../operators/churner/bin/: $(ls -la ../operators/churner/bin/ 2>/dev/null || echo 'directory not found')"
            exit 1
        fi
        set -a
        source $testpath/envs/churner.env
        set +a
        echo "DEBUG: Churner will listen on port: ${CHURNER_GRPC_PORT}"
        ../operators/churner/bin/server > $testpath/logs/churner.log 2>&1 &

        pid="$!"
        pids="$pids $pid"

        echo "DEBUG: Starting wait-for churner on port ${CHURNER_GRPC_PORT}"
        ./wait-for 0.0.0.0:${CHURNER_GRPC_PORT} -t 30 -- echo "Churner up" &
        waiters="$waiters $!"
    fi

    for FILE in $(ls $testpath/envs/dis*.env); do
        set -a
        source $FILE
        set +a
        id=$(basename $FILE | tr -d -c 0-9)
        ../disperser/bin/server > $testpath/logs/dis${id}.log 2>&1 &

        pid="$!"
        pids="$pids $pid"

        echo "DEBUG: Starting wait-for disperser on port ${DISPERSER_SERVER_GRPC_PORT}"
        ./wait-for 0.0.0.0:${DISPERSER_SERVER_GRPC_PORT} -t 30 -- echo "Disperser up" &
        waiters="$waiters $!"
    done

    # Check if encoder is provided by testinfra
    if [ -n "$ENCODER_URL" ]; then
        echo "DEBUG: Using encoder from testinfra at $ENCODER_URL"
    else
        for FILE in $(ls $testpath/envs/enc*.env); do
            set -a
            source $FILE
            set +a
            id=$(basename $FILE | tr -d -c 0-9)
            ../disperser/bin/encoder > $testpath/logs/enc${id}.log 2>&1 &

            pid="$!"
            pids="$pids $pid"

            echo "DEBUG: Starting wait-for encoder on port ${DISPERSER_ENCODER_GRPC_PORT}"
            ./wait-for 0.0.0.0:${DISPERSER_ENCODER_GRPC_PORT} -t 30 -- echo "Encoder up" &
            waiters="$waiters $!"
        done
    fi

    # Check if batcher is provided by testinfra
    if [ -n "$BATCHER_PROVIDED" ]; then
        echo "DEBUG: Using batcher from testinfra"
    else
        for FILE in $(ls $testpath/envs/batcher*.env); do
            set -a
            source $FILE
            set +a
            id=$(basename $FILE | tr -d -c 0-9)
            ../disperser/bin/batcher > $testpath/logs/batcher${id}.log 2>&1 &

            pid="$!"
            pids="$pids $pid"
        done
    fi

    for FILE in $(ls $testpath/envs/retriever*.env); do
        set -a
        source $FILE
        set +a
        ../retriever/bin/server > $testpath/logs/retriever.log 2>&1 &

        pid="$!"
        pids="$pids $pid"
    done

    for FILE in $(ls $testpath/envs/controller*.env); do
        set -a
        source $FILE
        set +a
        ../disperser/bin/controller > $testpath/logs/controller.log 2>&1 &

        pid="$!"
        pids="$pids $pid"
    done

    files=($(ls $testpath/envs/relay*.env))
    last_index=$(( ${#files[@]} - 1 ))
    for i in "${!files[@]}"; do
        FILE=${files[$i]}
        set -a
        source $FILE
        set +a
        id=$(basename $FILE | tr -d -c 0-9)
        ../relay/bin/relay > $testpath/logs/relay${id}.log 2>&1 &

        pid="$!"
        pids="$pids $pid"

        echo "DEBUG: Starting wait-for relay on port ${RELAY_GRPC_PORT}"
        ./wait-for 0.0.0.0:${RELAY_GRPC_PORT} -t 30 -- echo "Relay up" &
        waiters="$waiters $!"
    done

    # Check if operators are provided by testinfra
    if [ -n "$OPERATORS_PROVIDED" ]; then
        echo "DEBUG: Using operators from testinfra"
    else
        files=($(ls $testpath/envs/opr*.env))
        last_index=$(( ${#files[@]} - 1 ))

        for i in "${!files[@]}"; do
            if [ $i -eq $last_index ]; then
                sleep 10  # Sleep for 10 seconds before the last loop iteration
            fi
            FILE=${files[$i]}
            set -a
            source $FILE
            set +a
            id=$(basename $FILE | tr -d -c 0-9)
            ../node/bin/node > $testpath/logs/opr${id}.log 2>&1 &

            pid="$!"
            pids="$pids $pid"

            echo "DEBUG: Starting wait-for node on port ${NODE_DISPERSAL_PORT}"
            ./wait-for 0.0.0.0:${NODE_DISPERSAL_PORT} -t 30 -- echo "Node up" &
            waiters="$waiters $!"
        done
    fi

    echo $pids > $pid_file

    echo "DEBUG: Waiting for $(echo $waiters | wc -w) services to become available..."
    waiter_count=0
    for waiter in $waiters; do
        waiter_count=$((waiter_count + 1))
        echo "DEBUG: Waiting for service $waiter_count to be ready (PID: $waiter)"
        if wait $waiter; then
            echo "DEBUG: Service $waiter_count is ready"
        else
            echo "ERROR: Service $waiter_count failed to start (PID: $waiter)"
        fi
    done
    echo "DEBUG: All services startup checks completed"
}


function stop_detached {

    pid_file="$testpath/pids"
    if [[ -f "$pid_file" ]]; then
        pids=$(cat $pid_file)
        kill_processes
        rm -f $pid_file
    else
        echo "No PID file found, attempting fallback cleanup..."
        fallback_cleanup
    fi
}

function fallback_cleanup {
    echo "Performing fallback process cleanup..."
    
    # Kill processes by name patterns that are specific to EigenDA
    pkill -f "disperser/bin/server" 2>/dev/null || true
    pkill -f "disperser/bin/encoder" 2>/dev/null || true
    pkill -f "disperser/bin/batcher" 2>/dev/null || true
    pkill -f "disperser/bin/controller" 2>/dev/null || true
    pkill -f "node/bin/node" 2>/dev/null || true
    pkill -f "relay/bin/relay" 2>/dev/null || true
    pkill -f "retriever/bin/server" 2>/dev/null || true
    pkill -f "churner/bin/server" 2>/dev/null || true
    
    # Wait a moment for graceful shutdown
    sleep 2
    
    # Force kill any remaining processes
    pkill -9 -f "disperser/bin/server" 2>/dev/null || true
    pkill -9 -f "disperser/bin/encoder" 2>/dev/null || true
    pkill -9 -f "disperser/bin/batcher" 2>/dev/null || true
    pkill -9 -f "disperser/bin/controller" 2>/dev/null || true
    pkill -9 -f "node/bin/node" 2>/dev/null || true
    pkill -9 -f "relay/bin/relay" 2>/dev/null || true
    pkill -9 -f "retriever/bin/server" 2>/dev/null || true
    pkill -9 -f "churner/bin/server" 2>/dev/null || true
    
    echo "Fallback cleanup completed"
}

function force_stop {
    echo "Force stopping all EigenDA processes..."
    fallback_cleanup
    
    # Also clean up any remaining PID files
    find ./testdata -name "pids" -type f -delete 2>/dev/null || true
    rm -f ./anvil.pid 2>/dev/null || true
    
    echo "Force stop completed"
}

function start_anvil {

    echo "Starting anvil server ....."
    anvil --host 0.0.0.0 > /dev/null &
    anvil_pid=$!
    
    # Wait for anvil to be ready
    ./wait-for 0.0.0.0:8545 -- echo "Anvil ready"
    
    if [ $? -ne 0 ]; then
        echo "Failed to start anvil server"
        exit 1
    fi
    
    echo "Anvil server started ....."

    echo $anvil_pid > ./anvil.pid

}

function stop_anvil {

    pid_file="./anvil.pid"
    if [[ -f "$pid_file" ]]; then
        anvil_pid=$(cat $pid_file)
        if [[ -n "$anvil_pid" ]]; then
            kill $anvil_pid 2>/dev/null || true
        fi
        rm -f $pid_file
    else
        # Fallback: try to find and kill anvil process
        pkill -f "anvil --host 0.0.0.0" 2>/dev/null || true
    fi
}

function start_graph {

    echo "Starting graph node ....."
    pushd ./thegraph
        docker compose up -d
    popd
     ./wait-for http://0.0.0.0:8000 -- echo "GraphQL up"

     if [ $? -ne 0 ]; then
        echo "Failed to start graph node"
        exit 1
     fi

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
        
        Available commands:
        - start: Start all services with trap (interactive)
        - start-detached: Start all services in detached mode
        - stop: Stop all services gracefully
        - force-stop: Force kill all EigenDA processes
        - start-anvil: Start Anvil blockchain
        - stop-anvil: Stop Anvil blockchain
        - start-graph: Start Graph node
        - stop-graph: Stop Graph node
EOF
        ;;
    start)
        start_trap ${@:2} ;;
    start-detached)
        start_detached ${@:2} ;;
    stop)
        stop_detached ${@:2} ;;
    force-stop)
        force_stop ${@:2} ;;
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
