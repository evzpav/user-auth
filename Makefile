include help.mk

.PHONY: run-local git-config version clean install lint env env-stop test cover build image tag push deploy run run-docker remove-docker
.DEFAULT_GOAL := help

GITLAB_GROUP   	= evzpav
REGISTRY      	= myregistry
REGISTRY_GROUP 	= application

BUILD         	= $(shell git rev-parse --short HEAD)
DATE          	= $(shell date -uIseconds)
VERSION  	  	= $(shell git describe --always --tags)
NAME           	= $(shell basename $(CURDIR))
IMAGE          	= $(REGISTRY)/$(REGISTRY_GROUP)/$(NAME):$(BUILD)

MYSQL_NAME				= mysqldb_$(NAME)$(PIPELINE_ID)
NETWORK_NAME			= network_$(NAME)$(PIPELINE_ID)
NGINX_NAME				= nginx_$(NAME)$(PIPELINE_ID)
ACCEPTANCE_NETWORK		= acceptance_$(NETWORK_NAME)
ACCEPTANCE_MONGO		= acceptance_$(MYSQL_NAME)
ACCEPTANCE_APP_NAME		= app_$(NAME)$(PIPELINE_ID)
ACCEPTANCE_TESTS_IMAGE 	= acceptance_tests_$(NAME)$(PIPELINE_ID)

MYSQL_PASSWORD = mysqlpassword

git-config:
	git config --replace-all core.hooksPath .githooks

check-env-%:
	@ if [ "${${*}}" = ""  ]; then \
		echo "Variable '$*' not set"; \
		exit 1; \
	fi

version: ##@other Check version.
	@echo $(VERSION)

clean: ##@dev Remove folder vendor, public and coverage.
	rm -rf vendor public coverage

install: clean ##@dev Download dependencies via go mod.
	GO111MODULE=on go mod download
	GO111MODULE=on go mod vendor

lint: ##@check Run lint on docker.
	DOCKER_BUILDKIT=1 \
	docker build --progress=plain \
		--target=lint \
		--file=./build/package/dockerfile-lint .

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

# env: ##@environment Create network and run mysql container.
# 	# - docker network create $(NETWORK_NAME)	
# 	DOCKER_BUILDKIT=1 \
# 	docker run --rm -d \
# 		--name $(MYSQL_NAME) \
# 		-p 3306:3306 \
# 		-v $(HOME)/Documents/mysqldata:/var/lib/mysqldata/data \
# 		-e MYSQL_ROOT_PASSWORD=$(MYSQL_PASSWORD) \
# 		mysql:8.0


env-ip: ##@environment Return local MongoDB IP (from Docker container)
	@echo $$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(MYSQL_NAME))

env-stop: ##@environment Remove mongo container and remove network.
	-docker rm -vf $(MYSQL_NAME)
	# -docker network rm $(NETWORK_NAME)

test: ##@check Run tests and coverage.
	docker build --progress=plain \
		--network $(NETWORK_NAME) \
		--tag $(IMAGE) \
		--build-arg MONGO_URL=mongodb://${MYSQL_NAME}:27017 \
		--target=test \
		--file=./build/package/dockerfile-test .

	-mkdir coverage
	docker create --name $(NAME)-$(BUILD) $(IMAGE)
	docker cp $(NAME)-$(BUILD):/index.html ./coverage/.
	docker rm -vf $(NAME)-$(BUILD)

build: ##@build Build image.
	DOCKER_BUILDKIT=1 \
	docker build --progress=plain \
		--tag $(IMAGE) \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD=$(BUILD) \
		--build-arg DATE=$(DATE) \
		--target=build \
		--file=./build/package/dockerfile-build .

image: check-env-VERSION ##@build Create release docker image.
	DOCKER_BUILDKIT=1 \
	docker build --progress=plain \
		--tag $(IMAGE) \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD=$(BUILD) \
		--build-arg DATE=$(DATE) \
		--target=image \
		--file=./build/package/dockerfile-build .

tag: check-env-VERSION ##@build Add docker tag.
	docker tag $(IMAGE) \
		$(REGISTRY)/$(REGISTRY_GROUP)/$(NAME):$(VERSION)

push: check-env-VERSION ##@build Push docker image to registry.
	docker push $(REGISTRY)/$(REGISTRY_GROUP)/$(NAME):$(VERSION)

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
	MYSQL_URL=root:$(MYSQL_PASSWORD)@\(localhost:3306\)/user_auth?charset=utf8 \
	./user-auth

run-docker: remove-docker ##@docker Run docker container.
	docker run \
		--name $(NAME) \
		-p 8080:80 \
		$(IMAGE)

remove-docker:  ##@docker Remove docker container.
	-docker rm -vf $(NAME)

create-diagrams: ##@other Create diagrams.
	docker run --rm \
		-v $(shell pwd):/data \
		$(REGISTRY)/cli/plantuml:latest \
		sh /data/scripts/create-diagrams.sh

run-swagger:  ##@other Run Swagger OpenAPI doc server.
	docker run \
		-p 1234:8080 \
		-e BASE_URL=/ \
		-e SWAGGER_JSON=/api/api.yaml \
		-v ${PWD}/api:/api \
		swaggerapi/swagger-ui
	$(info Swagger docs running on http://localhost:1234)


env-acceptance-tests: ##@acceptance setup all acceptance tests environment
	docker network create $(ACCEPTANCE_NETWORK)

	./scripts/start_mongo.sh $(ACCEPTANCE_MONGO) $(ACCEPTANCE_NETWORK)

	docker run -d \
		--name $(NGINX_NAME) \
		--network $(ACCEPTANCE_NETWORK) \
		-v ${PWD}/test/acceptance/nginx/config/conf.d:/etc/nginx/conf.d \
		nginx:stable-alpine

	docker cp ${PWD}/test/acceptance/resources/ $(NGINX_NAME):/etc/nginx/

	docker run -d \
		--name $(ACCEPTANCE_APP_NAME) \
		--network $(ACCEPTANCE_NETWORK) \
 		-e LOGGER_LEVEL=debug \
		-e DATA_MENU_PORT=5003 \
 		-e MONGO_URL=mongodb://$(ACCEPTANCE_MONGO):27017 \
		$(IMAGE)

stop-env-acceptance-tests: ##@acceptance stop all acceptance tests environment
	- docker rm -vf $(NGINX_NAME)
	- docker rm -vf $(ACCEPTANCE_APP_NAME)
	- docker rm -vf $(ACCEPTANCE_MONGO)
	- docker network rm $(ACCEPTANCE_NETWORK)

do-acceptance-tests: ##@acceptance do acceptance tests
	docker build --progress=plain \
		--no-cache \
		--network $(ACCEPTANCE_NETWORK) \
		--tag $(ACCEPTANCE_TESTS_IMAGE) \
		--build-arg APP_HOST=$$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(ACCEPTANCE_APP_NAME)):5003 \
		--file=./test/acceptance/nginx/config/dockerfile-accept-tests .

acceptance-tests: stop-env-acceptance-tests ##@acceptance run acceptance tests locally
	make image
	make env-acceptance-tests
	make do-acceptance-tests
	make stop-env-acceptance-tests

