wallet = import_module("github.com/hugobyte/chain-package/services/icon/node-service/src/wallet.star")
setup_node = import_module("github.com/hugobyte/chain-package/services/icon/node-service/src/setup_icon_node.star")
icon_node_launcher = import_module("github.com/hugobyte/chain-package/services/icon/node-service/src/start_icon_node.star")


def main(plan,args):

    plan.print("Starting Kurtosis Package")
    node_service = icon_node_launcher.start_icon_node(plan,args)

    wallet_data = wallet.create_wallet(plan,node_service.node_service_response.name,"newwallet","newalletpassword")

    plan.print(wallet_data)

    uri = node_service.private_url
    plan.print(uri)
        

    adddress =  wallet.get_network_wallet_address(plan,node_service.node_service_response.name)
    plan.print(adddress)

    node = "node_hxb6b5791be0b5ef67063b3c10b840fb81514db2fd"
        
    plan.print(node)
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
   
        
