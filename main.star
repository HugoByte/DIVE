icon_setup_node = import_module("./services/jvm/icon/src/node-setup/setup_icon_node.star")
icon_node_launcher = import_module("./src/node-setup/start_icon_node.star")
eth_contract_service = import_module("./services/evm/eth/src/node-setup/contract-service.star")
eth_relay_setup = import_module("./services/evm/eth/src/relay-setup/contract_configuration.star")
eth_node = import_module("./services/evm/eth/eth.star")
icon_relay_setup = import_module("./services/jvm/icon/src/relay-setup/contract_configuration.star")
icon_service = import_module("./services/jvm/icon/icon.star")
btp_bridge = import_module("./services/bridges/btp/src/bridge.star")
input_parser = import_module("./package_io/input_parser.star")
cosmvm_node = import_module("./services/cosmvm/cosmvm.star")
cosmvm_relay = import_module("./services/bridges/ibc/src/bridge.star")
cosmvm_relay_setup = import_module("./services/cosmvm/archway/src/relay-setup/contract-configuration.star")
neutron_relay_setup = import_module("./services/cosmvm/neutron/src/relay-setup/contract-configuration.star")
btp_relay_setup = import_module("./services/bridges/btp/src/bridge.star")
ibc_relay_setup = import_module("./services/bridges/ibc/src/bridge.star")


def run(plan, command, node_name= None, custom_config = None, icon_service_config = None ,decentralize = false, chain_a = None, chain_b = None, service_config_a = None, service_config_b = None, bridge = false):
    """
    Parse the input and execute the specified action.

    Parameters:
        plan (Plan): The Kurtosis plan.
        command (str): The action to perform.
            - 'chain': Start a node.
            - 'bridge': Start a relay.
            - 'decentralize': decentralize already running icon node.

        node_name (str, optional): Name of the node to start.
            - Currently supported options: 'eth', 'hardhat', 'icon', 'neutron', 'archway'.

        custom_config (dict, optional): Custom configuration for node or relay. If empty, the node will start with default settings.
        bridge
        For ICON node with custom configuration, the following fields should be provided in custom_config_dict:
            - private_port (int): The private port for the node.
            - public_port (int): The public port for the node.
            - p2p_listen_address (str): The p2p listen address.
            - p2p_address (str): The p2p address.
            - cid (str): The CID (Chain ID) of the node.
            - genesis_file_path (str): The file path to the genesis file.
            - genesis_file_name (str): The name of the genesis file.

        For Cosmos (Archway/Neutron) node with custom configuration, the following fields should be provided in the custom config dict:
            - chain_id (str): The chain ID.
            - key (str): The key.
            - password (str): The password.
            - public_grpc (str): The public gRPC address.
            - public_http (str): The public HTTP address.
            - public_tcp (str): The public TCP address.
            - public_rpc (str): The public RPC address.

        icon_service_config (dict, optional): ServiceConfig, this field should be provided when wanna decentralize already running icon node

        decentralize (bool, optional): Flag indicating whether to decentralize the ICON node.
        relay_type (str, optional): The type of relay.
            - 'ibc': Start an IBC relay.
            - 'btp': Start a BTP bridge.

        chain_a (str): The source chain for relaying.
        chain_b (str): The destination chain for relaying.
        service_config_a (dict): Service configuration for chain A (source chain for relaying, Note: fields in dictonary should be same as output return after running node).
        service_config_b (dict): Service configuration for chain B (destination chain for relaying, Note: fields in dictonary should be same as output return after running node).
        bridge (bool): Flag indicating whether to use a BMV bridge.

    Returns:
        service_details (dict): Details about the service started.
    """
    return parse_input(plan, command, node_name, custom_config, icon_service_config ,decentralize, chain_a, chain_b, service_config_a, service_config_b, bridge)



