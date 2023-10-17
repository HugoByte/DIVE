RELAY_SERVICE_IMAGE = 'hugobyte/btp-relay'
RELAY_SERVICE_NAME = "btp-bridge"
RELAY_KEYSTORE_FILES_PATH = "/relay/keystores/"
RELAY_KEYSTORE_PATH = "../static-files/keystores/"
icon_relay_setup = import_module("../../../jvm/icon/src/relay-setup/contract_configuration.star")
icon_service = import_module("../../../jvm/icon/icon.star")
icon_setup_node = import_module("../../../jvm/icon/src/node-setup/setup_icon_node.star")
eth_contract_service = import_module("../../../evm/eth/src/node-setup/contract-service.star")
eth_relay_setup = import_module("../../../evm/eth/src/relay-setup/contract_configuration.star")
eth_node = import_module("../../../evm/eth/eth.star")
input_parser = import_module("../../../../package_io/input_parser.star")




def run_btp_setup(plan, src_chain, dst_chain, bridge):
    if src_chain == "icon" and dst_chain == "icon":
        data = icon_service.start_node_service_icon_to_icon(plan)
        src_chain_service_name = data.src_config["service_name"]
        dst_chain_service_name = data.dst_config["service_name"]

        icon_service.configure_icon_to_icon_node(plan, data.src_config, data.dst_config)

        config = start_btp_for_already_running_icon_nodes(plan, src_chain, dst_chain, data.src_config, data.dst_config, bridge)

        return config
    else:
        if (src_chain == "eth" or src_chain == "hardhat") and dst_chain == "icon":
            dst_chain = src_chain
            src_chain = "icon"

        if dst_chain == "eth" or dst_chain == "hardhat":
            src_chain_config = icon_service.start_node_service(plan)
            dst_chain_config = eth_node.start_eth_node_service(plan, dst_chain)

            src_chain_service_name = src_chain_config["service_name"]
            dst_chain_service_name = dst_chain_config["service_name"]

            icon_service.configure_icon_node(plan, src_chain_service_name, src_chain_config["endpoint"], src_chain_config["keystore_path"], src_chain_config["keypassword"], src_chain_config["nid"])

            config = start_btp_icon_to_eth_for_already_running_nodes(plan, src_chain, dst_chain, src_chain_config, dst_chain_config, bridge)

            return config

        else:
            fail("unsupported chain {0} - {1}".format(src_chain, dst_chain))





def start_btp_for_already_running_icon_nodes(plan, src_chain, dst_chain, src_chain_config, dst_chain_config, bridge):
    """
    Starts BTP for already running ICON nodes.

    Args:
        plan (str):  plan.
        src_chain (str): The source ICON chain name.
        dst_chain (str): The destination ICON chain name.
        src_chain_config (dict): Configuration for the source ICON chain.
        dst_chain_config (dict): Configuration for the destination ICON chain.
        bridge (bool): BMV bridge if true or false.

    Returns:
        dict: New configuration data for BTP.
    """
    # Deploy BMC ICON nodes
    src_bmc_address, dst_bmc_address = icon_service.deploy_bmc_icon(plan, src_chain, dst_chain, src_chain_config, dst_chain_config)

    # Deploy BMV ICON nodes
    response = icon_service.deploy_bmv_icon_to_icon(
        plan, 
        src_chain_config, 
        dst_chain_config, 
        src_bmc_address, 
        dst_bmc_address
    )

    # Deploy XCALL ICON nodes
    src_xcall_address, dst_xcall_address = icon_service.deploy_xcall_icon(
        plan,
        src_chain,
        dst_chain,
        src_chain_config,
        dst_chain_config,
        src_bmc_address,
        dst_bmc_address,
    )

    # Deploy DAPP ICON nodes
    src_dapp_address, dst_dapp_address = icon_service.deploy_dapp_icon(
        plan,
        src_chain,
        dst_chain,
        src_chain_config,
        dst_chain_config,
        src_xcall_address,
        dst_xcall_address
    )

    # Convert hexadecimal block heights to integers
    src_block_height = icon_setup_node.hex_to_int(plan, src_chain_config["service_name"], response.src_block_height)
    dst_block_height = icon_setup_node.hex_to_int(plan, dst_chain_config["service_name"], response.dst_block_height)

    # Create dictionaries for contract addresses
    src_contract_addresses = {
        "bmc": response.src_bmc,
        "bmv": response.src_bmv,
        "xcall": src_xcall_address,
        "dapp": src_dapp_address,
    }

    dst_contract_addresses = {
        "bmc": response.dst_bmc,
        "bmv": response.dst_bmv,
        "xcall": dst_xcall_address,
        "dapp": dst_dapp_address,
    }

    # Generate new configuration data for BTP
    config_data = input_parser.generate_new_config_data_for_btp(src_chain, dst_chain, src_chain_config["service_name"], dst_chain_config["service_name"], bridge)
    config_data["chains"][src_chain_config["service_name"]] = src_chain_config
    config_data["chains"][dst_chain_config["service_name"]] = dst_chain_config
    # Update network and contract information in the configuration data
    config_data["chains"][src_chain_config["service_name"]]["networkTypeId"] = response.src_network_type_id
    config_data["chains"][src_chain_config["service_name"]]["networkId"] = response.src_network_id
    config_data["chains"][dst_chain_config["service_name"]]["networkTypeId"] = response.dst_network_type_id
    config_data["chains"][dst_chain_config["service_name"]]["networkId"] = response.dst_network_id
    config_data["contracts"][src_chain_config["service_name"]] = src_contract_addresses
    config_data["contracts"][dst_chain_config["service_name"]] = dst_contract_addresses
    config_data["chains"][src_chain_config["service_name"]]["block_number"] = src_block_height
    config_data["chains"][dst_chain_config["service_name"]]["block_number"] = dst_block_height

    # Start BTP relayer
    start_btp_relayer(plan, response.src_bmc, response.dst_bmc, src_chain_config, dst_chain_config, bridge)

    # Set source and destination chain names in the configuration data
    config_data["links"]["src"] = src_chain_config["service_name"]
    config_data["links"]["dst"] = dst_chain_config["service_name"]

    return config_data



