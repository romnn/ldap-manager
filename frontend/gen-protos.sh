#!/bin/bash

set -e

export DIR=$(dirname $0)

cd $DIR
cmd="$@"

if [ -z "$SKIPPROTOCOMPILATION" ]
then
    echo "Installing compiler"
    pip install "grpc_web_proto_compile>=1.1.0"
    echo "Compiling protos"
    grpc_web_proto_compile -clear ../ ./src/generated
fi

exec $cmd