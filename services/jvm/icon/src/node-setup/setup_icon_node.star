wallet_config = import_module("github.com/hugobyte/chain-package/services/jvm/icon/src/node-setup/wallet.star")

BTP_VERSION = "21"

def get_main_preps(plan,service_name,uri):
    post_request = PostHttpRequestRecipe(
        port_id="rpc",
        endpoint="/api/v3/icon_dex",
        content_type="application/json",
        body = '{ "jsonrpc": "2.0", "id": 1, "method": "icx_call", "params": { "to": "cx0000000000000000000000000000000000000000", "dataType": "call", "data": { "method": "getMainPReps", "params": {  } } } }',
        extract={
            "preps" : '.result.preps'
        }
    )
    result = plan.wait(service_name=service_name,recipe=post_request,field="code",assertion="==",target_value=200)
    
    return result

def get_prep(plan,service_name,prep_address,uri):
    post_request = PostHttpRequestRecipe(
        port_id="rpc",
        endpoint="/api/v3/icon_dex",
        content_type="application/json",
        body='{"jsonrpc": "2.0","id": 1,"method": "icx_call","params": {"to": "cx0000000000000000000000000000000000000000", "dataType": "call","data": {"method": "getPRep", "params": {"address": "%s" }}}}' % prep_address,
        extract={
            "result_body" : '. | if .error != null then .error else .result end',
            "code" : '.| if .error.code != null then .error.code else 0 end | tonumber '
        }
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

def register_prep(plan,service_name,name,uri,keystorepath,keypassword,nid):


    method = "registerPRep"
    value = "0x6c6b935b8bbd400000"
    params = '{"name": "%s","country": "KOR", "city": "Seoul", "email": "test@example.com", "website": "https://test.example.com", "details": "https://test.example.com/details", "p2pEndpoint": "test.example.com:7100"}' % name

    

    exec_command = ["./bin/goloop","rpc","sendtx","call","--to","cx0000000000000000000000000000000000000000","--method",method,"--value",value,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","50000000000","--nid",nid]

    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")


    tx_result = get_tx_result(plan,service_name,tx_hash,uri)

    plan.assert(value=tx_result["extract.status"],assertion="==",target_value="0x1")

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

    return result

def set_stake(plan,service_name,amount,uri,keystorepath,keypassword,nid):
    
    method = "setStake"
    
    params = '{"value": "%s" }' % amount

    exec_command = ["./bin/goloop","rpc","sendtx","call","--to","cx0000000000000000000000000000000000000000","--method",method,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","50000000000","--nid",nid]

    plan.print(exec_command)
    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")
    tx_result = get_tx_result(plan,service_name,tx_hash,uri)

    plan.assert(value=tx_result["extract.status"],assertion="==",target_value="0x1")

    plan.print("Set Stake Completed")

def set_delegation(plan,service_name,address,amount,uri,keystorepath,keypassword,nid):
    method="setDelegation"
    params='{"delegations":[{"address":"%s","value":"%s"}]}' % (address,amount)

    exec_command = ["./bin/goloop","rpc","sendtx","call","--to","cx0000000000000000000000000000000000000000","--method",method,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","50000000000","--nid",nid]
    plan.print(exec_command)
    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")
    tx_result = get_tx_result(plan,service_name,tx_hash,uri)

    plan.assert(value=tx_result["extract.status"],assertion="==",target_value="0x1")

def set_bonder_list(plan,service_name,address,uri,keystorepath,keypassword,nid):
    method="setBonderList"
    params='{"bonderList":["%s"]}' % address

    exec_command = ["./bin/goloop","rpc","sendtx","call","--to","cx0000000000000000000000000000000000000000","--method",method,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","50000000000","--nid",nid]
    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")
    tx_result = get_tx_result(plan,service_name,tx_hash,uri)

    plan.assert(value=tx_result["extract.status"],assertion="==",target_value="0x1")

def set_bond(plan,service_name,address,amount,uri,keystorepath,keypassword,nid):

    method="setBond"
    params='{"bonds":[{"address":"%s","value":"%s"}]}' % (address,amount)

    exec_command = ["./bin/goloop","rpc","sendtx","call","--to","cx0000000000000000000000000000000000000000","--method",method,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","50000000000","--nid",nid]
    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")
    tx_result = get_tx_result(plan,service_name,tx_hash,uri)

    plan.assert(value=tx_result["extract.status"],assertion="==",target_value="0x1")

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

    plan.assert(value=tx_result["extract.status"],assertion="==",target_value="0x1")

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

    plan.assert(value=tx_result["extract.status"],assertion="==",target_value="0x1")

def ensure_decentralisation(plan,service_name,prep_address,uri,keystorepath,keypassword,nid):
    main_preps = get_main_preps(plan,service_name,uri)
    plan.print(main_preps)

    response = get_prep(plan,service_name,prep_address,uri)
    response_code = response.get("extract.code")


    plan.print("registerPRep")
    name =  "node_"+prep_address

    total_supply = get_total_supply(plan,service_name)
    min_delegated = get_min_delegated_amount(plan,service_name,total_supply)
    bond_amount = "0x152d02c7e14af6800000"
    stake = get_stake_amount(plan,service_name,bond_amount,min_delegated)

    response = register_prep(plan,service_name,name,uri,keystorepath,keypassword,nid)

    plan.print("ICON: setStake")

    set_stake(plan,service_name,stake,uri,keystorepath,keypassword,nid)

    plan.print("ICON: setDelegation")

    set_delegation(plan,service_name,prep_address,min_delegated,uri,keystorepath,keypassword,nid)

    plan.print("ICON: setBonderList")

    set_bonder_list(plan,service_name,prep_address,uri,keystorepath,keypassword,nid)

    plan.print("ICON: setBond")

    set_bond(plan,service_name,prep_address,bond_amount,uri,keystorepath,keypassword,nid)


    
    

def setup_node(plan,service_name,uri,keystorepath,keypassword,nid,prep_address):
    
    revision = get_revision(plan,service_name)
    
    plan.print("ICON: revision:%s " % revision)

    if revision != BTP_VERSION:
        plan.print("ICON: set revision to %s" % BTP_VERSION)

        set_revision(plan,service_name,uri,BTP_VERSION,keystorepath,keypassword,nid)

    pubKey = get_prep_node_public_key(plan,service_name,prep_address)

    plan.print(pubKey["body"])

    pubKey = wallet_config.get_network_wallet_public_key(plan,service_name)
    plan.print(pubKey)

    register_node_pubkey = register_prep_node_publickey(plan,service_name,prep_address,pubKey,uri,keystorepath,keypassword,nid)

    plan.print(register_node_pubkey)
    

def hex_to_int(plan,service_name,hex_number):
    exec_command = ["printf", "\"%u\"",hex_number,"|","jq tonumber"]
    result = plan.exec(service_name,recipe=ExecRecipe(command=exec_command))
    return result["output"].strip()

def get_min_delegated_amount(plan,service_name,total_supply):
    exec_command = ["python","-c","print(hex(int(%s / 500)))" % total_supply]
    result = plan.exec(service_name,recipe=ExecRecipe(exec_command))

    execute_cmd = ExecRecipe(command=["/bin/sh", "-c","echo \"%s\" | tr -d '\n\r'" % result["output"] ])
    result = plan.exec(service_name=service_name,recipe=execute_cmd)

    return result["output"].strip()

def get_stake_amount(plan,service_name,bond_amount,min_delegated):
    exec_command = ["python","-c","print(hex(int(%s) + int(%s)))" %(min_delegated,bond_amount)]
    result = plan.exec(service_name,recipe=ExecRecipe(exec_command))

    execute_cmd = ExecRecipe(command=["/bin/sh", "-c","echo \"%s\" | tr -d '\n\r'" % result["output"] ])
    result = plan.exec(service_name=service_name,recipe=execute_cmd)

    return result["output"].strip().replace("\n","")

def configure_node(plan,args):

    plan.print("Configuring ICON Node")

    icon_config_data = args["chains"]["icon"]
    

    service_name = icon_config_data["service_name"]
    uri = icon_config_data["endpoint"]
    keystorepath = icon_config_data["keystore_path"]
    keypassword = icon_config_data["keypassword"]
    nid = icon_config_data["nid"]

    prep_address =  wallet_config.get_network_wallet_address(plan,service_name)

    
    ensure_decentralisation(plan,service_name,prep_address,uri,keystorepath,keypassword,nid)

    plan.wait(service_name,recipe=ExecRecipe(command=["/bin/sh","-c","sleep 200s && echo 'success'"]),field="code",assertion="==",target_value=0,timeout="400s")

    main_preps = get_main_preps(plan,service_name,uri)
    plan.print(main_preps)

    setup_node(plan,service_name,uri,keystorepath,keypassword,nid,prep_address)

def open_btp_network(plan,service_name,args,uri,keystorepath,keypassword,nid):
    name = args["name"]
    owner = args["owner"]
    method="openBTPNetwork"
    params='{"networkTypeName":"eth","name":"%s","owner":"%s"}' % (name,owner)

    exec_command = ["./bin/goloop","rpc","sendtx","call","--to","cx0000000000000000000000000000000000000001","--method",method,"--params",params,"--uri",uri,"--key_store",keystorepath,"--key_password",keypassword,"--step_limit","50000000000","--nid",nid]
    result = plan.exec(service_name=service_name,recipe=ExecRecipe(command=exec_command))

    tx_hash = result["output"].replace('"',"")
    tx_result = filter_event(plan,service_name,tx_hash)

    plan.assert(value=tx_result["extract.status"],assertion="==",target_value="0x1")

    return tx_result

def get_last_block(plan,service_name):

    post_request = PostHttpRequestRecipe(

        port_id="rpc",
        endpoint="/api/v3/icon_dex",
        content_type="application/json",
        body='{"jsonrpc": "2.0","id": 1,"method": "icx_getLastBlock"}',
        extract={
            "height": ".result.height"
        }
    )

    response = plan.wait(service_name,recipe=post_request,field="code",assertion="==",target_value=200)

    return response["extract.height"]

def filter_event(plan,service_name,tx_hash):

    post_request = PostHttpRequestRecipe(
        port_id="rpc",
        endpoint="/api/v3/icon_dex",
        content_type="application/json",
        body='{ "jsonrpc": "2.0", "method": "icx_getTransactionResult", "id": 1, "params": { "txHash": %s } }' % tx_hash,
        extract={
            "status" : ".result.status",
            "network_type_id": '.result["eventLogs"] | .[1].indexed | .[1]',
            "network_id" :  '.result["eventLogs"] | .[1].indexed | .[2]'
        }
    )
   
    result = plan.wait(service_name=service_name,recipe=post_request,field="code",assertion="==",target_value=200)

    return result

def get_btp_network_info(plan,service_name,network_id):

    post_request = PostHttpRequestRecipe(
        port_id="rpc",
        endpoint="/api/v3/icon_dex",
        content_type="application/json",
        body='{ "jsonrpc": "2.0", "method": "btp_getNetworkInfo", "id": 1, "params": { "id": "%s" } }' % network_id,
        extract={
            "start_height" : '.result.startHeight',
        }
    )

    plan.print(post_request)

    result = plan.wait(service_name=service_name,recipe=post_request,field="code",assertion="==",target_value=200)

    exec_command = ["python","-c","print(hex(int(%s) + 1))" % result["extract.start_height"]]
    result = plan.exec(service_name,recipe=ExecRecipe(exec_command))

    execute_cmd = ExecRecipe(command=["/bin/sh", "-c","echo \"%s\" | tr -d '\n\r'" % result["output"] ])
    result = plan.exec(service_name=service_name,recipe=execute_cmd)

    return result["output"].strip()

def get_btp_header(plan,service_name,network_id,receipt_height):

    post_request = PostHttpRequestRecipe(
        port_id="rpc",
        endpoint="/api/v3/icon_dex",
        content_type="application/json",
        body='{ "jsonrpc": "2.0", "method": "btp_getHeader", "id": 1, "params": { "networkID": "%s" ,"height": "%s" } }' % (network_id,receipt_height),
        extract={
            "header" : '.result',
        }
    )

    result = plan.wait(service_name=service_name,recipe=post_request,field="code",assertion="==",target_value=200)

    command = ExecRecipe(command=["python", "-c","from base64 import b64encode, b64decode; print(b64decode('%s').hex())" % result["extract.header"]])

    first_header_hex = plan.exec(service_name,recipe=command)

    execute_cmd = ExecRecipe(command=["/bin/sh", "-c","echo \"%s\" | tr -d '\n\r'" % first_header_hex["output"] ])
    result = plan.exec(service_name=service_name,recipe=execute_cmd)

    return result["output"].strip()

def int_to_hex(plan,service_name,number):

    exec_command = ["python","-c","print(hex(int(%s)))" % number]
    result = plan.exec(service_name,recipe=ExecRecipe(exec_command))

    execute_cmd = ExecRecipe(command=["/bin/sh", "-c","echo \"%s\" | tr -d '\n\r'" % result["output"] ])
    result = plan.exec(service_name=service_name,recipe=execute_cmd)

    return result["output"]

