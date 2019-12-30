usage:
	@echo "make build            : builds binary runnable and compiled frontend for mpserver and mpfrontend"
	@echo "make build-mpserver   : builds binary for mpserver"
	@echo "make build-mpfrontend : builds compiled frontend from mpfrontend"
	@echo "make docker-dev       : build docker images for dev environment"
	@echo "make docker-prod      : build docker images for prod environment"
	@echo "make deploy           : runs mpserver, mpfrontend, and redis in dev environment"
	@echo "make undeploy         : tears down mpserver, mpfronted, and redis"
	@echo "make push             : pushes mpserver and mpfrontend images to dockerhub"
	@echo "make push-mpserver    : pushes the mpserver image to dockerhub"
	@echo "make push-mpfrontend  : pushes the mpfrontend image to dockerhub"

build: build-mpserver build-mpfrontend

build-mpserver:
	cd mpserver && GOOS=linux GOARCH=amd64 go build && cd ..

build-mpfrontend:
	cd mpfrontend && npm run build-dev

docker-dev:
	docker build -t aouyang1/mpserver:dev mpserver/
	docker build -t aouyang1/mpfrontend:dev mpfrontend/

docker-prod: docker-prod-mpserver docker-prod-mpfrontend

docker-prod-mpserver:
	./docker-prod-mpserver.sh

docker-prod-mpfrontend:
	./docker-prod-mpfrontend.sh

deploy: undeploy
	docker-compose up -d

undeploy:
	docker-compose down

push: push-mpserver push-mpfrontend

push-mpserver:
	./push-mpserver.sh

push-mpfrontend:
	./push-mpfrontend.sh
