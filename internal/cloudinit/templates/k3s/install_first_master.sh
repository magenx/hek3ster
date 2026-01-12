#!/bin/bash
# K3s installation script for first master node
# This script is executed on the first master to initialize the k3s cluster

curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION='{{ .K3sVersion }}' K3S_TOKEN='{{ .K3sToken }}' sh -s - server {{ .BaseArgs }}
