# Import required modules and constants
neutron_node_service = import_module("./src/node-setup/start_node.star")
constants = import_module("../../../package_io/constants.star")
neutron_private_ports = constants.NEUTRON_PRIVATE_PORTS
neutron_node1_config = constants.NEUTRON_NODE1_CONFIG
neutron_node2_config = constants.NEUTRON_NODE2_CONFIG

def start_node_services(plan, args):
    """
    Configure and start two Neutron node services, serving as the source and destination,
    to establish an IBC relay connection between them.

    Args:
        plan (plan): plan.
        args (dict): Arguments containing data for configuring the services.

    Returns:
        struct: Configuration information for the source and destination services.
    """

    data_src = args["src_config"]["data"]
    data_dst = args["dst_config"]["data"] 
    src_chain_config = ""
    dst_chain_config = ""

    if len(data_src) != 0:
        # Configure the service based on provided data for source chain
        chain_id = data_src["chainId"]
        key = data_src["key"]
        password = data_src["password"]
        public_grpc = data_src["public_grpc"]
        public_tcp = data_src["public_tcp"]
        public_http = data_src["public_http"]
        public_rpc = data_src["public_rpc"]

        src_chain_config = neutron_node_service.get_service_config(
            chain_id, key, password,
            neutron_private_ports.grpc, neutron_private_ports.http, neutron_private_ports.tcp, neutron_private_ports.rpc,
            public_grpc, public_http, public_tcp, public_rpc
        )
    else:
        # Use predefined port values for configuration for source chain
        src_chain_config = neutron_node_service.get_service_config(
            neutron_node1_config.chain_id, neutron_node1_config.key, neutron_node1_config.password,
            neutron_private_ports.grpc, neutron_private_ports.http,
            neutron_private_ports.tcp, neutron_private_ports.rpc,
            neutron_node1_config.grpc, neutron_node1_config.http,
            neutron_node1_config.tcp, neutron_node1_config.rpc
        )

    if len(data_dst) != 0:
        # Configure the service based on provided data for destination chain
        chain_id = data_dst["chainId"]
        key = data_dst["key"]
        password = data_dst["password"]
        public_grpc = data_dst["public_grpc"]
        public_tcp = data_dst["public_tcp"]
        public_http = data_dst["public_http"]
        public_rpc = data_dst["public_rpc"]

        dst_chain_config = neutron_node_service.get_service_config(
            chain_id, key, password,
            neutron_private_ports.grpc, neutron_private_ports.http, neutron_private_ports.tcp, neutron_private_ports.rpc,
            public_grpc, public_http, public_tcp, public_rpc
        )
    else:
        # Use predefined port values for configuration for destination chain
        dst_chain_config = neutron_node_service.get_service_config(
            neutron_node2_config.chain_id, neutron_node2_config.key, neutron_node2_config.password,
            neutron_private_ports.grpc, neutron_private_ports.http,
            neutron_private_ports.tcp, neutron_private_ports.rpc,
            neutron_node2_config.grpc, neutron_node2_config.http,
            neutron_node2_config.tcp, neutron_node2_config.rpc
        )
    
    # Start the source and destination Neutron node services
    src_chain_response = neutron_node_service.start_neutron_node(plan, src_chain_config)
    dst_chain_response = neutron_node_service.start_neutron_node(plan, dst_chain_config)

    # Create configuration dictionaries for both services
    src_service_config = {
        "service_name": src_chain_response.service_name,
        "endpoint": src_chain_response.endpoint,
        "endpoint_public": src_chain_response.endpoint_public,
        "chain_id": src_chain_response.chain_id,
        "chain_key": src_chain_response.chain_key
    }

    dst_service_config = {
        "service_name": dst_chain_response.service_name,
        "endpoint": dst_chain_response.endpoint,
        "endpoint_public": dst_chain_response.endpoint_public,
        "chain_id": dst_chain_response.chain_id,
        "chain_key": dst_chain_response.chain_key
    }

    return struct(
        src_config=src_service_config,
        dst_config=dst_service_config,
    )

def start_node_service(plan, args):
    """
    Start a Neutron node service with the provided configuration.

    Args:
        plan (plan): plan.
        args (dict): Arguments containing data for configuring the service.

    Returns:
        struct: The response from starting the Neutron node service.
    """

    data = args["data"]
    chain_config = ""

    if len(data) != 0:
        # Configure the service based on provided data
        chain_id = data["chainId"]
        key = data["key"]
        password = data["password"]
        public_grpc = data["public_grpc"]
        public_tcp = data["public_tcp"]
        public_http = data["public_http"]
        public_rpc = data["public_rpc"]

        chain_config = neutron_node_service.get_service_config(
            chain_id, key, password,
            neutron_private_ports.grpc, neutron_private_ports.http, neutron_private_ports.tcp, neutron_private_ports.rpc,
            public_grpc, public_http, public_tcp, public_rpc
        )
    else:
        # Use predefined port values for configuration
        chain_config = neutron_node_service.get_service_config(
            neutron_node1_config.chain_id, neutron_node1_config.key, neutron_node1_config.password,
            neutron_private_ports.grpc, neutron_private_ports.http,
            neutron_private_ports.tcp, neutron_private_ports.rpc,
            neutron_node1_config.grpc, neutron_node1_config.http,
            neutron_node1_config.tcp, neutron_node1_config.rpc
        )

    # Start the Neutron node service and return the response
    return neutron_node_service.start_neutron_node(plan, chain_config)
