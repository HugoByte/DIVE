
START = "github.com/hugobyte/dive/services/cosmvm/start-cosmos-0.sh"
START1 = "github.com/hugobyte/dive/services/cosmvm/start-cosmos-1.sh"
DEFAULT_CONTRACT_PATH = "github.com/hugobyte/dive/services/cosmvm/static_files/contracts"
SERVICE_NAME = "cosmos"
SERVICE_NAME1 = "cosmos1"
IMAGE = "archwaynetwork/archwayd:constantine"
PATH = "/start-scripts/"
CONTRACT_PATH = "/root/contracts/"
CHAIN_ID = "my-chain"
CHAIN_ID1 = "chain-1"

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
            "tcp"  : PortSpec(number=26656),
            "rpc"  : PortSpec(number=26657)
        },
        public_ports={
            "grpc" : PortSpec(number=9090),
            "http" : PortSpec(number=9091),
            "tcp"  : PortSpec(number=26656),
            "rpc"  : PortSpec(number=4564)
        },
        
        entrypoint=["/bin/sh","-c","cd ../../start-scripts && chmod +x start-cosmos-0.sh && ./start-cosmos-0.sh"]
    )

    node_service_response = plan.add_service(name=SERVICE_NAME, config= cosmwasm_node_config)

def get_service_config(service_name, cid):

    return struct(
        service_name = service_name,
        cid = cid
    )

def start_cosmos_node1(plan,args):

    plan.upload_files(src=START1, name="start-script1")

    cosmwasm_node_config1 = ServiceConfig(
        image=IMAGE,
        files={ 
            PATH: "start-script1",
            CONTRACT_PATH: "contract",
        },
        ports={
            "grpc" : PortSpec(number=9090),
            "http" : PortSpec(number=9091),
            "tcp"  : PortSpec(number=26656),
            "rpc"  : PortSpec(number=26657)
        },
        public_ports={
            "grpc" : PortSpec(number=9080),
            "http" : PortSpec(number=9092),
            "tcp"  : PortSpec(number=26658),
            "rpc"  : PortSpec(number=4566)
        },
        
        entrypoint=["/bin/sh","-c","cd ../../start-scripts && chmod +x start-cosmos-1.sh && ./start-cosmos-1.sh "]
    )

    node_service_response = plan.add_service(name=SERVICE_NAME1, config= cosmwasm_node_config1)

def get_service_config1(service_name, cid):

    return struct(
        service_name = service_name,
        cid = cid
    )
