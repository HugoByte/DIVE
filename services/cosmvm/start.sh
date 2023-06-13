#!/bin/sh

mkdir testnet && cd testnet
mkdir node1 && cd node1

# initialising the node
echo "Instialising the node"
archwayd init node1 --chain-id my-chain --home ./node1

# adding the keys
echo "Add keys"
archwayd keys add node1-account --home ./node1 | echo -n "enter keyring passphrase" | read -r phrase

# listing the keys
archwayd keys list

# adding the keys to genesis account
echo "Adding the keys to the genesis account"
archwayd add-genesis-account $(archwayd keys show node1-account -a --home ./node1) 1000000000stake --home ./node1 | echo -n 'enter keyring passphrase' | read -r phrase

# generate genesis transaction
echo "Generate transaction"
archwayd gentx node1-account 1000000000stake --chain-id my-chain --home ./node1 | echo -n "enter keyring passphrase" | read -r phrase

# collect geesis transaction
archwayd collect-gentxs --home ./node1

# starting the node
archwayd start --home ./node1