def parse_input(plan, action, node_name= None, custom_config = None, icon_service_config = None, decentralize = false, relay_type = None, chain_a = None, chain_b = None, service_config_a = None, service_config_b = None, bridge = false):
    # start the single node
    if action == "chain":
        if node_name != None:
            run_node(plan, node_name, decentralize, custom_config)
        else:
            fail("node_name parameter is missing, node_name require to start node.")


    elif action == "bridge":
        # Start btp relay between two nodes
        if relay_type == "btp":
            if chain_a != None and chain_b != None:
                if service_config_a == None and service_config_b == None:
                    data = btp_relay_setup.run_btp_setup(plan, chain_a, chain_b, bridge)
                elif service_config_a != None and  service_config_b != None:
                    if chain_a == "icon" and chain_b == "icon":
                        data = btp_relay_setup.start_btp_for_already_running_icon_nodes(plan, chain_a, chain_b, service_config_a, service_config_b, bridge)
                    elif chain_a == "icon" and chain_b in ["eth", "hardhat"]:
                        data = btp_relay_setup.start_btp_icon_to_eth_for_already_running_nodes(plan, chain_a, chain_b, service_config_a, service_config_b, bridge)
                    else:
                        fail("unsupported chain {0} - {1}".format(src_chain, dst_chain))
                else: 
                    fail("Add Service configs for both chain_a and chain_b")
                return data
            else:
                fail("chain_a and chain_b paramter are missing, Add chain_a and chain_b to start relay between them.")

        # Start ibc relay between two nodes
        elif args["relay"]["name"] == "ibc":
            if chain_a != None and chain_b != None:
                if service_config_a == None and service_config_b == None:
                    data = ibc_relay_setup.run_cosmos_ibc_setup(plan, chain_a, chain_b)
                elif service_config_a != None and service_config_b != None:
                    data = ibc_relay_setup.run_cosmos_ibc_relay_for_already_running_chains(plan, chain_a, chain_b, service_config_a, service_config_b)
                else: 
                    fail("Add Service configs for both chain_a and chain_b")
                return data
            else:
                fail("chain_a and chain_b paramter are missing, Add chain_a and chain_b to start relay between them.")

        else:
            fail("More Relay Support will be added soon")
    
    elif action == "decentralize":
        if icon_service_config != None:
            icon_setup_node.configure_node(icon_service_config["service_name"], icon_service_config["enpoint"], icon_service_config["keystore_path"], icon_service_config["keypassword"], icon_service_config["nid"])
        else: fail("icon_service_config paramter is missing, Add icon_service_config to decentralise the already running icon node")
    
    else: 
        fail("commands only support 'chain', 'bridge' and 'decentralize'")



def run_node(plan, node_name, decentralize, custom_config = None):
    if node_name == "icon":
        if custom_config == None:
            service_config = icon_service.start_node_service(plan)
        else: 
            service_config = icon_node_launcher.start_icon_node(plan, custom_config["private_port"], custom_config["private_port"], custom_config["p2p_listen_address"], custom_config["p2p_address"], custom_config["cid"], {}, custom_config["genesis_file_path"], custom_config["genesis_file_name"])

        if decentralize == true:
            icon_setup_node.configure_node(service_config["service_name"], service_config["enpoint"], service_config["keystore_path"], service_config["keypassword"], service_config["nid"])
        
        return service_config

    elif node_name == "eth" or node_name == "hardhat":
        return eth_node.start_eth_node_service(plan, node_name)

    elif node_name == "archway" or node_name == "neutron":
        if cutsom_config == None:
            return cosmvm_node.start_cosmvm_chains(plan, node_name)
        else:
            return cosmvm_node.start_cosmvm_chains(plan, node_name, custom_config["chain_id"], custom_config["key"], custom_config["password"], custom_config["public_grpc"], custom_config["public_http"], custom_config["public_tcp"], custom_config["public_rpc"])

    else:
        fail("Unknown Chain Type. Expected ['icon','eth','hardhat','cosmwasm']")

