RELAY_BIN = "../bin/relay"

def relay(plan,args,chain_config_path, deployment_path):
    if RELAY_BIN != "null":
        hello=plan.add_service(name="hello", config=ServiceConfig(
            image="ubuntu:latest",
        ))
        plan.exec(service_name="hello",recipe=ExecRecipe(
            command=["echo", "hello"],
        ))
    else:
        RELAY_BIN

    SRC_NETWORK = plan.exec(service_name="hello", recipe=ExecRecipe(command=["jq", "."+SRC+".network", deployment_path],))
    SRC_BMC_ADDRESS = plan.exec(service_name="hello", recipe=ExecRecipe(command=["jq", "."+SRC+".contracts.bmc"],))
    SRC_ADDRESS = "btp://"+SRC_NETWORK+"/"+SRC_BMC_ADDRESS
    SRC = plan.exec(service_name="hello", recipe=ExecRecipe(command=["jq", ".link.src",chain_config_path],))
    SRC_ENDPOINT = plan.exec(service_name="hello", recipe=ExecRecipe(command=["jq", ".chains."+SRC+".endpoint"],))
    SRC_KEY_STORE = plan.exec(service_name="hello", recipe=ExecRecipe(command=["jq", ".chains."+SRC+".keystore" ],))
    SRC_KEY_SECRET = plan.exec(service_name="hello", recipe=ExecRecipe(command=["jq", ".chains."+SRC+".keysecret" ],))
   
    if SRC_KEY_SECRET != "null":
        SRC_KEY_PASSWORD = plan.print(SRC_KEY_SECRET)
    else:
        SRC_KEY_PASSWORD = plan.exec(service_name="hello", recipe=ExecRecipe(command=["jq", "chain."+SRC+".keypass"],))

    DST_NETWORK = plan.exec(service_name="hello", recipe=ExecRecipe(command=["jq", "."+DST+".network", deployment_path],))
    DST_BMC_ADDRESS = plan.exec(service_name="hello", recipe=ExecRecipe(command=["jq", "."+DST+".contracts.bmc"],))
    DST_ADDRESS = "btp://"+DST_NETWORK+"/"+DST_BMC_ADDRESS
    DST = plan.exec(service_name="hello", recipe=ExecRecipe(command=["jq", ".link.dst",chain_config_path],))
    DST_ENDPOINT = plan.exec(service_name="hello", recipe=ExecRecipe(command=["jq", ".chains."+DST+".endpoint"],))
    DST_KEY_STORE = plan.exec(service_name="hello", recipe=ExecRecipe(command=["jq", ".chains."+DST+".keystore" ],))
    DST_KEY_SECRET = plan.exec(service_name="hello", recipe=ExecRecipe(command=["jq", ".chains."+DST+".keysecret" ],))

    if DST_KEY_SECRET != "null":
        DST_KEY_PASSWORD = plan.print(DST_KEY_SECRET)
    else:
        DST_KEY_PASSWORD = plan.exec(service_name="hello", recipe=ExecRecipe(command=["jq", "chain."+DST+".keypass"],))
   
    if  BMV_BRIDGE == True:
        plan.print("using bridge mode")
    else:
        plan.print("Using BTP Block node")
        BMV_BRIDGE = False

    exec_command = ExecRecipe(command=[
        RELAY_BIN,
        "--direction",
        "both",
        "--src.address",
        SRC_ADDRESS,
        "--src.endpoint",
        SRC_ENDPOINT,
        "--src.key_store",
        SRC_KEY_STORE,
        "--src.key_password",
        SRC_KEY_PASSWORD,
        "--src.bridge_mode",
        BMV_BRIDGE,
        "--dst.address",
        DST_ADDRESS,
        "--dst.endpoint",
        DST_ENDPOINT,
        "--dst.key_store"
        DST_KEY_STORE,
        "dst.key_password",
        DST_KEY_PASSWORD,
        "start"
    ],)
    result = plan.exec(service_name="hello", recipe=exec_command)

    return result["output"]


    