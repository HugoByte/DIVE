START_NODE = "github.com/hugobyte/chain-package/services/cosmvm/start.sh"
SERVICE_NAME = "cosmos"
IMAGE = "archwaynetwork/archwayd:constantine"
RPC_PORT_KEY = "rpc"
RPC_PRIVATE_PORT = 7070
RPC_PUBLIC_PORT = 7070
PATH = "/start-scripts/"

def run(plan,args):

    plan.print("Launching " +SERVICE_NAME+  " deployment service")

    plan.upload_files(src=START_NODE, name="start-script")

    cosmwasm_node_config = ServiceConfig(
        image=IMAGE,
        files={PATH: "start-script"},
        # ports={
        #     RPC_PORT_KEY : PortSpec(number=RPC_PUBLIC_PORT,transport_protocol="TCP",application_protocol="http")
        # },
        entrypoint=["/bin/sh","-c","sleep 999999999999"],
        cmd= ["/bin/sh", "-c", "chmod +x start.sh &&./start.sh"]
    )

    node_service_response = plan.add_service(name=SERVICE_NAME, config= cosmwasm_node_config)
