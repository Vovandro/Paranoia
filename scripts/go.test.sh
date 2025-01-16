#!/bin/bash

pkg=$(grep './' go.work | awk '{print $1}')

ROOT=$(pwd)

gotest() {
  cd $p

  go test ./... "$@"

  cd $ROOT
}

go test ./... "$@"

for p in $pkg; do
  if [ -d "$p" ]; then
    gotest "$@" &
  fi
done

wait