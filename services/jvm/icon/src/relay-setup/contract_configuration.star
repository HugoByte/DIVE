contract_deployment_service = import_module("github.com/hugobyte/dive/services/jvm/icon/src/node-setup/contract_deploy.star")
node_service = import_module("github.com/hugobyte/dive/services/jvm/icon/src/node-setup/setup_icon_node.star")

# Deploys BMC contract on ICON 
def deploy_bmc(plan,args):
    plan.print("Deploying BMC Contract")
    init_message = '{"_net":"%s"}' % args["network"]

    tx_hash = contract_deployment_service.deploy_contract(plan,"bmc",init_message,args)

    service_name = args["service_name"]
    score_address = contract_deployment_service.get_score_address(plan,service_name,tx_hash)
    return score_address

# Deploys xCall on ICON
def deploy_xcall(plan,bmc_address,args):

    plan.print("Deploying xCall Contract")
    init_message = '{"_bmc":"%s"}' % bmc_address

    tx_hash = contract_deployment_service.deploy_contract(plan,"xcall",init_message,args)
    service_name = args["service_name"]

    score_address = contract_deployment_service.get_score_address(plan,service_name,tx_hash)
    add_service(plan,bmc_address,score_address,args)
    return score_address   

# Adds services to BMC contract on ICON
def add_service(plan,bmc_address,xcall_address,args):

    plan.print("Adding xcall  to Bmc %s " % bmc_address)

    service_name = args["service_name"]
    uri = args["endpoint"]
    keystore_path = args["keystore_path"]
    keypassword = args["keypassword"]
    nid = args["nid"]

    method = "addService"
    params = '{"_svc":"xcall","_addr":"%s"}' % xcall_address

    exec_command = ["./bin/goloop","rpc","sendtx","call","--to",bmc_address,"--method",method,"--params",params,"--uri",uri,"--key_store",keystore_path,"--key_password",keypassword,"--step_limit","50000000000","--nid",nid]
    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")
    tx_result = node_service.get_tx_result(plan,service_name,tx_hash,uri)
    plan.assert(value=tx_result["extract.status"],assertion="==",target_value="0x1")



# Opens BTP Network on ICON
def open_btp_network(plan,service_name,src,dst,bmc_address,uri,keystorepath,keypassword,nid):
    plan.print("Opening BTP Network")

    last_block_height = node_service.get_last_block(plan,service_name)
    network_name = "{0}-{1}".format(dst,last_block_height)

    args = {"name":network_name,"owner":bmc_address}

    result = node_service.open_btp_network(plan,service_name,args,uri,keystorepath,keypassword,nid)
    return result

# Deploys BMV BTPBLOCK on ICON
def deploy_bmv_btpblock_java(plan,bmc_address,dst_network_id,dst_network_type_id,first_block_header,args):
    
    init_message = '{"bmc": "%s","srcNetworkID": "%s","networkTypeID": "%s", "blockHeader": "0x%s","seqOffset": "0x0"}' % (bmc_address,dst_network_id,dst_network_type_id,first_block_header)
    service_name = args["service_name"]

    tx_hash = contract_deployment_service.deploy_contract(plan,"bmv-btpblock",init_message,args)
    score_address = contract_deployment_service.get_score_address(plan,service_name,tx_hash)

    plan.print("BMV-BTPBlock: deployed")
    return score_address

# Deploys BMVBRIDGE on ICON
def deploy_bmv_bridge_java(plan,service_name,bmc_address,dst_network,offset,args):

    init_message = '{"_bmc": "%s","_net": "%s","_offset": "%s"}' %(bmc_address,dst_network,offset)
    tx_hash = contract_deployment_service.deploy_contract(plan,"bmv-bridge",init_message,args)

    score_address = contract_deployment_service.get_score_address(plan,service_name,tx_hash)
    plan.print("BMV-BTPBlock: deployed ")
    return score_address

