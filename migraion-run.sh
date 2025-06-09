#!/bin/bash

# Load environment variables from .env file
set -a
source .env
set +a

go run cmd/db/main.go