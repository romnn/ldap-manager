#!/bin/bash

set -e

export DIR=$(dirname $0)

cd $DIR
cmd="$@"

if [ -z "$SKIPPROTOCOMPILATION" ]
then
    echo "Installing compiler"
    pip install "grpc_web_proto_compile>=1.1.0"
    echo "Compiling protos (from $PWD)"
    grpc_web_proto_compile -clear --base_proto_parent_dir $(realpath $PWD/../) $(realpath $PWD/../) $(realpath $PWD/src/generated)
fi

exec $cmd