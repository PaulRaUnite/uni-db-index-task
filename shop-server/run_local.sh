#!/usr/bin/env bash
export KV_VIPER_FILE="config-local.yaml"
go run ./cmd/shop-server/ "$@"