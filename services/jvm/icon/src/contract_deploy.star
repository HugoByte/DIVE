DEFAULT_KEYSTORE_PATH = "config/keystore.json"
DEFAULT_KEY_SECRET = "gochain"
DEFAULT_STEP_LIMIT = "5000000000"

def deploy_contract(plan,service_name,args):

    contract_path = args.get("contract_name")
    init_message = args.get("init_message")
    keystore_path = args.get("keystore",DEFAULT_KEYSTORE_PATH)
    keystore_password = args.get("keypassword",DEFAULT_KEY_SECRET)
    uri = args.get("uri")
    setp_limit = args.get("step_limit",DEFAULT_STEP_LIMIT)

    execute_command = ["./bin/goloop","rpc","sendtx","deploy","contracts/"+contract_path,"--content_type","application/java","--key_store",keystore_path,"--key_password",keystore_password,"--step_limit",setp_limit,"--uri",uri,"--nid","0x3"]
    for i in init_message:
        execute_command.append("--param")
        execute_command.append("{0}={1}".format(i["key"],i["value"]))

    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=execute_command))

    return result["output"]

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