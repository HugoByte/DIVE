icon_setup_node = import_module("github.com/hugobyte/dive/services/jvm/icon/src/node-setup/setup_icon_node.star")
eth_contract_service = import_module("github.com/hugobyte/dive/services/evm/eth/src/node-setup/contract-service.star")
eth_relay_setup = import_module("github.com/hugobyte/dive/services/evm/eth/src/relay-setup/contract_configuration.star")
eth_node = import_module("github.com/hugobyte/dive/services/evm/eth/eth.star")
icon_relay_setup = import_module("github.com/hugobyte/dive/services/jvm/icon/src/relay-setup/contract_configuration.star")
icon_service = import_module("github.com/hugobyte/dive/services/jvm/icon/icon.star")
btp_bridge = import_module("github.com/hugobyte/dive/services/bridges/btp/src/bridge.star")
input_parser = import_module("github.com/hugobyte/dive/package_io/input_parser.star")
cosmvm_node = import_module("github.com/hugobyte/dive/services/cosmvm/cosmvm.star")
cosmvm_relay = import_module("github.com/hugobyte/dive/services/bridges/ibc/src/bridge.star")
cosmvm_relay_setup = import_module("github.com/hugobyte/dive/services/cosmvm/archway/src/relay-setup/contract-configuration.star")

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

        elif args["relay"]["name"] == "ibc":
            data = run_cosmos_ibc_setup(plan, args["relay"])

            return data

        else:
            fail("More Relay Support will be added soon")

def run_node(plan, node_name, args):
    if node_name == "icon":
        return icon_service.start_node_service(plan)

    elif node_name == "eth" or node_name == "hardhat":
        return eth_node.start_eth_node_serivce(plan, args, node_name)

    elif node_name == "archway":
        return cosmvm_node.start_cosmvm_chains(plan,node_name,args)

    else:
        fail("Unknown Chain Type. Expected ['icon','eth','hardhat','cosmwasm']")

def run_btp_setup(plan, args):
    links = args["links"]
    source_chain = links["src"]
    destination_chain = links["dst"]
    bridge = args["bridge"]

    if source_chain == "icon" and destination_chain == "icon":
        data = icon_service.start_node_service_icon_to_icon(plan)
        src_chain_service_name = data.src_config["service_name"]
        dst_chain_service_name = data.dst_config["service_name"]

        config_data = input_parser.generate_new_config_data(links, src_chain_service_name, dst_chain_service_name, bridge)
        config_data["chains"][src_chain_service_name] = data.src_config
        config_data["chains"][dst_chain_service_name] = data.dst_config

        icon_service.configure_icon_to_icon_node(plan, data.src_config, data.dst_config)

        config = start_btp_for_already_running_icon_nodes(plan, source_chain, destination_chain, config_data, data.src_config["service_name"], data.dst_config["service_name"])

        return config
    else:

        if (source_chain == "eth" or source_chain == "hardhat") and destination_chain == "icon":
            
            destination_chain = source_chain
            source_chain = "icon"
        
        if destination_chain == "eth" or destination_chain == "hardhat":
            src_chain_config = icon_service.start_node_service(plan)
            dst_chain_config = eth_node.start_eth_node_serivce(plan, args, destination_chain)

            src_chain_service_name = src_chain_config["service_name"]
            dst_chain_service_name = dst_chain_config["service_name"]

            config_data = input_parser.generate_new_config_data(links, src_chain_service_name, dst_chain_service_name, bridge)
            config_data["chains"][src_chain_service_name] = src_chain_config
            config_data["chains"][dst_chain_service_name] = dst_chain_config

            icon_service.configure_icon_node(plan, src_chain_config)
            config = start_btp_icon_to_eth_for_already_running_nodes(plan, source_chain,destination_chain, config_data,  src_chain_service_name, dst_chain_service_name)

            return config
            
        else:
            fail("unsupported chain {0} - {1}".format(source_chain,destination_chain))
    

def start_btp_for_already_running_icon_nodes(plan, src_chain, dst_chain, config_data, src_service_name, dst_service_name):
    src_bmc_address, dst_bmc_address = icon_service.deploy_bmc_icon(plan, src_chain, dst_chain,src_service_name, dst_service_name, config_data)

    response = icon_service.deploy_bmv_icon_to_icon(plan, src_service_name, dst_service_name, src_bmc_address, dst_bmc_address, config_data)

    src_xcall_address, dst_xcall_address = icon_service.deploy_xcall_icon(plan, src_chain, dst_chain, src_bmc_address, dst_bmc_address, config_data,src_service_name, dst_service_name)

    src_dapp_address, dst_dapp_address = icon_service.deploy_dapp_icon(plan, src_chain, dst_chain, src_xcall_address, dst_xcall_address, config_data,src_service_name, dst_service_name)

    src_block_height = icon_setup_node.hex_to_int(plan, src_service_name, response.src_block_height)
    dst_block_height = icon_setup_node.hex_to_int(plan, dst_service_name, response.dst_block_height)

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

    config_data["chains"][src_service_name]["networkTypeId"] = response.src_network_type_id
    config_data["chains"][src_service_name]["networkId"] = response.src_network_id
    config_data["chains"][dst_service_name]["networkTypeId"] = response.dst_network_type_id
    config_data["chains"][dst_service_name]["networkId"] = response.dst_network_id

    config_data["contracts"][src_service_name] = src_contract_addresses
    config_data["contracts"][dst_service_name] = dst_contract_addresses
    config_data["chains"][src_service_name]["block_number"] = src_block_height
    config_data["chains"][dst_service_name]["block_number"] = dst_block_height

    config_data = start_btp_relayer(plan, src_chain, dst_chain, config_data,src_service_name,dst_service_name)

    config_data["links"]["src"] = src_service_name
    config_data["links"]["dst"] = dst_service_name

    return config_data

