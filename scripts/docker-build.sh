#!/bin/bash
set -e

[ -z "$1" ] && echo "Usage: $0 <problem-dir-relative-path>" && exit 1

PROJECT_ROOT=$(pwd)
PROBLEM_PATH="$1"

PROJECT_NAME=$(basename "$PROBLEM_PATH")
DOCKERHUB_OWNER=$(yq -r '.docker.dockerhub.owner' config.yaml)

docker build \
  -f deployment/Dockerfile \
  -t "$DOCKERHUB_OWNER/$PROJECT_NAME" \
  --build-arg PROBLEM_PATH="$PROBLEM_PATH" \
  "$PROJECT_ROOT"

docker push "$DOCKERHUB_OWNER/$PROJECT_NAME:latest"