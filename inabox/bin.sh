#!/bin/bash

# -d returns the directories, -t sorts by modification time (newest first), head -1 takes the first line
# TODO: we should probably create a symlink "current" to the current testpath instead of relying on modification time
testpath=$(ls -td ./testdata/*/ | head -1)

# All processes started will have their PIDs stored here,
# and dumped into $testpath/pids when started in detached mode,
# so that they can be killed.
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

    # Churner is now started as a Docker container in the integration test
    # Skipping churner binary start

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

    mkdir -p $testpath/logs

    # Churner is now started as a Docker container in the integration test
    # Skipping churner binary start

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

        ./wait-for 0.0.0.0:${RELAY_GRPC_PORT} -- echo "Relay up" &
        waiters="$waiters $!"
    done

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

        ./wait-for 0.0.0.0:${NODE_DISPERSAL_PORT} -- echo "Node up" &
        waiters="$waiters $!"
    done

    echo $pids > $pid_file

    for waiter in $waiters; do
        wait $waiter
    done
}

# Stops all detached processes started with start_detached
# Meaning it walks through the pids file and kills those processes.
function stop_detached {

    pid_file="$testpath/pids"
    
    # Try to read PIDs from file if it exists
    if [[ -f "$pid_file" ]]; then
        pids=$(cat $pid_file)
        kill_processes
        rm -f $pid_file
    fi

    # We also call force_stop because the PID file approach is finicky
    # and might not catch all processes if the file is deleted or corrupted.
    force_stop
}

function force_stop {
    echo "Force stopping all EigenDA processes..."
    pkill -9 -f "disperser/bin/server" || true
    pkill -9 -f "disperser/bin/encoder" || true
    pkill -9 -f "disperser/bin/batcher" || true
    pkill -9 -f "disperser/bin/controller" || true
    pkill -9 -f "relay/bin/relay" || true
    pkill -9 -f "node/bin/node" || true
    pkill -9 -f "retriever/bin/server" || true
    rm -f $testpath/pids
    echo "All processes force killed"
}

help() {
    echo "Usage: $0 {start|start-detached|stop-detached|force-stop}"
    echo ""
    echo "Commands:"
    echo "  start              Start all services in the foreground with trap on SIGINT"
    echo "  start-detached     Start all services in the background and log output to files"
    echo "  stop-detached      Stop all background services started with start-detached"
    echo "  force-stop         Force kill all EigenDA related processes"
    echo ""
    echo "Logs are stored in $testpath/logs/"
    echo "PIDs of detached processes are stored in $testpath/pids"
}

case "$1" in
    start)
        start_trap ${@:2} ;;
    start-detached)
        start_detached ${@:2} ;;
    stop-detached)
        stop_detached ${@:2} ;;
    force-stop)
        force_stop ${@:2} ;;
    help)
        help;;
    *)
        help;;
esac
