#!/bin/bash
# Author: Pawananjani Kumar (pawanannjani.kumar@gmail.com)
# CreatedAt: 28 Mar 2024

echo "Env: $GRPC_APP_TAG"

if ! docker run --name gin_container --env-file envinronments.env --network=host -d -p 127.0.0.1:8081:8081 gin_image:v1; then
    echo "GIN container run failed."
    exit 1
  else
    echo "GIN container run successful."
  fi