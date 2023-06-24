WALLET = "github.com/hugobyte/dive/services/cosmvm/wallet.star"
START = "github.com/hugobyte/dive/services/cosmvm/start.sh"
DEFAULT_CONTRACT_PATH = "github.com/hugobyte/dive/services/cosmvm/static_files/contracts"
SERVICE_NAME = "cosmos"
IMAGE = "archwaynetwork/archwayd:constantine"
RPC_PORT_KEY = "rpc"
PATH = "/start-scripts/"
CONTRACT_PATH = "/root/contracts/"


def start_cosmos_node(plan,args):

    plan.print("Launching " +SERVICE_NAME+  " deployment service")

    plan.upload_files(src=START, name="start-script")
    plan.upload_files(src=DEFAULT_CONTRACT_PATH, name="contract")

    cosmwasm_node_config = ServiceConfig(
        image=IMAGE,
        files={ 
            PATH: "start-script",
            CONTRACT_PATH: "contract",
        },
        ports={
            "grpc" : PortSpec(number=9090),
            "http" : PortSpec(number=9091),
            "tcp" : PortSpec(number=26656),
            "rpc" : PortSpec(number=26657)
        },
        public_ports={
            "grpc" : PortSpec(number=9090 ),
            "http" : PortSpec(number=9091),
            "tcp"  :  PortSpec(number=26656),
            "rpc" : PortSpec(number=4564)
        },
        
        entrypoint=["/bin/sh","-c","cd ../../start-scripts && chmod +x start.sh && ./start.sh"]
    )

    node_service_response = plan.add_service(name=SERVICE_NAME, config= cosmwasm_node_config)
