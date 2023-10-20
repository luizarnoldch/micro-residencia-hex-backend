#!/bin/bash

FOLDERS=($(ls lambdas/*/))

export GOOS="linux"
export GOARCH="amd64"
export CGO_ENABLED="0"

build_lambda() {
  for folder in "${FOLDERS[@]}"; do
    (
      cd "lambdas/$folder" || exit
      go build -o bootstrap -tags lambda.norpc
      zip ../../bin/${folder}.zip bootstrap
      rm -rf bootstrap
    )
  done
}

build_lambda