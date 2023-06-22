## D.I.V.E.

### About

Dive deeply into the world of Blockchain and Web 3.0 using **D.I.V.E.** (Deployable Infrastructure for Virtually Effortless blockchain integration),The Dive package aim to implement its services and API for ICON Blockchain. The kurtosis services and API are designed to simplify the process of deploying various nodes and services for development and testing and enhance the overall user experience. Implementing kurtosis for the ICON blockchain can help ease the developers in the ecosystem to focus more on building the business logic without worrying about the setup which consumes a significant amount of time.

The vision is to making ICON the interoperable hub by easing the setup of BTP and IBC for ICON and the connecting chains.

### Setup and requirements

Before proceeding make sure to have

- [Docker installed and running](https://docs.kurtosis.com/install#i-install--start-docker)
- [Installed the kurtosis cli ](https://docs.kurtosis.com/install#ii-install-the-cli) or [(upgrading to the latest)](https://docs.kurtosis.com/upgrade)

### Integrating chain

 - ICON  
 - ETHEREUM

### Integrating node

- [**Icon node service package**](./jvm) - This package is responsible for running the ICON node and providing the configuration to the given services.
- [**Icon BTP Integration**](./jvm) - This provides the setup for Deploying BTP Smart Contracts and Relay
- [**Evm chain node package**](./evm/) - This package is responsible for running the EVM chain node and providing the configuration to the given services.
- [**Evm Util Package**](./evm/) - This package is responsible for Uploading and Interacting with Smart Contracts Deployed on EVM based chains.
- [**Evm BTP Integration**](./evm/) - This provides setup for Deploying BTP Smart Contracts and Relay Setup

### Running Kurtosis

```
kurtosis run . '{"links":{"src":" ","dst":" "}, "bridge":"true"}' --enclave <enclave_name>
```

- By running `kurtosis clean -a` , we can tear down the [encalve](https://docs.kurtosis.com/concepts-reference/enclaves/) and any of their artifacts.

Example.

* For Icon-Ethereum

```
kurtosis run . '{"links":{"src":"icon","dst":"eth"}, "bridge":"true"}' --enclave btp
```

- For Icon-Icon

```
kurtosis run . '{"links":{"src":"icon","dst":"icon"}, "bridge":"false"}' --enclave btp
```
 *Note:* The `bridge` should be false for Icon to Icon
