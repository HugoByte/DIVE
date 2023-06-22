icon_setup_node = import_module("github.com/hugobyte/chain-package/services/jvm/icon/src/node-setup/setup_icon_node.star")
eth_contract_service = import_module("github.com/hugobyte/chain-package/services/evm/eth/src/node-setup/contract-service.star")
eth_relay_setup = import_module("github.com/hugobyte/chain-package/services/evm/eth/src/relay-setup/contract_configuration.star")
eth_node = import_module("github.com/hugobyte/chain-package/services/evm/eth/eth.star")
icon_relay_setup = import_module("github.com/hugobyte/chain-package/services/jvm/icon/src/relay-setup/contract_configuration.star")
icon_service = import_module("github.com/hugobyte/chain-package/services/jvm/icon/icon.star")
btp_relay = import_module("github.com/hugobyte/chain-package/services/relay/btp_relay.star")




def run(plan,args):

    plan.print("Starting")

    links = args["links"]
    source_chain = links["src"]
    destination_chain = links["dst"]

    if destination_chain == "icon":
        destination_chain = "icon-1"

    bridge = args["bridge"]


    config_data = {
        "links": links,
        "chains" : {
            "%s" % source_chain : {},
            "%s" % destination_chain : {}
        },
        "contracts" : {
            "%s" % source_chain : {},
            "%s" % destination_chain : {}
        },
        "bridge" : bridge
    }

    

    if destination_chain == "icon-1":
        data = icon_service.start_node_service_icon_to_icon(plan)

        config_data["chains"][source_chain] = data.src_config
        config_data["chains"][destination_chain] = data.dst_config

        icon_service.configure_icon_to_icon_node(plan,config_data["chains"][source_chain],config_data["chains"][destination_chain])

        src_bmc_address , dst_bmc_address = icon_service.deploy_bmc_icon(plan,source_chain,destination_chain,config_data)

        response = icon_service.deploy_bmv_icon_to_icon(plan,source_chain,destination_chain,src_bmc_address,dst_bmc_address,config_data)

        src_xcall_address , dst_xcall_address = icon_service.deploy_xcall_icon(plan,source_chain,destination_chain,src_bmc_address,dst_bmc_address,config_data)

        src_dapp_address , dst_dapp_address = icon_service.deploy_dapp_icon(plan,source_chain,destination_chain,src_xcall_address,dst_xcall_address,config_data)


        src_block_height = icon_setup_node.hex_to_int(plan,data.src_config["service_name"],response.src_block_height)
        dst_block_height = icon_setup_node.hex_to_int(plan,data.src_config["service_name"],response.dst_block_height)

        src_contract_addresses = {
            "bmc": response.src_bmc,
            "bmv": response.src_bmv,
            "xcall": src_xcall_address,
            "dapp": src_dapp_address,
            "block_number" : src_block_height
        }

        dst_contract_addresses = {
            "bmc": response.dst_bmc,
            "bmv": response.dst_bmv,
            "xcall": dst_xcall_address,
            "dapp": dst_dapp_address,
            "block_number" : dst_block_height
        }

        config_data["chains"][source_chain]["networkTypeId"] = response.src_networkTypeId
        config_data["chains"][source_chain]["networkId"] = response.src_networkId
        config_data["chains"][destination_chain]["networkTypeId"] = response.dst_networkTypeId
        config_data["chains"][destination_chain]["networkId"] = response.dst_networkId

        config_data["contracts"][source_chain] = src_contract_addresses
        config_data["contracts"][destination_chain] = dst_contract_addresses




        
    if destination_chain == "eth":

        src_chain_config = icon_service.start_node_service(plan)

        dst_chain_config = eth_node.start_eth_node_serivce(plan,args)

        config_data["chains"][source_chain] = src_chain_config
        config_data["chains"][destination_chain] = dst_chain_config

        icon_service.configure_icon_node(plan,src_chain_config)

        eth_contract_service.start_deploy_service(plan,dst_chain_config)

        src_bmc_address , empty = icon_service.deploy_bmc_icon(plan,source_chain,destination_chain,config_data)

        dst_bmc_deploy_response = eth_relay_setup.deploy_bmc(plan,config_data)

        dst_bmc_address = dst_bmc_deploy_response.bmc


        dst_last_block_height_number = eth_contract_service.get_latest_block(plan,destination_chain,"localnet")

        dst_last_block_height_hex = icon_setup_node.int_to_hex(plan,src_chain_config["service_name"],dst_last_block_height_number)


        src_response = icon_service.deploy_bmv_icon(plan,source_chain,destination_chain,src_bmc_address ,dst_bmc_address,dst_last_block_height_hex,config_data)

        dst_bmv_address = eth_node.deploy_bmv_eth(plan,bridge,src_response,config_data)


        src_xcall_address , nil = icon_service.deploy_xcall_icon(plan,source_chain,destination_chain,src_bmc_address,dst_bmc_address,config_data)

        dst_xcall_address = eth_relay_setup.deploy_xcall(plan,config_data)

        src_dapp_address , nil = icon_service.deploy_dapp_icon(plan,source_chain,destination_chain,src_xcall_address,dst_xcall_address,config_data)

        dst_dapp_address = eth_relay_setup.deploy_dapp(plan,config_data)

        src_block_height = icon_setup_node.hex_to_int(plan,src_chain_config["service_name"],src_response.block_height)

        src_contract_addresses = {
            "bmc": src_response.bmc,
            "bmv": src_response.bmvbridge,
            "xcall": src_xcall_address,
            "dapp": src_dapp_address,
            "block_number" : src_block_height
        }

        dst_contract_addresses = {
            "bmc": dst_bmc_address,
            "bmcm" : dst_bmc_deploy_response.bmcm,
            "bmcs" : dst_bmc_deploy_response.bmcs,
            "bmv": dst_bmv_address,
            "xcall": dst_xcall_address,
            "dapp": dst_dapp_address,
            "block_number" : dst_last_block_height_number
        }


        config_data["contracts"][source_chain] = src_contract_addresses
        config_data["contracts"][destination_chain] = dst_contract_addresses
        config_data["chains"][source_chain]["networkTypeId"] = src_response.networkTypeId
        config_data["chains"][source_chain]["networkId"] = src_response.networkId


    src_network = config_data["chains"][source_chain]["network"]
    src_bmc = config_data["contracts"][source_chain]["bmc"]

    dst_network = config_data["chains"][destination_chain]["network"]
    dst_bmc = config_data["contracts"][destination_chain]["bmc"]

    src_btp_address = 'btp://{0}/{1}'.format(src_network,src_bmc)
    dst_btp_address = 'btp://{0}/{1}'.format(dst_network,dst_bmc)


    btp_relay.start_relayer(plan,source_chain,destination_chain,config_data,src_btp_address,dst_btp_address,bridge)


    plan.print(config_data)
    



    


    











   


