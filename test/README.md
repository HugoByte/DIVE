# End-to-End Testing Demo

This is a demo script from btp2 that demonstrates simple e2e testing scenarios.

## Prerequisites

To run the demo, the following software needs to be installed.

- Node.js 18 (LTS) \[[download](https://nodejs.org/en/download/)\]
- Docker compose (V2) \[[download](https://docs.docker.com/compose/install/)\]
- OpenJDK 11 or above \[[download](https://adoptium.net/)\]
- jq \[[download](https://github.com/stedolan/jq)\]
- go \[[download](https://go.dev/doc/install)\]

## Steps to run the script

- Step 1: Run the DIVE command that spins ups two chains for sending message using BTP

  ![img1](img/../../img/img1%202.png)

- Step 2: After you get the contract address from the output, update xCall and dApp address in deployment.json file

  ![img1](img/../../img/img2%202.png)

- Step 3: Update network and endpoint in the chain_config.json, deployments.json and hardhat.config.ts files

  ![img1](img/../../img/img3%202.png)

- Step 4: now run the command `make run-demo` to execute all scenarios.

  ![img1](img/../../img/Image3.png)

> **Note:**
> Running the demo script will copy all the dependencies required from the container. We can clean the dependencies by running `make clean-dep` command.

# End-to-End archway - archway Demo

## Steps to run the demo

- Step 1: Run the dive command that spins ups two archway chain for sending message using IBC. Wait for the setup to be completed

  ![img1](img/../../img/img4.png)

- Step 2: now run the command `make run-cosmos` to execute e2e demo which transfers token from one archway chain to another.

  ![img1](img/../../img/img5.png)

## Video

- D.I.V.E. package setup for testing the bridge between EVM and ICON using BTP2

  [![video1](img/../../img/video1.png)](https://www.youtube.com/watch?v=f3tMU-_E1a8&ab_channel=HugoByte)

- Setup EVM and JVM local nodes using the D.I.V.E package

  [![video1](img/../../img/video2.png)](https://www.youtube.com/watch?v=390s_uo19eA&t=25s&ab_channel=HugoByte)
