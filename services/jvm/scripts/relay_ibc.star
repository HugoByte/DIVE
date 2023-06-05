def relay(plan,args,gaia_data, relayer_conf, relayer_dir):
    plan.add_service(name="service_name", config=ServiceConfig(
        image="ubuntu:latest",
    ))
    plan.exec(service_name="service_name", recipe=ExecRecipe(
        command=["echo","relay"],
    ))

    execute_cmd = ExecRecipe(command=["which gaiad"],) 
    result = plan.exec(service_name="service_name", recipe=execute_cmd)
    if ! result:
        plan.print("Error: Gaiad is not installed, Try running 'make build-gaia' ")
    else:
        result
    
    plan.print("Gaia version info")
    execute_cmd = ExecRecipe(command=["gaiad version --long"],)
    result = plan.exec(service_name="service_name", recipe=execute_cmd)

    execute_cmd = ExecRecipe(command=["which jq"],)
    result = plan.exec(service_name="service_name", recipe=execute_cmd)
    if ! result:
        plan.print("jq (a tool for parsing json in the command line) is required...")
        plan.print("https://stedolan.github.io/jq/download/")
    else:
        result

    execute_cmd = ExecRecipe(command=["-d",gaia_data, "&&", "! $1" , " == ", "skip"],)
    result = plan.exec(service_name="service_name", recipe=execute_cmd)
    # plan.print()
    n = input("$0 will delete "+gaia_data+"and" +relayer_conf+ "folder. Do you wish to continue? (y/n): ")
    if n.lower() != "y":
        exit(1)

    exec = ExecRecipe(command=["rm", "-rf", gaia_data, "&>" , "/dev", "/null"],)
    res = plan.exec(service_name="service_name", recipe=exec)
    exec1 = ExecRecipe(command=["rm", "-rf", relayer_conf, "/dev/null"],)
    res1 = plan.exec(service_name="service_name", recipe="exec1")

    cmd = ExecRecipe(command=["killall", "gaiad", "&>", "/dev", "/null"],)
    res = plan.exec(service_name="service_name", recipe=cmd)

    chainid0 = "cosmos0"
    chainid1 = "cosmos1"

    plan.print("generating gaia configurations")
    execute_cmd = ExecRecipe(command=[
        "mkdir", "-p", gaia_data, "&&", "cd", gaia_data, "&&", 
        "cd ../",
        "./scripts",
        "/one-chain gaiad $chainid0",
        "./data 26657 26656 6060 9090 stake",
        "./scripts",
        "/one-chain gaiad $chainid1",
        "./data 26557 26556 6061 9091 rice beans",
     ],)
    res = plan.exec(service_name="service_name", recipe=execute_cmd)

    exe = ExecRecipe(commands=["cd", relayer_dir],)
    result = plan.exec(service_name="service_name", recipe=exe)

    plan.print("Building relayer...")
    execute = ExecRecipe(command=["make", "-c", "../../", "install"],)
    result = plan.exec(service_name="service_name", recipe=execute)

    execute = ExecRecipe(command=["pwd"],)
    res = plan.exec(service_name="service_name", recipe=execute)

    plan.print("generating rly configurations")
    excute_cmd = ExecRecipe(command=["rly","config","init"],)
    result = plan.exec(service_name="service_name", recipe=excute_cmd) 

    excute_cmd = ExecRecipe(command=["rly", "chains", "add-r", "configs/chains"],)
    result = plan.exec(service_name="service_name", recipe=excute_cmd)

    SEED0 = ExecRecipe(command=["jq", "-r", ".mnemonic", gaia_data, "/ibc-0", "/key_seed.json"],)
    plan.exec(service_name="service_name", recipe=SEED0)
    SEED1 = ExecRecipe(command=["jq", "-r", ".mnemonic", gaia_data, "/ibc-1", "/key_seed.json"],)
    plan.exec(service_name="service_name", recipe=SEED1)

    KEY1 = ExecRecipe(command=["rly", "keys", "restore", "ibc-0", "testkey", SEED0],)
    plan.exec(service_name="service_name", recipe=KEY1)
    KEY2 = ExecRecipe(command=["rly", "keys", "restore", "ibc-1", "testkey", SEED1],)
    plan.exec(service_name="service_name", recipe=KEY2)

    plan.print("key" +KEY1+ "imported from ibc-0 to relayer")
    plan.print("key" +KEY2+ "imported from ibc-1 to relayer")

    execute = ExecRecipe(command=["rly", "path", "add-dir", "configs", "/paths"],)
    plan.exec(service_name="service_name", recipe=execute)

def run(plan,args):
    relay(plan, args, gaia_data, relayer_conf, relayer_dir)
