# icon_node_launcher = import_module("github.com/hugobyte/chain-package/services/icon/node_launcher.star")
# wallet_creator = import_module("github.com/hugobyte/chain-package/services/icon/contract_deployer.star")
wallet = import_module("github.com/hugobyte/chain-package/services/icon/wallet.star")
setup_node = import_module("github.com/hugobyte/chain-package/services/icon/setup_icon_node.star")

icon_node_launcher = import_module("github.com/hugobyte/chain-package/services/icon/start_icon_node.star")

def run(plan,args):

    plan.print("Starting Kurtosis Package")

    if args["chain"] == "ICON":
        node_service = icon_node_launcher.start_icon_node(plan,args)

        wallet_data = wallet.create_wallet(plan,node_service.name,"newwallet","newalletpassword")

        plan.print(wallet_data)

        uri = "http://172.16.1.2:9080/api/v3/icon_dex"
        plan.print(uri)
        

        adddress =  wallet.get_network_wallet_address(plan,node_service.name)
        plan.print(adddress)

        node = "node_hxb6b5791be0b5ef67063b3c10b840fb81514db2fd"
        plan.print(node)        
        r = setup_node.ensure_decentralisation(plan,node_service.name,adddress,uri,"config/keystore.json","gochain","0x3")
        plan.print(r)

        # res = setup_node.get_total_supply(plan,node_service.name)


        # rs = setup_node.register_prep(plan,node_service.name,node,uri,"config/keystore.json","gochain","0x3")


        # ss = setup_node.set_stake(plan,node_service.name,"0xde0b6b3a7640000",uri,"config/keystore.json","gochain","0x3")

        # plan.print(ss)

        # sd = setup_node.set_delegation(plan,node_service.name,str(adddress),"0x2710",uri,"config/keystore.json","gochain","0x3")

        # sbl = setup_node.set_bonder_list(plan,node_service.name,str(adddress),uri,"config/keystore.json","gochain","0x3")

        # sbond = setup_node.set_bond(plan,node_service.name,str(adddress),"0x2710",uri,"config/keystore.json","gochain","0x3")

        
      
        # node_data = setup_node.get_main_preps(plan,node_service.name,uri)

        # plan.print(node_data)
       
        # d = setup_node.get_prep(plan,node_service.name,adddress,uri)

        # plan.print(d["code"])

        # rev = setup_node.get_revision(plan,node_service.name)

        # plan.print(rev)

        # if rev != "0x15":
        #     res = setup_node.set_revision(plan,node_service.name,uri,"0x15","config/keystore.json","gochain","0x3")
        #     plan.print(res)

        # res = setup_node.get_prep_node_public_key(plan,node_service.name,adddress)
        # # keystore_path = wallet_creator.create_wallet(plan,node_service.name,"newwallet","newalletpassword")

        # # contract_tx = wallet_creator.deploy_contract(plan,node_service.name,"BMC-0.1.0-optimized.jar",args,keystore_path,"newalletpassword","http://"+node_service.ip_address+":9080"+"/api/v3/icon_dex")

        # # plan.print(contract_tx)

        # pubkey = wallet.get_network_wallet_public_key(plan,node_service.name)

        # plan.print(pubkey)

        # result = setup_node.register_prep_node_publickey(plan,node_service.name,"hxb6b5791be0b5ef67063b3c10b840fb81514db2fd","0x03b3d972e61b4e8bf796c00e84030d22414a94d1830be528586e921584daadf934",uri,"config/keystore.json","gochain","0x3")

        # plan.print(result)

        plan.print("COMPLETED")

    else:
        plan.print("Not Configured")



# def run(plan, args):
   
#     plan.print("Starting Deployment Tool")

#     if args["chain"] == "ICON":
#         ip = icon_node_launcher.launch_icon_node(plan,args)
#         plan.print(ip)
#         response = plan.exec(service_name="icon",recipe=ExecRecipe(command=["../bin/goloop","rpc","lastblock","--uri","http://"+ip+"/api/v3"]),)
#         plan.print(response)

#         keystore_path = wallet_creator.create_wallet(plan,"icon","newwallet","newalletpassword")
# # deploy_contract(plan,service_name,contract_path,init_message,keystore_path,keystore_password,uri):
       
#         contract_tx = wallet_creator.deploy_contract(plan,"icon","BMC-0.1.0-optimized.jar",args,keystore_path,"newalletpassword","http://"+ip+"/api/v3")
#         plan.print(contract_tx)

#     else:
#         plan.print("Not Configured")


