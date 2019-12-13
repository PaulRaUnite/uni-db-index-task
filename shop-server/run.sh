#!/usr/bin/env bash
export KV_VIPER_FILE="config.yaml"
shop-server migrate up

_term() {
  echo "Caught SIGTERM signal!"
  kill -TERM "$child" 2>/dev/null
}

trap _term SIGTERM

(shop-server run )&

child=$!
wait "$child"