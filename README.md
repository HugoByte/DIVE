![DIVE](img/DIVE.png)

## D.I.V.E.

### About

Dive deeply into the world of Blockchain and Web 3.0 using **D.I.V.E.** (Deployable Infrastructure for Virtually Effortless blockchain integration),The Dive package aim to implement its services and API for ICON Blockchain. The kurtosis services and API are designed to simplify the process of deploying various nodes and services for development and testing and enhance the overall user experience. Implementing kurtosis for the ICON blockchain can help ease the developers in the ecosystem to focus more on building the business logic without worrying about the setup which consumes a significant amount of time.

The vision is to making ICON the interoperable hub by easing the setup of BTP and IBC for ICON and the connecting chains.

This repository uses [kurtosis package](https://docs.kurtosis.com/concepts-reference/packages)

### Setup and requirements

Before proceeding make sure to have

- [Docker installed and running](https://docs.kurtosis.com/install#i-install--start-docker)

- [Kurtosis installed and running ](https://docs.kurtosis.com/install#ii-install-the-cli) or [(upgrading to the latest)](https://docs.kurtosis.com/upgrade)

### Integrating chain

- ICON

- ETHEREUM

### Integrating node

- [**Icon node service package**](./jvm) - This package is responsible for running the ICON node and providing the configuration to the given services.

- [**Icon BTP Integration**](./jvm) - This provides the setup for Deploying BTP Smart Contracts and Relay

- [**Evm chain node package**](./evm/) - This package is responsible for running the EVM chain node and providing the configuration to the given services.

- [**Evm Util Package**](./evm/) - This package is responsible for Uploading and Interacting with Smart Contracts Deployed on EVM based chains.

- [**Evm BTP Integration**](./evm/) - This provides setup for Deploying BTP Smart Contracts and Relay Setup

### Running Dive

Dive-cli is a command line tool that will be used for starting the chain and crosschain communication between two different chains

### Installing Dive CLI

- Install using **`brew`**
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

### Commands

- **bridge** : For setting up communication between two different chains.This will setup the relayer to connect two different chains and pass messages between them.
  **Subcommand**: - **btp** : Starts connection between specified chains using [BTP](https://icon.community/assets/btp-litepaper.pdf)

- **chain**: For initialising and starting a specified blockchain node. By executing this command the node will be launched by enabling the network participation and ledger maintenance within the specified blockchain ecosystem
  **Subcommand**:

  - **eth** : Build, initialise and start a eth node.

  - **hardhat** : Build, initialise and start a hardhat node.

  - **icon** : Build, initialise and start a icon node.

- **clean**: Cleans the Enclave stared by Dive CLI
- **discord**: Redirect to the DIVE discord channel
- **tutorial**: Takes you to Dive tutorials
- **version**: Returns the current version of the CLI

### Usage

> Before proceeding, make sure the Kurtosis Engine is running in the background. If it's not already running, start it by executing the following command

    **```
    kurtosis engine start
    ```**

#### Setting up an Node

```
dive chain icon
```

After running the command, **DIVE CLI** will automatically start the ICON node and handle the necessary initialization processes. Please wait for the ICON node to fully initialize, which may take a few moments.

Once the initialization is complete, you can interact with the local ICON chain as needed.**DIVE CLI** sets up the ICON node on your local environment, enabling you to deploy and test smart contracts, explore transactions, and experiment with various ICON blockchain features

> Detailed output during execution can be enabled via Verbose flag

    Example: `dive chain icon --verbose=ture`
    Anyhow each Dive cli execution will be logged into log files under log folder in current working directory

After successful execution one can find service details in `services.json` created in current working directory.
Example `services.json`:

```javascript
{
"icon-0": {
	"service_name": "icon-node-0",
	"endpoint_public": "http://127.0.0.1:8090/api/v3/icon_dex",
	"endpoint": "http://172.16.0.2:9080/api/v3/icon_dex",
	"keypassword": "gochain",
	"keystore_path": "keystores/keystore.json",
	"network": "0x3.icon",
	"network_name": "icon-0",
	"nid": "0x3"
	}
}
```

#### Setting Bridge Between Two Chains

```bash
dive bridge btp --chainA icon --chainB eth
```

This command sets up btp bridge between icon and eth . After running this command **DIVE CLI** will automatically starts the ICON & ETH node , Deploys contract which is used for BTP and starts the realay to constanly exchange message between established connection.
After successful bridge setup all the neccessary details with respect to bridge will be added to `dive.json` file that will be present in current working directory.
Example `dive.json`:

```javascript
{
  "bridge": "false",
  "chains": {
    "icon": {
      "block_number": "235",
      "endpoint": "http://172.16.0.2:9080/api/v3/icon_dex",
      "endpoint_public": "http://127.0.0.1:8090/api/v3/icon_dex",
      "keypassword": "gochain",
      "keystore_path": "keystores/keystore.json",
      "network": "0x3.icon",
      "networkId": "0x1",
      "networkTypeId": "0x1",
      "network_name": "icon-0",
      "nid": "0x3",
      "service_name": "icon-node-0"
    },
    "icon-1": {
      "block_number": "233",
      "endpoint": "http://172.16.0.3:9081/api/v3/icon_dex",
      "endpoint_public": "http://127.0.0.1:8091/api/v3/icon_dex",
      "keypassword": "gochain",
      "keystore_path": "keystores/keystore.json",
      "network": "0x101.icon",
      "networkId": "0x1",
      "networkTypeId": "0x1",
      "network_name": "icon-1",
      "nid": "0x101",
      "service_name": "icon-node-1"
    }
  },
  "contracts": {
    "icon": {
      "bmc": "cx1755c5fe5012f3f56108a498723532314e003946",
      "bmv": "cx9a49018d108797f4dfe84cc930d33a3a5770a3a1",
      "dapp": "cx20053c926cc0218d0c0a607a5cffb96e207dbfe6",
      "xcall": "cxb895d6c1be173c155b682b4266e52e3165f38163"
    },
    "icon-1": {
      "bmc": "cx36a5fc49d6a77ec62f758577bdc1adb3d76982de",
      "bmv": "cx49c3d9a48a0606fb07639f5e6b56a039f64368c4",
      "dapp": "cx4fd6bdd547078398f255a241a14cab3796591028",
      "xcall": "cx7eab38422d1dedb292243cfed21b16148a42ec09"
    }
  },
  "links": {
    "dst": "icon",
    "src": "icon"
  }
}
```

#### Version

```bash
dive version
```

Prints out the current version of **DIVE CLI**

#### Cleaning

```bash
dive clean
```

Command cleans up the artifacts , services created on the Enclave during **DIVE** package execution

> Checkout More details on how to setup [BTP](https://www.xcall.dev/quickstart/setting-up-a-local-environment-with-dive-cli) bridge

### Testing

- Follow the instruction in [Test Folder](test/README.md#steps-to-run-the-script)

### Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request.

1. Fork the Project.

2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)

3. Commit your Changes (`git commit -m 'feature message'`)

4. Push to the Branch (`git push origin feature/AmazingFeature`)

5. Open a Pull Request.

## References

Special thanks to [Kurtosis-Tech](https://github.com/kurtosis-tech).

### License

Distributed under the Apache 2.0 License. See [LICENSE](./LICENSE) for more information.

### Feedback

We would happy to hear your thoughts on our project. Your feedback helps us improve and make it better for everyone. Please submit your valuable feedback [here](https://docs.google.com/forms/d/e/1FAIpQLScnesE-4IWPrFQ-W2FbRXHyQz8i_C0BVjIP_aWaxKe3myTgyw/viewform?usp=sharing)
