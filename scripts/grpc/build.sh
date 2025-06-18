#!/bin/bash
# Author: Pawananjani Kumar (pawanannjani.kumar@gmail.com)
# CreatedAt: 28 Mar 2024
echo "Env: $GRPC_IMAGE_TAG"
#pwd
#cat < Dockerfile_grpc_dev

if ! docker build -f Dockerfile_grpc_dev_builder --network=host -t grpc_image:v1 .  ; then
    echo "GRPC image build failed."
    exit 1
  else
    echo "GRPC image build successful."
  fi