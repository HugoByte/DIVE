ICON_SERVICE_NAME = "ICON_NODE"
ICON_NODE_IMAGE = "iconloop/goloop-icon:v1.3.5"
ICON_RPC_PORT = 9080
EXECUTABLE_PATH = "./bin/goloop"
ICON_BASE_CONFIG_FILES_PATH = "/goloop/config/"
ICON_CONTRACT_DIR = "/goloop/contracts/"
ICON_BASE_CONFIG_FILES_KEY = "base_config"
DEFAULT_ICON_BASE_CONFIG_FILES_PATH = "github.com/hugobte/chain-pacakge/articats/jvm/base-config/"
ICON_CONTRACT_DIR_KEY = "contracts"
DEFAULT_ICON_CONTRACT_DIR = "github.com/hugobte/chain-pacakge/articats/jvm/contracts/"
ICON_RPC_PORT_KEY = "rpc"

def start_icon_node(plan,args):

    plan.print("Launching"+ICON_SERVICE_NAME+"Deployment Service")

    icon_base_config_files = args.get(ICON_BASE_CONFIG_FILES_KEY,DEFAULT_ICON_BASE_CONFIG_FILES_PATH)
    contract_binaries = args.get(ICON_CONTRACT_DIR_KEY,DEFAULT_ICON_CONTRACT_DIR)

    plan.print("Uploading Files")

    plan.upload_files(src=icon_base_config_files,name="config-files")
    plan.upload_files(src=contract_binaries,name="contracs")

    icon_node_service_config = ServiceConfig(
        image=ICON_NODE_IMAGE,
        ports={
            ICON_RPC_PORT_KEY : PortSpec(number=ICON_RPC_PORT,transport_protocol="TCP")
        },
        files={
            ICON_BASE_CONFIG_FILES_PATH : "config-files",
            ICON_CONTRACT_DIR : "contracs"
        },
        env_vars={
            "GOLOOP_LOG_LEVEL": "trace",
            "GOLOOP_RPC_ADDR": ":9080",
            "GOLOOP_P2P_LISTEN": ":7080",
            "ICON_CONFIG": ICON_BASE_CONFIG_FILES_PATH+"icon_config.json"
        },
        cmd= ["/bin/sh","-c",ICON_BASE_CONFIG_FILES_PATH+"start.sh"]

    )

    icon_node_service = plan.add_service(name=ICON_SERVICE_NAME,config=icon_node_service_config)

    plan.exec(service_name=icon_node_service.name,recipe=ExecRecipe(command=["apk","add","jq"]))

    return icon_node_service