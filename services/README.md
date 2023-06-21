### About

In this repo, we aim to implement kurtosis services and API for ICON Blockchain. 

**Increased interoperability**: By integrating with ICON and Ethereum, Kurtosis will be able to more easily  interact with other projects and applications on these chains. This will increase the interoperability of Kurtosis and make it more valuable to users.

**Expanded user base**: ICON and Ethereum have a large user base, and by integrating with these chains,  
Kurtosis will be able to reach a wider audience. This will help to grow the Kurtosis ecosystem and make it more successful.

The kurtosis services and API are designed to simplify the process of deploying various nodes and services for development and testing and enhance the overall user experience. Implementing kurtosis for the ICON blockchain can help ease the developers in the ecosystem to focus more on building the business logic without worrying about the setup which consumes a significant amount of time.

- [**Icon and node service package**](./jvm) - This package is responsible for running the ICON node and providing the configuration to the given services.
- [**Icon BTP Integration**](./jvm) - This provides the setup for Deploying BTP Smart Contracts and Relay
- [**Evm chain node package**](./evm/) - This package is responsible for running the ICON node and providing the configuration to the given services.
- [**Evm Util Package**](./evm/) - This package is responsible for Uploading and Interacting with Smart Contracts Deployed on EVM based chains.
- [**Evm BTP Integration**](./evm/) - This provides setup for Deploying BTP Smart Contracts and Relay Setup

### Setup and requirements

Before proceeding make sure to have

- [Docker installed and running](https://docs.kurtosis.com/install#i-install--start-docker)
- [Installed the kurtosis cli ](https://docs.kurtosis.com/install#ii-install-the-cli) or [(upgrading to the latest)](https://docs.kurtosis.com/upgrade)

### Binary file

We can create a binary file where the task will be run to specify the options. The accepted file extensions are as follows

- `.ts` for Ethereum
- `.jar` for Icon

### Configuration File

- Icon

```
{
  "termPeriod": 100,
  "mainPRepCount": 1,
  "extraMainPRepCount": 0,
  "subPRepCount": 4
}
```

- EVM

```
{
  "compilerOptions": {
    "target": "es2020",
    "module": "commonjs",
    "esModuleInterop": true,
    "forceConsistentCasingInFileNames": true,
    "strict": true,
    "skipLibCheck": true,
    "resolveJsonModule": true
  }
}
```

### Running Kurtosis

```
kurtosis run . '{"links":{"src":" ","dst":" "}, "bridge":"true"}' --enclave btp
```

- If we want to tear down the encalve and any of their artifacts, by running `kurtosis clean -a`

While running this we have to give the source and destination path like Icon and EVM respectievely.







