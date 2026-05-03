#! /bin/bash
set -e

[ -z "$1" ] && echo "Usage: $0 project-name" && exit 1

TARGET_DIR="problems/$1"
mkdir -p "$TARGET_DIR"
touch "$TARGET_DIR/main.go"
