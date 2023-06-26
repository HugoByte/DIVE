ICON_NODE_IMAGE = "iconloop/goloop-icon:v1.3.7"
ICON_BASE_CONFIG_FILES_PATH = "/goloop/config/"
ICON_CONTRACT_DIR = "/goloop/contracts/"
ICON_BASE_CONFIG_FILES_KEY = "config_file_path"
ICON_BASE_CONFIG_FILES = "github.com/hugobyte/dive/services/jvm/icon/static-files/config/"
ICON_CONTRACT_DIR_KEY = "contract_file_path"
ICON_CONTRACT_FILES = "github.com/hugobyte/dive/services/jvm/icon/static-files/contracts/"
ICON_RPC_PORT_KEY = "rpc"
PUBLIC_IP_ADDRESS = "127.0.0.1"
ICON_RPC_ENDPOINT_PATH = "api/v3/icon_dex"


# Starts The Icon Node 
def start_icon_node(plan,service_config,id,start_file_name):

    service_name = service_config.service_name
    private_port = service_config.private_port
    public_port = service_config.public_port
    network_name = service_config.network_name
    p2p_listen_address = service_config.p2p_listen_address
    p2p_address = service_config.p2p_address
    cid = service_config.cid


    plan.print("Launching "+service_name+" Service")


    plan.print("Uploading Files")
    plan.upload_files(src=ICON_BASE_CONFIG_FILES,name="config-files-{0}".format(id))
    plan.upload_files(src=ICON_CONTRACT_FILES,name="contracts-{0}".format(id))

    icon_node_service_config = ServiceConfig(
        image=ICON_NODE_IMAGE,
        ports={
            ICON_RPC_PORT_KEY : PortSpec(number=private_port,transport_protocol="TCP",application_protocol="http")
        },
        public_ports = {
            ICON_RPC_PORT_KEY : PortSpec(number=public_port,transport_protocol="TCP",application_protocol="http")
        },
        files={
            ICON_BASE_CONFIG_FILES_PATH : "config-files-{0}".format(id),
            ICON_CONTRACT_DIR : "contracts-{0}".format(id),
        },
        env_vars={
            "GOLOOP_LOG_LEVEL": "trace",
            "GOLOOP_RPC_ADDR": ":%s" % private_port,
            "GOLOOP_P2P_LISTEN": ":%s" % p2p_listen_address,
            "GOLOOP_P2P": ":%s" % p2p_address,
            "ICON_CONFIG": ICON_BASE_CONFIG_FILES_PATH+"icon_config.json"
        },
        cmd= ["/bin/sh","-c",ICON_BASE_CONFIG_FILES_PATH+"%s" % start_file_name]

    )

    icon_node_service_response = plan.add_service(name=service_name,config=icon_node_service_config)
    plan.exec(service_name=service_name,recipe=ExecRecipe(command=["/bin/sh","-c","apk add jq"]))

    public_url = get_service_url(PUBLIC_IP_ADDRESS,icon_node_service_config.public_ports,ICON_RPC_ENDPOINT_PATH)
    private_url = get_service_url(icon_node_service_response.ip_address,icon_node_service_response.ports,ICON_RPC_ENDPOINT_PATH)

    chain_id = plan.exec(service_name=service_name,recipe=ExecRecipe(command=["/bin/sh","-c","./bin/goloop chain inspect %s --format {{.NID}} | tr -d '\n\r'" % cid]))

    network = "{0}.icon".format(chain_id["output"])
    
    
    return struct(
        service_name = service_name,
        network_name = network_name,
        network = network,
        nid = chain_id["output"],
        endpoint = private_url,
        endpoint_public = public_url,
        keystore_path = "config/keystore.json",
        keypassword= "gochain"
    )

"""
Returns URL
'ip_address' - Ip of the service running
'ports' - port on which service running
'path' - enpoint path
"""
def get_service_url(ip_address,ports,path):
    port_id = ports[ICON_RPC_PORT_KEY].number
    protocol = ports[ICON_RPC_PORT_KEY].application_protocol
    url = "{0}://{1}:{2}/{3}".format(protocol,ip_address,port_id,path)
    return url


# Retruns Service Config
def get_service_config(id,private_port,public_port,p2p_listen_address,p2p_address,cid):

    return struct(
        service_name = "icon-node-{0}".format(id),
        private_port = private_port,
        public_port = public_port,
        network_name = "icon-{0}".format(id),
        p2p_listen_address = p2p_listen_address,
        p2p_address = p2p_address,
        cid = cid
    )