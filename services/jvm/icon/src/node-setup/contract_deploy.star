DEFAULT_STEP_LIMIT = "500000000000"


def deploy_contract(plan,contract_name,init_message, service_name, uri, keystore_path, keystore_password, nid):

    contract = contract_name+".jar"

    execute_command = ["./bin/goloop","rpc","sendtx","deploy","contracts/"+contract,"--content_type","application/java","--params",init_message,"--key_store",keystore_path,"--key_password",keystore_password,"--step_limit",DEFAULT_STEP_LIMIT,"--uri",uri,"--nid",nid]

    plan.print(execute_command)
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

