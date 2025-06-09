#!/bin/bash

# Prompt the user for the migration name
read -p "Enter migration name: " migration_name

# Run the goose command with the provided name
export PATH="$(go env GOPATH)/bin:$PATH"
cd cmd/db/migration
goose -s create "$migration_name"