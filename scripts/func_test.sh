#!/bin/bash

LAMBDAS=("hello")

test_lambda() {
  for folder in "${LAMBDAS[@]}"; do
    (
      aws lambda invoke --function-name Cognito-Test-${folder} --payload file://events/${folder}/request/input.json --cli-binary-format raw-in-base64-out ./events/${folder}/response/output.json
      echo -e "\n"
      cat ./events/${folder}/response/output.json
      echo -e "\n"
    )
  done
}

test_lambda