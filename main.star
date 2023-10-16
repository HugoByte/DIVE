icon_setup_node = import_module("./services/jvm/icon/src/node-setup/setup_icon_node.star")
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
            data = btp_relay_setup.run_btp_setup(plan, args["relay"])

            return data

        elif args["relay"]["name"] == "ibc":
            data = run_cosmos_ibc_setup(plan, args["relay"])

            return data

        else:
            fail("More Relay Support will be added soon")
    if args["action"] == "start_relay":
        if args["relay"]["name"] == "ibc":
            service_name = args["relay"]["service"]
            cosmvm_relay.start_relay(plan, service_name)

def run_node(plan, node_name, args):
    if node_name == "icon":
        return icon_service.start_node_service(plan)

    elif node_name == "eth" or node_name == "hardhat":
        return eth_node.start_eth_node_service(plan, node_name)

    elif node_name == "archway" or node_name == "neutron":
        return cosmvm_node.start_cosmvm_chains(plan, node_name)

    else:
        fail("Unknown Chain Type. Expected ['icon','eth','hardhat','cosmwasm']")




def run_cosmos_ibc_setup(plan, args):
    links = args["links"]
    source_chain = links["src"]
    destination_chain = links["dst"]

    # Check if source and destination chains are both CosmVM-based chains (archway or neutron)
    if (source_chain in ["archway", "neutron"]) and (destination_chain in ["archway", "neutron"]):
        # Start IBC between two CosmVM chains
        data = cosmvm_node.start_ibc_between_cosmvm_chains(plan, source_chain, destination_chain)
        config_data = run_cosmos_ibc_relay_for_already_running_chains(plan, source_chain, destination_chain ,data.src_config, data.dst_config)
        return config_data

    if destination_chain in ["archway", "neutron"] and source_chain == "icon":
        # Start ICON node service
        src_chain_config = icon_service.start_node_service(plan)
        # Start CosmVM node service
        dst_chain_config = cosmvm_node.start_cosmvm_chains(plan, destination_chain)
        dst_chain_config = input_parser.struct_to_dict(dst_chain_config)
        # Get service names and new generate configuration data
        config_data = run_cosmos_ibc_relay_for_already_running_chains(plan,source_chain, destination_chain ,src_chain_config , dst_chain_config)
        return config_data



