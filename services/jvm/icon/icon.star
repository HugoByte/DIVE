wallet = import_module("./src/node-setup/wallet.star")
setup_node = import_module("./src/node-setup/setup_icon_node.star")
icon_node_launcher = import_module("./src/node-setup/start_icon_node.star")
icon_relay_setup = import_module("./src/relay-setup/contract_configuration.star")

START_FILE_FOR_ICON0 = "start-icon.sh"
START_FILE_FOR_ICON1 = "start-icon.sh"
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
ICON0_GENESIS_FILE_PATH = "../../static-files/config/genesis-icon-0.zip"
ICON1_GENESIS_FILE_PATH = "../../static-files/config/genesis-icon-1.zip"
ICON0_GENESIS_FILE_NAME = "genesis-icon-0.zip"
ICON1_GENESIS_FILE_NAME = "genesis-icon-1.zip"


def start_node_service_icon_to_icon(plan):
    """
    Spin up two ICON nodes, ICON-0 and ICON-1, and return their configuration.

    Args:
        plan (Plan): The kurtosis plan for ICON nodes.

    Returns:
        dict: A dictionary containing the service configuration for ICON-0 and ICON-1 nodes.

    Note:
        This function starts two ICON nodes, ICON-0 and ICON-1, and retrieves their configuration.
        It returns a dictionary with the configuration for both nodes.
    """
    source_chain_response = icon_node_launcher.start_icon_node(
        plan,
        ICON0_NODE_PRIVATE_RPC_PORT,
        ICON0_NODE_PUBLIC_RPC_PORT,
        ICON0_NODE_P2P_LISTEN_ADDRESS,
        ICON0_NODE_P2P_ADDRESS,
        ICON0_NODE_CID,
        {},
        ICON0_GENESIS_FILE_PATH,
        ICON0_GENESIS_FILE_NAME
    )

    destination_chain_response = icon_node_launcher.start_icon_node(
        plan,
        ICON1_NODE_PRIVATE_RPC_PORT,
        ICON1_NODE_PUBLIC_RPC_PORT,
        ICON1_NODE_P2P_LISTEN_ADDRESS,
        ICON1_NODE_P2P_ADDRESS,
        ICON1_NODE_CID,
        {},
        ICON1_GENESIS_FILE_PATH,
        ICON1_GENESIS_FILE_NAME
    )

    src_service_config = {
        "service_name": source_chain_response.service_name,
        "nid": source_chain_response.nid,
        "network": source_chain_response.network,
        "network_name": source_chain_response.network_name,
        "endpoint": source_chain_response.endpoint,
        "endpoint_public": source_chain_response.endpoint_public,
        "keystore_path": source_chain_response.keystore_path,
        "keypassword": source_chain_response.keypassword,
    }

    dst_service_config = {
        "service_name": destination_chain_response.service_name,
        "nid": destination_chain_response.nid,
        "network": destination_chain_response.network,
        "network_name": destination_chain_response.network_name,
        "endpoint": destination_chain_response.endpoint,
        "endpoint_public": destination_chain_response.endpoint_public,
        "keystore_path": destination_chain_response.keystore_path,
        "keypassword": destination_chain_response.keypassword,
    }

    return struct(
        src_config=src_service_config,
        dst_config=dst_service_config,
    )

def start_node_service(plan):
    """
    Spin up a single ICON node and return its configuration.

    Args:
        plan (Plan): The kurtosis plan for starting the ICON node.

    Returns:
        dict: A dictionary containing the configuration for the ICON node.

    Note:
        This function starts a single ICON node and retrieves its configuration.
        It returns a dictionary with the configuration for the node.
    """
    node_service_response = icon_node_launcher.start_icon_node(
        plan,
        ICON0_NODE_PRIVATE_RPC_PORT,
        ICON0_NODE_PUBLIC_RPC_PORT,
        ICON0_NODE_P2P_LISTEN_ADDRESS,
        ICON0_NODE_P2P_ADDRESS,
        ICON0_NODE_CID,
        {},
        ICON0_GENESIS_FILE_PATH,
        ICON0_GENESIS_FILE_NAME
    )

    chain_service_config = {
        "service_name": node_service_response.service_name,
        "nid": node_service_response.nid,
        "network": node_service_response.network,
        "network_name": node_service_response.network_name,
        "endpoint": node_service_response.endpoint,
        "endpoint_public": node_service_response.endpoint_public,
        "keystore_path": node_service_response.keystore_path,
        "keypassword": node_service_response.keypassword,
    }

    return chain_service_config


