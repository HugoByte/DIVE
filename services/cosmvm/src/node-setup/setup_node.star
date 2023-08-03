ICON_NODE = "http://localhost:9082/api/v3/"

def openBTPNetwork(plan,args, wallet, nid):

    plan.print("Opening Btp network of type eth")

    password = "gochain"
    method = "openBTPNetwork"
    name = args["name"]
    owner = args["owner"]
    params = '{"network_typeName":"eth", "name": "%s", "owner":"%s"}' % (name,owner)

    exec_command = ["./bin/goloop", "rpc", "sendtx", "call", "--uri", ICON_NODE, "--nid", nid, "--step-limit", "1000000000", "--to", "cx0000000000000000000000000000000000000001", "method", method, "--param", params, "--key_store", wallet, "key_password", password ]
    result = plan.exec(service_name="cosmos", recipe=ExecRecipe(command = exec_command))

    tx_hash = result["output"]

    tx_result = get_tx_result(plan, tx_hash )

    plan.assert(value=tx_result["extract.status"],assertion="==",target_value="200")

    plan.print("Completed")

def registerClient(plan,args, client_address, to_address, wallet, nid ):

    plan.print("registering the client")

    password = "gochain"
    method = "registerClient"
    params = '{"clientType":"07-tendermint","client":"%s"}' % (client_address)

    exec_command = ["./bin/goloop", "rpc", "sendtx", "call", "--uri", ICON_NODE, "--nid", nid, "--step_limit", "1000000000", "--to", to_address, "--method", method, "--param", params, "--key_store", wallet, "--key_password", password ]
    result = plan.exec(service_name= "cosmos", recipe=ExecRecipe(command = exec_command))

    tx_hash = result["output"]

    tx_result = get_tx_result(plan,tx_hash)

    plan.assert(value=tx_result["extract.status"], assertion="==", target_value="200")

    plan.print("Registered")

def bindPort(plan,args, mock_app_address, port_id, to_address, wallet, nid ):

    plan.print("Bind mock app to a port")

    password = "gochain"
    method = "bindPort"
    params = '{"moduleAddress":"%s","port_id":"%s"}' % (mock_app_address, port_id )

    exec_command = ["./bin/goloop", "rpc", "sendtx", "call", "--uri", ICON_NODE, "--nid", nid, "--step_limit", "1000000000", "--to", to_address, "--method", method, "--param", params, "--key_store", wallet, "--key_password", password ]
    result = plan.exec(service_name="cosmos", recipe=ExecRecipe(command = exec_command))

    tx_hash = result["output"]

    tx_result = get_tx_result(plan,tx_hash)

    plan.assert(value=tx_result["extract.status"], assertion="==", target_value="200")

    plan.print("Completed binding")

def get_tx_result(plan, tx_hash):

    post_request = PostHttpRequestRecipe(
        port_id="rpc",
        content_type="application/java",
        body='{ "jsonrpc": "2.0", "method": "getTransactionResult", "id": 1, "params": { "txHash": %s } }' % tx_hash,
        extract={
            "status":".result.status"
        }
    )
   
    result = plan.wait(service_name="cosmos",recipe=post_request,field="code",assertion="==",target_value=200)

    return result

def get_last_block(plan):

    post_request = PostHttpRequestRecipe(
        
        port_id="rpc",
        endpoint=ICON_NODE,
        content_type="application/java",
        body='{"jsonrpc": "2.0","id": 1,"method": "getLastBlock"}',
        extract={
            "height": ".result.height"
        }
    )

    response = plan.wait(service_name="cosmos",recipe=post_request,field="code",assertion="==",target_value=200)

    return response["extract.height"]

def get_btp_network_info(plan,service_name,network_id):

    post_request = PostHttpRequestRecipe(
        port_id="rpc",
        endpoint=ICON_NODE,
        content_type="application/java",
        body='{ "jsonrpc": "2.0", "method": "btp_getNetworkInfo", "id": 1, "params": { "id": "%s" } }' % network_id,
        extract={
            "start_height" : '.result.startHeight',
        }
    )
    result = plan.wait(service_name="cosmos",recipe=post_request,field="code",assertion="==",target_value=200)

    exec_command = ["python","-c","print(hex(int(%s) + 1))" % result["extract.start_height"]]
    result = plan.exec(service_name="cosmos",recipe=ExecRecipe(exec_command))

    execute_cmd = ExecRecipe(command=["/bin/sh", "-c","echo \"%s\" | tr -d '\n\r'" % result["output"] ])
    result = plan.exec(service_name="cosmos",recipe=execute_cmd)

    return result["output"]

def get_btp_header(plan,network_id,receipt_height):

    post_request = PostHttpRequestRecipe(
        port_id="rpc",
        endpoint=ICON_NODE,
        content_type="application/java",
        body='{ "jsonrpc": "2.0", "method": "btp_getHeader", "id": 1, "params": { "networkID": "%s" ,"height": "%s" } }' % (network_id,receipt_height),
        extract={
            "header" : '.result',
        }
    )

    result = plan.wait(service_name="cosmos",recipe=post_request,field="code",assertion="==",target_value=200)

    command = ExecRecipe(command=["python", "-c","from base64 import b64encode, b64decode; print(b64decode('%s').hex())" % result["extract.header"]])

    first_header_hex = plan.exec(service_name="cosmos",recipe=command)

    execute_cmd = ExecRecipe(command=["/bin/sh", "-c","echo \"%s\" | tr -d '\n\r'" % first_header_hex["output"] ])
    result = plan.exec(service_name="cosmos",recipe=execute_cmd)

    return result["output"]

def int_to_hex(plan,number):

    exec_command = ["python","-c","print(hex(int(%s)))" % number]
    result = plan.exec(service_name="cosmos",recipe=ExecRecipe(exec_command))

    execute_cmd = ExecRecipe(command=["/bin/sh", "-c","echo \"%s\" | tr -d '\n\r'" % result["output"] ])
    result = plan.exec(service_name="cosmos",recipe=execute_cmd)

    return result["output"]



    

    
