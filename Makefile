include help.mk
include .env

.PHONY: run-local git-config version clean install lint env env-stop test cover build image tag push deploy run run-docker remove-docker image
.DEFAULT_GOAL := help

BUILD         			= $(shell git rev-parse --short HEAD)
DATE          			= $(shell date -uIseconds)
VERSION  	  			= $(shell git describe --always --tags)
NAME           			= $(shell basename $(CURDIR))
IMAGE          			= $(NAME):$(BUILD)

MYSQL_NAME				= mysqldb_$(NAME)$(PIPELINE_ID)
NETWORK_NAME			= network_$(NAME)$(PIPELINE_ID)
MYSQL_PASSWORD 			= mysqlpassword

git-config:
	git config --replace-all core.hooksPath .githooks

check-env-%:
	@ if [ "${${*}}" = ""  ]; then \
		echo "Variable '$*' not set"; \
		exit 1; \
	fi

version: ##@other Check version.
	@echo $(VERSION)


build-mysql: ##@mysql build mysql docker image
	DOCKER_BUILDKIT=1 \
	docker build \
	--progress=plain \
	-t mysql_$(NAME):$(VERSION) \
	-f ./docker/mysql/Dockerfile \
	./docker/mysql/

run-mysql: build-mysql  ##@mysql run mysql on docker
	DOCKER_BUILDKIT=1 \
	docker run --rm -d \
		-v $(HOME)/Documents/mysqldata:/var/lib/mysqldata/data \
		-p 3306:3306 \
		--name mysql_$(NAME) \
		-e MYSQL_ROOT_PASSWORD=$(MYSQL_PASSWORD) \
		mysql_$(NAME):$(VERSION)

env-ip: ##@environment Return local MongoDB IP (from Docker container)
	@echo $$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(MYSQL_NAME))

env-stop: ##@environment Remove mongo container and remove network.
	-docker rm -vf $(MYSQL_NAME)

build-local: ##@dev Build binary locally
	-rm ./user-auth

	CGO_ENABLED=0 \
	GOOS=linux  \
	GOARCH=amd64  \
	go build -installsuffix cgo -o user-auth -ldflags \
	"-X main.version=${VERSION} -X main.build=${BUILD} -X main.date=${DATE}" \
	./cmd/server/main.go


run-local: build-local ##@dev Run locally.
	HOST=localhost \
	PORT=5001 \
	LOGGER_LEVEL=debug \
	EMAIL_FROM=$(EMAIL_FROM) \
	EMAIL_PASSWORD=$(EMAIL_PASSWORD) \
	GOOGLE_KEY=$(GOOGLE_KEY) \
	GOOGLE_SECRET=$(GOOGLE_SECRET) \
	MYSQL_URL=root:$(MYSQL_PASSWORD)@\(localhost:3306\)/user_auth?charset=utf8 \
	./user-auth

target: 
	DOCKER_BUILDKIT=1 \
	docker build --progress=plain \
		--tag $(IMAGE) \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD=$(BUILD) \
		--build-arg DATE=$(DATE) \
		--target=$(TARGET) \
		--file= .

build: ##@build Build image.
	make target TARGET=build

run-docker: build ##@docker Run docker container.
	docker run --rm \
		--name $(NAME) \
		--network=host \
		-e HOST=localhost \
		-e PORT=8080 \
		-e LOGGER_LEVEL=debug \
		-e EMAIL_FROM=$(EMAIL_FROM) \
		-e EMAIL_PASSWORD=$(EMAIL_PASSWORD) \
		-e GOOGLE_KEY=$(GOOGLE_KEY) \
		-e GOOGLE_SECRET=$(GOOGLE_SECRET) \
		-e MYSQL_URL=root:$(MYSQL_PASSWORD)@\(localhost:3306\)/user_auth?charset=utf8 \
		$(IMAGE)
