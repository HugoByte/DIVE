contract_deployment_service = import_module("github.com/hugobyte/chain-package/services/jvm/icon/src/node-setup/contract_deploy.star")
node_service = import_module("github.com/hugobyte/chain-package/services/jvm/icon/src/node-setup/setup_icon_node.star")

def deploy_bmc(plan,service_name,args):
    plan.print("Deploying BMC Contract")

    init_message = args.get("init_message")

    tx_hash = contract_deployment_service.deploy_contract(plan,service_name,"bmc",init_message,args)

    score_address = contract_deployment_service.get_score_address(plan,service_name,tx_hash)

    return score_address


def deploy_xcall(plan,service_name,bmc_address,args):

    plan.print("Deploying xCall Contract")

    init_message = {"key":"_bmc","value":"%s" % bmc_address}

    tx_hash = contract_deployment_service.deploy_contract(plan,service_name,"xcall",init_message,args)

    score_address = contract_deployment_service.get_score_address(plan,service_name,tx_hash)

    return score_address   

 
def open_btp_network(plan,service_name,src,dst,bmc_address,uri,keystorepath,keypassword,nid):
    plan.print("Opening BTP Network")

    last_block_height = node_service.get_last_block(plan,service_name)

    network_name = "{0}-{1}".format(dst,last_block_height)

    args = {"name":network_name,"owner":bmc_address}


    result = node_service.open_btp_network(plan,service_name,args,uri,keystorepath,keypassword,nid)

    return result

def get_first_btpblock_header(plan,service_name,network_id):

    receiptHeight = node_service.get_btp_network_info(plan,service_name,network_id)

    plan.print("receiptHeight %s" % receiptHeight)

    first_block_header = node_service.get_btp_header(plan,service_name,network_id,response)

    return first_block_header


def deploy_bmv_btpblock_java(plan,service_name,bmc_address,src_network_id,network_type_id,block_header,args):

    network_id = args["network_id"]

    first_block_header = get_first_btpblock_header(plan,service_name,network_id)
    init_message = {
      "bmc": bmc_address,
      "srcNetworkID": src_network_id,
      "networkTypeID": network_type_id,
      "blockHeader": first_block_header,
      "seqOffset": "0x0"
    }

    tx_hash = contract_deployment_service.deploy_contract(plan,service_name,"bmv-btpblock",init_message,args)

    score_address = contract_deployment_service.get_score_address(plan,service_name,tx_hash)

    plan.print("BMV-BTPBlock: deployed ")

    return score_address

def deploy_bmv_bridge_java(plan,service_name,bmc_address,dst_network,offset,args):

    init_message = {
        "_bmc": bmc_address,
        "_net": dst_network,
        "_offset": offset
    }

    tx_hash = contract_deployment_service.deploy_contract(plan,service_name,"bmv-bridge",init_message,args)

    score_address = contract_deployment_service.get_score_address(plan,service_name,tx_hash)

    plan.print("BMV-BTPBlock: deployed ")

    return score_address

def add_verifier(plan,service_name,bmc_address,dst_chain_network,bmv_address,uri,keystorepath,keypassword,nid):

    method = "addVerifier"
    params = '{"_net":"%s","_addr":"%s"}' % (dst_chain_network,bmv_address)

    exec_command = ["./bin/goloop","rpc","sendtx","call","--to",bmc_address,"--method",method,"--value",value,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","50000000000","--nid",nid]

    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")


    tx_result = get_tx_result(plan,service_name,tx_hash,uri)

    plan.assert(value=tx_result["extract.status"],assertion="==",target_value="0x1")

    return tx_result


def add_btp_link(plan,service_name,bmc_address,dst_bmc_address,src_network_id,uri,keystorepath,keypassword,nid):

    method = "addBTPLink"

    params = '{"_link":"%s","_networkId":"%s}' %(dst_bmc_address,src_network_id)

    exec_command = ["./bin/goloop","rpc","sendtx","call","--to",bmc_address,"--method",method,"--value",value,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","50000000000","--nid",nid]

    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")


    tx_result = get_tx_result(plan,service_name,tx_hash,uri)

    plan.assert(value=tx_result["extract.status"],assertion="==",target_value="0x1")

    return tx_result


def add_relay(plan,service_name,dst_bmc_address,relay_address,uri,keystorepath,keypassword,nid):

    method = "addRelay"
    params = '{"_link":"%s","_addr":"%s"}' % (dst_bmc_address,relay_address)

    exec_command = ["./bin/goloop","rpc","sendtx","call","--to",bmc_address,"--method",method,"--value",value,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","50000000000","--nid",nid]

    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")


    tx_result = get_tx_result(plan,service_name,tx_hash,uri)

    plan.assert(value=tx_result["extract.status"],assertion="==",target_value="0x1")


    return tx_result


def setup_link_icon(plan,service_name,bmc_address,dst_chain_network,dst_chain_bmc_address,src_chain_network_id,bmv_address,relay_address,args):

    dst_bmc_address = get_btp_address(dst_chain_network,dst_chain_bmc_address)

    plan.print(dst_bmc_address)

    uri = args["uri"]
    keystorepath = args["keystorepath"]
    keypassword = args["keypassword"]
    nid = args["nid"]

    response = add_verifier(plan,service_name,bmc_address,dst_chain_network,bmv_address,uri,keystorepath,keypassword,nid)

    plan.print(response)

    response = add_btp_link(plan,service_name,bmc_address,dst_bmc_address,src_chain_network_id,uri,keystorepath,keypassword,nid)

    plan.print(response)


    response =  add_relay(plan,service_name,dst_bmc_address,relay_address,uri,keystorepath,keypassword,nid)

    plan.print(response)


    plan.print("Icon Link Setup Completed")

def get_btp_address(network,dapp):

    return "btp://{0}/{1}".format(network,dapp)