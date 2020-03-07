[![Build Status](https://travis-ci.com/matrix-profile-foundation/matrixprofile-web.svg?branch=master)](https://travis-ci.com/matrix-profile-foundation/matrixprofile-web)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

# Requirements
* npm
* golang
* docker
* docker-compose

# Build and deploy locally
```sh
make build && make deploy
```

Go to `localhost:8080` in your browser

![screenshot](https://github.com/matrix-profile-foundation/matrix-profiles/blob/master/screenshot.png)

# Tear down
```sh
make undeploy
```
