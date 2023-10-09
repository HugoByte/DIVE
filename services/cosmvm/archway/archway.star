# Import necessary modules
archway_node_service = import_module("./src/node-setup/start_node.star")
constants = import_module("../../../package_io/constants.star")
archway_node_0_constant_config = constants.ARCHAY_NODE0_CONFIG
archway_node_1_constant_config = constants.ARCHAY_NODE1_CONFIG
archway_service_config = constants.ARCHWAY_SERVICE_CONFIG
archway_private_ports = constants.COMMON_ARCHWAY_PRIVATE_PORTS
parser = import_module("../../../package_io/input_parser.star")

# Function to start two Cosmos nodes for Archway
def start_nodes_services_archway(plan):
    """
    Configure and start two Cosmos nodes for Archway.

    Args:
        plan (plan): plan.

    Returns:
        struct: Configuration information for source and destination services.
    """
    service_name_src = "{0}-{1}".format(archway_service_config.service_name, archway_node_0_constant_config.chain_id)
    service_name_dst = "{0}-{1}".format(archway_service_config.service_name, archway_node_1_constant_config.chain_id)

    src_config = parser.struct_to_dict(archway_node_service.start_cosmos_node(
        plan,
        archway_node_0_constant_config.chain_id,
        archway_node_0_constant_config.key,
        archway_service_config.password,
        service_name_src,
        archway_private_ports.grpc,
        archway_private_ports.http,
        archway_private_ports.tcp,
        archway_private_ports.rpc,
        archway_node_0_constant_config.grpc,
        archway_node_0_constant_config.http,
        archway_node_0_constant_config.tcp,
        archway_node_0_constant_config.rpc,
    ))

    dst_config = parser.struct_to_dict(archway_node_service.start_cosmos_node(
        plan,
        archway_node_1_constant_config.chain_id,
        archway_node_1_constant_config.key,
        archway_service_config.password,
        service_name_dst,
        archway_private_ports.grpc,
        archway_private_ports.http,
        archway_private_ports.tcp,
        archway_private_ports.rpc,
        archway_node_1_constant_config.grpc,
        archway_node_1_constant_config.http,
        archway_node_1_constant_config.tcp,
        archway_node_1_constant_config.rpc,
    ))

    return struct(
        src_config = src_config,
        dst_config = dst_config,
    )

# Function to start a single Cosmos node for Archway
def start_node_service(plan, chain_id = None, key = None, password = None, public_grpc = None, public_http = None, public_tcp = None, public_rpc = None):
    """
    Configure and start a single Cosmos node for Archway.

    Args:
        plan (plan): Plan object for service deployment.
        chain_id (str): Chain ID.
        key (str): Key.
        password (str): Password.
        public_grpc (int, optional): Public gRPC port.
        public_http (int, optional): Public HTTP port.
        public_tcp (int, optional): Public TCP port.
        public_rpc (int, optional): Public RPC port.

    Returns:
        struct: Configuration information for the service.
    """
    chain_id = chain_id if chain_id != None else archway_node_0_constant_config.chain_id
    key = key if key != None else archway_node_0_constant_config.key
    password = password if password != None else archway_service_config.password
    public_grpc = public_grpc if public_grpc != None else archway_node_0_constant_config.grpc
    public_http = public_http if public_http != None else archway_node_0_constant_config.http
    public_tcp = public_tcp if public_tcp != None else archway_node_0_constant_config.tcp
    public_rpc = public_rpc if public_rpc != None else archway_node_0_constant_config.rpc
    service_name_src = "{0}-{1}".format(archway_service_config.service_name, chain_id)
    return archway_node_service.start_cosmos_node(
        plan,
        chain_id,
        key,
        password,
        service_name_src,
        archway_private_ports.grpc,
        archway_private_ports.http,
        archway_private_ports.tcp,
        archway_private_ports.rpc,
        public_grpc,
        public_http,
        public_tcp,
        public_rpc,
    )
