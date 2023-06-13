SERVICE_NAME = "cosmos"
IMAGE = "archwaynetwork/archwayd:constantine"
RPC_PORT_KEY = "rpc"
RPC_PRIVATE_PORT = 7070
RPC_PUBLIC_PORT = 7070
PATH = "/start-scripts/"

def run(plan,args):

    plan.print("Launching " +SERVICE_NAME+  " deployment service")

    cosmwasm_node_config = ServiceConfig(
        image=IMAGE,
        # files={PATH: "start-script"},
        # ports={
        #     RPC_PORT_KEY : PortSpec(number=RPC_PUBLIC_PORT,transport_protocol="TCP",application_protocol="http")
        # },
        entrypoint=["/bin/sh","-c","sleep 999999999999"],
         env_vars= {
            "PASSCODE" : "prathiksha"
        }
        # cmd= ["/bin/sh", "-c", "chmod +x start.sh &&./start.sh"]
    )

    node_service_response = plan.add_service(name=SERVICE_NAME, config= cosmwasm_node_config)

    plan.exec(service_name=node_service_response.name, recipe=ExecRecipe(command=["/bin/sh", "-c", "mkdir test"],))
    plan.exec(service_name=node_service_response.name, recipe=ExecRecipe(command=["/bin/sh", "-c", "cd test"],))
    plan.exec(service_name=node_service_response.name, recipe=ExecRecipe(command=["/bin/sh", "-c", "mkdir node1"],))
    plan.exec(service_name=node_service_response.name, recipe=ExecRecipe(command=["/bin/sh", "-c", "cd node1"],))

    # initialise the node
    plan.exec(service_name=node_service_response.name, recipe=ExecRecipe(command=["/bin/sh", "-c", "archwayd init node1 --chain-id my-chain --home ./node1"],))

    # adding the keys
    plan.exec(service_name=node_service_response.name, recipe=ExecRecipe(command=["/bin/sh", "-c", "(echo $PASSCODE; echo $PASSCODE) | archwayd keys add node1-account --home ./node1 "],))

    #listing the keys
    plan.exec(service_name=node_service_response.name, recipe=ExecRecipe(command=["archwayd", "keys", "list"],))

    # Adding the key to genesis account
    plan.exec(service_name=node_service_response.name, recipe=ExecRecipe(command=["/bin/sh", "-c",  "(echo $PASSCODE; echo $PASSCODE) | archwayd add-genesis-account $(archwayd keys show node1-account -a --home ./node1) 1000000000stake --home ./node1 "],))

    # Generate genesis transaction
    plan.exec(service_name=node_service_response.name, recipe=ExecRecipe(command=["/bin/sh", "-c", "(echo $PASSCODE; echo $PASSCODE) | archwayd gentx node1-account 1000000000stake --chain-id my-chain --home ./node1"],))

    # collect genesis transcation
    plan.exec(service_name=node_service_response.name, recipe=ExecRecipe(command=["/bin/sh", "-c", "archwayd collect-gentxs --home ./node1"],))

    # start the node
    plan.exec(service_name=node_service_response.name, recipe=ExecRecipe(command=["/bin/sh", "-c", "archwayd start --home ./node1 &"],))
  
