FROM golang:1.25-trixie AS builder

WORKDIR /app

COPY . .

ARG PROBLEM_PATH

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ${PROBLEM_PATH}

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 6969

CMD ["./main"]