wallet = import_module("github.com/hugobyte/dive/services/jvm/icon/src/node-setup/wallet.star")
setup_node = import_module("github.com/hugobyte/dive/services/jvm/icon/src/node-setup/setup_icon_node.star")
icon_node_launcher = import_module("github.com/hugobyte/dive/services/jvm/icon/src/node-setup/start_icon_node.star")
icon_relay_setup = import_module("github.com/hugobyte/dive/services/jvm/icon/src/relay-setup/contract_configuration.star")

START_FILE_FOR_ICON0 = "start-icon-0.sh"
START_FILE_FOR_ICON1 = "start-icon-1.sh"
ICON0_NODE_ID = 0
ICON1_NODE_ID = 1
ICON0_NODE_PRIVATE_RPC_PORT = 9080
ICON0_NODE_PUBLIC_RPC_PORT = 8090
ICON0_NODE_P2P_LISTEN_ADDRESS = 7080
ICON0_NODE_P2P_ADDRESS = 8080
ICON1_NODE_PRIVATE_RPC_PORT = 9081
ICON1_NODE_PUBLIC_RPC_PORT = 8091
ICON1_NODE_P2P_LISTEN_ADDRESS = 7081
ICON1_NODE_P2P_ADDRESS = 8081
ICON0_NODE_CID = "0xacbc4e"
ICON1_NODE_CID = "0x42f1f3"

# Spins up ICON Node ID 0
def start_icon_node_0(plan,service_config):

    plan.print("Starting Icon Node: Id 0")

    node_service = icon_node_launcher.start_icon_node(plan,service_config,ICON0_NODE_ID,START_FILE_FOR_ICON0)

    return node_service

# Spins up ICON Node ID 1
def start_icon_node_1(plan,service_config):

    plan.print("Starting Icon Node: Id 1")

    node_service = icon_node_launcher.start_icon_node(plan,service_config,ICON1_NODE_ID,START_FILE_FOR_ICON1)

    return node_service

# Spins up ICON Nodes {ICON-0 & ICON-1}
def start_node_service_icon_to_icon(plan):

    src_chain_config  = icon_node_launcher.get_service_config(ICON0_NODE_ID,ICON0_NODE_PRIVATE_RPC_PORT,ICON0_NODE_PUBLIC_RPC_PORT,ICON0_NODE_P2P_LISTEN_ADDRESS,ICON0_NODE_P2P_ADDRESS,ICON0_NODE_CID)
    dst_chain_config  = icon_node_launcher.get_service_config(ICON1_NODE_ID,ICON1_NODE_PRIVATE_RPC_PORT,ICON1_NODE_PUBLIC_RPC_PORT,ICON1_NODE_P2P_LISTEN_ADDRESS,ICON1_NODE_P2P_ADDRESS,ICON1_NODE_CID)

    source_chain_response = start_icon_node_0(plan,src_chain_config)

    destination_chain_response = start_icon_node_1(plan,dst_chain_config)

    src_service_config =  {
            
                "service_name" : source_chain_response.service_name,
                "nid" : source_chain_response.nid,
                "network" : source_chain_response.network,
                "network_name": source_chain_response.network_name,
                "endpoint": source_chain_response.endpoint ,
                "endpoint_public": source_chain_response.endpoint_public,
                "keystore_path" : source_chain_response.keystore_path,
                "keypassword": source_chain_response.keypassword

    }

    dst_service_config =  {
                "service_name" : destination_chain_response.service_name,
                "nid" : destination_chain_response.nid,
                "network" : destination_chain_response.network,
                "network_name": destination_chain_response.network_name,
                "endpoint": destination_chain_response.endpoint ,
                "endpoint_public": destination_chain_response.endpoint_public,
                "keystore_path" : destination_chain_response.keystore_path,
                "keypassword": destination_chain_response.keypassword
            }
        


    return struct(
        src_config = src_service_config,
        dst_config = dst_service_config
    )

# Spins up single ICON node
def start_node_service(plan):

    chain_config = icon_node_launcher.get_service_config(ICON0_NODE_ID,ICON0_NODE_PRIVATE_RPC_PORT,ICON0_NODE_PUBLIC_RPC_PORT,ICON0_NODE_P2P_LISTEN_ADDRESS,ICON0_NODE_P2P_ADDRESS,ICON0_NODE_CID)

    node_service_response = start_icon_node_0(plan,chain_config)

    chain_service_config =  {
            
                "service_name" : node_service_response.service_name,
                "nid" : node_service_response.nid,
                "network" : node_service_response.network,
                "network_name": node_service_response.network_name,
                "endpoint": node_service_response.endpoint ,
                "endpoint_public": node_service_response.endpoint_public,
                "keystore_path" : node_service_response.keystore_path,
                "keypassword": node_service_response.keypassword

    }

    return chain_service_config

