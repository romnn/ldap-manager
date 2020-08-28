#!/bin/bash

# openldap chart is in the @stable repository
helm repo add stable https://kubernetes-charts.storage.googleapis.com/
helm dependency update