# Configures ICON Nodes setup
def configure_icon_to_icon_node(plan, src_chain_config, dst_chain_config):
    plan.print("Configuring ICON Nodes")
    setup_node.configure_node(plan, src_chain_config["service_name"], src_chain_config["endpoint"], src_chain_config["keystore_path"], src_chain_config["keypassword"], src_chain_config["nid"])
    setup_node.configure_node(plan, dst_chain_config["service_name"], dst_chain_config["endpoint"], dst_chain_config["keystore_path"], dst_chain_config["keypassword"], dst_chain_config["nid"])


# Configures ICON Node setup
def configure_icon_node(plan, service_name, uri, keystorepath, keypassword, nid):
    plan.print("configure ICON Node")
    setup_node.configure_node(plan, service_name, uri, keystorepath, keypassword, nid)


def deploy_bmc_icon(
    plan, 
    src_chain, 
    dst_chain,
    src_chain_config,
    dst_chain_config
):
    """
    Deploy a BMC (BTP Message center contract) on the source chain and optionally on the destination chain if both are ICON chains.

    Args:
        plan (str): The name of the plan to deploy the BMC contract.
        src_chain (str): The name of the source chain.
        dst_chain (str): The name of the destination chain.
        src_chain_config (dict): Source chain configuration, a dictionary containing the following parameters:
            - "service_name": The name of the source chain's service.
            - "network": The source chain's network.
            - "endpoint": The endpoint URL for the source chain.
            - "keystore_path": The path to the keystore file for the source chain.
            - "keypassword": The password for the keystore.
            - "nid": The Network ID for the source chain.

        dst_chain_config (dict): Destination chain configuration, a dictionary containing the following parameters:
            - "service_name": The name of the destination chain's service.
            - "network": destination chain's network.
            - "endpoint": The endpoint URL for the destination chain.
            - "keystore_path": The path to the keystore file for the destination chain.
            - "keypassword": The password for the keystore.
            - "nid": The Network ID for the destination chain.

    Returns:
        str: The address of the deployed BMC contract on the source chain.
        str (optional): The address of the deployed BMC contract on the destination chain if both are ICON chains.
    """
   
    src_bmc_address = icon_relay_setup.deploy_bmc(plan, src_chain_config["network"], src_chain_config["service_name"], src_chain_config["endpoint"], src_chain_config["keystore_path"], src_chain_config["keypassword"], src_chain_config["nid"])

    if src_chain == "icon" and dst_chain == "icon":
        dst_bmc_address = icon_relay_setup.deploy_bmc(plan, dst_chain_config["network"], dst_chain_config["service_name"], dst_chain_config["endpoint"], dst_chain_config["keystore_path"], dst_chain_config["keypassword"], dst_chain_config["nid"])

        return src_bmc_address, dst_bmc_address

    return src_bmc_address



