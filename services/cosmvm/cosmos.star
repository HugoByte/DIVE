cosmvm_node = import_module("github.com/hugobyte/dive/services/cosmvm/src/node-setup/start_node.star")
cosmvm_setup = import_module("github.com/hugobyte/dive/services/cosmvm/src/node-setup/setup_node.star")


SERVICE_NAME = "cosmos"
SERVICE_NAME1 = "cosmos1"
COSMOS_NODE1_CID = "my-chain"
COSMOS_NODE2_CID = "chain-2"

# spins up the cosmos node1

def start_cosmos_0(plan,args):

    plan.print("Starting the cosmos node ID 0")

    node_service = cosmvm_node.start_cosmos_node(plan,args)

    return node_service

def start_cosmos_1(plan,args):

    plan.print("Starting the cosmos node ID 1")

    node_service = cosmvm_node.start_cosmos_node_1(plan,args)

    return node_service

# spins up the 2 comsos nodes
def start_node_service_cosmos_to_cosmos(plan):

    src_chain_config = cosmvm_node.get_service_config(SERVICE_NAME, COSMOS_NODE1_CID )
    dst_chain_config = cosmvm_node.get_service_config_1(SERVICE_NAME1, COSMOS_NODE2_CID )

    source_chain_response = start_cosmos_0(plan,src_chain_config)
    destination_chain_response = start_cosmos_1(plan,dst_chain_config)

    src_service_config =  {
            
                "service_name" : source_chain_response.service_name,
                "endpoint": source_chain_response.endpoint ,
                "endpoint_public": source_chain_response.endpoint_public
    }

    dst_service_config =  {
                "service_name" : destination_chain_response.service_name,
                "endpoint": destination_chain_response.endpoint ,
                "endpoint_public": destination_chain_response.endpoint_public
            }
        
    return struct(
        src_config = src_service_config,
        dst_config = dst_service_config
    )

# spins up the single cosmos node

def start_node_service(plan,args):

    chain_config = cosmvm_node.get_service_config(SERVICE_NAME, COSMOS_NODE1_CID)

    node_service_response = start_cosmos_0(plan,chain_config)


# configures the cosmos node setup

def configure_cosmos_to_cosmos_node(plan,src_chain_config, dst_chain_config):

    plan.print("configuring the nodes")

    cosmvm_setup.configure_node(plan,src_chain_config)
    cosmvm_setup.configure_node(plan,dst_chain_config)

# Configures ICON Node setup
def configure_cosmos_node(plan,chain_config):

    plan.print("configure cosmos Node")

    cosmvm_setup.configure_node(plan,chain_config) 
