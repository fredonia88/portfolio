#!/bin/bash

sudo yum clean all
sudo rm -f /var/log/*.log*
sudo rm -rf /tmp/*

# Check for --prune-volumes argument
if [[ "$1" == "--prune-volumes" ]]; then
    echo "Pruning Docker system, including volumes..."
    docker system prune -a --volumes
else
    echo "Pruning Docker system without volumes..."
    docker system prune -a
fi