# Function to start the BTP from ICON to Ethereum for already running nodes.
def start_btp_icon_to_eth_for_already_running_nodes(plan, src_chain, dst_chain, src_chain_config, dst_chain_config, bridge):
    """
    Starts BTP from ICON to Ethereum for already running nodes.

    Args:
        plan (str): plan.
        src_chain (str): The source ICON chain name.
        dst_chain (str): The destination Ethereum chain name.
        src_chain_config (dict): Configuration for the source ICON chain.
        dst_chain_config (dict): Configuration for the destination Ethereum chain.
        bridge (str): The bridge configuration.

    Returns:
        dict: New configuration data for BTP.
    """
    # Start the Ethereum contract service
    eth_contract_service.start_deploy_service(plan, dst_chain_config["endpoint"])

    # Deploy BMC ICON node on the source ICON chain
    src_bmc_address = icon_service.deploy_bmc_icon(plan, src_chain, dst_chain, src_chain_config, dst_chain_config)

    # Deploy BMC Ethereum node on the destination Ethereum chain
    dst_bmc_deploy_response = eth_relay_setup.deploy_bmc(plan, dst_chain, dst_chain_config["network"], dst_chain_config["network_name"])
    dst_bmc_address = dst_bmc_deploy_response.bmc

    # Get the latest block height on the destination Ethereum chain
    dst_last_block_height_number = eth_contract_service.get_latest_block(plan, dst_chain, "localnet")
    dst_last_block_height_hex = icon_setup_node.int_to_hex(plan, src_chain_config["service_name"], dst_last_block_height_number)

    # Deploy BMV ICON node on the source ICON chain
    src_response = icon_service.deploy_bmv_icon(plan, src_bmc_address, dst_bmc_address, src_chain_config, dst_chain_config, dst_last_block_height_hex)

    # Deploy BMV Ethereum node on the destination Ethereum chain
    dst_bmv_address = eth_node.deploy_bmv_eth(plan, bridge, src_response, dst_chain_config["network"], dst_chain_config["network_name"], dst_chain)

    # Deploy XCALL ICON node on the source ICON chain
    src_xcall_address = icon_service.deploy_xcall_icon(
        plan,
        src_chain,
        dst_chain,
        src_chain_config,
        dst_chain_config,
        src_bmc_address,
        dst_bmc_address
    )

    # Deploy XCALL Ethereum node on the destination Ethereum chain
    dst_xcall_address = eth_relay_setup.deploy_xcall(plan, dst_chain, dst_chain_config["network"], dst_chain_config["network_name"])

    # Deploy DAPP ICON node on the source ICON chain
    src_dapp_address = icon_service.deploy_dapp_icon(plan, src_chain, dst_chain, src_chain_config, dst_chain_config, src_xcall_address, dst_xcall_address)

    # Deploy DAPP Ethereum node on the destination Ethereum chain
    dst_dapp_address = eth_relay_setup.deploy_dapp(plan, dst_chain, dst_chain_config["network"], dst_chain_config["network_name"])

    # Convert hexadecimal block height to integer
    src_block_height = icon_setup_node.hex_to_int(plan, src_chain_config["service_name"], src_response.block_height)

    # Create dictionaries for contract addresses
    src_contract_addresses = {
        "bmc": src_response.bmc,
        "bmv": src_response.bmvbridge,
        "xcall": src_xcall_address,
        "dapp": src_dapp_address,
    }

    dst_contract_addresses = {
        "bmc": dst_bmc_address,
        "bmcm": dst_bmc_deploy_response.bmcm,
        "bmcs": dst_bmc_deploy_response.bmcs,
        "bmv": dst_bmv_address,
        "xcall": dst_xcall_address,
        "dapp": dst_dapp_address,
    }

    # Generate new configuration data for BTP
    config_data = input_parser.generate_new_config_data_for_btp(src_chain, dst_chain, src_chain_config["service_name"], dst_chain_config["service_name"], bridge)
    config_data["chains"][src_chain_config["service_name"]] = src_chain_config
    config_data["chains"][dst_chain_config["service_name"]] = dst_chain_config
 
    config_data["contracts"][src_chain_config["service_name"]] = src_contract_addresses
    config_data["contracts"][dst_chain_config["service_name"]] = dst_contract_addresses
    config_data["chains"][src_chain_config["service_name"]]["networkTypeId"] = src_response.network_type_id
    config_data["chains"][src_chain_config["service_name"]]["networkId"] = src_response.network_id
    config_data["chains"][src_chain_config["service_name"]]["block_number"] = src_block_height
    config_data["chains"][dst_chain_config["service_name"]]["block_number"] = dst_last_block_height_number

    # Start BTP relayer
    start_btp_relayer(plan, src_response.bmc, dst_bmc_address, src_chain_config, dst_chain_config, bridge)

    # Set source and destination chain names in the configuration data
    config_data["links"]["src"] = src_chain_config["service_name"]
    config_data["links"]["dst"] = dst_chain_config["service_name"]

    return config_data



