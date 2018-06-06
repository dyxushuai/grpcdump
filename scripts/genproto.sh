#!/usr/bin/env bash
PbDepsPath=$GOPATH/src
# gen corepb
PbPath=./grpc_example/helloworld/helloworld
PbFile=helloworld.proto
protoc  -I $PbDepsPath \
        -I $PbPath \
        --go_out=plugins=grpc:$PbPath $PbPath/$PbFile