# Deploys BMV for ICON to ICON setup
def deploy_bmv_icon_to_icon(
    plan,
    src_chain_config,
    dst_chain_config,
    src_bmc_address,
    dst_bmc_address
): 

    """
    Deploy a BMV contract between two ICON networks.

    Args:
        plan: The deployment plan.
        src_chain_config (dict): Source chain configuration, a dictionary containing the following parameters:
            - "service_name": The name of the source chain's service.
            - "network_name": The name of the source chain's network.
            - "endpoint": The endpoint URL for the source chain.
            - "keystore_path": The path to the keystore file for the source chain.
            - "keypassword": The password for the keystore.
            - "nid": The Network ID for the source chain.

        dst_chain_config (dict): Destination chain configuration, a dictionary containing the following parameters:
            - "service_name": The name of the destination chain's service.
            - "network_name": The name of the destination chain's network.
            - "endpoint": The endpoint URL for the destination chain.
            - "keystore_path": The path to the keystore file for the destination chain.
            - "keypassword": The password for the keystore.
            - "nid": The Network ID for the destination chain.

        src_bmc_address (str): Source BMC (Blockchain Management Contract) address.
        dst_bmc_address (str): Destination BMC address.

    Returns:
        A struct containing information about the deployment.
    """

    src_last_block_height = setup_node.get_last_block(plan, src_chain_config["service_name"])
    dst_last_block_height = setup_node.get_last_block(plan, dst_chain_config["service_name"])
    src_network_name = "{0}-{1}".format(src_chain_config["network_name"], src_last_block_height)
    dst_network_name = "{0}-{1}".format(dst_chain_config["network_name"], dst_last_block_height)
    src_data = {
        "name": src_network_name,
        "owner": src_bmc_address,
    }
    dst_data = {
        "name": dst_network_name,
        "owner": dst_bmc_address,
    }
    src_open_btp_network_response = setup_node.open_btp_network(plan, src_chain_config["service_name"], src_data, src_chain_config["endpoint"], src_chain_config["keystore_path"], src_chain_config["keypassword"], src_chain_config["nid"])
    dst_open_btp_network_response = setup_node.open_btp_network(plan, dst_chain_config["service_name"], dst_data, dst_chain_config["endpoint"], dst_chain_config["keystore_path"], dst_chain_config["keypassword"], dst_chain_config["nid"])
    src_btp_network_info = setup_node.get_btp_network_info(plan, src_chain_config["service_name"], src_open_btp_network_response["extract.network_id"])
    src_first_block_header = setup_node.get_btp_header(plan, src_chain_config["service_name"], src_open_btp_network_response["extract.network_id"], src_btp_network_info)
    dst_btp_network_info = setup_node.get_btp_network_info(plan, dst_chain_config["service_name"], dst_open_btp_network_response["extract.network_id"])
    dst_first_block_header = setup_node.get_btp_header(plan, dst_chain_config["service_name"], dst_open_btp_network_response["extract.network_id"], dst_btp_network_info)
    src_bmv_address = icon_relay_setup.deploy_bmv_btpblock_java(plan, src_bmc_address, dst_chain_config["network"], dst_open_btp_network_response["extract.network_type_id"], dst_first_block_header, src_chain_config["service_name"], src_chain_config["endpoint"], src_chain_config["keystore_path"], src_chain_config["keypassword"], src_chain_config["nid"])
    dst_bmv_address = icon_relay_setup.deploy_bmv_btpblock_java(plan, dst_bmc_address, src_chain_config["network"], src_open_btp_network_response["extract.network_type_id"], src_first_block_header, dst_chain_config["service_name"], dst_chain_config["endpoint"], dst_chain_config["keystore_path"], dst_chain_config["keypassword"], dst_chain_config["nid"])
    src_relay_address = wallet.get_network_wallet_address(plan, src_chain_config["service_name"])
    dst_relay_address = wallet.get_network_wallet_address(plan, dst_chain_config["service_name"])
    icon_relay_setup.setup_link_icon(plan, src_chain_config["service_name"], src_bmc_address, dst_chain_config["network"], dst_bmc_address, src_open_btp_network_response["extract.network_id"], src_bmv_address, src_relay_address, src_chain_config["endpoint"], src_chain_config["keystore_path"], src_chain_config["keypassword"], src_chain_config["nid"])
    icon_relay_setup.setup_link_icon(plan, dst_chain_config["service_name"], dst_bmc_address, src_chain_config["network"], src_bmc_address, dst_open_btp_network_response["extract.network_id"], dst_bmv_address, dst_relay_address, dst_chain_config["endpoint"], dst_chain_config["keystore_path"], dst_chain_config["keypassword"], dst_chain_config["nid"])
    
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
def deploy_xcall_icon(plan, src_chain, dst_chain, src_chain_config, dst_chain_config, src_bmc_address, dst_bmc_address):

    src_xcall_address = icon_relay_setup.deploy_xcall(plan, src_bmc_address, src_chain_config["service_name"], src_chain_config["endpoint"], src_chain_config["keystore_path"], src_chain_config["keypassword"], src_chain_config["nid"])

    if src_chain == "icon" and dst_chain == "icon":
        dst_xcall_address = icon_relay_setup.deploy_xcall(plan, dst_bmc_address, dst_chain_config["service_name"], dst_chain_config["endpoint"], dst_chain_config["keystore_path"], dst_chain_config["keypassword"], dst_chain_config["nid"])

        return src_xcall_address, dst_xcall_address

    return src_xcall_address