def start_btp_relayer(plan, src_bmc, dst_bmc, src_chain_config, dst_chain_config ,bridge):
  
    src_btp_address = "btp://{0}/{1}".format(src_chain_config["network"], src_bmc)
    dst_btp_address = "btp://{0}/{1}".format(dst_chain_config["network"], dst_bmc)
    start_relayer(plan, src_btp_address, dst_btp_address, src_chain_config, dst_chain_config ,bridge)



def start_relayer(plan, src_btp_address, dst_btp_address, src_chain_config, dst_chain_config, bridge):
    """
    Start a BTP Relay Service.

    Args:
        plan (Plan): The kurtosi plan for starting the relay service.
        src_btp_address (str): The source BTP address.
        dst_btp_address (str): The destination BTP address.
        src_chain_config (dict): Configuration for the source chain.
        dst_chain_config (dict): Configuration for the destination chain.
        bridge (str): The bridge mode.

    Note:
        This function starts a BTP Relay Service to relay messages between two blockchain networks.
        It uploads keystore files, configures the service, and starts the relay with the provided parameters.
    """
    plan.print("Starting BTP Relay Service")

    plan.upload_files(src=RELAY_KEYSTORE_PATH, name="keystores")
    
    relay_service = ServiceConfig(
        image=RELAY_SERVICE_IMAGE,
        files={
            RELAY_KEYSTORE_FILES_PATH: "keystores"
        },
        cmd=[
            "/bin/sh",
            "-c",
            "./bin/relay --direction both --log_writer.filename log/relay.log --src.address %s --src.endpoint %s --src.key_store %s --src.key_password %s --src.bridge_mode=%s --dst.address %s --dst.endpoint %s --dst.key_store %s --dst.key_password %s start " % (src_btp_address, src_chain_config["endpoint"], src_chain_config["keystore_path"], src_chain_config["keypassword"], bridge, dst_btp_address, dst_chain_config["endpoint"], dst_chain_config["keystore_path"], dst_chain_config["keypassword"])
        ]
    )

    plan.add_service(name=RELAY_SERVICE_NAME, config=relay_service)
