#!/bin/bash

set -e

GOOS=linux GOARCH=amd64 go build -o server
docker build . --platform linux/amd64 --tag vslinko/secret-auth-server:latest
docker push vslinko/secret-auth-server:latest
