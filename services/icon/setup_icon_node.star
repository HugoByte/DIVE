BTP_VERSION = 21

def get_main_preps(plan,service_name,uri):
    post_request = PostHttpRequestRecipe(
        port_id="rpc",
        endpoint="/api/v3/icon_dex",
        content_type="application/json",
        body = '{ "jsonrpc": "2.0", "id": 1, "method": "icx_call", "params": { "to": "cx0000000000000000000000000000000000000000", "dataType": "call", "data": { "method": "getMainPReps", "params": {  } } } }',
        extract={
            "preps" : ".result.preps"
        }
    )
    result = plan.wait(service_name=service_name,recipe=post_request,field="code",assertion="==",target_value=200)
    
    return result["extract.preps"]

def get_prep(plan,service_name,prep_address,uri):

    post_request = PostHttpRequestRecipe(
        port_id="rpc",
        endpoint="/api/v3/icon_dex",
        content_type="application/json",
        body='{"jsonrpc": "2.0","id": 1,"method": "icx_call","params": {"to": "cx0000000000000000000000000000000000000000", "dataType": "call","data": {"method": "getPRep", "params": {"address": %s }}}}' % prep_address
    )
    result = plan.wait(service_name=service_name,recipe=post_request,field="code",assertion=">=",target_value=200)
    return result

def get_total_supply(plan,service_name):

    post_request= PostHttpRequestRecipe(
        port_id="rpc",
        endpoint="/api/v3/icon_dex",
        content_type="application/json",
        body='{ "jsonrpc": "2.0", "method": "icx_getTotalSupply", "id": 1 }',
        extract={
            "supply":".result[2:]| explode | reverse | map(if . > 96  then . - 87 else . - 48 end) | reduce .[] as $c ([1,0]; (.[0] * 16) as $b | [$b, .[1] + (.[0] * $c)])| .[1] | tonumber"
        }
    )    
    result = plan.wait(service_name=service_name,recipe=post_request,field="code",assertion="==",target_value=200)
    return result["extract.supply"]

