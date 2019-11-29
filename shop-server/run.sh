#!/usr/bin/env bash
export KV_VIPER_FILE="config.yaml"
shop-server migrate up
shop-server run