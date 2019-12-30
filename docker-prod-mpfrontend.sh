#!/bin/bash

TAG=latest
if [ "${TRAVIS_TAG}" != "" ]; then
  TAG=${TRAVIS_TAG}
fi
docker build --no-cache -t aouyang1/mpfrontend:${TAG} -f mpfrontend/Dockerfile_prod mpfrontend/
