ICON_SERVICE_NAME = "ICON"
ICON_NODE_IMAGE = "iconloop/goloop-icon:v1.3.5"
ICON_RPC_PRIVATE_PORT = 9080
ICON_RPC_PUBLIC_PORT = 8090
EXECUTABLE_PATH = "./bin/goloop"
ICON_BASE_CONFIG_FILES_PATH = "/goloop/config/"
ICON_CONTRACT_DIR = "/goloop/contracts/"
ICON_BASE_CONFIG_FILES_KEY = "config_file_path"
DEFAULT_ICON_BASE_CONFIG_FILES_PATH = "github.com/hugobte/chain-pacakge/services/jvm/icon/static-files/config/"
ICON_CONTRACT_DIR_KEY = "contract_file_path"
DEFAULT_ICON_CONTRACT_DIR = "github.com/hugobte/chain-pacakge/services/jvm/icon/static-files/contracts/"
ICON_RPC_PORT_KEY = "rpc"
PUBLIC_IP_ADDRESS = "127.0.0.1"
ICON_RPC_ENDPOINT_PATH = "api/v3/icon_dex"

def start_icon_node(plan,args):

    plan.print("Launching "+ICON_SERVICE_NAME+" Deployment Service")

    service_params = args.get("service_params")

    icon_base_config_files = service_params.get(ICON_BASE_CONFIG_FILES_KEY,DEFAULT_ICON_BASE_CONFIG_FILES_PATH)
    icon_contract_files = service_params.get(ICON_CONTRACT_DIR_KEY,DEFAULT_ICON_CONTRACT_DIR)

    plan.print("Uploading Files")
    plan.upload_files(src=icon_base_config_files,name="config-files")
    plan.upload_files(src=icon_contract_files,name="contracts")

    icon_node_service_config = ServiceConfig(
        image=ICON_NODE_IMAGE,
        ports={
            ICON_RPC_PORT_KEY : PortSpec(number=ICON_RPC_PRIVATE_PORT,transport_protocol="TCP",application_protocol="http")
        },
        public_ports = {
            ICON_RPC_PORT_KEY : PortSpec(number=ICON_RPC_PUBLIC_PORT,transport_protocol="TCP",application_protocol="http")
        },
        files={
            ICON_BASE_CONFIG_FILES_PATH : "config-files",
            ICON_CONTRACT_DIR : "contracts"
        },
        env_vars={
            "GOLOOP_LOG_LEVEL": "trace",
            "GOLOOP_RPC_ADDR": ":9080",
            "GOLOOP_P2P_LISTEN": ":7080",
            "ICON_CONFIG": ICON_BASE_CONFIG_FILES_PATH+"icon_config.json"
        },
        cmd= ["/bin/sh","-c",ICON_BASE_CONFIG_FILES_PATH+"start.sh"]

    )

    icon_node_service_response = plan.add_service(name=ICON_SERVICE_NAME,config=icon_node_service_config)
    plan.exec(service_name=icon_node_service_response.name,recipe=ExecRecipe(command=["apk","add","jq"]))

    public_url = get_service_url(PUBLIC_IP_ADDRESS,icon_node_service_config.public_ports,ICON_RPC_ENDPOINT_PATH)
    private_url = get_service_url(icon_node_service_response.ip_address,icon_node_service_response.ports,ICON_RPC_ENDPOINT_PATH)
    
    response = struct(
        service_config = icon_node_service_config,
        node_service_response = icon_node_service_response,
        public_url = public_url,
        private_url = private_url
    )
    
    return response

def get_service_url(ip_address,ports,path):
    port_id = ports[ICON_RPC_PORT_KEY].number
    protocol = ports[ICON_RPC_PORT_KEY].application_protocol
    url = "{0}://{1}:{2}/{3}".format(protocol,ip_address,port_id,path)
    return url