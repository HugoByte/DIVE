icon_setup_node = import_module("github.com/hugobyte/chain-package/services/jvm/icon/src/node-setup/setup_icon_node.star")
eth_contract_service = import_module("github.com/hugobyte/chain-package/services/evm/eth/src/node-setup/contract-service.star")
eth_relay_setup = import_module("github.com/hugobyte/chain-package/services/evm/eth/src/relay-setup/contract_configuration.star")
eth_node = import_module("github.com/hugobyte/chain-package/services/evm/eth/eth.star")
icon_relay_setup = import_module("github.com/hugobyte/chain-package/services/jvm/icon/src/relay-setup/contract_configuration.star")
RELAY_SERVICE_IMAGE = 'relay'
RELAY_SERVICE_NAME = "btp-relay"
RELAY_CONFIG_FILES_PATH = "/relay/config/"
icon_service = import_module("github.com/hugobyte/chain-package/services/jvm/icon/icon.star")

def start_relay(plan,src_chain,dst_chain,args,src_btp_address,dst_btp_address,bridge):

    plan.print("Starting Relay Service")

    src_config = args["chains"][src_chain]

    src_service_name = src_config["service_name"]

    src_endpoint = src_config["endpoint"]
    src_keystore = src_config["keystore_path"]
    src_keypassword =src_config["keypassword"]

    dst_config = args["chains"][dst_chain]

    dst_service_name = dst_config["service_name"]

    dst_endpoint = dst_config["endpoint"]
    dst_keystore = dst_config["keystore_path"]
    dst_keypassword =dst_config["keypassword"]


    relay_service = ServiceConfig(
        image=RELAY_SERVICE_IMAGE,
        files={
            RELAY_CONFIG_FILES_PATH: "config-files-0"
        },
        cmd=["/bin/sh","-c","./bin/relay --direction both --log_writer.filename log/relay.log --src.address %s --src.endpoint %s --src.key_store %s --src.key_password %s  --src.bridge_mode=%s --dst.address %s --dst.endpoint %s --dst.key_store %s --dst.key_password %s start " % (src_btp_address,src_endpoint,src_keystore,src_keypassword,bridge, dst_btp_address, dst_endpoint, dst_keystore, dst_keypassword)]

    )

    plan.add_service(name=RELAY_SERVICE_NAME,config=relay_service)





def run(plan,args):

    plan.print("Starting")

    links = args["links"]
    source_chain = links["src"]
    destination_chain = links["dst"]

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

        src_contract_addresses = {
            "bmc": response.src_bmc,
            "bmv": response.src_bmv,
            "xcall": src_xcall_address,
            "dapp": src_dapp_address
        }

        dst_contract_addresses = {
            "bmc": response.dst_bmc,
            "bmv": response.dst_bmv,
            "xcall": dst_xcall_address,
            "dapp": dst_dapp_address
        }

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

        dst_bmc_address = eth_relay_setup.deploy_bmc(plan,config_data)

        dst_last_block_height_number = eth_contract_service.get_latest_block(plan,destination_chain,"localnet")

        dst_last_block_height_hex = icon_setup_node.int_to_hex(plan,src_chain_config["service_name"],dst_last_block_height_number)


        src_response = icon_service.deploy_bmv_icon(plan,source_chain,destination_chain,src_bmc_address ,dst_bmc_address,dst_last_block_height_hex,config_data)

        dst_bmv_address = eth_node.deploy_bmv_eth(plan,bridge,src_response,config_data)


        src_xcall_address , nil = icon_service.deploy_xcall_icon(plan,source_chain,destination_chain,src_bmc_address,dst_bmc_address,config_data)

        dst_xcall_address = eth_relay_setup.deploy_xcall(plan,config_data)

        src_dapp_address , nil = icon_service.deploy_dapp_icon(plan,source_chain,destination_chain,src_xcall_address,dst_xcall_address,config_data)

        dst_dapp_address = eth_relay_setup.deploy_dapp(plan,config_data)

        src_contract_addresses = {
            "bmc": src_response.bmc,
            "bmv": src_response.bmvbridge,
            "xcall": src_xcall_address,
            "dapp": src_dapp_address
        }

        dst_contract_addresses = {
            "bmc": dst_bmc_address,
            "bmv": dst_bmv_address,
            "xcall": dst_xcall_address,
            "dapp": dst_dapp_address
        }


        config_data["contracts"][source_chain] = src_contract_addresses
        config_data["contracts"][destination_chain] = dst_contract_addresses


    src_network = config_data["chains"][source_chain]["network"]
    src_bmc = config_data["contracts"][source_chain]["bmc"]

    dst_network = config_data["chains"][destination_chain]["network"]
    dst_bmc = config_data["contracts"][destination_chain]["bmc"]

    src_btp_address = 'btp://{0}/{1}'.format(src_network,src_bmc)
    dst_btp_address = 'btp://{0}/{1}'.format(dst_network,dst_bmc)


    start_relay(plan,source_chain,destination_chain,config_data,src_btp_address,dst_btp_address,bridge)


    plan.print(config_data)
    



    


    











   


