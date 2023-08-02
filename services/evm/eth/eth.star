eth_node = import_module("github.com/hugobyte/dive/services/evm/eth/src/node-setup/start-eth-node.star")
eth_relay_setup = import_module("github.com/hugobyte/dive/services/evm/eth/src/relay-setup/contract_configuration.star")

def start_eth_node_serivce(plan,args,node_type):

    node_service_data = eth_node.start_node_service(plan,args,node_type)

    config_data = {
                "service_name" : node_service_data.service_name,
                "nid" : node_service_data.nid,
                "network" : node_service_data.network,
                "network_name": node_service_data.network_name,
                "endpoint":  node_service_data.endpoint ,
                "endpoint_public": node_service_data.endpoint_public ,
                "keystore_path" : node_service_data.keystore_path,
                "keypassword": node_service_data.keypassword
            }

    return config_data


def deploy_bmv_eth(plan,bridge,data,args,chain_name):

    if bridge == "true":

        address = eth_relay_setup.deploy_bmv_bridge(plan,args,data.block_height,data.bmc,data.network,chain_name)
        return address

    else :
        address = eth_relay_setup.deploy_bmv(plan,args,data.block_header,data.bmc,data.network,data.network_type_id,chain_name)

        return address