# Configures ICON Nodes setup
def configure_icon_to_icon_node(plan,src_chain_config,dst_chain_config):

    plan.print("Configuring ICON Nodes")

    setup_node.configure_node(plan,src_chain_config) 
    setup_node.configure_node(plan,dst_chain_config)

# Configures ICON Node setup
def configure_icon_node(plan,chain_config):

    plan.print("configure ICON Node")

    setup_node.configure_node(plan,chain_config) 

# Deploys BMC on ICON
def deploy_bmc_icon(plan,src_chain,dst_chain,args):

    src_config = args["chains"][src_chain]
    
    src_bmc_address = icon_relay_setup.deploy_bmc(plan,src_config)

    if dst_chain == "icon-1":
        dst_config = args["chains"][dst_chain]
        dst_bmc_address = icon_relay_setup.deploy_bmc(plan,dst_config)

        return src_bmc_address , dst_bmc_address

    return src_bmc_address

# Deploys BMV for ICON to ICON setup
def deploy_bmv_icon_to_icon(plan,src_chain,dst_chain,src_bmc_address,dst_bmc_address,args):

    src_chain_config = args["chains"][src_chain]
    dst_chain_config = args["chains"][dst_chain]

    src_chain_service = src_chain_config["service_name"]
    src_chain_network = src_chain_config["network"]
    src_chain_network_name = src_chain_config["network_name"]
    src_chain_keystore_path = src_chain_config["keystore_path"]
    src_chain_keypassword = src_chain_config["keypassword"]
    src_chain_nid = src_chain_config["nid"]
    src_chain_endpoint = src_chain_config["endpoint"]


    dst_chain_service = dst_chain_config["service_name"]
    dst_chain_network = dst_chain_config["network"]
    dst_chain_network_name = dst_chain_config["network_name"]
    dst_chain_keystore_path = dst_chain_config["keystore_path"]
    dst_chain_keypassword = dst_chain_config["keypassword"]
    dst_chain_nid = dst_chain_config["nid"]
    dst_chain_endpoint = dst_chain_config["endpoint"]

    src_last_block_height = setup_node.get_last_block(plan,src_chain_service)
    dst_last_block_height = setup_node.get_last_block(plan,dst_chain_service)

    src_network_name = "{0}-{1}".format(src_chain_network_name,src_last_block_height)
    dst_network_name = "{0}-{1}".format(dst_chain_network_name,dst_last_block_height)

    src_data = {
        "name" : src_network_name,
        "owner" : src_bmc_address
    }

    dst_data = {
        "name": dst_network_name,
        "owner": dst_bmc_address
    }


    src_open_btp_network_response = setup_node.open_btp_network(plan,src_chain_service,src_data,src_chain_endpoint,src_chain_keystore_path,src_chain_keypassword,src_chain_nid)

    dst_open_btp_network_response = setup_node.open_btp_network(plan,dst_chain_service,dst_data,dst_chain_endpoint,dst_chain_keystore_path,dst_chain_keypassword,dst_chain_nid)

   

    src_btp_network_info = setup_node.get_btp_network_info(plan,src_chain_service,src_open_btp_network_response["extract.network_id"])

    src_first_block_header = setup_node.get_btp_header(plan,src_chain_service,src_open_btp_network_response["extract.network_id"],src_btp_network_info)

    dst_btp_network_info = setup_node.get_btp_network_info(plan,dst_chain_service,dst_open_btp_network_response["extract.network_id"])

    dst_first_block_header = setup_node.get_btp_header(plan,dst_chain_service,dst_open_btp_network_response["extract.network_id"],dst_btp_network_info)


    src_bmv_address = icon_relay_setup.deploy_bmv_btpblock_java(plan,src_bmc_address,dst_chain_network,dst_open_btp_network_response["extract.network_type_id"],dst_first_block_header,src_chain_config)

    dst_bmv_address = icon_relay_setup.deploy_bmv_btpblock_java(plan,dst_bmc_address,src_chain_network,src_open_btp_network_response["extract.network_type_id"],src_first_block_header,dst_chain_config)


    src_relay_address = wallet.get_network_wallet_address(plan,src_chain_service)
    dst_relay_address = wallet.get_network_wallet_address(plan,dst_chain_service)

    icon_relay_setup.setup_link_icon(plan,src_chain_service,src_bmc_address,dst_chain_network,dst_bmc_address,src_open_btp_network_response["extract.network_id"],src_bmv_address,src_relay_address,src_chain_config)

    icon_relay_setup.setup_link_icon(plan,dst_chain_service,dst_bmc_address,src_chain_network,src_bmc_address,dst_open_btp_network_response["extract.network_id"],dst_bmv_address,dst_relay_address,dst_chain_config)



    return struct(
        src_bmc = src_bmc_address,
        src_bmv = src_bmv_address,
        dst_bmc = dst_bmc_address,
        dst_bmv = dst_bmv_address,
        src_block_height = src_last_block_height,
        dst_block_height = dst_last_block_height,
        src_network_type_id = src_open_btp_network_response["extract.network_type_id"],
        src_network_id = src_open_btp_network_response["extract.network_id"],
        dst_network_type_id = dst_open_btp_network_response["extract.network_type_id"],
        dst_network_id = dst_open_btp_network_response["extract.network_id"],


    )

