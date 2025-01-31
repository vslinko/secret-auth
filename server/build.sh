#!/bin/bash

set -e

GOOS=linux GOARCH=amd64 go build -o server
docker build . --platform linux/amd64 --tag docker-registry.vslinko.xyz/vslinko/secret-auth-server:latest
docker push docker-registry.vslinko.xyz/vslinko/secret-auth-server:latest
