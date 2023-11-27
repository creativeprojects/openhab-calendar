GOCMD=env go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
GOGET=$(GOCMD) get

BINARY=openhab-calendar
CONFIG=config*.json
DEPLOY_TEST_SERVER=
DEPLOY_PROD_SERVER=
DEPLOY_DIR=~/openhab-calendar/
BUILD_PROD=./build
TESTS=./...
COVERAGE_FILE=coverage.out

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: all test build coverage clean build-linux deploy-test deploy-prod deploy-config

all: test build

build:
		$(GOBUILD) -o $(BINARY) -v

test:
		$(GOTEST) -v $(TESTS)

coverage:
		$(GOTEST) -coverprofile=$(COVERAGE_FILE) $(TESTS)
		$(GOTOOL) cover -html=$(COVERAGE_FILE)

clean:
		$(GOCLEAN)
		rm -f $(BINARY) $(COVERAGE_FILE) ${BUILD_PROD}/$(BINARY)

build-linux:
		GOOS="linux" GOARCH="amd64" $(GOBUILD) -o ${BUILD_PROD}/$(BINARY) -v

deploy-test: build-linux
		rsync -av ${BUILD_PROD}/$(BINARY) $(DEPLOY_TEST_SERVER):$(DEPLOY_DIR)
		ssh $(DEPLOY_TEST_SERVER) "sudo systemctl restart $(BINARY)"

deploy-prod: build-linux
		rsync -av ${BUILD_PROD}/$(BINARY) $(DEPLOY_PROD_SERVER):$(DEPLOY_DIR)

deploy-config:
		# rsync -av $(CONFIG) $(DEPLOY_TEST_SERVER):$(DEPLOY_DIR)
		rsync -av $(CONFIG) $(DEPLOY_PROD_SERVER):$(DEPLOY_DIR)
