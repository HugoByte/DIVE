# Import required modules and constants
neutron_node_service = import_module("./src/node-setup/start_node.star")
constants = import_module("../../../package_io/constants.star")
neutron_private_ports = constants.NEUTRON_PRIVATE_PORTS
neutron_node1_config = constants.NEUTRON_NODE1_CONFIG
neutron_node2_config = constants.NEUTRON_NODE2_CONFIG
neutron_service_config = constants.NEUTRON_SERVICE_CONFIG

def start_node_services(plan):
    """
    Configure and start two Neutron node services, serving as the source and destination to establish an IBC relay connection between them.

    Args:
        plan (plan): plan.

    Returns:
        struct: Configuration information for the source and destination services.
    """

    # Start the source and destination Neutron node services
    service_name_src = "{0}-{1}".format(neutron_service_config.service_name, neutron_node1_config.chain_id)
    service_name_dst = "{0}-{1}".format(neutron_service_config.service_name, neutron_node2_config.chain_id)
    src_chain_response = neutron_node_service.start_neutron_node(plan, neutron_node1_config.chain_id, neutron_node1_config.key, neutron_node1_config.password, service_name_src, neutron_private_ports.http, neutron_private_ports.rpc, neutron_private_ports.tcp, neutron_private_ports.grpc, neutron_node1_config.http, neutron_node1_config.rpc, neutron_node1_config.tcp, neutron_node1_config.grpc)
    dst_chain_response = neutron_node_service.start_neutron_node(plan, neutron_node2_config.chain_id, neutron_node2_config.key, neutron_node2_config.password, service_name_dst, neutron_private_ports.http, neutron_private_ports.rpc, neutron_private_ports.tcp, neutron_private_ports.grpc, neutron_node2_config.http, neutron_node2_config.rpc, neutron_node2_config.tcp, neutron_node2_config.grpc)

    # Create configuration dictionaries for both services
    src_service_config = {
        "service_name": src_chain_response.service_name,
        "endpoint": src_chain_response.endpoint,
        "endpoint_public": src_chain_response.endpoint_public,
        "chain_id": src_chain_response.chain_id,
        "chain_key": src_chain_response.chain_key,
    }

    dst_service_config = {
        "service_name": dst_chain_response.service_name,
        "endpoint": dst_chain_response.endpoint,
        "endpoint_public": dst_chain_response.endpoint_public,
        "chain_id": dst_chain_response.chain_id,
        "chain_key": dst_chain_response.chain_key,
    }

    return struct(
        src_config = src_service_config,
        dst_config = dst_service_config,
    )

def start_node_service(plan, chain_id = None, key = None, password = None, public_grpc = None, public_http = None, public_tcp = None, public_rpc = None):
    """
    Start a Neutron node service with the provided configuration.

    Args:
        plan: Plan
        chain_id: Chain Id of the chain to be started.
        key: Key used for creating account.
        password: Password for Key.
        public_grpc: GRPC Endpoint for chain to run.
        public_http: HTTP Endpoint for chain to run .
        public_tcp: TCP Endpoint for chain to run.
        public_rpc: RPC Endpoint for chain to run.

    Returns:
        struct: The response from starting the Neutron node service.

    """

    # Start the Neutron node service with default configuration and return the response
    chain_id = chain_id if chain_id != None else neutron_node1_config.chain_id
    key = key if key != None else neutron_node1_config.key
    password = password if key != None else neutron_node1_config.password
    public_http = public_http if public_http != None else neutron_node1_config.http
    public_rpc = public_rpc if public_rpc != None else neutron_node1_config.rpc
    public_tcp = public_tcp if public_tcp != None else neutron_node1_config.tcp
    public_grpc = public_grpc if public_grpc != None else neutron_node1_config.grpc
    service_name = "{0}-{1}".format(neutron_service_config.service_name, chain_id)

    return neutron_node_service.start_neutron_node(
        plan,
        chain_id,
        key,
        password,
        service_name,
        neutron_private_ports.http,
        neutron_private_ports.rpc,
        neutron_private_ports.tcp,
        neutron_private_ports.grpc,
        public_http,
        public_rpc,
        public_tcp,
        public_grpc,
    )
