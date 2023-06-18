icon = import_module("github.com/hugobyte/chain-package/services/jvm/icon/main.star")
icon_setup_node = import_module("github.com/hugobyte/chain-package/services/jvm/icon/src/node-setup/setup_icon_node.star")
wallet = import_module("github.com/hugobyte/chain-package/services/jvm/icon/src/node-setup/wallet.star")
eth_contract_service = import_module("github.com/hugobyte/chain-package/services/evm/eth/src/node-setup/contract-service.star")
eth_relay_service = import_module("github.com/hugobyte/chain-package/services/evm/eth/src/relay-setup/contract_configuration.star")
eth_node = import_module("github.com/hugobyte/chain-package/services/evm/eth/src/node-setup/start-eth-node.star")
icon_relay_service = import_module("github.com/hugobyte/chain-package/services/jvm/icon/src/relay-setup/contract_configuration.star")
RELAY_SERVICE_IMAGE = 'alpine'


def run_node_service(plan,args):

    plan.print("Statring Node Service")

    links = args["links"]
    source_chain = links["src"]
    destination_chain = links["dst"]

    bridge = args["bridge"]

    source_chain_response = icon.node_service(plan,args)

    destination_chain_response = eth_node.start_eth_node(plan,args)

    config_data = {
        "links": links,
        "chains" : {
            "%s" % source_chain : {
                "service_name" : source_chain_response.service_name,
                "nid" : source_chain_response.nid,
                "network" : source_chain_response.network,
                "network_name": source_chain_response.network_name,
                "endpoint": source_chain_response.endpoint ,
                "endpoint_public": source_chain_response.endpoint_public,
                "keystore_path" : source_chain_response.keystore_path,
                "keypassword": source_chain_response.keypassword

            },
            "%s" % destination_chain : {
                "service_name" : destination_chain_response.service_name,
                "nid" : destination_chain_response.nid,
                "network" : destination_chain_response.network,
                "network_name": destination_chain_response.network_name,
                "endpoint": destination_chain_response.endpoint ,
                "endpoint_public": "",
                "keystore_path" : "config/eth_keystore.json",
                "keypassword": "password"
            }
        },
        "bridge": bridge
    }

    return config_data

def configure_nodes(plan,config_data):

    plan.print("Configuring Nodes")

    icon_setup_node.configure_node(plan,config_data)

    eth_contract_service.start_deploy_service(plan,config_data)


def deploy_relay_contracts(plan,args):

    plan.print("contracts")


    src_bmc_address = icon_relay_service.deploy_bmc(plan,args)

    dst_bmc_address = eth_relay_service.deploy_bmc(plan,args)


    src_address,dst_address = deploy_bmv_contract(plan,args,src_bmc_address,dst_bmc_address)

    src_xcall_address = icon_relay_service.deploy_xcall(plan,src_bmc_address,args)

    dst_xcall_address = eth_relay_service.deploy_xcall(plan,args)


    src_dapp_address = icon_relay_service.deploy_dapp(plan,src_xcall_address,args)

    dst_dapp_address = eth_relay_service.deploy_dapp(plan,args)


    start_relay(plan,args,src_address,dst_address,args["bridge"])


def deploy_bmv_contract(plan,args,source_bmc_address,dst_bmc_address):

    icon_config_data = args["chains"]["icon"]
    icon_service_name = icon_config_data["service_name"]
    icon_network = icon_config_data["network"]
    icon_keystore_path = icon_config_data["keystore_path"]
    icon_keypassword = icon_config_data["keypassword"]
    icon_nid = icon_config_data["nid"]
    icon_endpoint = icon_config_data["endpoint"]

    bridge = args["bridge"]


    eth_config_data = args["chains"]["eth"]
    dst_name = eth_config_data["network_name"]
    dts_network = eth_config_data["network"]

    src_last_block_height = icon_setup_node.get_last_block(plan,icon_service_name)

    plan.print("Src Chain Block Height %s" % src_last_block_height)

    network_name = "{0}-{1}".format(dst_name,src_last_block_height)

    data = {
        "name": network_name,
        "owner": source_bmc_address
    }

    open_btp_net_response = icon_setup_node.open_btp_network(plan,icon_service_name,data,icon_endpoint,icon_keystore_path,icon_keypassword,icon_nid)

    dst_last_block_height_number = eth_contract_service.get_latest_block(plan,dst_name,"localnet")

    dst_last_block_height_hex = icon_setup_node.int_to_hex(plan,icon_service_name,dst_last_block_height_number)

    plan.print("Dst Chain Block Height %s" % dst_last_block_height_hex)

    src_btp_network_info = icon_setup_node.get_btp_network_info(plan,icon_service_name,open_btp_net_response["extract.network_id"])

    src_first_block_header = icon_setup_node.get_btp_header(plan,icon_service_name,open_btp_net_response["extract.network_id"],src_btp_network_info)

    icon_bmv_address = icon_relay_service.deploy_bmv_bridge_java(plan,icon_service_name,source_bmc_address,dts_network,dst_last_block_height_hex,icon_config_data)

    if bridge == "true":

        eth_relay_service.deploy_bmv_bridge(plan,args,src_last_block_height,source_bmc_address,icon_network)

    else :
        eth_relay_service.deploy_bmv(plan,args,src_first_block_header,source_bmc_address,icon_network,open_btp_net_response["extract.network_type_id"])


    # setup link 

    relay_address = wallet.get_network_wallet_address(plan,icon_service_name)

    icon_relay_service.setup_link_icon(plan,icon_service_name,source_bmc_address,dts_network,dst_bmc_address,open_btp_net_response["extract.network_id"],icon_bmv_address,relay_address,args)


    src_btp_address = 'btp://{0}/{1}'.format(icon_network,source_bmc_address)
    dst_btp_address = 'btp://{0}/{1}'.format(dts_network,dst_bmc_address)

    return src_btp_address,dst_btp_address


