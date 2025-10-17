#!/bin/bash

pkg=$(grep './' go.work | awk '{print $1}')

ROOT=$(pwd)

goupdate() {
  cd $p

  go get -u ./... "$@"

  go mod tidy

  cd $ROOT
}

go get -u ./... "$@"

for p in $pkg; do
  if [ -d "$p" ]; then
    goupdate "$@" &
  fi
done

wait