# Adds Verifier to BMC contract on ICON
def add_verifier(plan,service_name,bmc_address,dst_chain_network,bmv_address,uri,keystorepath,keypassword,nid):

    method = "addVerifier"
    params = '{"_net":"%s","_addr":"%s"}' % (dst_chain_network,bmv_address)

    exec_command = ["./bin/goloop","rpc","sendtx","call","--to",bmc_address,"--method",method,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","50000000000","--nid",nid]
    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")
    tx_result = node_service.get_tx_result(plan,service_name,tx_hash,uri)

    plan.assert(value=tx_result["extract.status"],assertion="==",target_value="0x1")
    return tx_result

# Adds BTP Link to BMC contract on ICON
def add_btp_link(plan,service_name,bmc_address,dst_bmc_address,src_network_id,uri,keystorepath,keypassword,nid):

    method = "addBTPLink"

    params = '{"_link":"%s","_networkId":"%s"}' %(dst_bmc_address,src_network_id)

    exec_command = ["./bin/goloop","rpc","sendtx","call","--to",bmc_address,"--method",method,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","50000000000","--nid",nid]
    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")


    tx_result = node_service.get_tx_result(plan,service_name,tx_hash,uri)
    plan.assert(value=tx_result["extract.status"],assertion="==",target_value="0x1")

    return tx_result

# Adds Relay Address to BMC contract on ICON
def add_relay(plan,service_name,bmc_address,dst_bmc_address,relay_address,uri,keystorepath,keypassword,nid):

    method = "addRelay"
    params = '{"_link":"%s","_addr":"%s"}' % (dst_bmc_address,relay_address)

    exec_command = ["./bin/goloop","rpc","sendtx","call","--to",bmc_address,"--method",method,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","500000000000","--nid",nid]
    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")
    tx_result = node_service.get_tx_result(plan,service_name,tx_hash,uri)
    plan.assert(value=tx_result["extract.status"],assertion="==",target_value="0x1")
    return tx_result

# Configures Link in BMC on ICON
def setup_link_icon(plan,service_name,bmc_address,dst_chain_network,dst_chain_bmc_address,src_chain_network_id,bmv_address,relay_address,args):

    dst_bmc_address = get_btp_address(dst_chain_network,dst_chain_bmc_address)

    uri = args["endpoint"]
    keystore_path = args["keystore_path"]
    keypassword = args["keypassword"]
    nid = args["nid"]

    add_verifier(plan,service_name,bmc_address,dst_chain_network,bmv_address,uri,keystore_path,keypassword,nid)
    add_btp_link(plan,service_name,bmc_address,dst_bmc_address,src_chain_network_id,uri,keystore_path,keypassword,nid)
    add_relay(plan,service_name,bmc_address,dst_bmc_address,relay_address,uri,keystore_path,keypassword,nid)

    plan.print("Icon Link Setup Completed")

# Returns BTP address
def get_btp_address(network,dapp):
    return "btp://{0}/{1}".format(network,dapp)

# Deploys dAPP on ICON
def deploy_dapp(plan,xcall_address,args):

    plan.print("Deploying dapp Contract")
    init_message = '{"_callService":"%s"}' % xcall_address

    tx_hash = contract_deployment_service.deploy_contract(plan,"dapp-sample",init_message,args)
    service_name = args["service_name"]
    
    score_address = contract_deployment_service.get_score_address(plan,service_name,tx_hash)
    return score_address   


# Deploy ibc_hndler
def deploy_ibc_handler(plan,args):

    plan.print("IBC handler")

    init_message = '{}' 

    tx_hash = contract_deployment_service.deploy_contract(plan,"ibc-0.1.0-optimized",init_message, args)
    plan.print(tx_hash)
    service_name = args["service_name"]

    score_address = contract_deployment_service.get_score_address(plan,service_name,tx_hash)

    plan.print("deployed ibc handler")

    return score_address

# deploy light_client 
def deploy_light_client_for_icon(plan,args, ibc_handler_address):

    plan.print("deploy tendermint lightclient")

    init_message = '{"ibcHandler":"%s"}' % ibc_handler_address

    tx_hash = contract_deployment_service.deploy_contract(plan, "tendermint-0.1.0-optimized", init_message, args)
    service_name = args["service_name"]
    score_address = contract_deployment_service.get_score_address(plan,service_name,tx_hash)

    plan.print("deployed light client")

    return score_address

def deploy_xcall_connection(plan,args,xcall_address,ibc_address):

    plan.print("deploy xcall connection")
    plan.print(xcall_address)
    
    init_message= '{"_xCall": "%s","_ibc": "%s","port": "xcall"}' % (xcall_address,ibc_address)

   
    tx_hash = contract_deployment_service.deploy_contract(plan, "xcall-connection-0.1.0-optimized", init_message, args)

    service_name = args["service_name"]
    score_address = contract_deployment_service.get_score_address(plan,service_name,tx_hash)

    return score_address


def deploy_xcall_for_ibc(plan,args):

    plan.print("Deploying xCall Contract for IBC")
    init_message = '{"networkId":"%s"}' % args["network"]

    tx_hash = contract_deployment_service.deploy_contract(plan,"xcall-0.1.0-optimized",init_message,args)
    service_name = args["service_name"]

    score_address = contract_deployment_service.get_score_address(plan,service_name,tx_hash)
    
    return score_address  

def deploy_xcall_dapp(plan,args,xcall_address):
    
    plan.print("Deploying Xcall Dapp Contract")

    params = '{"_callService":"%s"}' % (xcall_address)

    tx_hash = contract_deployment_service.deploy_contract(plan,"dapp-multi-protocol-0.1.0-optimized",params,args)
    service_name = args["service_name"]

    score_address = contract_deployment_service.get_score_address(plan,service_name,tx_hash)
    
    return score_address  

def add_connection_xcall_dapp(plan,xcall_dapp_address,wasm_network_id,java_xcall_connection_address,wasm_xcall_connection_address,service_name,uri,keystorepath,keypassword,nid):

    plan.print("Configure dapp connection")
    method = "addConnection"
    params = '{"nid":"%s","source":"%s","destination":"%s"}' % (wasm_network_id,java_xcall_connection_address,wasm_xcall_connection_address)

    #execute
    exec_command = ["./bin/goloop","rpc","sendtx","call","--to",xcall_dapp_address,"--method",method,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","500000000000","--nid",nid]
    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")
    tx_result = node_service.get_tx_result(plan,service_name,tx_hash,uri)
    plan.assert(value=tx_result["extract.status"],assertion="==",target_value="0x1")
    return tx_result

def configure_xcall_connection(plan,xcall_connection_address,connection_id,counterparty_port_id,counterparty_nid,client_id,service_name,uri,keystorepath,keypassword,nid):

    plan.print("Configure Xcall Connection")

    method = "configureConnection"
    params = '{"connectionId":"%s","counterpartyPortId":"%s","counterpartyNid":"%s","clientId":"%s","timeoutHeight":1000000}' % (connection_id,counterparty_port_id,counterparty_nid,client_id)
    
    exec_command = ["./bin/goloop","rpc","sendtx","call","--to",xcall_connection_address,"--method",method,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","500000000000","--nid",nid]
    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")
    tx_result = node_service.get_tx_result(plan,service_name,tx_hash,uri)
    plan.assert(value=tx_result["extract.status"],assertion="==",target_value="0x1")
    return tx_result



def set_default_connection_xcall(plan,xcall_address,wasm_network_id,xcall_connection_address,service_name,uri,keystorepath,keypassword,nid):

    plan.print("Setting Up  Xcall Default connection")
    method = "setDefaultConnection"
    params = '{"nid":"%s","connection":"%s"}' % (wasm_network_id,xcall_connection_address)

    exec_command = ["./bin/goloop","rpc","sendtx","call","--to",xcall_address,"--method",method,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","500000000000","--nid",nid]
    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")
    tx_result = node_service.get_tx_result(plan,service_name,tx_hash,uri)
    plan.assert(value=tx_result["extract.status"],assertion="==",target_value="0x1")
    return tx_result

def setup_contracts_for_ibc_java(plan,args):
    
    plan.print("Setting Contracts")

    ibc_core_address = deploy_ibc_handler(plan,args)

    xcall_address = deploy_xcall_for_ibc(plan,args)

    light_client_address = deploy_light_client_for_icon(plan,args,ibc_core_address)

    xcall_connection_address = deploy_xcall_connection(plan,args,xcall_address,ibc_core_address)

    contracts = {
        "ibc_core": ibc_core_address,
        "xcall" : xcall_address,
        "light_client" : light_client_address,
        "xcall_connection" : xcall_connection_address
    }

    return contracts

def configure_connection_for_java(plan,args,xcall_address,xcall_connection_address,wasm_network_id,connection_id,counterparty_port_id, counterparty_nid, client_id, service_name, uri, keystorepath, keypassword, nid):

    plan.print("configure conection fopr channel")

    configure_xcal_connection_result = configure_xcall_connection(plan,xcall_connection_address,connection_id,counterparty_port_id, counterparty_nid, client_id, service_name, uri, keystorepath, keypassword, nid)

    set_xcall_connection_result = set_default_connection_xcall(plan,xcall_address,wasm_network_id,xcall_connection_address,service_name,uri,keystorepath,keypassword,nid)

    return set_xcall_connection_result

def deploy_and_configure_dapp_java(plan,args,xcall_address,wasm_network_id,java_xcall_connection_address,wasm_xcall_connection_address,service_name,uri,keystorepath,keypassword,nid):

    plan.print("Deploy and Configure Dapp")

    xcall_dapp_address = deploy_xcall_dapp(plan,args,xcall_address)

    add_connection_result = add_connection_xcall_dapp(plan,xcall_dapp_address,wasm_network_id,java_xcall_connection_address,wasm_xcall_connection_address,service_name,uri,keystorepath,keypassword,nid)

    result = {
        "xcall_dapp" : xcall_dapp_address,
        "add_connection_result" : add_connection_result
    }

    return result

def registerClient(plan,service_name, args, light_client_address, keystorepath,keystore_password ,nid, uri,ibc_core_address ):

    plan.print("registering the client")

    method = "registerClient"
    params = '{"clientType":"07-tendermint","client":"%s"}' % (light_client_address)

    exec_command = ["./bin/goloop", "rpc", "sendtx", "call", "--uri", uri, "--nid", nid, "--step_limit", "5000000000", "--to", ibc_core_address, "--method", method, "--params", params, "--key_store", keystorepath, "--key_password", keystore_password ]
    plan.print(exec_command)
    result = plan.exec(service_name=service_name, recipe=ExecRecipe(command = exec_command))

    tx_hash = result["output"]
    # tx_result = get_tx_result(plan,tx_hash,service_name,)
    tx_result = node_service.get_tx_result(plan,service_name,tx_hash,uri)

    plan.assert(value=tx_result["extract.status"],assertion="==",target_value="0x1")

    return tx_hash

def bindPort(plan,service_name,args,xcall_conn_address,keystorepath,keystore_password,nid,uri,ibc_core_address,port_id):

    plan.print("Bind Port")

    password = "gochain"
    method = "bindPort"
    params = '{"portId":"%s", "moduleAddress":"%s"}' % (port_id,xcall_conn_address)

    exec_command = ["./bin/goloop", "rpc", "sendtx", "call", "--uri", uri, "--nid", nid, "--step_limit", "5000000000", "--to", ibc_core_address, "--method", method, "--params", params, "--key_store", keystorepath, "--key_password", keystore_password ]
    
    result = plan.exec(service_name=service_name, recipe=ExecRecipe(command = exec_command))

    tx_hash = result["output"]
    tx_result = node_service.get_tx_result(plan,service_name,tx_hash,uri)

    plan.assert(value=tx_result["extract.status"],assertion="==",target_value="0x1")

    return tx_hash