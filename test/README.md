# End-to-End Testing Demo

This is a demo script from btp2 that demonstrates simple e2e testing scenarios. 

## Prerequisites

To run the demo, the following software needs to be installed.

* Node.js 18 (LTS) \[[download](https://nodejs.org/en/download/)\]
* Docker compose (V2) \[[download](https://docs.docker.com/compose/install/)\]
* OpenJDK 11 or above \[[download](https://adoptium.net/)\]
* jq \[[download](https://github.com/stedolan/jq)\]
* go \[[download](https://go.dev/doc/install)\]

## Steps to run the script
* Step 1: Run the DIVE package that spins ups two chains for sending message using BTP
* Step 2: After you get the contract address from the output, update xCall and dApp address in deployment.json file
* Step 3: Update network and endpoint in the chain_config.json, deployments.json and hardhat.config.ts files
* Step 4: now run the command 'make run-demo' to execute all scenarios.