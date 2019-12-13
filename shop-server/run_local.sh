#!/usr/bin/env bash
export KV_VIPER_FILE="config-local.yaml"
go run ./cmd/shop-server/ migrate up
go run ./cmd/shop-server/ run