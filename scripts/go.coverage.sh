#!/bin/bash

ROOT="$(pwd)"

COVERAGE_TEST="-cover -covermode=atomic -test.gocoverdir=$ROOT/coverage"

mkdir -p $ROOT/coverage

./scripts/go.test.sh $COVERAGE_TEST "$@"

go tool covdata textfmt -i="$ROOT/coverage" -o "$ROOT/coverage.out"

rm $ROOT/coverage/*