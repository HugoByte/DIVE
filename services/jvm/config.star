def bmc(plan, service_name, node_uri, contract_path, keystore_path, key_secret, network_address):
    exec_recipe = ExecRecipe(command=[
        "./bin/goloop",
        "rpc",
        "--uri",
        node_uri,
        "sendtx",
        "deploy",
        "contracts/"+contract_path
        "--key_store ",
        keystore_path
        "--key_secret",
        key_secret
        "--nid",
        "0x3",
        "--content-type"
        "application/java",
        "--step-limit"
        "13610920001",
        "--param"
        network_address
    ],)

    bmc = plan.exec(
        service_name=service_name,
        recipe=exec_recipe,
    )
    return bmc["output"]


def bts(plan, service_name, node_uri, contract_path, keystore_path, key_secret, get_score_address_bmc, javascore_serialized, network_address):
    exec_recipe = ExecRecipe(command=[
        "./bin/goloop",
        "rpc",
        "--uri",
        node_uri,
        "sendtx",
        "deploy",
        "contracts/"+contract_path
        "key_store",
        keystore_path
        "key_secret",
        key_secret
        "--nid",
        "0x3",
        "--content_type"
        "application/java",
        "--step_limit"
        "13610920001",
        "--param"
        get_score_address_bmc,
        "--param"
        javascore_serialized,
        "--param"
        network_address,
        "--param"
        "_decimals=0x12",
        "--param"
        "_feeNumerator=0x64",
        "--param"
        "_fixedFee=0x3bacab37b62e0000"
    ],)

    bts = plan.exec(
        service_name=service_name,
        recipe=exec_recipe,
    )
    return bts["result"]


def bsr(plan, service_name, node_url, contract_path, keystore_path, key_secret):
    exec_recipe = ExecRecipe(command=[
        "./bin/goloop",
        "rpc",
        "--uri",
        node_url,
        "sendtx"
        "deploy",
        "contracts/"+contract_path
        "--key_store",
        keystore_path
        "key_secret",
        key_secret
        "--nid",
        "0x3",
        "--content_type"
        "application/java",
        "step_limit"
        "13610920001",
        "--param"
        "_name= "TICX"",
        "--param"
        "_decimals=0x12",
        "--param "
        "_initialSupply=0x186a0"
    ],)

    bsr = plan.exec(
        service_name=service_name,
        recipe=exec_recipe,
    )
    return bsr["output"]
