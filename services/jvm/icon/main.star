wallet = import_module("github.com/hugobyte/chain-package/services/jvm/icon/src/node-setup/wallet.star")
setup_node = import_module("github.com/hugobyte/chain-package/services/jvm/icon/src/node-setup/setup_icon_node.star")
icon_node_launcher = import_module("github.com/hugobyte/chain-package/services/jvm/icon/src/node-setup/start_icon_node.star")
contract_service = import_module("github.com/hugobyte/chain-package/services/jvm/icon/src/node-setup/contract_deploy.star")

# relay_setup = import_module("github.com/hugobyte/chain-package/services/jvm/icon/src/relay-setup/contract_configuration.star")

def node_service(plan,args):

    node_service = icon_node_launcher.start_icon_node(plan,args)

    plan.print("Icon Node Started")

    return node_service
   
        
def contract_deployer(plan,args):

    plan.print("Running Contract Deployer")

    service_params = args["service_params"]

    service_name = service_params.get("service_name","ICON")

    tx_hash = contract_service.deploy_contract(plan,service_name,service_params)

    score_address = contract_service.get_score_address(plan,service_name,tx_hash)

    return score_address

def icon_service(plan,args):

    if args["service_method"] == "deploy_node":

      response = node_service(plan,args)

      plan.print("Private_URL " + response.private_url)

      plan.print("Public_URL  "+response.public_url)

      

    elif args["service_method"] == "deploy_contract":

        response = contract_deployer(plan,args)

        plan.print("ScoreAddress %s" % response)

    # else: 
    #     relay_setup.open_btp_network(plan,"ICON","src","dst","cx0000000000000000000000000000000000000000","http://172.16.1.2:9080/api/v3/icon_dex","config/keystore.json","gochain","0x3")