def register_prep(plan,service_name,prep_address,uri,keystorepath,keypassword,nid):
    plan.print("registerPRep")

    name =  prep_address
    method = "registerPRep"
    value = "0x6c6b935b8bbd400000"
    params = '{"name": "%s","country": "KOR", "city": "Seoul", "email": "test@example.com", "website": "https://test.example.com", "details": "https://test.example.com/details", "p2pEndpoint": "test.example.com:7100"}' % name

    

    exec_command = ["./bin/goloop","rpc","sendtx","call","--to","cx0000000000000000000000000000000000000000","--method",method,"--value",value,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","50000000000","--nid",nid]

    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")


    tx_result = get_tx_result(plan,service_name,tx_hash,uri)

    plan.assert(value=tx_result,assertion="==",target_value="0x1")

    plan.print("Completed RegisterPrep")


def get_tx_result(plan,service_name,tx_hash,uri):

    post_request = PostHttpRequestRecipe(
        port_id="rpc",
        endpoint="/api/v3/icon_dex",
        content_type="application/json",
        body='{ "jsonrpc": "2.0", "method": "icx_getTransactionResult", "id": 1, "params": { "txHash": %s } }' % tx_hash,
        extract={
            "status":".result.status"
        }
    )
   
    result = plan.wait(service_name=service_name,recipe=post_request,field="code",assertion="==",target_value=200)

    return result["extract.status"]

def set_stake(plan,service_name,amount,uri,keystorepath,keypassword,nid):
    
    method = "setStake"
    
    params = '{"value": "%s" }' % amount

    exec_command = ["./bin/goloop","rpc","sendtx","call","--to","cx0000000000000000000000000000000000000000","--method",method,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","50000000000","--nid",nid]
    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")
    tx_result = get_tx_result(plan,service_name,tx_hash,uri)

    plan.assert(value=tx_result,assertion="==",target_value="0x1")

    plan.print("Set Stake Completed")

def set_delegation(plan,service_name,address,amount,uri,keystorepath,keypassword,nid):
    method="setDelegation"
    params='{"delegations":[{"address":%s,"value":"%s"}]}' % (address,amount)

    exec_command = ["./bin/goloop","rpc","sendtx","call","--to","cx0000000000000000000000000000000000000000","--method",method,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","50000000000","--nid",nid]
    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")
    tx_result = get_tx_result(plan,service_name,tx_hash,uri)

    plan.assert(value=tx_result,assertion="==",target_value="0x1")

def set_bonder_list(plan,service_name,address,uri,keystorepath,keypassword,nid):
    method="setBonderList"
    params='{"bonderList":[%s]}' % address

    exec_command = ["./bin/goloop","rpc","sendtx","call","--to","cx0000000000000000000000000000000000000000","--method",method,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","50000000000","--nid",nid]
    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")
    tx_result = get_tx_result(plan,service_name,tx_hash,uri)

    plan.assert(value=tx_result,assertion="==",target_value="0x1")

def set_bond(plan,service_name,address,amount,uri,keystorepath,keypassword,nid):

    method="setBond"
    params='{"bonds":[{"address":%s,"value":"%s"}]}' % (address,amount)

    exec_command = ["./bin/goloop","rpc","sendtx","call","--to","cx0000000000000000000000000000000000000000","--method",method,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","50000000000","--nid",nid]
    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")
    tx_result = get_tx_result(plan,service_name,tx_hash,uri)

    plan.assert(value=tx_result,assertion="==",target_value="0x1")

def get_revision(plan,service_name):

    post_request = PostHttpRequestRecipe(
        port_id="rpc",
        endpoint="/api/v3/icon_dex",
        content_type="application/json",
        body='{"jsonrpc": "2.0","id": 1,"method": "icx_call","params": {"to": "cx0000000000000000000000000000000000000000", "dataType": "call","data": {"method": "getRevision", "params": { }}}}',
        extract={
            # "rev_number" : '.result'
             "rev_number": '.result[2:]| explode | reverse | map(if . > 96  then . - 87 else . - 48 end) | reduce .[] as $c ([1,0]; (.[0] * 16) as $b | [$b, .[1] + (.[0] * $c)])| .[1] | tonumber '
        }
    )
    result = plan.wait(service_name=service_name,recipe=post_request,field="code",assertion="==",target_value=200)

    plan.print(result["extract.rev_number"])

    return result["extract.rev_number"]

def set_revision(plan,service_name,uri,code,keystorepath,keypassword,nid):

    method="setRevision"
    params='{"code":"%s"}' % code

    exec_command = ["./bin/goloop","rpc","sendtx","call","--to","cx0000000000000000000000000000000000000001","--method",method,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","50000000000","--nid",nid]
    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")
    tx_result = get_tx_result(plan,service_name,tx_hash,uri)

    plan.assert(value=tx_result,assertion="==",target_value="0x1")

def get_prep_node_public_key(plan,service_name,address):
    post_request = PostHttpRequestRecipe(
        port_id="rpc",
        endpoint="/api/v3/icon_dex",
        content_type="application/json",
        body='{"jsonrpc": "2.0","id": 1,"method": "icx_call","params": {"to": "cx0000000000000000000000000000000000000000", "dataType": "call","data": {"method": "getPRepNodePublicKey", "params": { "address": %s}}}}' % address,
       
        
    )
    result = plan.wait(service_name=service_name,recipe=post_request,field="code",assertion=">=",target_value=200)

    return result

def register_prep_node_publickey(plan,service_name,address,pubkey,uri,keystorepath,keypassword,nid):
    method="registerPRepNodePublicKey"
    
    params="{\"address\":\"%s\",\"pubKey\":\"%s\"}" % (address,pubkey)


    exec_command = ["./bin/goloop","rpc","sendtx","call","--to","cx0000000000000000000000000000000000000000","--method",method,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","50000000000","--nid",nid]
    plan.print(exec_command)
    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")
    tx_result = get_tx_result(plan,service_name,tx_hash,uri)

    plan.assert(value=tx_result,assertion="==",target_value="0x1")

def ensure_decentralisation(plan,service_name,prep_address,uri,keystorepath,keypassword,nid):

    plan.print("Setting Up Icon Node")

    main_preps = get_main_preps(plan,service_name,uri)

    prep = get_prep(plan,service_name,prep_address,uri)

    plan.print(prep["code"])

    setup_node(plan,service_name,uri,keystorepath,keypassword,nid,prep_address)
    

def setup_node(plan,service_name,uri,keystorepath,keypassword,nid,prep_address):
    
    revision = get_revision(plan,service_name)
    
    plan.print("ICON: revision:%s " % revision)

    # if revision < BTP_VERSION:
    #     plan.print("ICON: set revision to %s" % BTP_VERSION)

    #     set_revision(plan,service_name,uri,BTP_VERSION,keystorepath,keypassword,nid)

    pubKey = get_prep_node_public_key(plan,service_name,prep_address)

    plan.print(pubKey["body"])
    

def hex_to_int(plan,service_name,hex_number):
    exec_command = ["printf", "\"%u\"",hex_number,"|","jq tonumber"]
    result = plan.exec(service_name,recipe=ExecRecipe(command=exec_command))
    return result["output"].strip()