#!/bin/bash
set -e

[ -z "$1" ] && echo "Usage: $0 project-root-dir" && exit 1

PROJECT_NAME=$(basename "$1")
DOCKERHUB_OWNER=$(yq -r '.docker.dockerhub.owner' config.yaml)

docker build -f "$1/deployment/Dockerfile" -t "$DOCKERHUB_OWNER/$PROJECT_NAME" "$1"