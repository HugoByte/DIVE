cosmvm_deploy = import_module("github.com/hugobyte/dive/services/cosmvm/archway/src/node-setup/deploy.star")
PASSCODE="password"

node_constants = import_module("github.com/hugobyte/dive/package_io/constants.star")
password = node_constants.ARCHWAY_SERVICE_CONFIG.password

def deploy_core(plan,service_name,chain_id,chain_key):

    plan.print("Deploying ibc-core contract")

    message = '{}'

    
    contract_addr_ibc_core = cosmvm_deploy.deploy(plan,chain_id,chain_key, "cw_ibc_core", message,service_name,password)

    return contract_addr_ibc_core

def deploy_xcall(plan,service_name,chain_id,chain_key,network_id,denom):

    plan.print("Deploying xcall contract")

    message = '{"network_id":"%s" , "denom":"%s"}' % (network_id,denom)

    contract_addr_xcall = cosmvm_deploy.deploy(plan,chain_id,chain_key, "cw_xcall", message,service_name,password)

    return contract_addr_xcall

def deploy_light_client(plan,service_name,chain_id,chain_key,ibc_host_address):

    plan.print("Deploying the light client")

    message = '{"ibc_host":"%s"}' % (ibc_host_address)

    contract_addr_light_client = cosmvm_deploy.deploy(plan,chain_id,chain_key,"cw_icon_light_client", message,service_name,password)
    

    return contract_addr_light_client

def deploy_xcall_connection(plan,service_name,chain_id,chain_key,xcall_address,ibc_host,port_id,denom):

    plan.print("Deploying the xcall ibc connection")

    message = '{"ibc_host":"%s","port_id":"%s","xcall_address":"%s", "denom":"%s"}' % (ibc_host,port_id,xcall_address,denom)

    contract_addr_xcall_connection = cosmvm_deploy.deploy(plan,chain_id,chain_key,"cw_xcall_ibc_connection", message,service_name,password)
    
    return contract_addr_xcall_connection


def bindPort(plan,service_name,chain_id,chain_key,ibc_address,conn_address):

    plan.print("bind mock app to the port")

    exec = ExecRecipe(command=["/bin/sh", "-c", "echo '%s' | archwayd tx wasm execute %s '{\"bind_port\":{\"address\":\"%s\", \"port_id\":\"xcall\"}}' --from %s --keyring-backend test --chain-id %s --output json -y" % (PASSCODE,ibc_address, conn_address,chain_key,chain_id)])
    plan.print(exec)
    result = plan.exec(service_name=service_name, recipe=exec)
   
    tx_hash = result["output"] 

    return tx_hash

def registerClient(plan,service_name,chain_id,chain_key,ibc_address,client_address):

    plan.print("registering the client")

    exec = ExecRecipe(command=["/bin/sh", "-c", "echo '%s' | archwayd tx wasm execute \"%s\" '{\"register_client\":{\"client_type\":\"iconclient\",\"client_address\":\"%s\"}}' --from %s --keyring-backend test --chain-id %s --output json -y" % (PASSCODE,ibc_address,client_address,chain_key,chain_id)])
    result = plan.exec(service_name=service_name, recipe=exec)

    tx_hash = result["output"]

    return tx_hash

def deploy_xcall_dapp(plan,service_name,chain_id,chain_key,xcall_address):

    plan.print("Deploying the xcall dapp")

    message = '{"address":"%s"}' % (xcall_address)

    xcall_dapp_address = cosmvm_deploy.deploy(plan,chain_id,chain_key,"cw_xcall_ibc_connection", message,service_name)
    
    return xcall_dapp_address

def add_connection_xcall_dapp(plan,service_name,chain_id,chain_key,xcall_dapp_address,wasm_xcall_connection_address,xcall_connection_address,java_network_id):

    plan.print("Configure xcall dapp")

    params = '{"add_connection":{"src_endpoint":"%s","dest_endpoint":"%s","network_id":"%s"}}' % (wasm_xcall_connection_address,xcall_connection_address,java_network_id)

    exec = ExecRecipe(command=["/bin/sh", "-c", "echo '%s' | archwayd tx wasm execute \"%s\" %s --from %s --keyring-backend test --chain-id %s --output json -y" % (PASSCODE,xcall_dapp_address,params,chain_key,chain_id)])
    result = plan.exec(service_name=service_name, recipe=exec)

    tx_hash = result["output"]

    return tx_hash