# Deploys dApp Contract on ICON nodes
def deploy_dapp_icon(plan, src_chain, dst_chain, src_chain_config, dst_chain_config ,src_xcall_address, dst_xcall_address):
    """
    Deploy DApp contract on ICON networks.

    Args:
        plan (str): The deployment plan.
        src_chain (str): The source chain name.
        dst_chain (str): The destination chain name.
        plan: The deployment plan.
        src_chain_config (dict): Source chain configuration, a dictionary containing the following parameters:
            - "service_name": The name of the source chain's service.
            - "network_name": The name of the source chain's network.
            - "endpoint": The endpoint URL for the source chain.
            - "keystore_path": The path to the keystore file for the source chain.
            - "keypassword": The password for the keystore.
            - "nid": The Network ID for the source chain.

        dst_chain_config (dict): Destination chain configuration, a dictionary containing the following parameters:
            - "service_name": The name of the destination chain's service.
            - "network_name": The name of the destination chain's network.
            - "endpoint": The endpoint URL for the destination chain.
            - "keystore_path": The path to the keystore file for the destination chain.
            - "keypassword": The password for the keystore.
            - "nid": The Network ID for the destination chain.

        src_xcall_address (str): The source XCALL contract address.
        dst_xcall_address (str): The destination XCALL contract address.
    
    Returns:
        tuple or str: If both source and destination chains are "icon," returns a tuple of source and destination DApp contract addresses.
                      If only the source chain is "icon," returns the source DApp contract address as a string.
    """
    src_dapp_address = icon_relay_setup.deploy_dapp(plan, src_xcall_address, src_chain_config["service_name"], src_chain_config["endpoint"], src_chain_config["keystore_path"], src_chain_config["keypassword"], src_chain_config["nid"])

    if src_chain == "icon" and dst_chain == "icon":
        dst_dapp_address = icon_relay_setup.deploy_dapp(plan, dst_xcall_address, dst_chain_config["service_name"], dst_chain_config["endpoint"], dst_chain_config["keystore_path"], dst_chain_config["keypassword"], dst_chain_config["nid"])

        return src_dapp_address, dst_dapp_address

    return src_dapp_address


# Deploy BMV on ICON Node
def deploy_bmv_icon(
    plan, 
    src_bmc_address, 
    dst_bmc_address, 
    src_chain_config,
    dst_chain_config,
    dst_last_block_height
):
    """
    Deploy BMV (BTP Multi-Validator) from one ICON network to another ICON network.

    Args:
        plan (str): The deployment plan.
        src_bmc_address (str): The source BMC (Blockchain Management Contract) address.
        dst_bmc_address (str): The destination BMC address.
        plan: The deployment plan.
        src_chain_config (dict): Source chain configuration, a dictionary containing the following parameters:
            - "service_name": The name of the source chain's service.
            - "network_name": The name of the source chain's network.
            - "endpoint": The endpoint URL for the source chain.
            - "keystore_path": The path to the keystore file for the source chain.
            - "keypassword": The password for the keystore.
            - "nid": The Network ID for the source chain.

        dst_chain_config (dict): Destination chain configuration, a dictionary containing the following parameters:
            - "service_name": The name of the destination chain's service.
            - "network_name": The name of the destination chain's network.
            - "endpoint": The endpoint URL for the destination chain.
            - "keystore_path": The path to the keystore file for the destination chain.
            - "keypassword": The password for the keystore.
            - "nid": The Network ID for the destination chain.

        dst_last_block_height (str): The destination chain's last block height.

    Returns:
        dict: A dictionary containing information about the deployment.
    """
    src_chain_last_block_height = setup_node.get_last_block(plan, src_chain_config["service_name"])

    plan.print("source block height %s" % src_chain_last_block_height)

    network_name = "{0}-{1}".format(dst_chain_config["network_name"], src_chain_last_block_height)

    src_data = {
        "name": network_name,
        "owner": src_bmc_address,
    }

    src_open_btp_net_response = setup_node.open_btp_network(plan, src_chain_config["service_name"], src_data, src_chain_config["endpoint"], src_chain_config["keystore_path"], src_chain_config["keypassword"], src_chain_config["nid"])

    src_btp_network_info = setup_node.get_btp_network_info(plan, src_chain_config["service_name"], src_open_btp_net_response["extract.network_id"])

    src_first_block_header = setup_node.get_btp_header(plan, src_chain_config["service_name"], src_open_btp_net_response["extract.network_id"], src_btp_network_info)

    src_bmv_address = icon_relay_setup.deploy_bmv_bridge_java(plan, src_chain_config["service_name"], src_bmc_address, dst_chain_config["network"], dst_last_block_height, src_chain_config["endpoint"], src_chain_config["keystore_path"], src_chain_config["keypassword"], src_chain_config["nid"])

    relay_address = wallet.get_network_wallet_address(plan, src_chain_config["service_name"])

    icon_relay_setup.setup_link_icon(plan, src_chain_config["service_name"], src_bmc_address, dst_chain_config["network"], dst_bmc_address, src_open_btp_net_response["extract.network_id"], src_bmv_address, relay_address, src_chain_config["endpoint"], src_chain_config["keystore_path"], src_chain_config["keypassword"], src_chain_config["nid"])

    return struct(
        bmc = src_bmc_address,
        bmvbridge = src_bmv_address,
        network_type_id = src_open_btp_net_response["extract.network_type_id"],
        network_id = src_open_btp_net_response["extract.network_id"],
        block_header = src_first_block_header,
        block_height = src_chain_last_block_height,
        network = src_chain_config["network"],
    )
