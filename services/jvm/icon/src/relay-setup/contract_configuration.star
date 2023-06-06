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



def deploy_bmv(plan,args):
    plan.print("Deploying BMV contract")

    src = args["src"]
    dst = args["dst"]
    bmc_address = args["bmc_address"]
    

    response = open_btp_network(plan,)





 
def open_btp_network(plan,service_name,src,dst,bmc_address,uri,keystorepath,keypassword,nid):
    plan.print("Opening BTP Network")

    last_block_height = node_service.get_last_block(plan,service_name)

    network_name = "{0}-{1}".format(dst,last_block_height)

    args = {"name":network_name,"owner":bmc_address}


    result = node_service.open_btp_network(plan,service_name,args,uri,keystorepath,keypassword,nid)

    return result 


