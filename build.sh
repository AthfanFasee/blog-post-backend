#!/bin/sh

set -e

echo "creating docker images"
docker-compose build --build-arg VERSION=$(git describe --always --dirty --tags --long)

echo "running docker containers in same network"
docker compose up

