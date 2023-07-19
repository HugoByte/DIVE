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

The available commands are -

1. `bridge`: 
  
   For connecting the two different chains. This will create an relay to connect two different chains and pass any messages between them
   
      subcommands

      - `btp` : Starts BTP bridge between Chain A and Chain B
  
2. `chain` : 
   
   For building, initialising and starting a specified blockchain node. By executing this command the node will be launched by enabling the network participation and ledger maintenance within the specified blockchain ecosystem

      subcommands

    - `eth` : Build, initialise and start a eth node.
    - `hardhat`: Build, initialise and start a hardhat node.
    - `icon`: Build, initialise and start a icon node.


3. `clean`: 
   
   For cleaning the kurtosis enclave
   
4. `discord`:
   
    Redirect to the DIVE discord channel
   
5. `tutorial`: 
  
   Redirect to the DIVE youtube channel
   
6. `version`: 
  
   For getting the current version of the CLI

**Example**

- For building,initialising and starting the icon chain

```
dive-cli chain icon
```

- icon-icon cross chain communication

```
dive-cli bridge btp --chainA icon --chainB icon
```

- icon-eth cross chain communication

```
dive-cli bridge btp --chainA icon --chainB eth
```

- For cleaning the enclave

```
dive-cli clean
```

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
