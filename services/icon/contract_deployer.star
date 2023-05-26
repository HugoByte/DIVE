DEAFULT_SERVICE_NAME="icon"
SERVICE_NAME_KEY="service"



def deploy_contract(plan,service_name,contract_path,init_message,keystore_path,keystore_password,uri):

    plan.print("Deploying A SmartContract")

    execute_command = ["./bin/goloop","rpc","sendtx","deploy","contracts/"+contract_path,"--content_type","application/java","--key_store",keystore_path,"--key_password",keystore_password,"--step_limit","500000000","--uri",uri,"--nid","0x3"]
    for i in init_message["params"]:
        execute_command.append("--param")
        execute_command.append("{0}={1}".format(i["key"],i["value"]))

    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=execute_command))

    return result["output"]

# def get_nid(plan):
#     result = plan.exec(service_name="icon",recipe=ExecRecipe(command=["cat","nid.icon"]))
#     return result["output"]


# ./bin/goloop rpc sendtx deploy /goloop/contracts/BMC-0.1.0-optimized.jar --key_store newwallet --key_password newalletpassword --step_limit 500000000 --content_type application/java --uri http://172.16.1.2:9080/api/v3/icon_dex --nid 0x3 --param _net=btp://0x1.icon