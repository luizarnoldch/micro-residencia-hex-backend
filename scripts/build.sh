#!/bin/bash

FOLDERS=($(ls -d lambdas/*/))

export GOOS="linux"
export GOARCH="amd64"
export CGO_ENABLED="0"

build_lambda() {
  for folder in "${FOLDERS[@]}"; do
    (
      folder_name=$(basename "${folder}")
      cd "lambdas/$folder_name" || exit
      go build -o bootstrap -tags lambda.norpc
      zip ../../bin/${folder_name}.zip bootstrap
      rm -rf bootstrap
    )
  done
}

build_lambda