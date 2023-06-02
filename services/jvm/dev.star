
def dev(plan,args,gaia_conf, relayer_conf, relayer_dir, path_name):
    plan.add_service(name="service_name", config=ServiceConfig(
        image="ubuntu:latest",
    ))
    plan.exec(service_name="service_name", recipe=ExecRecipe(
        command=["echo","relay"],
    ))

if (relayer_conf in globals() or gaia_conf in globals()) and not (arg1 == "skip"):
    n = input("$0 will delete "+relayer_conf+"and" +gaia_conf+ "folder. Do you wish to continue? (y/n): ")
    if n.lower() != "y":
        exit(1)

    exec = ExecRecipe(command=["cd", relayer_dir],)
    plan.exec(service_name="service_name", recipe=exec)    

    exec = ExecRecipe(command=["rm", "-rf", relayer_conf, "&> /dev", "/null"],)
    plan.exec(service_name="service-name", recipe=exec)

    #
    exec = ExecRecipe(command=["scripts/relay_ibc", "skip"],)
    plan.exec(service_name="service_name", recipe=exec)
    sleep(3)

    exec = ExecRecipe(command=["rly","tx","link",path-name,"-d", "-t", "3s"],)
    plan.exec(service_name="service_name", recipe=exec)
    sleep(2)

    plan.print("Initial balances")
    bal1 = ExecRecipe(command=["rly","q", "bal", "cosmos0"],)
    plan.exec(service_name="service_name", recipe=bal1)
    bal2 = ExecRecipe(command=["rly", "q", "bal", "cosmos-1"],)
    plan.exec(service_name="service_name", exec=bal2)

    plan.print("balance0" +bal1+ )
    plan.print("balance1" +bal2+)

    keys = ExecRecipe(command=["rly", "keys", "show", "cosmos-1" ], )
    plan.exec(service_name="service_name", recipe=keys)
    plan.print("Sending IBC transactions")
    exec = ExecRecipe(command=[
        "rly","tx","transfer","cosmos-0", "cosmos-1", "100000samoleans", keys , "channel-0", "-d"
        ],)
    sleep(5)

    execute = ExecRecipe(command=[
        "rly", "tx", "relay-packets", path-name, "channel-0", "-d"
    ],)
    plan.exec(service_name="service_name", recipe=execute)
    sleep(5)

    execute = ExecRecipe(command=[
        "rly", "tx", "relay-acknowledgemnts", path-name, "channel-0", "-d"
    ],)
    plan.exec(service_name="service_name", recipe=execute)
    sleep(5)
    plan.print("-- Balances after packets are sent --")
    bal1 = ExecRecipe(command=["rly","q", "bal", "cosmos0"],)
    plan.exec(service_name="service_name", recipe=bal1)
    bal2 = ExecRecipe(command=["rly", "q", "bal", "cosmos-1"],)
    plan.exec(service_name="service_name", exec=bal2)

    plan.print("balance0" +bal1+ )
    plan.print("balance1" +bal2+)

    keys1 = ExecRecipe(command=["rly", "keys", "show", "cosmos-0" ], )
    exec_cmd = ExecRecipe(command=[
        "rly","tx","transfer", "cosmos-1", "cosmos-0", "100000transfer", "/channel-0", "/samoleans", keys1 , "channel-0", "-d"
    ],)
    plan.exec(service_name="service_name", exec=exec_cmd)
    sleep(5)

    execute = ExecRecipe(command=[
        "rly", "tx", "relay-packets", path-name, "channel-0", "-d"
    ],)
    plan.exec(service_name="service_name", recipe=execute)
    sleep(5)

     execute = ExecRecipe(command=[
        "rly", "tx", "relay-acknowledgemnts", path-name, "channel-0", "-d"
    ],)
    plan.exec(service_name="service_name", recipe=execute)
    sleep(5)

    echo "-- Balances after packets are sent --"
    echo "balance 0 $(rly q bal ibc-0)"
    echo "balance 1 $(rly q bal ibc-1)"
    




