constants = import_module("github.com/hugobyte/dive/package_io/constants.star")

def relay(plan,args,seed0,seed1):

    plan.print("starting the cosmos-cosmos relay setup")

    plan.exec(service_name="cosmos-relay",recipe=ExecRecipe(command=["/bin/sh", "-c", "rly config init"]))

    plan.print("Adding the chain1")

    plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "rly chains add --file ../../script/chains/archway1.json '%s' " % (constants.COSMOS_NODE_CLIENT.chain_id)]))

    plan.print("Adding the chain2")

    plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "rly chains add --file ../../script/chains/archway2.json '%s' " % (constants.COSMOS_NODE_CLIENT.chain_id_1)]))

    plan.print("Adding the keys")

    plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "rly keys restore '%s' fd '%s' " % (constants.COSMOS_NODE_CLIENT.chain_id,seed0["output"])]))

    plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "rly keys restore '%s' fd1 '%s' " % (constants.COSMOS_NODE_CLIENT.chain_id_1,seed1["output"])]))

    plan.print("Adding the paths")

    plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "rly paths new my-chain chain-2 demo"]))

    plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "rly transact connection demo"]))

    plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "rly tx link demo -d -t 3s "]))

    plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "rly paths list"]))

    plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "rly start &"]))