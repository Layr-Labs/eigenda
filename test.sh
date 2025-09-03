#!/bin/bash

result=0

go clean -testcache
# expand all arguments to script at the end of this line
CI=true go test -short ./... "$@"
result=$?
