#!/bin/sh

set -e

echo "creating docker images with build time variables"
docker-compose build --build-arg VERSION=$(git describe --always --dirty --tags --long) --build-arg CURRENT_TIME=$(date +"%y-%m-%d")

echo "running docker containers in same network"
docker compose up

