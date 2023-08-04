#!/bin/sh

# this script instantiates localnet required genesis files

set -e

echo clearing $HOME/.archway
rm -rf $HOME/.archway
echo initting new chain
# init config files
archwayd init archwayd-id --chain-id my-chain

# create accounts
archwayd keys add fd --keyring-backend test --output json > ../../start-scripts/key_seed.json 2>&1
(echo "password"; echo "password") |archwayd keys add node1-account 
(echo "password"; echo "password") |archwayd keys add test-account

apk add jq

addr=$(archwayd keys show fd -a --keyring-backend=test)
addres=$(echo "password"| archwayd keys show node1-account -a)
test_address=$(echo "password"| archwayd keys show test-account -a)
val_addr=$(archwayd keys show fd  --keyring-backend=test --bech val -a)

chmod o+x /root/.archway/config/genesis.json

sed -i -e "s|\"accounts\": *\[\]|\"accounts\": [{\"@type\": \"/cosmos.auth.v1beta1.BaseAccount\",\"address\": \"$addres\", \"pub_key\": null,\"account_number\": \"0\", \"sequence\": \"0\"},{\"@type\": \"/cosmos.auth.v1beta1.BaseAccount\",\"address\": \"$test_address\",\"pub_key\": null,\"account_number\": \"0\",\"sequence\": \"0\"}]|" /root/.archway/config/genesis.json
# give the accounts some money
archwayd add-genesis-account "$addr" 1000000000000stake --keyring-backend=test

# save configs for the daemon
archwayd gentx fd 10000000stake --chain-id my-chain --keyring-backend=test

# input genTx to the genesis file
archwayd collect-gentxs
# verify genesis file is fine
archwayd validate-genesis
echo changing network settings
sed -i 's/127.0.0.1/0.0.0.0/g' $HOME/.archway/config/config.toml

echo test account address: "$addr"
echo test account private key: "$(yes | archwayd keys export fd --unsafe --unarmored-hex --keyring-backend=test)"
echo account for --from flag "fd"

echo starting network...
archwayd start 