def configure_xcall_connection(plan,service_name,chain_id,chain_key,xcall_connection_address,connection_id,counterparty_port_id,counterparty_nid,client_id):

    plan.print("Configure Xcall Connections Connection ")

    params = '{"configure_connection":{"connection_id":"%s","counterparty_port_id":"%s","counterparty_nid":"%s","client_id":"%s","timeout_height":30000}}' % (connection_id,counterparty_port_id,counterparty_nid,client_id)

    exec_cmd = ["/bin/sh", "-c","echo '%s'| archwayd tx wasm execute %s %s --from %s --keyring-backend test --chain-id %s --output json -y" % (PASSCODE,xcall_connection_address,params,chain_key,chain_id)]

    result = plan.exec(service_name=service_name, recipe=ExecRecipe(command=exec_cmd))

    tx_result = check_tx_result(params,result["output"],service_name)

    tx_hash = result["output"]


def set_default_connection_xcall(plan,service_name,chain_id,chain_key,network_id,xcall_connection_address,xcall_address):
    plan.print("Set Xcall default connection ")  
    params = '{"set_default_connection":{"nid":"%s","address":"%s"}}' % (network_id,xcall_connection_address)

    exec_cmd = ["/bin/sh", "-c","echo '%s'| archwayd tx wasm execute %s %s --from %s --keyring-backend test --chain-id %s --output json -y" % (PASSCODE,xcall_address,params,chain_key,chain_id)]

    result = plan.exec(service_name=service_name, recipe=ExecRecipe(command=exec_cmd))

    tx_result = check_tx_result(params,result["output"],service_name)

    tx_hash = result["output"]


def check_tx_result(plan,tx_hash,service_name):

    plan.print("Check Tx Result")

    # exec_cmd = ["/bin/sh","-c","archwayd query tx %s  --chain-id %s | jq .code" % (tx_hash,chain_id)]

    post_request = PostHttpRequestRecipe(
        port_id="rpc",
        endpoint="",
        content_type="application/json",
        body='{ "jsonrpc": "2.0", "method": "tx", "id": 1, "params": { "hash": %s } }' % tx_hash,
        extract={
            "status" : ".result.code",
        }
    )
   
    result = plan.wait(service_name=service_name,recipe=post_request,field="extract.status",assertion="==",target_value=0)

    return result

def setup_contracts_for_ibc_wasm(plan,service_name,chain_id,chain_key,network_id,denom,port_id):
    plan.print("Deploying Contracts for IBC Setup")

    ibc_core_address = deploy_core(plan,service_name,chain_id,chain_key)

    light_client_address = deploy_light_client(plan,service_name,chain_id,chain_key,ibc_core_address)

    xcall_address = deploy_xcall(plan,service_name,chain_id,chain_key,network_id,denom)

    xcall_connection_address = deploy_xcall_connection(plan,service_name,chain_id,chain_key,xcall_address,ibc_core_address,port_id,denom)

    contracts = {
        "ibc_core" : ibc_core_address,
        "xcall" : xcall_address,
        "light_client": light_client_address,
        "xcall_connection" : xcall_connection_address
    }

    return contracts

def configure_connection_for_wasm(plan,service_name,chain_id,chain_key,xcall_connection_address,connection_id,counterparty_port_id, counterparty_nid, client_id,network_id,xcall_address):

    plan.print("Configure Connection for Channel Steup IBC")

    configure_xcall_connection_result  = configure_xcall_connection(plan,service_name,chain_id,chain_key,xcall_connection_address,connection_id,counterparty_port_id,counterparty_nid,client_id)

    plan.print(configure_xcall_connection_result)

    configure_xcall_result = set_default_connection_xcall(plan,service_name,chain_id,chain_key,network_id,xcall_connection_address,xcall_address)

    plan.print(configure_xcall_result)

def deploy_and_configure_xcall_dapp(plan,service_name,chain_id,chain_key,xcall_address,wasm_xcall_connection_address,xcall_connection_address,network_id):

    plan.print("Configure Xcall Dapp")


    xcall_dapp_address = deploy_xcall_dapp(plan,service_name,chain_id,chain_key,xcall_address)

    add_connection_result = add_connection_xcall_dapp(plan,service_name,chain_id,chain_key,xcall_dapp_address,wasm_xcall_connection_address,xcall_connection_address,network_id)

    result = {
        "xcall_dapp" : xcall_dapp_address,
        "add_connection_result" : add_connection_result
    }

    return result