#!/bin/bash

if [ -d "./certs" ]; then
  echo "Installing certificates ..."
  for cert in $(ls ./certs/*.cert)
    do
      /opt/cprocsp/bin/amd64/certmgr  -inst -store uMy -silent -file $cert
    done
fi

echo "Starting mainnet CSP miner ..."
# Execute the main binary and pass all script arguments
exec /mir/build/bin/mir "$@"
