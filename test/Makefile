PWD = $(abspath ./)

DOCKER_CONTAINER_ID=$(shell docker ps -aqf "name=eth-contract-deployer")

HARDHAT = npx hardhat
NETWORK = --network localnet
FOLDERS = solidity
DEP := solidity node_modules typechain-types

.DEFAULT_GOAL := run-demo

run-demo:$(FOLDERS)
	@ echo ">>> Run btp demo_dapp.ts " ; \
	$(HARDHAT) $(NETWORK) run scripts/demo_dapp.ts 

.PHONY: $(FOLDERS)

$(FOLDERS):
	@if [ -n "$(DOCKER_CONTAINER_ID)" ]; then \
		if [ ! -d "$@" ]; then \
			echo "Copying dependencies to local..." && \
			docker cp -a $(DOCKER_CONTAINER_ID):/static-files/contracts/. solidity/ && \
			docker cp -a $(DOCKER_CONTAINER_ID):/static-files/node_modules/. node_modules/ && \
			docker cp -a $(DOCKER_CONTAINER_ID):/static-files/typechain-types/. typechain-types/ ; \
		else \
			echo "Skipping to copy dependencies as it already exists locally."; \
		fi \
	else \
		npm install; \
	fi



.PHONY: clean $(DEP)
clean-dep:
	rm -rf $(DEP)


run-smoke: BINARY
	@echo "Running Smoke test cases"
	@cd $(CURDIR)/functional && ginkgo --focus="Smoke Tests" -v -p 

run-bridge: BINARY
	@echo "Running all bridge test cases"
	@cd $(CURDIR)/functional && ginkgo --focus="Bridge command Test" -v -p -timeout=3h

run-icon: BINARY
	@echo "Running icon chain commands test cases"
	@cd $(CURDIR)/functional && ginkgo --focus="Icon chain commands" -v -p -timeout=3h

run-eth: BINARY
	@echo "Running eth chain commands test cases"
	@cd $(CURDIR)/functional && ginkgo --focus="Eth chain commands" -v -p -timeout=3h

run-hardhat: BINARY
	@echo "Running hardhat chain commands test cases"
	@cd $(CURDIR)/functional && ginkgo --focus="Hardhat chain commands" -v -p -timeout=3h

run-archway: BINARY
	@echo "Running archway chain commands test cases"
	@cd $(CURDIR)/functional && ginkgo --focus="Archway chain commands" -v -p -timeout=3h

run-relaychain-kusama: BINARY
	@echo "Running relaychain chain commands test cases for kusama"
	@cd $(CURDIR)/functional && relayChain=kusama ginkgo --focus="Relaychain commands" -v -p -timeout=3h

run-relaychain-polkadot: BINARY
	@echo "Running relaychain chain commands test cases for polkadot"
	@cd $(CURDIR)/functional && relayChain=polkadot ginkgo --focus="Relaychain commands" -v -p -timeout=3h

run-relaychain-all: BINARY
	@echo "Running relaychain chain commands test cases for polkadot and kusama"
	@cd $(CURDIR)/functional && ginkgo --focus="Relaychain commands" -v -p -timeout=3h

run-parachain-kusama: BINARY
	@echo "Running parachain chain commands test cases for kusama"
	@cd $(CURDIR)/functional && relayChain=kusama ginkgo --focus="Parachain commands" -v -p -timeout=3h

run-parachain-polkadot: BINARY
	@echo "Running parachain chain commands test cases for polkadot"
	@cd $(CURDIR)/functional && relayChain=polkadot ginkgo --focus="Parachain commands" -v -p -timeout=3h

run-parachain-all: BINARY
	@echo "Running parachain chain commands test cases for polkadot and kusama"
	@cd $(CURDIR)/functional && ginkgo --focus="Parachain commands" -v -p -timeout=3h

run-relay&parachain-kusama: BINARY
	@echo "Running relay and parachain chain commands test cases for kusama and its parachains"
	@cd $(CURDIR)/functional && relayChain=kusama ginkgo --focus="Relaychain and parachain commands" -v -p -timeout=3h

run-relay&parachain-polkadot: BINARY
	@echo "Running relay and parachain chain commands test cases for polkadot and its parachains"
	@cd $(CURDIR)/functional && relayChain=polkadot ginkgo --focus="Relaychain and parachain commands" -v -p -timeout=3h

run-relay&parachain-all: BINARY
	@echo "Running relay and parachain chain commands test cases for polkadot and kusama along its respective parachains"
	@cd $(CURDIR)/functional && ginkgo --focus="Relaychain and parachain commands" -v -p -timeout=3h

run-misc: BINARY
	@echo "Running other test cases"
	@cd $(CURDIR)/functional && ginkgo --focus="Other commands" -v -p -timeout=3h

BINARY: GINKGO
	@cd $(CURDIR)/../cli && go build -o dive

GINKGO:
	@if ! command -v ginkgo version; then \
		echo "Installing ginkgo..."; \
		go install github.com/onsi/ginkgo/v2/ginkgo; \
	fi

remove-bin:
	@cd $(CURDIR)/../cli && rm -rf dive

run-archway-e2e:
	@echo "Running e2e demo cosmos"
	@cd $(CURDIR)/scripts/cosmos && npx ts-node cosmos_e2e_demo.ts archway

run-neutron-e2e:
	@echo "Running e2e demo cosmos"
	@cd $(CURDIR)/scripts/cosmos && npx ts-node cosmos_e2e_demo.ts neutron

run-icon-archway-e2e:
	@echo "Running e2e demo cosmos"
	@cd $(CURDIR)/scripts/cosmos && npx ts-node icon_cosmos_e2e_demo.ts archway

run-icon-neutron-e2e:
	@echo "Running e2e demo cosmos"
	@cd $(CURDIR)/scripts/cosmos && npx ts-node icon_cosmos_e2e_demo.ts neutron




