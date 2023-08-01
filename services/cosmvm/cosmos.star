cosmvm_node = import_module("github.com/hugobyte/dive/services/cosmvm/src/node-setup/start_node.star")
cosmvm_setup = import_module("github.com/hugobyte/dive/services/cosmvm/src/node-setup/setup_node.star")


SERVICE_NAME = "cosmos"
SERVICE_NAME1 = "cosmos1"
COSMOS_NODE1_CID = "my-chain"
COSMOS_NODE2_CID = "chain-2"

# spins up the coamos node1

def start_cosmos0(plan,args):

    plan.print("Starting the cosmos node ID 0")

    node_service = cosmvm_node.start_cosmos_node(plan,args)

    return node_service

def start_cosmos1(plan,args):

    plan.print("Starting the cosmos node ID 1")

    node_service = cosmvm_node.start_cosmos_node1(plan,args)

    return node_service


# spins up the 2 comsos node

def start_node_service_cosmos_to_cosmos(plan):

    src_chain_config = cosmvm_node.get_service_config(SERVICE_NAME, COSMOS_NODE1_CID )
    dst_chain_config = cosmvm_node.get_service_config1(SERVICE_NAME1, COSMOS_NODE2_CID )

    source_chain_response = start_cosmos0(plan,src_chain_config)
    destination_chain_response = start_cosmos1(plan,dst_chain_config)


# spins up the single cosmos node

def start_node_service(plan,args):

    chain_config = cosmvm_node.get_service_config(SERVICE_NAME, COSMOS_NODE1_CID)

    node_service_response = start_cosmos0(plan,chain_config)


# configures the cosmos node setup

def configure_cosmos_to_cosmos_node(plan,src_chain_config, dst_chain_config):

    plan.print("configuring the nodes")

    cosmvm_setup.configure_node(plan,src_chain_config)
    cosmvm_setup.configure_node(plan,dst_chain_config)

# Configures ICON Node setup
def configure_cosmos_node(plan,chain_config):

    plan.print("configure cosmos Node")

    cosmvm_setup.configure_node(plan,chain_config) 









# # spins up the cosmos_node0
# def start_cosmos0(plan,args):

#     plan.print("Starting the cosmos node0")

#     node_service = cosmvm_node.start_cosmos_node(plan, args, SERVICE_NAME, CID)

#     return node_service

# # spins up the cosmos_node1
# def start_cosmos1(plan,args):

#     plan.print("starting cosmos node1")

#     node_service = cosmvm_node.start_cosmos_node(plan,args, SERVICE_NAME, CID)

#     return node_service

# # spins up the cosmos_nodes {cosmos0 & cosmos1}
# def start_node_service_cosmos_to_cosmos(plan):

#     src_chain_config = cosmvm_node.get_service_config(COSMOS0,SERVICE_NAME, CID)
#     dst_chain_config = cosmvm_node.get_service_config(COSMOS1,SERVICE_NAME, CID)

#     source_chain_reponse = start_cosmos0(plan,args)

#     destination_chain_response = start_cosmos1(plan,args)

#     src_service_config = {
#         "service_name" : source_chain_reponse.service_name,
#         "nid" : source_chain_reponse.nid
#     }

#     dst_service_config = {
#         "service_name" : destination_chain_response.service_name,
#         "nid" : destination_chain_response.nid
#     }

#     return struct(
#         src_config = src_service_config,
#         dst_config = dst_service_config
#     )

# # Configures COSMOS nodes setup
# def configure_cosmos_to_cosmos(plan,src_chain_config,dst_chain_config):

#     plan.print("Configuring COSM Nodes")

#     setup_node.configure_node(plan,src_chain_config) 
#     setup_node.configure_node(plan,dst_chain_config)

# # Configures COSMO Node setup
# def configure_cosmos_node(plan,chain_config):

#     plan.print("configure ICON Node")

#     setup_node.configure_node(plan,chain_config) 




