def bmc(plan, service_name, node_url, sendtx_deploy, network_id, network_address):
    exec_recipe = ExecRecipe(command = [
        "/bin/goloop", 
        "rpc", 
        "--uri", 
        node_url , 
        sendtx_deploy, 
        "--key_store ", 
        "--key_secret", 
        "--nid",
        network_id,
        "--content-type application/java",
        "--step-limit 13610920001" ,
        "--param"
        network_address
    ],)

bmc = plan.exec(
    service_name = service_name,
    recipe = exec_recipe,
)

def bts(plan, service_name, node_url, sendtx_deploy, network_id, get_score_address_bmc, javascore_serialized, network_address):
    exec_recipe = ExecRecipe(command = [
        "/bin/goloop",
        "rpc",
        "--uri",
        node_url,
        sendtx_deploy,
        "key_store",
        "key_secret",
        "--nid",
        network_id,
        "--content_type application/java",
        "--step_limit 13610920001",
        "--param"
        get_score_address_bmc,
        "--param"
        javascore_serialized,
        "--param"
        network_address,
        "--param _decimals=0x12",
        "--param _feeNumerator=0x64",
        "--param _fixedFee=0x3bacab37b62e0000"
    ],)

bts = plan.exec(
    service_name = service_name,
    recipe = exec_recipe,
)

def bsr(plan, service_name, node_url, sendtx_deploy, network_id):
    exec_recipe = ExecRecipe(command = [
        "/bin/goloop",
        "rpc",
        "--uri",
        node_url,
        sendtx_deploy,
        "--key_store",
        "key_secret",
        "--nid",
        network_id,
        "--content_type application/java",
        "step_limit 13610920001",
        "--param _name= "TICX"",
        "--param _decimals=0x12",
        "--param _initialSupply=0x186a0"
    ],)

bsr = plan.exec(
    service_name = service_name,
    recipe = exec_recipe,
)
