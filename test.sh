#!/bin/bash

result=0
function tearDown {
    docker stop localstack-test
    exit $result
}
trap tearDown EXIT

go run ./inabox/deploy/cmd -localstack-port=4570 -deploy-resources=false localstack
go clean -testcache 
# expand all arguments to script at the end of this line
LOCALSTACK_PORT=4570 DEPLOY_LOCALSTACK=false go test -short ./... "$@"
result=$?
