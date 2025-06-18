#!/bin/bash
# Author: Pawananjani Kumar (pawanannjani.kumar@gmail.com)
# CreatedAt: 28 Mar 2024

echo "Env: $GRPC_APP_TAG"

if ! docker run --name grpc_container --env-file envinronments.env --network=host -d -p 127.0.0.1:8080:8080 grpc_image:v1; then
    echo "GRPC container run failed."
    exit 1
  else
    echo "GRPC container run successful."
  fi