def run_cosmos_ibc_relay_for_already_running_chains(plan, src_chain, dst_chain, src_chain_config, dst_chain_config):
    source_chain = src_chain
    destination_chain = dst_chain
    if src_chain in ["archway", "neutron"] and dst_chain in ["archway", "neutron"]:
        src_chain_service_name = src_chain_config["service_name"]
        dst_chain_service_name = dst_chain_config["service_name"]
        src_chain_id = src_chain_config["chain_id"]
        src_chain_key = src_chain_config["chain_key"]
        dst_chain_id = dst_chain_config["chain_id"]
        dst_chain_key = dst_chain_config["chain_key"]

        config_data = input_parser.generate_new_config_data_for_ibc(src_chain, dst_chain, src_chain_service_name, dst_chain_service_name)
        config_data["chains"][src_chain_service_name] = src_chain_config
        config_data["chains"][dst_chain_service_name] = dst_chain_config
        cosmvm_relay.start_cosmos_relay(plan, src_chain, dst_chain, src_chain_config, dst_chain_config)

    elif src_chain == "icon" and dst_chain in ["archway", "neutron"]:
        src_chain_service_name = src_chain_config["service_name"]
        dst_chain_service_name = dst_chain_config["service_name"]
        config_data = input_parser.generate_new_config_data_for_ibc(src_chain,dst_chain, src_chain_service_name, dst_chain_service_name)
        # Add chain configurations to the configuration data
        config_data["chains"][src_chain_service_name] = src_chain_config
        config_data["chains"][dst_chain_service_name] = dst_chain_config

        # Setup ICON contracts for IBC
        deploy_icon_contracts = icon_relay_setup.setup_contracts_for_ibc_java(plan, src_chain_config["service_name"], src_chain_config["endpoint"], src_chain_config["keystore_path"], src_chain_config["keypassword"], src_chain_config["nid"], src_chain_config["network"])
        icon_register_client = icon_relay_setup.registerClient(plan, src_chain_service_name, deploy_icon_contracts["light_client"], src_chain_config["keystore_path"], src_chain_config["keypassword"], src_chain_config["nid"], src_chain_config["endpoint"], deploy_icon_contracts["ibc_core"])

        # Configure ICON node
        icon_setup_node.configure_node(plan, src_chain_config["service_name"], src_chain_config["endpoint"], src_chain_config["keystore_path"], src_chain_config["keypassword"], src_chain_config["nid"])

        src_chain_last_block_height = icon_setup_node.get_last_block(plan, src_chain_service_name)

        plan.print("source block height %s" % src_chain_last_block_height)

        network_name = "{0}-{1}".format("dst_chain_network_name", src_chain_last_block_height)

        src_data = {
            "name": network_name,
            "owner": deploy_icon_contracts["ibc_core"],
        }

        # Open BTP network on ICON chain
        tx_result_open_btp_network = icon_setup_node.open_btp_network(plan, src_chain_service_name, src_data, src_chain_config["endpoint"], src_chain_config["keystore_path"], src_chain_config["keypassword"], src_chain_config["nid"])

        icon_bind_port = icon_relay_setup.bindPort(plan, src_chain_service_name, deploy_icon_contracts["xcall_connection"], src_chain_config["keystore_path"], src_chain_config["keypassword"], src_chain_config["nid"], src_chain_config["endpoint"], deploy_icon_contracts["ibc_core"], "xcall")

        # Depending on the destination chain (archway or neutron), set up Cosmos contracts
        if dst_chain == "archway":
            deploy_cosmos_contracts = cosmvm_relay_setup.setup_contracts_for_ibc_wasm(plan, dst_chain_service_name, dst_chain_config["chain_id"], dst_chain_config["chain_key"], dst_chain_config["chain_id"], "stake", "xcall")
            cosmvm_relay_setup.registerClient(plan, dst_chain_service_name, dst_chain_config["chain_id"], dst_chain_config["chain_key"], deploy_cosmos_contracts["ibc_core"], deploy_cosmos_contracts["light_client"])
            plan.wait(service_name = dst_chain_service_name, recipe = ExecRecipe(command = ["/bin/sh", "-c", "sleep 10s && echo 'success'"]), field = "code", assertion = "==", target_value = 0, timeout = "200s")
            cosmvm_relay_setup.bindPort(plan, dst_chain_service_name, dst_chain_config["chain_id"], dst_chain_config["chain_key"], deploy_cosmos_contracts["ibc_core"], deploy_cosmos_contracts["xcall_connection"])
        elif dst_chain == "neutron":
            deploy_cosmos_contracts = neutron_relay_setup.setup_contracts_for_ibc_wasm(plan, dst_chain_service_name, dst_chain_config["chain_id"], dst_chain_config["chain_key"], dst_chain_config["chain_id"], "stake", "xcall")
            neutron_relay_setup.registerClient(plan, dst_chain_service_name, dst_chain_config["chain_id"], dst_chain_config["chain_key"], deploy_cosmos_contracts["ibc_core"], deploy_cosmos_contracts["light_client"])
            plan.wait(service_name = dst_chain_service_name, recipe = ExecRecipe(command = ["/bin/sh", "-c", "sleep 10s && echo 'success'"]), field = "code", assertion = "==", target_value = 0, timeout = "200s")
            neutron_relay_setup.bindPort(plan, dst_chain_service_name, dst_chain_config["chain_id"], dst_chain_config["chain_key"], deploy_cosmos_contracts["ibc_core"], deploy_cosmos_contracts["xcall_connection"])

        # Add contract information to configuration data
        config_data["contracts"][src_chain_service_name] = deploy_icon_contracts
        config_data["contracts"][dst_chain_service_name] = deploy_cosmos_contracts

        src_chain_id = src_chain_config["network_name"].split('-', 1)[1]

        network_id = icon_setup_node.hex_to_int(plan, src_chain_service_name, src_chain_config["nid"])
        btp_network_id = icon_setup_node.hex_to_int(plan, src_chain_service_name, tx_result_open_btp_network["extract.network_id"])
        btp_network_type_id = icon_setup_node.hex_to_int(plan, src_chain_service_name, tx_result_open_btp_network["extract.network_type_id"])

        src_chain_data = {
            "chain_id": src_chain_id,
            "rpc_address": src_chain_config["endpoint"],
            "ibc_address": deploy_icon_contracts["ibc_core"],
            "password": src_chain_config["keypassword"],
            "network_id": network_id,
            "btp_network_id": btp_network_id,
            "btp_network_type_id": btp_network_type_id
        }

        dst_chain_data = {
            "chain_id": dst_chain_config["chain_id"],
            "key": dst_chain_config["chain_key"],
            "rpc_address": dst_chain_config["endpoint"],
            "ibc_address": deploy_cosmos_contracts["ibc_core"],
            "service_name": dst_chain_config["service_name"],
        }

        # Start the Cosmos relay for ICON to Cosmos communication
        relay_service_response = cosmvm_relay.start_cosmos_relay_for_icon_to_cosmos(plan, src_chain, dst_chain ,src_chain_data, dst_chain_data)
        path_name = cosmvm_relay.setup_relay(plan, src_chain_data, dst_chain_data)

        relay_data = cosmvm_relay.get_relay_path_data(plan, relay_service_response.service_name, path_name)

        dapp_result_java = icon_relay_setup.deploy_and_configure_dapp_java(plan, src_chain_config, deploy_icon_contracts["xcall"], dst_chain_config["chain_id"], deploy_icon_contracts["xcall_connection"], deploy_cosmos_contracts["xcall_connection"])

        # Depending on the destination chain (archway or neutron), deploy and configure the DApp for Wasm
        if dst_chain == "archway":
            dapp_result_wasm = cosmvm_relay_setup.deploy_and_configure_xcall_dapp(plan, dst_chain_service_name, dst_chain_config["chain_id"], dst_chain_config["chain_key"], deploy_cosmos_contracts["xcall"], deploy_cosmos_contracts["xcall_connection"], deploy_icon_contracts["xcall_connection"], src_chain_config["network"])
            cosmvm_relay_setup.configure_connection_for_wasm(plan, dst_chain_service_name, dst_chain_config["chain_id"], dst_chain_config["chain_key"], deploy_cosmos_contracts["xcall_connection"], relay_data.dst_connection_id, "xcall", src_chain_config["network"], relay_data.dst_client_id, deploy_cosmos_contracts["xcall"])
        elif dst_chain == "neutron":
            dapp_result_wasm = neutron_relay_setup.deploy_and_configure_xcall_dapp(plan, dst_chain_service_name, dst_chain_config["chain_id"], dst_chain_config["chain_key"], deploy_cosmos_contracts["xcall"], deploy_cosmos_contracts["xcall_connection"], deploy_icon_contracts["xcall_connection"], src_chain_config["network"])
            neutron_relay_setup.configure_connection_for_wasm(plan, dst_chain_service_name, dst_chain_config["chain_id"], dst_chain_config["chain_key"], deploy_cosmos_contracts["xcall_connection"], relay_data.dst_connection_id, "xcall", src_chain_config["network"], relay_data.dst_client_id, deploy_cosmos_contracts["xcall"])

        icon_relay_setup.configure_connection_for_java(plan, deploy_icon_contracts["xcall"], deploy_icon_contracts["xcall_connection"], dst_chain_config["chain_id"], relay_data.src_connection_id, "xcall", dst_chain_config["chain_id"], relay_data.src_client_id, src_chain_service_name, src_chain_config["endpoint"], src_chain_config["keystore_path"], src_chain_config["keypassword"], src_chain_config["nid"])
            
        config_data["contracts"][src_chain_service_name]["dapp"] = dapp_result_java["xcall_dapp"]
        config_data["contracts"][dst_chain_service_name]["dapp"] = dapp_result_wasm["xcall_dapp"]

        # Start relay channel
        cosmvm_relay.start_channel(plan, relay_service_response.service_name, path_name, "xcall", "xcall")


    return config_data
