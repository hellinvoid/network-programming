#! /bin/bash
set -e

[ -z "$1" ] && echo "Usage: $0 project-name" && exit 1

TARGET_DIR="problems/$1"

mkdir -p "$TARGET_DIR/deployment"
cat <<EOF > "$TARGET_DIR/.dockerignore"
deployment
EOF

touch "$TARGET_DIR/main.go"

cat <<EOF > "$TARGET_DIR/deployment/Dockerfile"
FROM golang:1.25-trixie AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 6969

CMD ["./main"]
EOF

GITHUB_ACCOUNT="github.com/$(yq -r '.git.github.owner' config.yaml)"
GITHUB_REPO=$(yq -r '.git.github.repo' config.yaml)

go -C "$TARGET_DIR" mod init "$GITHUB_ACCOUNT/$GITHUB_REPO/$TARGET_DIR"