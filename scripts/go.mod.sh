#!/bin/bash

DIR=$(pwd)

cd ../

pkg=$(grep './' go.work | awk '{print $1}')

ROOT=$(pwd)

gotest() {
  cd $p

  go mod download "$@"

  cd $ROOT
}

go mod download "$@"

for p in $pkg; do
  if [ -d "$p" ]; then
    gotest "$@" &
  fi
done

wait