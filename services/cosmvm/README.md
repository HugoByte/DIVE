### Cosmvm Node Configuration

- Deploying the node

 PATH - contracts wasm path should be given 
  
  ```
  archwayd tx wasm store  /PATH/  --from node1-account --chain-id my-chain --gas auto --gas-adjustment 1.3 -y --output json -b block | jq -r '.logs[0].events[-1].attributes[0].value'
  ```

- Running

```
kurtosis clean -a
```

```
kurtosis run . --enclave chain '{"contract_name":"NAME_OF_CONTRACT", "message":{"KEY":"VALUE"}}'
```