def start_btp_icon_to_eth_for_already_running_nodes(plan, src_chain, dst_chain, config_data, src_service_name, dst_service_name):
    dst_chain_config = config_data["chains"][dst_service_name]
    src_chain_config = config_data["chains"][src_service_name]

    eth_contract_service.start_deploy_service(plan, dst_chain_config)

    src_bmc_address = icon_service.deploy_bmc_icon(plan, src_chain, dst_chain,src_service_name, dst_service_name,config_data)

    dst_bmc_deploy_response = eth_relay_setup.deploy_bmc(plan, config_data, dst_chain,dst_service_name)

    dst_bmc_address = dst_bmc_deploy_response.bmc

    dst_last_block_height_number = eth_contract_service.get_latest_block(plan, dst_chain, "localnet")

    dst_last_block_height_hex = icon_setup_node.int_to_hex(plan, src_service_name, dst_last_block_height_number)

    src_response = icon_service.deploy_bmv_icon(plan, src_service_name, dst_service_name, src_bmc_address, dst_bmc_address, dst_last_block_height_hex, config_data)

    dst_bmv_address = eth_node.deploy_bmv_eth(plan, config_data["bridge"], src_response, config_data, dst_chain,dst_service_name)

    src_xcall_address = icon_service.deploy_xcall_icon(plan, src_chain, dst_chain, src_bmc_address, dst_bmc_address, config_data,src_service_name, dst_service_name)

    dst_xcall_address = eth_relay_setup.deploy_xcall(plan, config_data, dst_chain,dst_service_name)

    src_dapp_address = icon_service.deploy_dapp_icon(plan, src_chain, dst_chain, src_xcall_address, dst_xcall_address, config_data,src_service_name, dst_service_name)

    dst_dapp_address = eth_relay_setup.deploy_dapp(plan, config_data, dst_chain,dst_service_name)

    src_block_height = icon_setup_node.hex_to_int(plan, src_service_name, src_response.block_height)

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

    config_data["contracts"][src_service_name] = src_contract_addresses
    config_data["contracts"][dst_service_name] = dst_contract_addresses
    config_data["chains"][src_service_name]["networkTypeId"] = src_response.network_type_id
    config_data["chains"][src_service_name]["networkId"] = src_response.network_id
    config_data["chains"][src_service_name]["block_number"] = src_block_height
    config_data["chains"][dst_service_name]["block_number"] = dst_last_block_height_number

    config_data = start_btp_relayer(plan, src_chain, dst_chain, config_data,src_service_name,dst_service_name)
    
    config_data["links"]["src"] = src_service_name
    config_data["links"]["dst"] = dst_service_name

    return config_data

def start_btp_relayer(plan, src_chain, dst_chain, config_data,src_service_name,dst_service_name):
    src_network = config_data["chains"][src_service_name]["network"]
    src_bmc = config_data["contracts"][src_service_name]["bmc"]

    dst_network = config_data["chains"][dst_service_name]["network"]
    dst_bmc = config_data["contracts"][dst_service_name]["bmc"]

    src_btp_address = "btp://{0}/{1}".format(src_network, src_bmc)
    dst_btp_address = "btp://{0}/{1}".format(dst_network, dst_bmc)

    btp_bridge.start_relayer(plan, src_service_name, dst_service_name, config_data, src_btp_address, dst_btp_address, config_data["bridge"])

    return config_data

