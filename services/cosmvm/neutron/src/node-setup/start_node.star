# Import constants from an external module
constants = import_module("github.com/hugobyte/dive/package_io/constants.star")
neutron_node_constants = constants.NEUTRON_SERVICE_CONFIG
network_port_keys_and_ip = constants.NETWORK_PORT_KEYS_AND_IP_ADDRESS

def start_neutron_node(plan, args):
    """
    Start a Neutron node service with the provided configuration.

    Args:
        plan (Plan): The deployment plan.
        args (dict): Arguments for configuring the service.

    Returns:
        None
    """

    chain_id = args["cid"]
    key = args["key"]
    password = args["password"]
    service_name = args["service_name"]

    plan.print("Launching " + service_name + " deployment service")

    start_script_file = "start-script-%s" % chain_id
    init_script_file = "init-script-%s" % chain_id
    init_neutrond_script_file = "init-neutrond-script-%s" % chain_id

    plan.upload_files(src = neutron_node_constants.init_script, name = init_script_file)
    plan.upload_files(src = neutron_node_constants.init_nutrond_script, name = init_neutrond_script_file)
    plan.upload_files(src = neutron_node_constants.start_script, name = start_script_file)

    # Define Neutron node configuration
    neutron_node_config = ServiceConfig(
        image=neutron_node_constants.image,
        files = {
            neutron_node_constants.path + "start": start_script_file,
            neutron_node_constants.path + "init": init_script_file,
            neutron_node_constants.path + "init-neutrond": init_neutrond_script_file,
        },
        ports={
            network_port_keys_and_ip.http: PortSpec(
                number=args["private_http"],
                transport_protocol=network_port_keys_and_ip.tcp.upper(),
                application_protocol=network_port_keys_and_ip.http,
                wait="2m"
            ),
            network_port_keys_and_ip.rpc: PortSpec(
                number = args["private_rpc"], 
                transport_protocol =network_port_keys_and_ip.tcp.upper(),
                application_protocol =network_port_keys_and_ip.http, 
                wait = "2m"
            ),
            network_port_keys_and_ip.tcp: PortSpec(
                number = args["private_tcp"], 
                transport_protocol =network_port_keys_and_ip.tcp.upper(), 
                application_protocol =network_port_keys_and_ip.http, 
                wait = "2m"
            ),
            network_port_keys_and_ip.grpc: PortSpec(
                number = args["private_grpc"],
                transport_protocol =network_port_keys_and_ip.tcp.upper(), 
                application_protocol =network_port_keys_and_ip.http, 
                wait = "2m"
            ),

        },
        public_ports={
            network_port_keys_and_ip.http: PortSpec(
                number=args["public_http"],
                transport_protocol=network_port_keys_and_ip.tcp.upper(),
                application_protocol=network_port_keys_and_ip.http,
                wait="2m"
            ),
            network_port_keys_and_ip.rpc: PortSpec(
                number = args["public_rpc"], 
                transport_protocol =network_port_keys_and_ip.tcp.upper(), 
                application_protocol =network_port_keys_and_ip.http, 
                wait = "2m"
            ),
            network_port_keys_and_ip.tcp: PortSpec(
                number = args["public_tcp"], 
                transport_protocol =network_port_keys_and_ip.tcp.upper(), 
                application_protocol =network_port_keys_and_ip.http, 
                wait = "2m"
            ),
            network_port_keys_and_ip.grpc: PortSpec(
                number = args["public_grpc"], 
                transport_protocol =network_port_keys_and_ip.tcp.upper(), 
                application_protocol =network_port_keys_and_ip.http, 
                wait = "2m"
            ),
        },
        entrypoint=["/bin/sh", "-c"],
        cmd = ["chmod +x ../..%s/init/init.sh && chmod +x ../..%s/start/start.sh && chmod +x ../..%s/init-neutrond/init-neutrond.sh && key=%s password=\"%s\" CHAINID=%s ../..%s/init/init.sh && key=%s password=\"%s\" CHAINID=%s ../..%s/init-neutrond/init-neutrond.sh && key=%s password=\"%s\" CHAINID=%s ../..%s/start/start.sh" % (neutron_node_constants.path, neutron_node_constants.path, neutron_node_constants.path, key, password, chain_id, neutron_node_constants.path, key, password, chain_id,neutron_node_constants.path, key, password, chain_id, neutron_node_constants.path)],
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


def get_service_config(cid, key, password, private_grpc, private_http, private_tcp, private_rpc, public_grpc, public_http, public_tcp, public_rpc):
    """
    Get the service configuration based on provided port values.

    Args:
        private_grpc (int): Private gRPC port.
        private_http (int): Private HTTP port.
        private_tcp (int): Private TCP port.
        private_rpc (int): Private RPC port.
        public_http (int): Public HTTP port.
        public_tcp (int): Public TCP port.
        public_rpc (int): Public RPC port.
        public_rpc (int): Public RPC port.
        service_name (str): Service name.


    Returns:
        dict: Service configuration dictionary.
    """

    service_name = "{0}-{1}".format(neutron_node_constants.service_name, cid)
    config = {
        "public_grpc": public_grpc,
        "public_http": public_http,
        "public_tcp": public_tcp,
        "public_rpc": public_rpc,
        "private_http": private_http,
        "private_tcp": private_tcp,
        "private_rpc": private_rpc,
        "private_grpc": private_grpc,
        "service_name": service_name,
        "cid" : cid,
        "key" : key,
        "password" : password,
    }
    return config
