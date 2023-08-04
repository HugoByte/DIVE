cosmvm_deploy = import_module("github.com/hugobyte/dive/services/cosmvm/src/node-setup/deploy.star")
cosmvm_node = import_module("github.com/hugobyte/dive/services/cosmvm/src/node-setup/start_node.star")
contract_deployment_service = import_module("github.com/hugobyte/dive/services/jvm/icon/src/node-setup/contract_deploy.star")
icon_relay_setup = import_module("github.com/hugobyte/dive/services/jvm/icon/src/relay-setup/contract_configuration.star")

def deploy_core(plan,args):

    plan.print("Deploying ibc-core contract")

    message = '{}'

    contract_addr_ibc_core = cosmvm_deploy.deploy(plan,args, "cw_ibc_core", message)

    return contract_addr_ibc_core

def deploy_xcall(plan,args, contract_addr_ibc_core):

    plan.print("Deploying xcall contract")

    message = '{"timeout_height":10 , "ibc_host":"%s"}' % contract_addr_ibc_core 

    contract_addr_xcall = cosmvm_deploy.deploy(plan,args, "cw_xcall", message)

    return contract_addr_xcall

def deploy_light_client(plan,args):

    plan.print("Deploying the light client")

    message = '{}' 

    contract_addr_light_client = cosmvm_deploy.deploy(plan,args,"cw_icon_light_client", message)

    return contract_addr_light_client

def cosmwasm(plan, args):
    
    value = cosmvm_node.start_cosmos_node(plan,args)

    ibc_core = deploy_core(plan,args)
    plan.print(ibc_core)
    
    xcall = deploy_xcall(plan,args, ibc_core)
    plan.print(xcall)
   
    light_client = deploy_light_client(plan,args)
    plan.print(light_client)

    return value

# Deploy ibc_hndler
def ibc_handler(plan,args):

    plan.print("IBC handler")

    init_message = '{}' 

    tx_hash = contract_deployment_service.deploy_contract(plan,"ibc-0.1.0-optimized",init_message, args)

    plan.print("deployed ibc handler")

    return tx_hash

# Deploy bmc
def deploy_bmc(plan,args):
    plan.print("Deploying BMC Contract")
    init_message = '{"_net":"%s"}' % args["network"]

    bmc_hash = contract_deployment_service.deploy_contract(plan,"bmc",init_message,args)

    return bmc_hash

# deploy light_client 
def light_client_for_icon(plan,args):

    plan.print("deploy tendermint lightclient")

    init_message = '{}' 

    tx_hash = contract_deployment_service.deploy_contract(plan, "tendermint-0.1.0-optimized", init_message, args)

    plan.print("deploy light client")

    return tx_hash


