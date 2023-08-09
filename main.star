icon_setup_node = import_module("github.com/hugobyte/dive/services/jvm/icon/src/node-setup/setup_icon_node.star")
eth_contract_service = import_module("github.com/hugobyte/dive/services/evm/eth/src/node-setup/contract-service.star")
eth_relay_setup = import_module("github.com/hugobyte/dive/services/evm/eth/src/relay-setup/contract_configuration.star")
eth_node = import_module("github.com/hugobyte/dive/services/evm/eth/eth.star")
icon_relay_setup = import_module("github.com/hugobyte/dive/services/jvm/icon/src/relay-setup/contract_configuration.star")
icon_service = import_module("github.com/hugobyte/dive/services/jvm/icon/icon.star")
btp_bridge = import_module("github.com/hugobyte/dive/services/bridges/btp/src/bridge.star")
input_parser = import_module("github.com/hugobyte/dive/package_io/input_parser.star")
cosmvm_node = import_module("github.com/hugobyte/dive/services/cosmvm/cosmos.star")
cosmvm_relay = import_module("github.com/hugobyte/dive/services/bridges/ibc/src/bridge.star")

def run(plan, args):
    return parse_input(plan, args)

def parse_input(plan, args):
    if args["action"] == "start_node":
        node_name = args["node_name"]

        run_node(plan, node_name, args)

    if args["action"] == "start_nodes":
        nodes = args["nodes"]

        if len(nodes) == 1:
            if nodes[0] == "icon":
                data = icon_service.start_node_service_icon_to_icon(plan)

                return data
            else:
                fail("Only Icon Supported for MutliRun")

        if len(nodes) > 2:
            fail("Maximum two nodes are allowed")

        if nodes[0] == "icon" and nodes[1] == "icon":
            data = icon_service.start_node_service_icon_to_icon(plan)
            return data

        else:
            node_0 = run_node(plan, nodes[0], args)
            node_1 = run_node(plan, nodes[1], args)

            return node_0, node_1

    if args["action"] == "setup_relay":
        if args["relay"]["name"] == "btp":
            data = run_btp_setup(plan, args["relay"])

            return data

        elif args["relay"]["name"] == "cosmos":
            data = run_cosmos_setup(plan, args["relay"])

            return data

        else:
            fail("More Relay Support will be added soon")

def run_node(plan, node_name, args):
    if node_name == "icon":
        return icon_service.start_node_service(plan)

    elif node_name == "eth" or node_name == "hardhat":
        return eth_node.start_eth_node_serivce(plan, args, node_name)

    elif node_name == "cosmos":
        return cosmvm_node.start_node_service(plan, args)

    else:
        fail("Unknown Chain Type. Expected ['icon','eth','hardhat','cosmwasm']")
        return

def run_btp_setup(plan, args):
    args_data = input_parser.get_args_data(args)

    config_data = input_parser.generate_config_data(args)
    if args_data.dst == "icon-1":
        data = icon_service.start_node_service_icon_to_icon(plan)

        config_data["chains"][args_data.src] = data.src_config
        config_data["chains"][args_data.dst] = data.dst_config

        icon_service.configure_icon_to_icon_node(plan, config_data["chains"][args_data.src], config_data["chains"][args_data.dst])

        config = start_btp_for_already_running_icon_nodes(plan, args_data.src, args_data.dst, config_data, data.src_config["service_name"], data.dst_config["service_name"])

        return config

    if args_data.dst == "eth" or args_data.dst == "hardhat":
        src_chain_config = icon_service.start_node_service(plan)

        dst_chain_config = eth_node.start_eth_node_serivce(plan, args, args_data.dst)

        config_data["chains"][args_data.src] = src_chain_config
        config_data["chains"][args_data.dst] = dst_chain_config

        icon_service.configure_icon_node(plan, src_chain_config)

        config = start_btp_icon_to_eth_for_already_running_nodes(plan, args_data.src, args_data.dst, config_data, src_chain_config["service_name"], dst_chain_config["service_name"])

        return config

