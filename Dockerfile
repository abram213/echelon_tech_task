# Start from golang base image
FROM golang:alpine as builder

ENV GO111MODULE=on

# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

RUN mkdir /app
WORKDIR /app
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Build the go app
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file, .env and views from the previous stage.
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/views ./views

EXPOSE 8080

ENTRYPOINT ./main