START_COSMOS = "github.com/hugobyte/dive/services/cosmvm/start-cosmos-0.sh"
START_COSMOS_1 = "github.com/hugobyte/dive/services/cosmvm/start-cosmos-1.sh"
DEFAULT_CONTRACT_PATH = "github.com/hugobyte/dive/services/cosmvm/static_files/contracts"
SERVICE_NAME = "cosmos"
SERVICE_NAME_1 = "cosmos1"
IMAGE = "archwaynetwork/archwayd:constantine"
PATH = "/start-scripts/"
CONTRACT_PATH = "/root/contracts/"
CHAIN_ID = "my-chain"
CHAIN_ID_1 = "chain-1"
PUBLIC_IP_ADDRESS = "127.0.0.1"

def start_cosmos_node(plan,args):

    plan.print("Launching " +SERVICE_NAME+  " deployment service")

    plan.upload_files(src=START_COSMOS, name="start-script")
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

    plan.print(node_service_response)

    public_url = get_service_url(PUBLIC_IP_ADDRESS,cosmwasm_node_config.public_ports)
    private_url = get_service_url(node_service_response.ip_address,node_service_response.ports)

    return struct(
        service_name = SERVICE_NAME,
        endpoint = private_url,
        endpoint_public = public_url
    )

# returns url
def get_service_url(ip_address,ports):
    port_id = ports["rpc"].number
    url = "{0}:{1}".format(ip_address,port_id)
    return url

def get_service_config(service_name, cid):

    return struct(
        service_name = service_name,
        cid = cid
    )

def start_cosmos_node_1(plan,args):

    plan.upload_files(src=START_COSMOS_1, name="start-script1")

    cosmwasm_node_config_1 = ServiceConfig(
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

    node_service_response = plan.add_service(name=SERVICE_NAME_1, config= cosmwasm_node_config_1)

    public_url = get_service_url_1(PUBLIC_IP_ADDRESS,cosmwasm_node_config_1.public_ports)
    private_url = get_service_url_1(node_service_response.ip_address,node_service_response.ports)

    return struct(
        service_name = SERVICE_NAME_1,
        endpoint = private_url,
        endpoint_public = public_url
    )

# returns url
def get_service_url_1(ip_address,ports):
    port_id = ports["rpc"].number
    url = "{0}:{1}".format(ip_address,port_id)
    return url

def get_service_config_1(service_name, cid):

    return struct(
        service_name = service_name,
        cid = cid
    )
