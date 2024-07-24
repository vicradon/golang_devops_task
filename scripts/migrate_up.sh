#!/bin/bash

set -e

DB_URL="postgres://postgres:password@localhost:5432/golang_internet_clipboard?sslmode=disable"

goose -dir db/migrations postgres "$DB_URL" up
