<<<<<<< HEAD
<<<<<<< HEAD
ICON_NODE_IMAGE = "iconloop/goloop-icon:v1.3.5"
ICON_BASE_CONFIG_FILES_PATH = "/goloop/config/"
ICON_CONTRACT_DIR = "/goloop/contracts/"
ICON_BASE_CONFIG_FILES_KEY = "config_file_path"
ICON_BASE_CONFIG_FILES = "github.com/hugobyte/dive/services/jvm/icon/static-files/config/"
ICON_CONTRACT_DIR_KEY = "contract_file_path"
ICON_CONTRACT_FILES = "github.com/hugobyte/dive/services/jvm/icon/static-files/contracts/"
ICON_RPC_PORT_KEY = "rpc"
PUBLIC_IP_ADDRESS = "127.0.0.1"
ICON_RPC_ENDPOINT_PATH = "api/v3/icon_dex"

=======
constants = import_module("github.com/hugobyte/dive/package_io/constants.star")
>>>>>>> main
=======
constants = import_module("github.com/hugobyte/dive/package_io/constants.star")
>>>>>>> origin/development

# Starts The Icon Node 
def start_icon_node(plan,service_config,id,start_file_name):

    icon_node_constants = constants.ICON_NODE_CLIENT

    service_name = service_config.service_name
    private_port = service_config.private_port
    public_port = service_config.public_port
    network_name = service_config.network_name
    p2p_listen_address = service_config.p2p_listen_address
    p2p_address = service_config.p2p_address
    cid = service_config.cid


    plan.print("Launching "+service_name+" Service")


    plan.print("Uploading Files for %s Service" % service_name) 
    plan.upload_files(src=icon_node_constants.config_files_path,name="config-files-{0}".format(id))
    plan.upload_files(src=icon_node_constants.contract_files_path,name="contracts-{0}".format(id))
    plan.upload_files(src=icon_node_constants.keystore_files_path,name="kesytore-{0}".format(id) )

    icon_node_service_config = ServiceConfig(
        image=icon_node_constants.node_image,
        ports={
            icon_node_constants.port_key : PortSpec(number=private_port,transport_protocol="TCP",application_protocol="http")
        },
        public_ports = {
            icon_node_constants.port_key : PortSpec(number=public_port,transport_protocol="TCP",application_protocol="http")
        },
        files={
            icon_node_constants.config_files_directory : "config-files-{0}".format(id),
            icon_node_constants.contracts_directory : "contracts-{0}".format(id),
            icon_node_constants.keystore_directory : "kesytore-{0}".format(id)
        },
        env_vars={
            "GOLOOP_LOG_LEVEL": "trace",
            "GOLOOP_RPC_ADDR": ":%s" % private_port,
            "GOLOOP_P2P_LISTEN": ":%s" % p2p_listen_address,
            "GOLOOP_P2P": ":%s" % p2p_address,
            "ICON_CONFIG": icon_node_constants.config_files_directory+"icon_config.json"
        },
        cmd= ["/bin/sh","-c",icon_node_constants.config_files_directory+"%s" % start_file_name]

    )

    icon_node_service_response = plan.add_service(name=service_name,config=icon_node_service_config)
    plan.exec(service_name=service_name,recipe=ExecRecipe(command=["/bin/sh","-c","apk add jq"]))

    public_url = get_service_url(icon_node_constants.public_ip_address,icon_node_service_config.public_ports,icon_node_constants.rpc_endpoint_path)
    private_url = get_service_url(icon_node_service_response.ip_address,icon_node_service_response.ports,icon_node_constants.rpc_endpoint_path)

    chain_id = plan.exec(service_name=service_name,recipe=ExecRecipe(command=["/bin/sh","-c","./bin/goloop chain inspect %s --format {{.NID}} | tr -d '\n\r'" % cid]))

    network = "{0}.icon".format(chain_id["output"])
    
    
    return struct(
        service_name = service_name,
        network_name = network_name,
        network = network,
        nid = chain_id["output"],
        endpoint = private_url,
        endpoint_public = public_url,
        keystore_path = "keystores/keystore.json",
        keypassword= "gochain"
    )

"""
Returns URL
'ip_address' - Ip of the service running
'ports' - port on which service running
'path' - enpoint path
"""
def get_service_url(ip_address,ports,path):
    port_id = ports[constants.ICON_NODE_CLIENT.port_key].number
    protocol = ports[constants.ICON_NODE_CLIENT.port_key].application_protocol
    url = "{0}://{1}:{2}/{3}".format(protocol,ip_address,port_id,path)
    return url


# Retruns Service Config
def get_service_config(id,private_port,public_port,p2p_listen_address,p2p_address,cid):

    return struct(
        service_name = "{0}{1}".format(constants.ICON_NODE_CLIENT.service_name,id),
        private_port = private_port,
        public_port = public_port,
        network_name = "icon-{0}".format(id),
        p2p_listen_address = p2p_listen_address,
        p2p_address = p2p_address,
        cid = cid
    )