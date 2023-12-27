#!/bin/bash

go mod tidy
go build -o ../node/plugin/bin/nodeplugin ../node/plugin/cmd
../node/plugin/bin/nodeplugin
 

