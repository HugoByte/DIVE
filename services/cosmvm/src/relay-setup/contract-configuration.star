cosmvm_deploy = import_module("github.com/hugobyte/dive/services/cosmvm/src/node-setup/deploy.star")
cosmvm_node = import_module("github.com/hugobyte/dive/services/cosmvm/src/node-setup/start_node.star")

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


