#!/bin/sh

srcChainId=$1
dstChainId=$2
srcChainKey=$3
dstChainKey=$4
srcChainEndpoint=$5
dstChainEndpoint=$6
srcChainSeed=$7
dstChainSeed=$8
relayConnectionPath=${srcChainId}-${dstChainId}

sed -i -e 's|\"rpc-addr\": \"\"|\"rpc-addr\": \"'${srcChainEndpoint}'\"|' ../script/${srcChainId}/cosmos-${srcChainId}.json

sed -i -e 's|\"rpc-addr\": \"\"|\"rpc-addr\": \"'${dstChainEndpoint}'\"|' ../script/${dstChainId}/cosmos-${dstChainId}.json


echo "Init Relay"

rly config init

echo "Adding Chain ${srcChainId}"

rly chains add --file ../script/${srcChainId}/cosmos-${srcChainId}.json ${srcChainId}

echo "Adding Chain ${dstChainId}"

rly chains add --file ../script/${dstChainId}/cosmos-${dstChainId}.json ${dstChainId}

echo "Adding Keys"

rly keys restore ${srcChainId} ${srcChainKey} "${srcChainSeed}"

rly keys restore ${dstChainId} ${dstChainKey} "${dstChainSeed}"

echo "Adding the paths"

rly paths new ${srcChainId} ${dstChainId} ${relayConnectionPath}

rly transact connection ${relayConnectionPath}

rly tx link ${relayConnectionPath} -d -t 3s 

rly paths list

rly start