# Deploys xCall Contract on ICON nodes
def deploy_xcall_icon(plan,src_chain,dst_chain,src_bmc_address,dst_bmc_address,args):

    src_config = args["chains"][src_chain]
    dst_config = args["chains"][dst_chain]

    src_xcall_address = icon_relay_setup.deploy_xcall(plan,src_bmc_address,src_config)

    if dst_chain == "icon-1":

        dst_xcall_address = icon_relay_setup.deploy_xcall(plan,dst_bmc_address,dst_config)

        return src_xcall_address, dst_xcall_address


    return src_xcall_address

# Deploys dApp Contract on ICON nodes
def deploy_dapp_icon(plan,src_chain,dst_chain,src_xcall_address,dst_xcall_address,args):

    src_config = args["chains"][src_chain]
    dst_config = args["chains"][dst_chain]

    src_dapp_address = icon_relay_setup.deploy_dapp(plan,src_xcall_address,src_config)

    if dst_chain == "icon-1":

        dst_dapp_address = icon_relay_setup.deploy_dapp(plan,dst_xcall_address,dst_config)

        return src_dapp_address,dst_dapp_address

    return src_dapp_address


# Deploy BMV on ICON Node
def deploy_bmv_icon(plan,src_chain,dst_chain,src_bmc_address,dst_bmc_address,dst_last_block_height,args):

    src_chain_config = args["chains"][src_chain]

    src_chain_service = src_chain_config["service_name"]
    src_chain_network = src_chain_config["network"]
    src_chain_network_name = src_chain_config["network_name"]
    src_chain_keystore_path = src_chain_config["keystore_path"]
    src_chain_keypassword = src_chain_config["keypassword"]
    src_chain_nid = src_chain_config["nid"]
    src_chain_endpoint = src_chain_config["endpoint"]

    dst_chain_config = args["chains"][dst_chain]

    dst_chain_service = dst_chain_config["service_name"]
    dst_chain_network = dst_chain_config["network"]
    dst_chain_network_name = dst_chain_config["network_name"]
    dst_chain_keystore_path = dst_chain_config["keystore_path"]
    dst_chain_keypassword = dst_chain_config["keypassword"]
    dst_chain_nid = dst_chain_config["nid"]
    dst_chain_endpoint = dst_chain_config["endpoint"]

    src_chain_last_block_height = setup_node.get_last_block(plan,src_chain_service)



    plan.print("source block height %s" % src_chain_last_block_height)

    network_name = "{0}-{1}".format(dst_chain_network_name,src_chain_last_block_height)

    src_data = {
        "name"  : network_name,
        "owner" : src_bmc_address 
    }

    src_open_btp_net_response = setup_node.open_btp_network(plan,src_chain_service,src_data,src_chain_endpoint,src_chain_keystore_path,src_chain_keypassword,src_chain_nid)

    src_btp_network_info = setup_node.get_btp_network_info(plan,src_chain_service,src_open_btp_net_response["extract.network_id"])

    src_first_block_header = setup_node.get_btp_header(plan,src_chain_service,src_open_btp_net_response["extract.network_id"],src_btp_network_info)

    src_bmv_address = icon_relay_setup.deploy_bmv_bridge_java(plan,src_chain_service,src_bmc_address,dst_chain_network,dst_last_block_height,src_chain_config)

    relay_address = wallet.get_network_wallet_address(plan,src_chain_service)

    icon_relay_setup.setup_link_icon(plan,src_chain_service,src_bmc_address,dst_chain_network,dst_bmc_address,src_open_btp_net_response["extract.network_id"],src_bmv_address,relay_address,src_chain_config)



    return struct(
        bmc =  src_bmc_address,
        bmvbridge = src_bmv_address ,
        network_type_id = src_open_btp_net_response["extract.network_type_id"],
        network_id = src_open_btp_net_response["extract.network_id"],
        block_header = src_first_block_header,
        block_height = src_chain_last_block_height,
        network = src_chain_network

    )


