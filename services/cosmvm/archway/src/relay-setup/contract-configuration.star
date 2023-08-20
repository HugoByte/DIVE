cosmvm_deploy = import_module("github.com/hugobyte/dive/services/cosmvm/archway/src/node-setup/deploy.star")
PASSCODE="password"

def deploy_core(plan,args):

    plan.print("Deploying ibc-core contract")

    message = '{}'

    contract_addr_ibc_core = cosmvm_deploy.deploy(plan,args, "cw_ibc_core", message)

    return contract_addr_ibc_core

def deploy_xcall(plan,args,network_id,denom):

    plan.print("Deploying xcall contract")

    message = '{"network_id":"%s" , "denom":"%s"}' % (network_id,denom)

    contract_addr_xcall = cosmvm_deploy.deploy(plan,args, "cw_xcall", message)

    return contract_addr_xcall

def deploy_light_client(plan,args,ibc_host_address):

    plan.print("Deploying the light client")

    message = '{"ibc_host":"%s"}' % (ibc_host_address)

    contract_addr_light_client = cosmvm_deploy.deploy(plan,args,"cw_icon_light_client", message)
    

    return contract_addr_light_client

def deploy_xcall_connection(plan,args,xcall_address,ibc_host,port_id,denom):

    plan.print("Deploying the xcall ibc connection")

    message = '{"ibc_host":"%s","port_id":"%s","xcall_address":"%s", "denom":"%s"}' % (ibc_host,port_id,xcall_address,denom)

    contract_addr_xcall_connection = cosmvm_deploy.deploy(plan,args,"cw_xcall_ibc_connection", message)
    
    return contract_addr_xcall_connection


def bindPort(plan,args,ibc_address,conn_address):

    plan.print("bind mock app to the port")

    exec = ExecRecipe(command=["/bin/sh", "-c", "echo '%s' | archwayd tx wasm execute %s '{\"bind_port\":{\"address\":\"%s\", \"port_id\":\"xcall\"}}' --from fd --chain-id my-chain --output json -y" % (PASSCODE,ibc_address, conn_address )])
    plan.print(exec)
    result = plan.exec(service_name="cosmos", recipe=exec)
   
    tx_hash = result["output"] 

    return tx_hash

def registerClient(plan,args,ibc_address,client_address):

    plan.print("registering the client")

    exec = ExecRecipe(command=["/bin/sh", "-c", "echo '%s' | archwayd tx wasm execute \"%s\" '{\"register_client\":{\"client_type\":\"iconclient\",\"client_address\":\"%s\"}}' --from fd --chain-id my-chain --output json -y" % (PASSCODE,ibc_address,client_address)])
    result = plan.exec(service_name="cosmos", recipe=exec)

    tx_hash = result["output"]

    return tx_hash

def deploy_xcall_dapp(plan,xcall_address):

    plan.print("Deploying the xcall dapp")

    message = '{"address":"%s"}' % (xcall_address)

    xcall_dapp_address = cosmvm_deploy.deploy(plan,args,"cw_xcall_ibc_connection", message)
    
    return xcall_dapp_address

def add_connection_xcall_dapp(plan,xcall_dapp_address,wasm_xcall_connection_address,xcall_connection_address,network_id):

    plan.print("Configure xcall dapp")

    params = '{"add_connection":{"src_endpoint":"%s","dest_endpoint":"%s","network_id":"%s"}}' % (wasm_xcall_connection_address,xcall_connection_address,network_id)

    exec = ExecRecipe(command=["/bin/sh", "-c", "echo '%s' | archwayd tx wasm execute \"%s\" %s --from fd --chain-id my-chain --output json -y" % (PASSCODE,xcall_dapp_address,params)])
    result = plan.exec(service_name="cosmos", recipe=exec)

    tx_hash = result["output"]

    return tx_hash

def configure_xcall_connection(plan,args,connection_id,counterparty_port_id,counterparty_nid,client_id):

    plan.print("Configure Xcall Connections Connection ")

    params = '{"configure_connection":{"connection_id":"%s","counterparty_port_id":"%s","counterparty_nid":"%s","client_id":"%s","timeout_height":30000}}' % (connection_id,counterparty_port_id,counterparty_nid,client_id)


def set_default_connection_xcall(paln,network_id,xcall_connection_address,xcall_address):
    plan.print("Set Xcall default connection ")  
    params = '{"set_default_connection":{"nid":"%s","address":"%s"}}' % (network_id,xcall_connection_address)
