#!/bin/bash

# This script waits for the database container to begin accepting connections,
# because we don't want the tests to fail just because the container wasn't
# fully started yet.

set -e

HOST=db
PORT=5432
TIMEOUT=5 # in seconds

echo -n "Waiting for TCP connection to $HOST:$PORT..."

try_count=0

while ! echo | nc -w 1 $HOST $PORT 2> /dev/null; do
  echo -n .
  sleep 1
  let "try_count = $try_count + 1"
  if [[ $try_count -ge $TIMEOUT ]]; then
    echo
    echo "Timed out after $TIMEOUT seconds"
    exit 1
  fi
done

echo 'ok'
