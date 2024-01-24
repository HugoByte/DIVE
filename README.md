![DIVE](img/DIVE.png)

**Dive deeply into the world of Blockchain and Web 3.0 using **D.I.V.E.** (Deployable Infrastructure for Virtually Effortless blockchain integration)**

[![Discord](https://img.shields.io/badge/Discord-hugobyte-2FC252?style=flat&logo=discord&labelColor=black)](https://discord.gg/GyRQSBN3Cu)

# D.I.V.E.

[![run smoke testcases](https://github.com/HugoByte/DIVE/actions/workflows/smoke-test.yaml/badge.svg)](https://github.com/HugoByte/DIVE/actions/workflows/smoke-test.yaml)

# Table of Contents
- [Introduction](#introduction)
  - [About](#about)
  - [Prerequisites](#prerequisites)
- [Installing Dive CLI](#installing-dive-cli)
- [Commands](#commands)
    - [Available Flags](#available-flags)
- [Usage](#usage)
    - [Setting up a Node](#setting-up-a-node)
    - [Setting Up Bridge Between Two Chains Which Are Already Running](#setting-up-bridge-between-two-chains-which-are-already-running)
    - [Setting Up Bridge Between Two Chains](#setting-up-bridge-between-two-chains)
    - [Setting Up Relay and Para Chains](#setting-up-relay-and-para-chains)
- [Configuration File Guidelines](#configuration-file-guidelines)
    - [Configuration Fields](#configuration-fields)
    - [Config Templates](#config-templates)
- [Service Details](#service-details)
    - [Chain Command](#chain-command)
    - [Bridge Command](#bridge-command)
    - [Relay and Para Chains](#relay-and-para-chains)
- [Logs](#logs)
- [Version](#version)
- [Cleaning](#cleaning)
- [Enclaves](#enclaves)
- [Socials](#socials)
- [Testing](#testing)
- [Known Issues](#known-issues)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)
- [Feedback](#feedback)

# Introduction

<p align='justify'>
DIVE CLI – a powerful tool designed to streamline the entire process of node setup, network configuration, and bridge creation.
With DIVE CLI, developers can easily connect and interact with various blockchain networks, paving the way for seamless cross-chain communication and smart contract deployment.
Serving as an all-in-one solution, DIVE CLI eliminates the hassle of manually configuring nodes, allowing developers to effortlessly set up nodes for the bridge networks with just a few simple commands. The tool provides a user-friendly interface that makes the process accessible even to those new to blockchain development.
</p>

## About

The Dive CLI aim to implement its services and API for various Blockchains. The kurtosis services and API are designed to simplify the process of deploying various nodes and services for development and testing and enhance the overall user experience. Implementing kurtosis for the ICON blockchain can help ease the developers in the ecosystem to focus more on building the business logic without worrying about the setup which consumes a significant amount of time.

The vision is to making ICON the interoperable hub by easing the setup of BTP and IBC for ICON and the connecting chains.

## Prerequisites

Ensure the following prerequisites are met before using the Dive-package:
- [Docker](https://www.docker.com/): Make sure Docker is installed on your machine. You can [install it here](https://www.docker.com/).

- [Kurtosis](https://www.kurtosis.com/): Ensure Kurtosis is installed on your machine. You can [install it here](https://www.kurtosis.com/).


# Installing Dive CLI

- Install on **`MacOS`**
  ```
  brew install hugobyte/tap/dive-cli
  ```
- Install on **`linux`**

  Please find the latest release [here](https://github.com/HugoByte/DIVE/releases)

  Run below command to install DIVE CLI by mentioning latest release version and machine arch where dive is getting installed:

  ```shell
  curl -L https://github.com/HugoByte/DIVE/releases/download/{latest-version}/dive-cli_{latest-version}_linux_{arch}.tar.gz | sudo  tar -xzv -C /usr/local/bin/ dive
  ```

  For example, if the latest version is v0.0.13-beta and system's architecture is amd64, the command will be:

  ```shell
  curl -L https://github.com/HugoByte/DIVE/releases/download/v0.0.13-beta/dive-cli_v0.0.13-beta_linux_amd64.tar.gz | sudo tar xzv -C /usr/local/bin/ dive
  ```

  Verify installation by running below command which should print out the dive version installed

  ```shell
  dive version
  ```

- Install on **`windows`**:
  ```bash
  Invoke-WebRequest -Uri "https://github.com/HugoByte/DIVE/releases/download/{latest-version}/dive-cli_{latest-version}_windows_{arch}.tar.gz" -OutFile dive.tar.gz
  tar -xvzf dive.tar.gz
  dive.exe
  ```

# Commands

- **bridge:** For setting up communication between two different chains.This will setup the relayer to connect two different chains and pass messages between them.

  **Subcommand:**

  - **btp:** Starts connection between specified chains using [BTP](https://icon.community/assets/btp-litepaper.pdf).
  - **ibc:** Starts connection between specified chains using [IBC](https://www.ibcprotocol.dev/).

- **chain:** For initialising and starting a specified blockchain node. By executing this command the node will be launched by enabling the network participation and ledger maintenance within the specified blockchain ecosystem.

  **Subcommand:**

  - **archway:** Build, initialise and start a archway node.
  - **eth:** Build, initialise and start a eth node.
  - **hardhat:** Build, initialise and start a hardhat node.
  - **icon:** Build, initialise and start a icon node.
  - **kusama:**      Build, initialize and start a Kusama node
  - **neutron:** Build, initialise and start a neutron node.
  - **polkadot:**    Build, initialize and start a Polkadot node

- **clean:** Destroys and removes any running encalves.
- **discord:** Redirect to the DIVE discord channel.
- **enclaves:** Prints info about kurtosis enclaves.
- **tutorial:** Redirects to the Dive tutorials playlist.
- **twitter:** Opens official HugoByte twitter home page.
- **version:** Returns the current version of the CLI.

## Available Flags

**Global Flags**

- **enclaveName (string):** Provide an enclave name to run services inside. Default enclave name is `dive`
- **verbose (bool):** Print out logs to Stdout.
- **h or help (bool):** Help.

**Flags for Bridges**

- **b or bmvbridge (bool):** Whether to use BMV bridge or not. (Only for BTP bridge)
- **chainA (string):** Name of the source chain.
- **chainAServiceName (string):** Service Name of the source chain from the service details file.
- **chainB (string):** Name of the destination chain.
- **chainBServiceName (string):** Service Name of the destination chain from the service details file.

**Flags for Chains**

For archway/neutron:
    
- **c or config (string):** Path to custom config json file.

For Icon:

- **c or config (string):** Path to custom config json file.
- **d or decentralization (bool):** Decentralize Icon Node.
- **g or genesis (string):** Path to custom genesis file.

For Kusama/Polkadot:

- **c or config (string):** Path to custom config json file to start kusama relaychain and parachain nodes.
- **explorer (bool):** Specify the bool flag if you want to start polkadot js explorer service
- **metrics (bool):** Specify the bool flag if you want to start prometheus and grafana metrics service
- **n or network (string):** Specify the network to run (localnet/testnet/mainnet). Default will be localnet.
- **no-relay (bool):** Specify the bool flag to run parachain only (only for testnet and mainnet)
- **p or parachain (strings):** Specify the list of parachains to spawn parachain node

**Flags for Clean**

- **a (bool):** Clean all running enclaves.


# Usage 

> Before proceeding, make sure the Kurtosis Engine is running in the background. If it's not already running, start it by executing the following command: `kurtosis engine start`

## Setting up a Node

To set up an individual node, simply pass the name of the chain to the dive chain command:

```bash
dive chain archway
```

After running the command, **DIVE CLI** will automatically start the Archway node and handle the necessary initialization processes. Please wait for the Archway node to fully initialize, which may take a few moments.

Once the initialization is complete, you can interact with the local Archway chain as needed.**DIVE CLI** sets up the Archway node on your local environment, enabling you to deploy and test smart contracts, explore transactions, and experiment with various Archway blockchain features

> Detailed output during execution can be enabled via Verbose flag.

> example: `dive chain archway --verbose`

> Each Dive cli execution will be logged into log files under log folder in current working directory.

After successful execution service details can be found in `services_xxxx_xxxx123xx.json`. For more details, refer to the [Service Details](#service-details)

You can also pass your custom config using the `-c` flag:

```bash
dive chain icon -c path/to/config/file
```

For detailed instructions on writing the configuration file, refer to the [Configuration File Guidelines](#configuration-file-guidelines)

## Setting Up Bridge Between Two Chains Which Are Already Running

To set up a bridge between two chains which are already running, you must have the service names of both the chains.

For ex: To start an ibc bridge between icon and archway, we need to have service details of icon and archway as follows:

```json
{
    "icon-node-0xacbc4e": {
        "service_name": "icon-node-0xacbc4e",
        "endpoint_public": "http://127.0.0.1:44444/api/v3/icon_dex",
        "endpoint": "http://172.16.4.6:9080/api/v3/icon_dex",
        "keypassword": "gochain",
        "keystore_path": "keystores/keystore.json",
        "network": "0x3.icon",
        "network_name": "icon-0xacbc4e",
        "nid": "0x3"
    },
    "node-service-constantine-3": {
        "service_name": "node-service-constantine-3",
        "endpoint_public": "http://127.0.0.1:9431",
        "endpoint": "http://172.16.4.5:26657",
        "chain_id": "constantine-3",
        "chain_key": "constantine-3-key"
    }
}
```

Run this command to start a bridge between the two nodes:

```bash
dive bridge ibc --chainA icon --chainB archway --chainAServiceName icon-node-0xacbc4e --chainBServiceName node-service-constantine-3
```
> Note: You can also pass a single running chain instead of passing both. In such case, you can only pass either one of the service name flags that corresponds to the chain.

## Setting Up Bridge Between Two Chains

To set up a bridge between two chains, run this command:

```bash
dive bridge btp --chainA icon --chainB eth -b
```

> `-b` flag is used to specify the type of bmv contract to be deployed for btp setup.

This command sets up btp bridge between icon and eth . After running this command **DIVE CLI** will automatically starts the ICON & ETH node, deploy contracts which is used for BTP and starts the relay to constantly exchange messages between the established connection.

After successful bridge setup all the neccessary details with respect to bridge will be added to `dive_xxxx_xxxx123xx.json`. For more details, refer to the [Service Details](#service-details)

> Checkout More details on how to setup [BTP](https://www.xcall.dev/quickstart/setting-up-a-local-environment-with-dive-cli) bridge

## Setting Up Relay and Para Chains

To set up a relaychain, run this command:

```bash
dive chain kusama
```

To set up a parachain in polkadot along with the relaychain, run this command:

```bash
dive chain polkadot -p frequency
```

To set up only the parachain in kusama, run this command:

```bash
dive chain kusama -p encointer -n testnet --no-relay
```

To set up explorer, pass this flag:

```bash
dive chain kusama -p encointer --explorer
```

To set up metrics, pass this flag:

```bash
dive chain polkadot -p frequency --metrics
```

To specify the network, use the `n` flag:

```bash
dive chain polkadot -n mainnet
```

To pass a custom config, use the `c` flag:

```bash
dive chain kusama -c path/to/config/file
```

For detailed instructions on writing the configuration file, refer to the [Configuration File Guidelines](#configuration-file-guidelines)

> Note: You can run a parachain without a relaychain only in testnet and mainnet networks.

> Note: The default network type is localnet.

# Configuration File Guidelines

You can also pass custom config using the -c flag for chains that support custom config. 

For **cosmos chain (Archway/Neutron)**, the config file is as follows:

```json
{
    "chain_id": "archway-node-0",
    "key": "archway-node-0-key",
    "password": "password"
}
```

For **Icon chain**, the config file is as follows:

```json
{
  "p2p_listen_address": "7080",
  "p2p_address": "8080",
  "cid": "0xacbc4e"
}
```

For **Polkadot/Kusama**, the config file is as follows:

```json
{
  "chain_type": "localnet",
  "relaychain": {
    "name": "rococo-local",
    "nodes": [
      {
        "name": "alice",
        "node_type": "validator",
        "prometheus": false
      },
      {
        "name": "bob",
        "node_type": "validator",
        "prometheus": true

      }
    ]
  },

  "parachains": [
    {
      "name":"acala",
      "nodes": [
        {
          "name": "alice",
          "node_type": "collator",
          "prometheus": false

        },
        {
          "name": "bob",
          "node_type": "full",
          "prometheus": true
        }
      ]
    }
  ],
  "explorer": true
}
```

## Configuration Fields

For **cosmos chain (Archway/Neutron):**

- **chain_id:** The Chain ID of the chain.
- **key:** The Key to use to spawn the node.
- **password:** The Password to use to spawn the node.

For **Icon chain:**

- **p2p_listen_address:** The p2p listen address.
- **p2p_address:** The p2p address.
- **cid:** The CID (Chain ID) of the node.

> Note: The cid for ICON chain must be dervied from the genesis file.

For  **Polkadot/Kusama:**

- **chain_type:** Specifies the type of the network (e.g., "localnet","testnet", "mainnet").
- **relaychain:** Configuration for the relaychain. (When chain_type is "testnet" or "mainenet", the "relaychain" can be an empty dictonary).
  - **name:** Name of the relaychain (e.g., "rococo-local", "rococo", "polkadot" or "kusama").
  - **nodes:** List of nodes on the relaychain, each with:
    - **name:** Node name (e.g., "alice").
    - **node_type:** Node type, can be "validator" or "full".
    - **prometheus:** Whether Prometheus monitoring is enabled (true/false).
- **parachains:** List of parachains, each with:
  - **name:** Parachain name (e.g., "kilt").
  - **nodes:** List of nodes on the parachain, similar to relaychain nodes.
    - **name:** Node name (e.g., "alice").
    - **node_type:** Node type, can be "collator" or "full".
    - **prometheus:** Whether Prometheus monitoring is enabled (true/false).
- **explorer:** Whether Polkadot js explorer is enabled (true/false).

> Note: The polkadot/kusama command start two nodes in relaychain and one node in parachain by default in localnet. In testnet and mainnet, only one node is started for both by default. 

## Config Templates

Feel free to modify this example configuration file based on your specific network requirements. [Here](./cli/sample-jsons) is a link to the official templates that you can edit and use.

# Service Details

The service details are all stored in the output folder in the current working directory. The output is further divided into sub-folders named after the enclave. The sample service details for various commands are given below.

## Chain Command
The service details returned after running a chain command is as follows:

```json
{
    "icon-node-0xacbc4e": {
        "service_name": "icon-node-0xacbc4e",
        "endpoint_public": "http://127.0.0.1:44444/api/v3/icon_dex",
        "endpoint": "http://172.16.4.6:9080/api/v3/icon_dex",
        "keypassword": "gochain",
        "keystore_path": "keystores/keystore.json",
        "network": "0x3.icon",
        "network_name": "icon-0xacbc4e",
        "nid": "0x3"
    }
}
```

## Bridge Command

The service details returned after running a bridge command is as follows:

```json
{
    "ibc-bridge-icon-archway": {
        "chains": {
            "icon-node-0xacbc4e": {
                "endpoint": "http://172.16.4.6:9080/api/v3/icon_dex",
                "endpoint_public": "http://127.0.0.1:44444/api/v3/icon_dex",
                "keypassword": "gochain",
                "keystore_path": "keystores/keystore.json",
                "network": "0x3.icon",
                "network_name": "icon-0xacbc4e",
                "nid": "0x3",
                "service_name": "icon-node-0xacbc4e"
            },
            "node-service-constantine-3": {
                "chain_id": "constantine-3",
                "chain_key": "constantine-3-key",
                "endpoint": "http://172.16.4.5:26657",
                "endpoint_public": "http://127.0.0.1:9431",
                "service_name": "node-service-constantine-3"
            }
        },
        "contracts": {
            "icon-node-0xacbc4e": {
                "dapp": "cxdc423964a82cb08ce561b35162b1206ade3199b9",
                "ibc_core": "cxf60b8dfcf5745df1ca81832bbd0281fb1c961413",
                "light_client": "cx22edaace91d092b3bd62008a57cef77fb8cc458c",
                "xcall": "cxc34f0537d11e3c26ee4bbcb6c181daba3a84d0cd",
                "xcall_connection": "cx13c69e008ae87e5d0c90ea72f3aa3202e068fe3b"
            },
            "node-service-constantine-3": {
                "dapp": "archway1eyfccmjm6732k7wp4p6gdjwhxjwsvje44j0hfx8nkgrm8fs7vqfshgatxw",
                "ibc_core": "archway14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sy85n2u",
                "light_client": "archway1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrqgj33g6",
                "xcall": "archway17p9rzwnnfxcjp32un9ug7yhhzgtkhvl9jfksztgw5uh69wac2pgssf05p7",
                "xcall_connection": "archway1ghd753shjuwexxywmgs4xz7x2q732vcnkm6h2pyv9s6ah3hylvrqvlzdpl"
            }
        },
        "links": {
            "dst": "node-service-constantine-3",
            "src": "icon-node-0xacbc4e"
        }
    }
}
```
## Relay and Para Chains

The service details returned after running a relaychain/parachain is as follows:

```json
{
    "frequency-alice-localnet": {
        "service_name": "frequency-alice-localnet",
        "endpoint_public": "ws://127.0.0.1:31584",
        "endpoint": "ws://172.16.0.10:9946",
        "endpoint_prometheus": "tcp://127.0.0.1:13713",
        "prometheus": true,
        "ip_address": "172.16.0.10",
        "prometheus_port": 9615,
        "prometheus_public_port": 13713
    },
    "grafana": {
        "service_name": "grafana",
        "endpoint_public": "http://127.0.0.1:64304",
        "endpoint": "http://172.16.0.13:3000"
    },
    "polkadot-js-explorer": {
        "service_name": "polkadot-js-explorer",
        "endpoint_public": "http://127.0.0.1:80",
        "endpoint": "http://172.16.0.14:80"
    },
    "prometheus": {
        "service_name": "prometheus",
        "endpoint_public": "http://127.0.0.1:25553",
        "endpoint": "http://172.16.0.12:9090"
    },
    "rococo-local-alice": {
        "service_name": "rococo-local-alice",
        "endpoint_public": "ws://127.0.0.1:54384",
        "endpoint": "ws://172.16.0.5:9944",
        "endpoint_prometheus": "tcp://127.0.0.1:60957",
        "prometheus": true,
        "ip_address": "172.16.0.5",
        "prometheus_port": 9615,
        "prometheus_public_port": 60957
    },
    "rococo-local-bob": {
        "service_name": "rococo-local-bob",
        "endpoint_public": "ws://127.0.0.1:21996",
        "endpoint": "ws://172.16.0.4:9944",
        "endpoint_prometheus": "tcp://127.0.0.1:9351",
        "prometheus": true,
        "ip_address": "172.16.0.4",
        "prometheus_port": 9615,
        "prometheus_public_port": 9351
    }
}
```

> Note: The service files are named as `services_xxxx_xxxx123xx.json` and `dive_xxxx_xxxx123xx.json` for chain and bridge commands respectively.

> Note: The file names for both chain and bridge commands contain the name of the enclave and the short UUID of the enclave. 

# Logs

The logs are located within the 'logs' folder, further organized into individual folders named after each enclave.

You can find the 'logs' folder in the current working directory.

Each folder named after the enclave, has two files:

- **dive.log** : It is created when the execution of a command starts. It contains the execution logs.
- **error.log** : It is created when an error occurs during the execution of the command. It contains the error logs.

The file names also contain the timestamp which can be helpful to find a particular log file.

> Note: The logs folder is not deleted when dive clean is run.

# Version

To check the current version of **DIVE CLI**, run this command:

```bash
dive version
```

# Cleaning

To clean a specific enclave, use:

```bash
dive clean --enclaveName 'enclave'
```

To clean all running enclaves, use:

```bash
dive clean -a
```

> Note: Using the clean command will remove the service files in the output folder. Using `-a` flag removes the output folder.

# Enclaves

To get the details of all the running enclaves, use:

```bash
dive enclaves
```

# Socials

To access any of HugoByte's official social media, run:

```bash
dive discord
```

```bash
dive tutorial
```

# Testing

For guidelines on testing, please refer [here](test/README.md).

## Known Issues

[Here](https://github.com/HugoByte/DIVE/issues) is a list of known issues and their status that our team is currently working to resolve. 

## Troubleshooting

If you encounter issues while using the Dive-packages, refer to the following troubleshooting tips to resolve common problems:

- Clean kurtosis engine using:
```bash
kurtosis clean -a
```

- Restart kurtosis engine using:
```bash
kurtosis engine restart
```

- Check if your docker is installed and accessible:
```bash
docker --version
```

- Check if your DIVE-CLI is installed and accessible:
```bash
dive version
```

- Upgrade or Re-install your DIVE-CLI:

Refer [here](#installing-dive-cli) to upgrade or re-install your cli.

If you still experience issues after following these troubleshooting tips, please [open an issue](https://github.com/HugoByte/DIVE/issues) to get further assistance.


## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. We welcome contributions to enhance and expand the functionality of the DIVE. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request.

1. Fork the Project.

2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)

3. Commit your Changes (`git commit -m 'feature message'`)

4. Push to the Branch (`git push origin feature/AmazingFeature`)

5. Open a Pull Request.

## References

- This repository uses [dive-packages](https://github.com/HugoByte/dive-packages)
- This repository uses [polkadot-kurtosis-packages](https://github.com/HugoByte/polkadot-kurtosis-package)

## License

Distributed under the Apache 2.0 License. See [LICENSE](./LICENSE) for more information.

## Feedback

We would happy to hear your thoughts on our project. Your feedback helps us improve and make it better for everyone. Please submit your valuable feedback [here](https://docs.google.com/forms/d/e/1FAIpQLScnesE-4IWPrFQ-W2FbRXHyQz8i_C0BVjIP_aWaxKe3myTgyw/viewform?usp=sharing)
