#!/usr/bin/env bash

SCRIPT_PATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
FULL_PATH="$SCRIPT_PATH/smpp-gateway"

echo $FULL_PATH
command() {
  $FULL_PATH >output.log 2>error.log &
  echo $! >> smpp-gateway.pid
}

command
