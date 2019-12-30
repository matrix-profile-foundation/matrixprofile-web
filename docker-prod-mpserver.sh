#!/bin/bash

TAG=latest
if [ "${TRAVIS_TAG}" != "" ]; then
  TAG=${TRAVIS_TAG}
fi
docker build --no-cache -t aouyang1/mpserver:${TAG} -f mpserver/Dockerfile_prod mpserver/
