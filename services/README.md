### About

This repo contains how we are integrating the Icon and Evm chain with the kurtois.

**Increased interoperability**: By integrating with ICON and Ethereum, Kurtosis will be able to more easily  interact with other projects and applications on these chains. This will increase the interoperability of Kurtosis and make it more valuable to users.

**Expanded user base**: ICON and Ethereum have a large user base, and by integrating with these chains,  
Kurtosis will be able to reach a wider audience. This will help to grow the Kurtosis ecosystem and make it more successful.

- By configuring the icon and Ethereum node
- Deploying the configured node
- Getting the node endpoints of the deployed node
- Integrating the icon with btp by contract configuration and setuping the relay
- Integrating the evm with btp by contract configuration and setuping the relay

### Setup and requirements

Before proceeding make sure to have

- [Docker installed and running](https://docs.kurtosis.com/install#i-install--start-docker)
- [Installed the kurtosis cli ](https://docs.kurtosis.com/install#ii-install-the-cli) or [(upgrading to the latest)](https://docs.kurtosis.com/upgrade)

### Binary file

We can create a binary file where the task will be run to specify the options. The accepted file extensions are as follows

- `.ts` for Ethereum
- `.jar` for Icon

### Running Kurtosis

```
kurtosis run . '{"links":{"src":" ","dst":" "}, "bridge":"true"}' --enclave btp
```

- If we want to tear down the encalve and any of their artifacts, by running `kurtosis clean -a`

While running this we have to give the source and destination path like Icon and EVM respectievely.







