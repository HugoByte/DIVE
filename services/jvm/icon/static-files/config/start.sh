#!/bin/sh

start_chain() {
  while true; do
    RES=$(goloop system info 2>&1)
    if [ "$?" == "0" ]; then
      break
    fi
    sleep 1
  done
  echo $RES

  
  CID=${1}
  if [ ! -e ${GOLOOP_NODE_DIR}/${CID} ]; then
    # join chain
    GENESIS=/goloop/genesis/${2}
    goloop chain join \
        --platform icon \
        --channel icon_dex \
        --genesis ${GENESIS} \
        --tx_timeout 10000 \
        --node_cache small \
        --normal_tx_pool 1000 \
        --db_type rocksdb \
        --role 3
    goloop system config rpcIncludeDebug true
  fi
  goloop chain start ${CID}
}

# start chain in backgound
start_chain "$1" "$2" &

# start goloop server
exec goloop server start
