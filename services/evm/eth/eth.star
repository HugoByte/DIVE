eth_node = import_module("github.com/hugobyte/chain-package/services/evm/eth/src/node-setup/start-eth-node.star")
eth_relay_setup = import_module("github.com/hugobyte/chain-package/services/evm/eth/src/relay-setup/contract_configuration.star")

def start_eth_node_serivce(plan,args):

    node_service_data = eth_node.start_eth_node(plan,args)

    config_data = {
                "service_name" : node_service_data.service_name,
                "nid" : node_service_data.nid,
                "network" : node_service_data.network,
                "network_name": node_service_data.network_name,
                "endpoint": "http://%s" % node_service_data.endpoint ,
                "endpoint_public": "",
                "keystore_path" : "config/eth_keystore.json",
                "keypassword": "password"
            }

    return config_data


def deploy_bmv_eth(plan,bridge,data,args):

    if bridge == "true":

        address = eth_relay_setup.deploy_bmv_bridge(plan,args,data.block_height,data.bmc,data.network)
        return address

    else :
        address = eth_relay_setup.deploy_bmv(plan,args,data.block_header,data.bmc,data.network,data.networkTypeId)

        return address


