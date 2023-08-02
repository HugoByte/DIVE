constants = import_module("github.com/hugobyte/dive/package_io/constants.star")

# Starts The Icon Node 
def start_icon_node(plan,service_config,id,uploaded_genesis,genesis_file_path,genesis_file_name):

    plan.print(uploaded_genesis)

    icon_node_constants = constants.ICON_NODE_CLIENT

    service_name = service_config["service_name"]
    private_port = service_config["private_port"]
    public_port = service_config["public_port"]
    network_name = service_config["network_name"]
    p2p_listen_address = service_config["p2p_listen_address"]
    p2p_address = service_config["p2p_address"]
    cid = service_config["cid"]


    plan.print("Launching "+service_name+" Service")


    plan.print("Uploading Files for %s Service" % service_name) 
    plan.upload_files(src=icon_node_constants.config_files_path,name="config-files-{0}".format(id))
    plan.upload_files(src=icon_node_constants.contract_files_path,name="contracts-{0}".format(id))
    plan.upload_files(src=icon_node_constants.keystore_files_path,name="kesytore-{0}".format(id) )

    file_path = ""
    file_name = ""
    if len(uploaded_genesis) == 0:
       plan.upload_files(src=genesis_file_path,name=genesis_file_name)
       file_path = genesis_file_name
       file_name = genesis_file_name
    else:
        file_path = uploaded_genesis["file_path"]
        file_name = uploaded_genesis["file_name"]

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
            icon_node_constants.keystore_directory : "kesytore-{0}".format(id),
            icon_node_constants.genesis_file_path : file_path

        },
        env_vars={
            "GOLOOP_LOG_LEVEL": "trace",
            "GOLOOP_RPC_ADDR": ":%s" % private_port,
            "GOLOOP_P2P_LISTEN": ":%s" % p2p_listen_address,
            "GOLOOP_P2P": ":%s" % p2p_address,
            "ICON_CONFIG": icon_node_constants.config_files_directory+"icon_config.json"
        },
        cmd= ["/bin/sh","-c",icon_node_constants.config_files_directory+"start.sh %s %s" % (cid,file_name)]

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

    config = {
        "service_name" : "{0}{1}".format(constants.ICON_NODE_CLIENT.service_name,id),
        "private_port" : private_port,
        "public_port" : public_port,
        "network_name" : "icon-{0}".format(id),
        "p2p_listen_address" : p2p_listen_address,
        "p2p_address" : p2p_address,
        "cid":cid
    }


    return config