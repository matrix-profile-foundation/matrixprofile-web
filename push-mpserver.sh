#!/bin/bash

TAG=latest
if [ "${TRAVIS_TAG}" != "" ]; then
  TAG=${TRAVIS_TAG}
fi
docker push aouyang1/mpserver:${TAG}
