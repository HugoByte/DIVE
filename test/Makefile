PWD = $(abspath ./)

DOCKER_CONTAINER_ID=$(shell docker ps -aqf "name=eth-contract-deployer")

HARDHAT = npx hardhat
NETWORK = --network localnet
FOLDERS = solidity
DEP := solidity node_modules typechain-types

.DEFAULT_GOAL := run-demo

run-demo:$(FOLDERS)
	@ echo ">>> Run demo_dapp.ts " ; \
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


