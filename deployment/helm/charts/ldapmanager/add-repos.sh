#!/bin/bash
set -e

export DIR=$(dirname $0)

# openldap chart is in the @stable repository
helm repo add stable https://charts.helm.sh/stable

echo "Updating dependencies for $(realpath $DIR)..."
helm dependency update $(realpath $DIR)
