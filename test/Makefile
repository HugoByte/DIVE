PWD = $(abspath ./)

DOCKER_CONTAINER_ID=$(shell docker ps -aqf "name=eth-contract-deployer")

JAVASCORE_PATH = $(PWD)/../javascore
GRADLE = ./gradlew

HARDHAT = npx hardhat
TS_NODE = npx ts-node
NETWORK = --network localnet
SOLIDITY_CONTRACTS = ./solidity/contracts

.DEFAULT_GOAL := all
all:
	@echo $(PWD)

.PHONY: build clean deploy

run-demo:
	mkdir -p solidity/build/
	npm install
	docker cp -a $(DOCKER_CONTAINER_ID):/static-files/contracts/. solidity/
	@ echo ">>> Run demo_dapp.ts " ; \
	$(HARDHAT) $(NETWORK) run scripts/demo_dapp.ts
