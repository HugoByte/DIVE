icon_setup_node = import_module("github.com/hugobyte/dive/services/jvm/icon/src/node-setup/setup_icon_node.star")
eth_contract_service = import_module("github.com/hugobyte/dive/services/evm/eth/src/node-setup/contract-service.star")
eth_relay_setup = import_module("github.com/hugobyte/dive/services/evm/eth/src/relay-setup/contract_configuration.star")
eth_node = import_module("github.com/hugobyte/dive/services/evm/eth/eth.star")
icon_relay_setup = import_module("github.com/hugobyte/dive/services/jvm/icon/src/relay-setup/contract_configuration.star")
icon_service = import_module("github.com/hugobyte/dive/services/jvm/icon/icon.star")
btp_bridge = import_module("github.com/hugobyte/dive/services/bridges/btp/src/bridge.star")
input_parser = import_module("github.com/hugobyte/dive/package_io/input_parser.star")



def run(plan,args):

    plan.print("Starting")

    return parse_input(plan,args)


def parse_input(plan,args):

    if args["action"] == "start_node":

        node_name = args["node_name"]

        run_node(plan,node_name,args)

   

    if args["action"] == "start_nodes":

        nodes = args["nodes"]

        if len(nodes) == 1:

            if  nodes[0] == "icon":

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
            node_0 = run_node(plan,nodes[0],args)
            node_1 = run_node(plan,nodes[1],args)

            return node_0,node_1


    

    if args["action"] == "setup_relay":

        if args["relay"]["name"] == "btp":
            data = run_btp_setup(plan,args["relay"])
            
            return data

        else:
            fail("More Relay Support will be added soon")


def run_node(plan,node_name,args):

    if node_name == "icon":

        return icon_service.start_node_service(plan)
        
    elif node_name == "eth" or node_name == "hardhat":

        return eth_node.start_eth_node_serivce(plan,args,node_name)

    else :
        fail("Unknown Chain Type. Expected ['icon','eth','hardhat']")


def run_btp_setup(plan,args):

    args_data = input_parser.get_args_data(args)
    
    config_data = input_parser.generate_config_data(args)    

    if args_data.dst == "icon-1":
        data = icon_service.start_node_service_icon_to_icon(plan)

        config_data["chains"][args_data.src] = data.src_config
        config_data["chains"][args_data.dst] = data.dst_config

        icon_service.configure_icon_to_icon_node(plan,config_data["chains"][args_data.src],config_data["chains"][args_data.dst])

        src_bmc_address , dst_bmc_address = icon_service.deploy_bmc_icon(plan,args_data.src,args_data.dst,config_data)

        response = icon_service.deploy_bmv_icon_to_icon(plan,args_data.src,args_data.dst,src_bmc_address,dst_bmc_address,config_data)

        src_xcall_address , dst_xcall_address = icon_service.deploy_xcall_icon(plan,args_data.src,args_data.dst,src_bmc_address,dst_bmc_address,config_data)

        src_dapp_address , dst_dapp_address = icon_service.deploy_dapp_icon(plan,args_data.src,args_data.dst,src_xcall_address,dst_xcall_address,config_data)


        src_block_height = icon_setup_node.hex_to_int(plan,data.src_config["service_name"],response.src_block_height)
        dst_block_height = icon_setup_node.hex_to_int(plan,data.src_config["service_name"],response.dst_block_height)

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

        config_data["chains"][args_data.src]["networkTypeId"] = response.src_network_type_id
        config_data["chains"][args_data.src]["networkId"] = response.src_network_id
        config_data["chains"][args_data.dst]["networkTypeId"] = response.dst_network_type_id
        config_data["chains"][args_data.dst]["networkId"] = response.dst_network_id

        config_data["contracts"][args_data.src] = src_contract_addresses
        config_data["contracts"][args_data.dst] = dst_contract_addresses
        config_data["chains"][args_data.src]["block_number"] = src_block_height
        config_data["chains"][args_data.dst]["block_number"] =  dst_block_height




        
    if args_data.dst == "eth" or args_data.dst == "hardhat":

        src_chain_config = icon_service.start_node_service(plan)

        dst_chain_config = eth_node.start_eth_node_serivce(plan,args,args_data.dst)

        config_data["chains"][args_data.src] = src_chain_config
        config_data["chains"][args_data.dst] = dst_chain_config

        icon_service.configure_icon_node(plan,src_chain_config)

        eth_contract_service.start_deploy_service(plan,dst_chain_config)

        src_bmc_address = icon_service.deploy_bmc_icon(plan,args_data.src,args_data.dst,config_data)

        dst_bmc_deploy_response = eth_relay_setup.deploy_bmc(plan,config_data,args_data.dst)

        dst_bmc_address = dst_bmc_deploy_response.bmc


        dst_last_block_height_number = eth_contract_service.get_latest_block(plan,args_data.dst,"localnet")

        dst_last_block_height_hex = icon_setup_node.int_to_hex(plan,src_chain_config["service_name"],dst_last_block_height_number)


        src_response = icon_service.deploy_bmv_icon(plan,args_data.src,args_data.dst,src_bmc_address ,dst_bmc_address,dst_last_block_height_hex,config_data)

        dst_bmv_address = eth_node.deploy_bmv_eth(plan,args_data.bridge,src_response,config_data,args_data.dst)


        src_xcall_address  = icon_service.deploy_xcall_icon(plan,args_data.src,args_data.dst,src_bmc_address,dst_bmc_address,config_data)

        dst_xcall_address = eth_relay_setup.deploy_xcall(plan,config_data,args_data.dst)

        src_dapp_address = icon_service.deploy_dapp_icon(plan,args_data.src,args_data.dst,src_xcall_address,dst_xcall_address,config_data)

        dst_dapp_address = eth_relay_setup.deploy_dapp(plan,config_data,args_data.dst)

        src_block_height = icon_setup_node.hex_to_int(plan,src_chain_config["service_name"],src_response.block_height)

        src_contract_addresses = {
            "bmc": src_response.bmc,
            "bmv": src_response.bmvbridge,
            "xcall": src_xcall_address,
            "dapp": src_dapp_address,
        }

        dst_contract_addresses = {
            "bmc": dst_bmc_address,
            "bmcm" : dst_bmc_deploy_response.bmcm,
            "bmcs" : dst_bmc_deploy_response.bmcs,
            "bmv": dst_bmv_address,
            "xcall": dst_xcall_address,
            "dapp": dst_dapp_address,
        }


        config_data["contracts"][args_data.src] = src_contract_addresses
        config_data["contracts"][args_data.dst] = dst_contract_addresses
        config_data["chains"][args_data.src]["networkTypeId"] = src_response.network_type_id
        config_data["chains"][args_data.src]["networkId"] = src_response.network_id
        config_data["chains"][args_data.src]["block_number"] = src_block_height
        config_data["chains"][args_data.dst]["block_number"] =  dst_last_block_height_number


    src_network = config_data["chains"][args_data.src]["network"]
    src_bmc = config_data["contracts"][args_data.src]["bmc"]

    dst_network = config_data["chains"][args_data.dst]["network"]
    dst_bmc = config_data["contracts"][args_data.dst]["bmc"]

    src_btp_address = 'btp://{0}/{1}'.format(src_network,src_bmc)
    dst_btp_address = 'btp://{0}/{1}'.format(dst_network,dst_bmc)


    btp_bridge.start_relayer(plan,args_data.src,args_data.dst,config_data,src_btp_address,dst_btp_address,args_data.bridge)


    return config_data


    











   