def start_btp_for_already_running_icon_nodes(plan, src_chain, dst_chain, config_data, src_service_name, dst_src_name):
    src_bmc_address, dst_bmc_address = icon_service.deploy_bmc_icon(plan, src_chain, dst_chain, config_data)

    response = icon_service.deploy_bmv_icon_to_icon(plan, src_chain, dst_chain, src_bmc_address, dst_bmc_address, config_data)

    src_xcall_address, dst_xcall_address = icon_service.deploy_xcall_icon(plan, src_chain, dst_chain, src_bmc_address, dst_bmc_address, config_data)

    src_dapp_address, dst_dapp_address = icon_service.deploy_dapp_icon(plan, src_chain, dst_chain, src_xcall_address, dst_xcall_address, config_data)

    src_block_height = icon_setup_node.hex_to_int(plan, src_service_name, response.src_block_height)
    dst_block_height = icon_setup_node.hex_to_int(plan, dst_src_name, response.dst_block_height)

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

    config_data["chains"][src_chain]["networkTypeId"] = response.src_network_type_id
    config_data["chains"][src_chain]["networkId"] = response.src_network_id
    config_data["chains"][dst_chain]["networkTypeId"] = response.dst_network_type_id
    config_data["chains"][dst_chain]["networkId"] = response.dst_network_id

    config_data["contracts"][src_chain] = src_contract_addresses
    config_data["contracts"][dst_chain] = dst_contract_addresses
    config_data["chains"][src_chain]["block_number"] = src_block_height
    config_data["chains"][dst_chain]["block_number"] = dst_block_height

    config_data = start_btp_relayer(plan, src_chain, dst_chain, config_data)

    return config_data

def start_btp_icon_to_eth_for_already_running_nodes(plan, src_chain, dst_chain, config_data, src_service_name, dst_src_name):
    dst_chain_config = config_data["chains"][dst_chain]
    src_chain_config = config_data["chains"][src_chain]
    eth_contract_service.start_deploy_service(plan, dst_chain_config)

    src_bmc_address = icon_service.deploy_bmc_icon(plan, src_chain, dst_chain, config_data)

    dst_bmc_deploy_response = eth_relay_setup.deploy_bmc(plan, config_data, dst_chain)

    dst_bmc_address = dst_bmc_deploy_response.bmc

    dst_last_block_height_number = eth_contract_service.get_latest_block(plan, dst_chain, "localnet")

    dst_last_block_height_hex = icon_setup_node.int_to_hex(plan, src_chain_config["service_name"], dst_last_block_height_number)

    src_response = icon_service.deploy_bmv_icon(plan, src_chain, dst_chain, src_bmc_address, dst_bmc_address, dst_last_block_height_hex, config_data)

    dst_bmv_address = eth_node.deploy_bmv_eth(plan, config_data["bridge"], src_response, config_data, dst_chain)

    src_xcall_address = icon_service.deploy_xcall_icon(plan, src_chain, dst_chain, src_bmc_address, dst_bmc_address, config_data)

    dst_xcall_address = eth_relay_setup.deploy_xcall(plan, config_data, dst_chain)

    src_dapp_address = icon_service.deploy_dapp_icon(plan, src_chain, dst_chain, src_xcall_address, dst_xcall_address, config_data)

    dst_dapp_address = eth_relay_setup.deploy_dapp(plan, config_data, dst_chain)

    src_block_height = icon_setup_node.hex_to_int(plan, src_chain_config["service_name"], src_response.block_height)

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

    config_data["contracts"][src_chain] = src_contract_addresses
    config_data["contracts"][dst_chain] = dst_contract_addresses
    config_data["chains"][src_chain]["networkTypeId"] = src_response.network_type_id
    config_data["chains"][src_chain]["networkId"] = src_response.network_id
    config_data["chains"][src_chain]["block_number"] = src_block_height
    config_data["chains"][dst_chain]["block_number"] = dst_last_block_height_number

    config_data = start_btp_relayer(plan, src_chain, dst_chain, config_data)

    return config_data

def start_btp_relayer(plan, src_chain, dst_chain, config_data):
    src_network = config_data["chains"][src_chain]["network"]
    src_bmc = config_data["contracts"][src_chain]["bmc"]

    dst_network = config_data["chains"][dst_chain]["network"]
    dst_bmc = config_data["contracts"][dst_chain]["bmc"]

    src_btp_address = "btp://{0}/{1}".format(src_network, src_bmc)
    dst_btp_address = "btp://{0}/{1}".format(dst_network, dst_bmc)

    btp_bridge.start_relayer(plan, src_chain, dst_chain, config_data, src_btp_address, dst_btp_address, config_data["bridge"])

    return config_data

# starts cosmos relay setup

def run_cosmos_setup(plan, args):
    args_data = input_parser.get_args_data(args)

    config_data = input_parser.generate_config_data(args)

    if args_data.dst == "cosmwasm1":
        data, src_service_config, dst_service_config = cosmvm_node.start_node_service_cosmos_to_cosmos(plan)

        config_data["chains"][args_data.src] = data.src_config
        config_data["chains"][args_data.dst] = data.dst_config

        plan.print(config_data)

        cosmvm_relay.start_cosmos_relay(plan, src_service_config.key, src_service_config.cid, dst_service_config.key, dst_service_config.cid, data.src_config, data.dst_config)

    return config_data
