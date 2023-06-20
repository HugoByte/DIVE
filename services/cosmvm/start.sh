#!/bin/sh

mkdir testnet && cd testnet
mkdir node1 && cd node1

# initialising the node
echo "Instialising the node"

archwayd init node1 --chain-id my-chain

PASSCODE="password"

# adding the keys
echo "Add keys"
(echo $PASSCODE; echo $PASSCODE) | archwayd keys add node1-account
(echo $PASSCODE; echo $PASSCODE) | archwayd keys add node2 


# listing the keys
(echo $PASSCODE; echo $PASSCODE) | archwayd keys list 

# adding the keys to genesis account
echo "Adding the keys to the genesis account"
(echo $PASSCODE; echo $PASSCODE) | archwayd add-genesis-account $(archwayd keys show node1-account -a ) 1000000000stake

# generate genesis transaction
echo "Generate transaction"
(echo $PASSCODE; echo $PASSCODE) | archwayd gentx node1-account 1000000000stake --chain-id my-chain 
# collect geesis transaction
archwayd collect-gentxs

# starting the node
archwayd start