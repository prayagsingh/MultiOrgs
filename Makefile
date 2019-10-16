.PHONY: all dev clean build env-up env-down run

all: clean build env-up run

dev: build run

##### BUILD
build:
	@echo "Build ..."
	@go build
	@echo "Build done"

##### ENV
env-up:
	@echo "Start environment ..."
	@cd networks && docker-compose up --force-recreate -d
	@echo "checking present working directory ..."
	@pwd
	@echo "Environment up"

env-down:
	@echo "Stop environment ..."
	@cd networks && docker-compose down
	@echo "Environment down"

##### RUN
run:
	@echo "Start app ..."
	@./MultiOrgs

##### CLEAN
clean: env-down
	@echo "Clean up ..."
	@echo "cleaning identities from wallet .."
	@rm -f /c/Projects/Go/src/github.com/MultiOrgs/wallet/org1/cert/*
	@rm -rf /c/Projects/Go/src/github.com/MultiOrgs/wallet/org1/key/*
	@rm -f /c/Projects/Go/src/github.com/MultiOrgs/wallet/org2/cert/* 
	@rm -rf /c/Projects/Go/src/github.com/MultiOrgs/wallet/org2/key/* 
	@rm -f ${GOPATH}/src/github.com/MultiOrgs/wallet/org1/cert/*
	@rm -rf ${GOPATH}/src/github.com/MultiOrgs/wallet/org1/key/*
	@rm -f ${GOPATH}/src/github.com/MultiOrgs/wallet/org2/cert/* 
	@rm -rf ${GOPATH}/src/github.com/MultiOrgs/wallet/org2/key/*
	@echo "cleaning identities completed. Now check the folders"
	@tree /c/Projects/Go/src/github.com/MultiOrgs/wallet/
	@echo ""
	@tree ${GOPATH}/src/github.com/MultiOrgs/wallet/
	@rm -rf /tmp/MultiOrgs-* MultiOrgs
	@docker rm -f -v `docker ps -a --no-trunc | grep "dev-peer" | cut -d ' ' -f 1` 2>/dev/null || true
	@docker rmi `docker images --no-trunc | grep "MultiOrgs" | cut -d ' ' -f 1` 2>/dev/null || true
	@echo "checking CC docker images cleared or not"
	@docker images
	@echo "Clearing docker volume and network"
	@docker volume prune -f
	@docker network prune -f
	@docker system prune -f
	@echo "Clean up done"