# starts cosmos ibc relay setup
def run_cosmos_ibc_setup(plan, args):
    links = args["links"]
    source_chain = links["src"]
    destination_chain = links["dst"]

    if source_chain == "archway" and destination_chain == "archway":
        data = cosmvm_node.start_ibc_between_cosmvm_chains(plan,source_chain,destination_chain)
                
        config_data = run_cosmos_ibc_relay_for_already_running_chains(plan,links,data.src_config,data.dst_config)
        return config_data

    if destination_chain == "archway":

        src_chain_config = icon_service.start_node_service(plan)
        data = {"data":{}}
        dst_chain_config = cosmvm_node.start_cosmvm_chains(plan, destination_chain,data) 

        src_chain_service_name = src_chain_config["service_name"]
        dst_chain_service_name = dst_chain_config.service_name

        config_data = input_parser.generate_new_config_data(links, src_chain_service_name, dst_chain_service_name,"")

        config_data["chains"][src_chain_service_name] = src_chain_config
        config_data["chains"][dst_chain_service_name] = dst_chain_config

        deploy_icon_contracts = icon_relay_setup.setup_contracts_for_ibc_java(plan,src_chain_config)

        icon_register_client = icon_relay_setup.registerClient(plan,src_chain_service_name,args,deploy_icon_contracts["light_client"],src_chain_config["keystore_path"],src_chain_config["keypassword"],src_chain_config["nid"],src_chain_config["endpoint"],deploy_icon_contracts["ibc_core"])

        icon_bind_port = icon_relay_setup.bindPort(plan,src_chain_service_name,args,deploy_icon_contracts["xcall_connection"],src_chain_config["keystore_path"],src_chain_config["keypassword"],src_chain_config["nid"],src_chain_config["endpoint"],deploy_icon_contracts["ibc_core"],"xcall")

        icon_setup_node.configure_node(plan,src_chain_config)

        src_chain_last_block_height = icon_setup_node.get_last_block(plan,src_chain_service_name)

        plan.print("source block height %s" % src_chain_last_block_height)

        network_name = "{0}-{1}".format("dst_chain_network_name",src_chain_last_block_height)

        src_data = {
            "name"  : network_name,
            "owner" : deploy_icon_contracts["ibc_core"] 
        }

        icon_setup_node.open_btp_network(plan,src_chain_service_name,src_data,src_chain_config["endpoint"],src_chain_config["keystore_path"], "gochain",src_chain_config["nid"])

        deploy_archway_contracts = cosmvm_relay_setup.setup_contracts_for_ibc_wasm(plan,dst_chain_service_name,dst_chain_config.chain_id,dst_chain_config.chain_key,dst_chain_config.chain_id,"stake","xcall")

        cosmvm_relay_setup.registerClient(plan,dst_chain_service_name,dst_chain_config.chain_id,dst_chain_config.chain_key,deploy_archway_contracts["ibc_core"],deploy_archway_contracts["light_client"])

        cosmvm_relay_setup.bindPort(plan,dst_chain_service_name,dst_chain_config.chain_id,dst_chain_config.chain_key,deploy_archway_contracts["ibc_core"],deploy_archway_contracts["xcall_connection"])

        src_contract_address = {
            "contracts" : deploy_icon_contracts,
        }

        dst_contract_address = {
            "contracts_archway" : deploy_archway_contracts,
        }

        config_data["contracts"][src_chain_service_name] = src_contract_address
        config_data["contracts"][dst_chain_service_name] = dst_contract_address

        cosmos = cosmvm_relay.start_cosmos_relay_for_icon_to_cosmos(plan,args)

        SEED0 = plan.exec(service_name=dst_chain_service_name, recipe=ExecRecipe(command=["/bin/sh", "-c", "jq -r '.mnemonic' ../../start-scripts/key_seed.json | tr -d '\n\r'"]))

        plan.exec(service_name="ibc-relayer", recipe=ExecRecipe(command=["/bin/sh", "-c", "sed -i -e 's|\"ibc-handler-address\":\"\"|\"ibc-handler-address\": \"'%s'\"|' ../script/icon.json" % (deploy_icon_contracts["ibc_core"])]))

        plan.exec(service_name="ibc-relayer", recipe=ExecRecipe(command=["/bin/sh", "-c", "sed -i -e 's|\"ibc-handler-address\":\"\"|\"ibc-handler-address\": \"'%s'\"|' ../script/archway1.json" % (deploy_archway_contracts["ibc_core"])]))

        #  plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "sed -i -e 's|\"rpc-addr\": \"\"|\"rpc-addr\": \"http://'%s'\"|' ../script/chains/archway1.json" % (["endpoint"])]) )
        cosmvm_relay.setup_relay(plan,args,SEED0)


        # config_data = run_cosmos_ibc_relay_for_already_running_chains(plan,links,src_chain_config,dst_chain_config)
        return config_data

    
        
def run_cosmos_ibc_relay_for_already_running_chains(plan,links,src_config,dst_config):

    src_chain_service_name = src_config["service_name"]
    dst_chain_service_name = dst_config["service_name"]
    src_chain_id = src_config["chain_id"]
    src_chain_key = src_config["chain_key"]
    dst_chain_id = dst_config["chain_id"]
    dst_chain_key = dst_config["chain_key"]

    config_data = input_parser.generate_new_config_data_cosmvm_cosmvm(links, src_chain_service_name, dst_chain_service_name)
    config_data["chains"][src_chain_service_name] = src_config
    config_data["chains"][dst_chain_service_name] = dst_config
    cosmvm_relay.start_cosmos_relay(plan, src_chain_key, src_chain_id, dst_chain_key, dst_chain_id, src_config, dst_config)

    return config_data
