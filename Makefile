usage:
	@echo "make build            : builds docker images for mpserver and mpfrontend"
	@echo "make build-mpserver   : builds docker image for mpserver"
	@echo "make build-mpfrontend : builds docker image from mpfrontend"
	@echo "make deploy           : runs mpserver, mpfrontend, and redis"
	@echo "make undeploy         : tears down mpserver, mpfronted, and redis"
	@echo "make push : pushes mpserver and mpfrontend images to dockerhub"
	@echo "make push-mpserver    : pushes the mpserver image to dockerhub"
	@echo "make push-mpfrontend  : pushes the mpfrontend image to dockerhub"

build: build-mpserver build-mpfrontend

build-mpserver:
	docker build --no-cache -t aouyang1/mpserver mpserver/

build-mpfrontend:
	docker build --no-cache -t aouyang1/mpfrontend -f mpfrontend/Dockerfile_prod mpfrontend/

build-mpfrontend-dev:
	docker build --no-cache -t aouyang1/mpfrontend mpfrontend/

deploy: undeploy
	docker-compose up --build -d

undeploy:
	docker-compose down

push: push-mpserver push-mpfrontend

push-mpserver:
	docker push aouyang1/mpserver

push-mpfrontend:
	docker push aouyang1/mpfrontend
