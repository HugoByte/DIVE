![DIVE](img/DIVE.png)

# D.I.V.E.

[![run smoke testcases](https://github.com/HugoByte/DIVE/actions/workflows/smoke-test.yaml/badge.svg)](https://github.com/HugoByte/DIVE/actions/workflows/smoke-test.yaml)

## About

Dive deeply into the world of Blockchain and Web 3.0 using **D.I.V.E.** (Deployable Infrastructure for Virtually Effortless blockchain integration), The Dive package aims to implement its services and API for ICON Blockchain. The kurtosis services and API are designed to simplify the process of deploying various nodes and services for development and testing and enhance the overall user experience. Implementing kurtosis for the ICON blockchain can help ease the developers in the ecosystem to focus more on building the business logic without worrying about the setup which consumes a significant amount of time.

The vision is to make ICON the interoperable hub by easing the setup of BTP and IBC for ICON and the connecting chains.

This repository uses [kurtosis package](https://docs.kurtosis.com/concepts-reference/packages)

## Setup and requirements

Before proceeding make sure to have

- [Docker installed and running](https://docs.kurtosis.com/install#i-install--start-docker)

- [Kurtosis installed and running ](https://docs.kurtosis.com/install#ii-install-the-cli) or [(upgrading to the latest)](https://docs.kurtosis.com/upgrade)

## Integrating chain

- ICON

- ETHEREUM

### Integrating node

- [**Icon node service package**](./jvm) - This package is responsible for running the ICON node and providing the configuration to the given services.

- [**Icon BTP Integration**](./jvm) - This provides the setup for Deploying BTP Smart Contracts and Relay

- [**Evm chain node package**](./evm/) - This package is responsible for running the EVM chain node and providing the configuration to the given services.

- [**Evm Util Package**](./evm/) - This package is responsible for Uploading and Interacting with Smart Contracts Deployed on EVM based chains.

- [**Evm BTP Integration**](./evm/) - This provides setup for Deploying BTP Smart Contracts and Relay Setup

## DIVE CLI

<p align='justify'>
DIVE CLI – a powerful tool designed to streamline the entire process of node setup, network configuration, and BTP bridge creation.
With DIVE CLI, developers can easily connect and interact with various blockchain networks, paving the way for seamless cross-chain communication and smart contract deployment.
Serving as an all-in-one solution, DIVE CLI eliminates the hassle of manually configuring nodes, allowing developers to effortlessly set up nodes for the BTP network with just a few simple commands. The tool provides a user-friendly interface that makes the process accessible even to those new to blockchain development.
</p>

## Installing Dive CLI

- Install on **`MacOS`**
  ```
  brew install hugobyte/tap/dive-cli
  ```
- Install on **`linux`**

  ```bash
   $ curl -L https://github.com/HugoByte/DIVE/releases/download/{latest-version}/dive-cli_{latest-version}_linux_{arch}.tar.gz
   $ tar -xvzf dive-cli_{latest-version}_linux_{arch}.tar.gz
   $ mv dive /usr/local/bin
   $ rm dive-cli_{latest-version}_linux_{arch}.tar.gz

  ```

- Install on **`windows`**:
  ```bash
  Invoke-WebRequest -Uri "https://github.com/HugoByte/DIVE/releases/download/{latest-version}/dive-cli_{latest-version}_windows_{arch}.tar.gz" -OutFile dive.tar.gz
  tar -xvzf dive.tar.gz
  dive.exe
  ```

## Commands

- **bridge** : For setting up communication between two different chains.This will setup the relayer to connect two different chains and pass messages between them.

  **Subcommand**:

  - **btp** : Starts connection between specified chains using [BTP](https://icon.community/assets/btp-litepaper.pdf)

- **chain**: For initialising and starting a specified blockchain node. By executing this command the node will be launched by enabling the network participation and ledger maintenance within the specified blockchain ecosystem

  **Subcommand**:

  - **eth** : Build, initialise and start a eth node.

  - **hardhat** : Build, initialise and start a hardhat node.

  - **icon** : Build, initialise and start a icon node.

- **clean**: Cleans the Enclave stared by Dive CLI
- **discord**: Redirect to the DIVE discord channel
- **tutorial**: Takes you to Dive tutorials
- **version**: Returns the current version of the CLI

## Usage

> Before proceeding, make sure the Kurtosis Engine is running in the background. If it's not already running, start it by executing the following command: `kurtosis engine start`

### Setting up a Node

```
dive chain icon
```

After running the command, **DIVE CLI** will automatically start the ICON node and handle the necessary initialization processes. Please wait for the ICON node to fully initialize, which may take a few moments.

Once the initialization is complete, you can interact with the local ICON chain as needed.**DIVE CLI** sets up the ICON node on your local environment, enabling you to deploy and test smart contracts, explore transactions, and experiment with various ICON blockchain features

> Detailed output during execution can be enabled via Verbose flag.
> example: `dive chain icon --verbose`
> Anyhow each Dive cli execution will be logged into log files under log folder in current working directory

After successful execution one can find service details in `services.json` created in current working directory.
Example `services.json`:

```javascript
{
"icon-node-0xacbc4e": { # this key specifies service name
	"block_number": "206",
	"endpoint": "http://172.16.0.2:9080/api/v3/icon_dex",
	"endpoint_public": "http://127.0.0.1:8090/api/v3/icon_dex",
	"keypassword": "gochain",
	"keystore_path": "keystores/keystore.json",
	"network": "0x3.icon",
	"networkId": "0x1",
	"networkTypeId": "0x1",
	"network_name": "icon-0xacbc4e",
	"nid": "0x3",
	"service_name": "icon-node-0xacbc4e"
 }
}
```

### Setting Up Bridge Between Two Chains Which Is Already Running

- Starting ICON

  ```bash
  dive chain icon -d #This spins up icon and decentralise for btp
  ```

- Starting ETH

  ```bash
  dive chain eth  --verbose=true #This spins up Eth
  ```

  > `--verbose=true` can be used to see details execution logs

  Once chains are running you can find a services.json file in current working directory. Example services.json can be found below.

  ```javascript
  {

  	"icon-node-0xacbc4e": {
  		"block_number": "206",
  		"endpoint": "http://172.16.0.2:9080/api/v3/icon_dex",
  		"endpoint_public": "http://127.0.0.1:8090/api/v3/icon_dex",
  		"keypassword": "gochain",
  		"keystore_path": "keystores/keystore.json",
  		"network": "0x3.icon",
  		"networkId": "0x1",
  		"networkTypeId": "0x1",
  		"network_name": "icon-0xacbc4e",
  		"nid": "0x3",
  		"service_name": "icon-node-0xacbc4e"
  	},
  	"el-1-geth-lighthouse": {
  		"block_number": "24",
  		"endpoint": "http://172.16.0.7:8545",
  		"endpoint_public": "http://",
  		"keypassword": "password",
  		"keystore_path": "keystores/eth_keystore.json",
  		"network": "0x301824.eth",
  		"network_name": "eth",
  		"nid": "0x301824",
  		"service_name": "el-1-geth-lighthouse"
  	}
  }
  ```

  Now you can start bridge just by running

  ```bash
  dive bridge btp --chainA icon --chainB eth --chainAServiceName icon-node-0xacbc4e  --chainBServiceName el-1-geth-lighthouse
  ```

### Setting Bridge Between Two Chains

Run below command to start btp connection between any supported chain

```bash
dive bridge btp --chainA icon --chainB eth -b
```

> `-b` flag is used to specify the type of bmv contract to be deployed for btp setup.

This command sets up btp bridge between icon and eth . After running this command **DIVE CLI** will automatically starts the ICON & ETH node , Deploys contract which is used for BTP and starts the relay to constantly exchange message between established connection.
After successful bridge setup all the necessary details with respect to bridge will be added to `dive.json` file that will be present in current working directory.
Example `dive.json`:

```javascript
{
	"bridge": "true",
	"chains": {
		"el-1-geth-lighthouse": {
			"block_number": "24",
			"endpoint": "http://172.16.0.7:8545",
			"endpoint_public": "http://",
			"keypassword": "password",
			"keystore_path": "keystores/eth_keystore.json",
			"network": "0x301824.eth",
			"network_name": "eth",
			"nid": "0x301824",
			"service_name": "el-1-geth-lighthouse"
		},
		"icon-node-0xacbc4e": {
			"block_number": "206",
			"endpoint": "http://172.16.0.2:9080/api/v3/icon_dex",
			"endpoint_public": "http://127.0.0.1:8090/api/v3/icon_dex",
			"keypassword": "gochain",
			"keystore_path": "keystores/keystore.json",
			"network": "0x3.icon",
			"networkId": "0x1",
			"networkTypeId": "0x1",
			"network_name": "icon-0xacbc4e",
			"nid": "0x3",
			"service_name": "icon-node-0xacbc4e"
		}
	},
	"contracts": {
		"el-1-geth-lighthouse": {
			"bmc": "0xB9D7a3554F221B34f49d7d3C61375E603aFb699e",
			"bmcm": "0xAb2A01BC351770D09611Ac80f1DE076D56E0487d",
			"bmcs": "0xBFF5cD0aA560e1d1C6B1E2C347860aDAe1bd8235",
			"bmv": "0x765E6b67C589A4b40184AEd9D9ae7ba40E32F8d4",
			"dapp": "0x9bE03fF3E1888A216f9e48c68B587A89c5b94CD6",
			"xcall": "0x5911A6541729C227fAda7D5187ee7518B47fB237"
		},
		"icon-node-0xacbc4e": {
			"bmc": "cx3f9b7aa2a7fa0334a0068a324c9020b9138363f1",
			"bmv": "cxe308f1f14febfff0f906df89d4bd191ba11b4689",
			"dapp": "cxfe3cdbe04e78ff3747b076cb7122c4f7ba58cf49",
			"xcall": "cx4b33b94cb04bf2c179cda1af81c1d1eb639c5e98"
		}
	},
	"links": {
		"dst": "el-1-geth-lighthouse", # service name for eth chain
		"src": "icon-node-0xacbc4e"  # Service name for ICON chain
	}
}
```

### Version

```bash
dive version
```

Prints out the current version of **DIVE CLI**

### Cleaning

```bash
dive clean
```

Command cleans up the artifacts, services created on the Enclave during **DIVE** package execution

> Checkout More details on how to setup [BTP](https://www.xcall.dev/quickstart/setting-up-a-local-environment-with-dive-cli) bridge

## Testing

- Follow the instruction in [Test Folder](test/README.md#steps-to-run-the-script)

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request.

1. Fork the Project.

2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)

3. Commit your Changes (`git commit -m 'feature message'`)

4. Push to the Branch (`git push origin feature/AmazingFeature`)

5. Open a Pull Request.

## References

Special thanks to [Kurtosis-Tech](https://github.com/kurtosis-tech).

## License

Distributed under the Apache 2.0 License. See [LICENSE](./LICENSE) for more information.

## Feedback

We would be happy to hear your thoughts on our project. Your feedback helps us improve and make it better for everyone. Please submit your valuable feedback [here](https://docs.google.com/forms/d/e/1FAIpQLScnesE-4IWPrFQ-W2FbRXHyQz8i_C0BVjIP_aWaxKe3myTgyw/viewform?usp=sharing)