def start_relay(plan,args,src_btp_address,dst_btp_address,bridge):

    plan.print("Starting Relay Service")

    src_config = args["chains"]["icon"]

    src_service_name = src_config["service_name"]

    src_endpoint = src_config["endpoint"]
    src_keystore = src_config["keystore_path"]
    src_keypassword =src_config["keypassword"]

    dst_config = args["chains"]["eth"]

    dst_service_name = dst_config["service_name"]

    dst_endpoint = dst_config["endpoint"]
    dst_keystore = dst_config["keystore_path"]
    dst_keypassword =dst_config["keypassword"]


    exec_command = ["./bin/relay","--direction","both","--src.address",src_btp_address,"--src.endpoint",src_endpoint,"--src.key_store",src_keystore,"--src.key_password",src_keypassword,"--src.bridge_mode=%s" % bridge,"--dst.address", dst_btp_address, "--dst.endpoint","http://%s" % dst_endpoint, "--dst.key_store",dst_keystore, "--dst.key_password",dst_keypassword,"start"]

    plan.exec(service_name=src_service_name,recipe=ExecRecipe(command=exec_command))


# def run_icon_node_setup(plan,data):

#     plan.print("Setting Up Icon Node")

#     prep_adddress =  wallet.get_network_wallet_address(plan,node_service.service_name)

#     data = {
#              "service_name": icon_node_service_response.service_name,
#              "prep_address":prep_adddress,
#              "uri":icon_node_service_response.endpoint,
#              "keystorepath":icon_node_service_response.keystore_path ,
#              "keypassword":icon_node_service_response.keypassword,
#              "nid":icon_node_service_response.nid  
#     }        
#     setup_node.configure_node(plan,data)



def run(plan,args):

    plan.print("Starting")

    
    config_data = run_node_service(plan,args)

    configure_nodes(plan,config_data)

    deploy_relay_contracts(plan,config_data)

    


    


    
    




    

    


# def run(plan,args):

#     plan.print("Starting Kurtosis Package")

#     url = eth_node.start_eth_node(plan,args)
    
#     response = eth_contract_service.start_deploy_service(plan,args,url)

#     plan.print(response)

#     contract_address = eth_executor_service.deploy_contract(plan,response.name,"bmc",'{\"link\":\"eth\",\"chainNetwork\":\"0x543.eth\"}',"localnet","")

#     plan.print(contract_address)

#     contract_address = eth_executor_service.deploy_contract(plan,response.name,"bmv",'{ "src": { "name": "eth", "network": "0x543.eth", "networkTypeId": "0x1", "lastBlockHeight": "30375" }, "dst": { "name": "icon", "network": "0x3.icon", "firstBlockHeader": "0xf8468301037f00a0d090304264eeee3c3562152f2dc355601b0b423a948824fd0a012c11c3fc2fb4c00e01f80000f80097d6d594b040bff300eee91f7665ac8dcf89eb0871015306", "bmcAddress": "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9\" } }',"localnet",'bridge="true"')
   
#     contract_address = eth_executor_service.deploy_contract(plan,response.name,"xcall",'{"name":"eth"}',"localnet",'')

#     contract_address = eth_executor_service.deploy_contract(plan,response.name,"dapp",'{"name":"eth"}',"localnet",'')




   
        

