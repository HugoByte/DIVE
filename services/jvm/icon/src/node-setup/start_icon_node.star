constants = import_module("../../../../../package_io/constants.star")
network_keys_and_public_address = constants.NETWORK_PORT_KEYS_AND_IP_ADDRESS

def start_icon_node(plan, private_port, public_port, p2p_listen_address, p2p_address, cid, uploaded_genesis, genesis_file_path, genesis_file_name):
    """
    Function to start an ICON node.

    Args:
        plan: kurtosis plan.
        private_port: The private port for the ICON node.
        public_port: The public port for the ICON node.
        p2p_listen_address: The P2P listen address for the ICON node.
        p2p_address: The P2P address for the ICON node.
        cid: The chain ID for the ICON network.
        uploaded_genesis: A dictionary containing uploaded genesis file data.
        genesis_file_path: The path to the genesis file.
        genesis_file_name: The name of the genesis file.

    Returns:
        Configuration data for the started ICON node service as a dictionary.
    """
    plan.print(uploaded_genesis)

    icon_node_constants = constants.ICON_NODE_CLIENT

    service_name = "{0}{1}".format(constants.ICON_NODE_CLIENT.service_name, cid)
    network_name = "icon-{0}".format(cid)

    plan.print("Launching " + service_name + " Service")

    plan.print("Uploading Files for %s Service" % service_name)
    plan.upload_files(src=icon_node_constants.config_files_path, name="config-files-{0}".format(cid))
    plan.upload_files(src=icon_node_constants.contract_files_path, name="contracts-{0}".format(cid))
    plan.upload_files(src=icon_node_constants.keystore_files_path, name="keystore-{0}".format(cid))

    file_path = ""
    file_name = ""
    if len(uploaded_genesis) == 0:
        plan.upload_files(src=genesis_file_path, name=genesis_file_name)
        file_path = genesis_file_name
        file_name = genesis_file_name
    else:
        file_path = uploaded_genesis["file_path"]
        file_name = uploaded_genesis["file_name"]

    icon_node_service_config = ServiceConfig(
        image=icon_node_constants.node_image,
        ports={
            network_keys_and_public_address.rpc: PortSpec(
                number=private_port, transport_protocol=network_keys_and_public_address.tcp.upper(),
                application_protocol=network_keys_and_public_address.http
            ),
        },
        public_ports={
            network_keys_and_public_address.rpc: PortSpec(
                number=public_port, transport_protocol=network_keys_and_public_address.tcp.upper(),
                application_protocol=network_keys_and_public_address.http
            ),
        },
        files={
            icon_node_constants.config_files_directory: "config-files-{0}".format(cid),
            icon_node_constants.contracts_directory: "contracts-{0}".format(cid),
            icon_node_constants.keystore_directory: "keystore-{0}".format(cid),
            icon_node_constants.genesis_file_path: file_path,
        },
        env_vars={
            "GOLOOP_LOG_LEVEL": "trace",
            "GOLOOP_RPC_ADDR": ":%s" % private_port,
            "GOLOOP_P2P_LISTEN": ":%s" % p2p_listen_address,
            "GOLOOP_P2P": ":%s" % p2p_address,
            "ICON_CONFIG": icon_node_constants.config_files_directory + "icon_config.json",
        },
        cmd=["/bin/sh", "-c", icon_node_constants.config_files_directory + "start.sh %s %s" % (cid, file_name)],
    )

    icon_node_service_response = plan.add_service(name=service_name, config=icon_node_service_config)
    plan.exec(service_name=service_name, recipe=ExecRecipe(command=["/bin/sh", "-c", "apk add jq"]))

    public_url = get_service_url(
        network_keys_and_public_address.public_ip_address, icon_node_service_config.public_ports,
        icon_node_constants.rpc_endpoint_path
    )
    private_url = get_service_url(icon_node_service_response.ip_address, icon_node_service_response.ports, icon_node_constants.rpc_endpoint_path)

    chain_id = plan.exec(
        service_name=service_name, recipe=ExecRecipe(
            command=["/bin/sh", "-c", "./bin/goloop chain inspect %s --format {{.NID}} | tr -d '\n\r'" % cid]
        )
    )

    network = "{0}.icon".format(chain_id["output"])

    return struct(
        service_name=service_name,
        network_name=network_name,
        network=network,
        nid=chain_id["output"],
        endpoint=private_url,
        endpoint_public=public_url,
        keystore_path="keystores/keystore.json",
        keypassword="gochain",
    )


def get_service_url(ip_address, ports, path):
    """_summary_

    Args:
        ip_address (string): 
        ports (int): _description_
        path (int): _description_

    Returns:
        string: service url
    """
    port_id = ports[network_keys_and_public_address.rpc].number
    protocol = ports[network_keys_and_public_address.rpc].application_protocol
    url = "{0}://{1}:{2}/{3}".format(protocol, ip_address, port_id, path)
    return url