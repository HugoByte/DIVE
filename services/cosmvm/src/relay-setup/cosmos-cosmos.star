def relay(plan,args,SEED0,SEED1):

    plan.print("starting the cosmos-cosmos relay setup")

    plan.exec(service_name="cosmos-relay",recipe=ExecRecipe(command=["/bin/sh", "-c", "rly config init"]))

    plan.print("Adding the chain1")

    plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "rly chains add --file ../../script/chains/archway1.json my-chain"]))

    plan.print("Adding the chain2")

    plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "rly chains add --file ../../script/chains/archway2.json chain-2"]))

    plan.print("Adding the keys")

    plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "rly keys restore my-chain fd '%s' " % (SEED0["output"])]))

    plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "rly keys restore chain-2 fd1 '%s' " % (SEED1["output"])]))

    plan.print("Adding the paths")

    plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "rly paths new my-chain chain-2 demo"]))

    plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "rly transact connection demo"]))

    plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "rly tx link demo -d -t 3s "]))

    plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "rly paths list"]))

    plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "rly start &"]))