DEFAULT_STEP_LIMIT = "500000000000"

"""
Deploys Contract on Icon Chain
'contract_name' - Name of the Contract to be deployed
'args' - Dict of params for deployment
"""
def deploy_contract(plan,contract_name,init_message,args):

    contract = contract_name+".jar"
    service_name = args["service_name"]
    uri = args["endpoint"]
    keystore_path = args["keystore_path"]
    keystore_password = args["keypassword"]
    nid = args["nid"]


    execute_command = ["./bin/goloop","rpc","sendtx","deploy","contracts/"+contract,"--content_type","application/java","--params",init_message,"--key_store",keystore_path,"--key_password",keystore_password,"--step_limit",DEFAULT_STEP_LIMIT,"--uri",uri,"--nid",nid]

    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=execute_command))

    return result["output"]

"""
Returns Contract Address
'tx_hash' - transaction hash
"""
def get_score_address(plan,service_name,tx_hash):

    post_request = PostHttpRequestRecipe(
        port_id="rpc",
        endpoint="/api/v3/icon_dex",
        content_type="application/json",
        body='{ "jsonrpc": "2.0", "method": "icx_getTransactionResult", "id": 1, "params": { "txHash": %s } }' % tx_hash,
        extract={
            "score_address":".result.scoreAddress"
        }
    )
   
    result = plan.wait(service_name=service_name,recipe=post_request,field="code",assertion="==",target_value=200)

    return result["extract.score_address"]

