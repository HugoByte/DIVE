{
"type": "wasm",
"value": {
"key-directory": "/root/.relayer/keys/my-chain",
"key": "{{.KEY}}",
"chain-id": "{{.CHAINID}}",
"rpc-addr": "{{.RPCADDRESS}}",
"account-prefix": "archway",
"keyring-backend": "test",
"gas-adjustment": 1.5,
"gas-prices": "0.025stake",
"min-gas-amount": 1000000,
"debug": true,
"timeout": "20s",
"block-timeout": "",
"output-format": "json",
"sign-mode": "direct",
"extra-codecs": [],
"coin-type": 0,
"broadcast-mode": "batch",
"ibc-handler-address":"{{.IBCADDRESS}}",
"first-retry-block-after": 0,
"start-height": 0,
"block-interval": 6000
}
}