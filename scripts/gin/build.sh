#!/bin/bash
# Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
# CreatedAt: 28 Mar 2024
echo "Env: $TAG"

#if ! go build -tags "$TAG" -o ./src/build/gin/app ./src/cmd/gin/main.go ; then
#    echo "Build failed."
#    exit 1
#  else
#    echo "Build successful."
#  fi

if ! docker build -f Dockerfile_gin_dev --network=host -t gin_image:v1 .  ; then
    echo "GIN image build failed."
    exit 1
  else
    echo "GIN image build successful."
  fi