#!/bin/bash -e

REPO_PATH="registrator"

export GOPATH=${PWD}/gopath

rm -f $GOPATH/src/${REPO_PATH}
mkdir -p $GOPATH/src
ln -s ${PWD} $GOPATH/src/${REPO_PATH}

eval $(go env)

go get $(go list -f "{{range .Imports}}{{ .  }} {{end}}")

CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "-s" -o bin/registrator ${REPO_PATH}
