# Import constants from an external module
constants = import_module("../../../../../package_io/constants.star")
neutron_node_constants = constants.NEUTRON_SERVICE_CONFIG
network_port_keys_and_ip = constants.NETWORK_PORT_KEYS_AND_IP_ADDRESS

def start_neutron_node(plan, chain_id, key, password, service_name, private_http_port, private_rcp_port, private_tcp_port, private_grpc_port, public_http_port, public_rcp_port, public_tcp_port, public_grpc_port):
    """
    Start a Neutron node service with the provided configuration.

    Args:
        plan (plan): plan.
        chain_id (str): Chain ID.
        key (str): Key.
        password (str): Mnemonic password between 20 - 34 characters.
        service_name (str): Service name.
        private_http_port (int): Private HTTP port.
        private_rpc_port (int): Private RPC port.
        private_tcp_port (int): Private TCP port.
        private_grpc_port (int): Private gRPC port.
        public_http_port (int): Public HTTP port.
        public_rpc_port (int): Public RPC port.
        public_tcp_port (int): Public TCP port.
        public_grpc_port (int): Public gRPC port.

    Returns:
        struct: Service configuration
    """
    plan.print("Launching " + service_name + " deployment service")

    start_script_file = "start-script-%s" % chain_id
    init_script_file = "init-script-%s" % chain_id
    init_neutrond_script_file = "init-neutrond-script-%s" % chain_id
    contract_files = "contract-%s" % chain_id

    plan.upload_files(src = neutron_node_constants.init_script, name = init_script_file)
    plan.upload_files(src = neutron_node_constants.init_nutrond_script, name = init_neutrond_script_file)
    plan.upload_files(src = neutron_node_constants.start_script, name = start_script_file)
    plan.upload_files(src= neutron_node_constants.default_contract_path, name=contract_files)

    # Define Neutron node configuration
    neutron_node_config = ServiceConfig(
        image=neutron_node_constants.image,
        files = {
            neutron_node_constants.path + "start": start_script_file,
            neutron_node_constants.path + "init": init_script_file,
            neutron_node_constants.path + "init-neutrond": init_neutrond_script_file,
            neutron_node_constants.contract_path: contract_files
        },
        ports={
            network_port_keys_and_ip.http: PortSpec(
                number = private_http_port,
                transport_protocol=network_port_keys_and_ip.tcp.upper(),
                application_protocol=network_port_keys_and_ip.http,
                wait="2m"
            ),
            network_port_keys_and_ip.rpc: PortSpec(
                number = private_rcp_port, 
                transport_protocol =network_port_keys_and_ip.tcp.upper(),
                application_protocol =network_port_keys_and_ip.http, 
                wait = "2m"
            ),
            network_port_keys_and_ip.tcp: PortSpec(
                number = private_tcp_port, 
                transport_protocol =network_port_keys_and_ip.tcp.upper(), 
                application_protocol =network_port_keys_and_ip.http, 
                wait = "2m"
            ),
            network_port_keys_and_ip.grpc: PortSpec(
                number = private_grpc_port,
                transport_protocol =network_port_keys_and_ip.tcp.upper(), 
                application_protocol =network_port_keys_and_ip.http, 
                wait = "2m"
            ),

        },
        public_ports={
            network_port_keys_and_ip.http: PortSpec(
                number = public_http_port,
                transport_protocol=network_port_keys_and_ip.tcp.upper(),
                application_protocol=network_port_keys_and_ip.http,
                wait="2m"
            ),
            network_port_keys_and_ip.rpc: PortSpec(
                number = public_rcp_port, 
                transport_protocol =network_port_keys_and_ip.tcp.upper(), 
                application_protocol =network_port_keys_and_ip.http, 
                wait = "2m"
            ),
            network_port_keys_and_ip.tcp: PortSpec(
                number = public_tcp_port, 
                transport_protocol =network_port_keys_and_ip.tcp.upper(), 
                application_protocol =network_port_keys_and_ip.http, 
                wait = "2m"
            ),
            network_port_keys_and_ip.grpc: PortSpec(
                number = public_grpc_port, 
                transport_protocol =network_port_keys_and_ip.tcp.upper(), 
                application_protocol =network_port_keys_and_ip.http, 
                wait = "2m"
            ),
        },
        entrypoint=["/bin/sh", "-c"],
        cmd = ["chmod +x ../..%s/init/init.sh && chmod +x ../..%s/start/start.sh && chmod +x ../..%s/init-neutrond/init-neutrond.sh && key=%s password=\"%s\" CHAINID=%s ../..%s/init/init.sh && CHAINID=%s ../..%s/init-neutrond/init-neutrond.sh && CHAINID=%s ../..%s/start/start.sh" % (neutron_node_constants.path, neutron_node_constants.path, neutron_node_constants.path, key, password, chain_id, neutron_node_constants.path, chain_id,neutron_node_constants.path, chain_id, neutron_node_constants.path)],
        env_vars={
            "RUN_BACKGROUND": "0",
        },
    )

    # Add the service to the plan
    node_service_response = plan.add_service(name=service_name, config=neutron_node_config)
    plan.print(node_service_response)

    # Get public and private url, (private IP returned by kurtosis service)
    public_url = get_service_url(network_port_keys_and_ip.public_ip_address, neutron_node_config.public_ports)
    private_url = get_service_url(node_service_response.ip_address, node_service_response.ports)

    #return service name  and endpoints
    return struct(
        service_name = service_name,
        endpoint = private_url,
        endpoint_public = public_url,
        chain_id = chain_id,
        chain_key = key
    )


def get_service_url(ip_address, ports):
    """
    Get the service URL based on IP address and ports.

    Args:
        ip_address (str): IP address of the service.
        ports (dict): Dictionary of service ports.

    Returns:
        str: The constructed service URL.
    """

    port_id = ports[network_port_keys_and_ip.rpc].number
    protocol = ports[network_port_keys_and_ip.rpc].application_protocol
    url = "{0}://{1}:{2}".format(protocol, ip_address, port_id)
    return url