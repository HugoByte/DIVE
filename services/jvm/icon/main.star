wallet = import_module("github.com/hugobyte/chain-package/services/jvm/icon/src/wallet.star")
setup_node = import_module("github.com/hugobyte/chain-package/services/jvm/icon/src/setup_icon_node.star")
icon_node_launcher = import_module("github.com/hugobyte/chain-package/services/jvm/icon/src/start_icon_node.star")
contract_service = import_module("github.com/hugobyte/chain-package/services/jvm/icon/src/contract_deploy.star")

def node_service(plan,args):

    plan.print("Starting Kurtosis Package")
    node_service = icon_node_launcher.start_icon_node(plan,args)

    uri = node_service.private_url
        
    adddress =  wallet.get_network_wallet_address(plan,node_service.node_service_response.name)

    data = {
             "service_name": node_service.node_service_response.name,
             "prep_address":adddress,
             "uri":uri,
             "keystorepath": "config/keystore.json",
             "keypassword":"gochain",
             "nid":"0x3"  
    }        
    r = setup_node.configure_node(plan,data)
    plan.print(r)

    plan.print("COMPLETED")

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

