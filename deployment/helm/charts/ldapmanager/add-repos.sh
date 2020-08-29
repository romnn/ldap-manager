#!/bin/bash
set -e

export DIR=$(dirname $0)

# openldap chart is in the @stable repository
helm repo add stable https://kubernetes-charts.storage.googleapis.com/

echo "Updating dependencies for $(realpath $DIR)..."
helm dependency update $(realpath $DIR)