#!/bin/bash

sudo yum clean all
sudo rm -f /var/log/*.log*
sudo rm -rf /tmp/*
docker system prune -a --volumes