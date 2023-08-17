archway_node_service = import_module("github.com/hugobyte/dive/services/cosmvm/archway/src/node-setup/start_node.star")
constants = import_module("github.com/hugobyte/dive/package_io/constants.star")
archway_node_0_constant_config = constants.ARCHAY_NODE0_CONFIG
archway_node_1_constant_config = constants.ARCHAY_NODE1_CONFIG
archway_service_config = constants.ARCHWAY_SERVICE_CONFIG
archway_private_ports = constants.COMMON_ARCHWAY_PRIVATE_PORTS

# spins up the 2 comsos nodes
def start_nodes_services_archway(plan):
    src_chain_config = archway_node_service.get_service_config(archway_node_0_constant_config.chain_id, archway_node_0_constant_config.key, archway_private_ports.grpc, archway_private_ports.http, archway_private_ports.tcp, archway_private_ports.rpc,archway_node_0_constant_config.grpc,archway_node_0_constant_config.http,archway_node_0_constant_config.tcp,archway_node_0_constant_config.rpc, archway_service_config.password)
    dst_chain_config = archway_node_service.get_service_config(archway_node_1_constant_config.chain_id, archway_node_1_constant_config.key, archway_private_ports.grpc, archway_private_ports.http, archway_private_ports.tcp, archway_private_ports.rpc, archway_node_1_constant_config.grpc, archway_node_1_constant_config.http, archway_node_1_constant_config.tcp, archway_node_1_constant_config.rpc, archway_service_config.password)

    source_chain_response = archway_node_service.start_cosmos_node(plan, src_chain_config)
    destination_chain_response = archway_node_service.start_cosmos_node(plan, dst_chain_config)

    src_service_config = {
        "service_name": source_chain_response.service_name,
        "endpoint": source_chain_response.endpoint,
        "endpoint_public": source_chain_response.endpoint_public,
        "chain_id": source_chain_response.chain_id,
        "chain_key": source_chain_response.chain_key
    }

    dst_service_config = {
        "service_name": destination_chain_response.service_name,
        "endpoint": destination_chain_response.endpoint,
        "endpoint_public": destination_chain_response.endpoint_public,
        "chain_id": destination_chain_response.chain_id,
        "chain_key": destination_chain_response.chain_key
    }

    return struct(
        src_config = src_service_config,
        dst_config = dst_service_config,
    )

# spins up the single cosmos node

def start_node_service(plan, args):
    data = args["data"]
    chain_config = ""
    if len(data) != 0:
        # Private ports
        private_grpc = data["private_grpc"]
        private_tcp = data["private_tcp"]
        private_http = data["private_http"]
        private_rpc = data["private_rpc"]

        # Public Ports

        public_grpc = data["public_grpc"]
        public_tcp = data["public_tcp"]
        public_http = data["public_http"]
        public_rpc = data["public_rpc"]
        chain_config = archway_node_service.get_service_config(archway_node_0_constant_config.chain_id, archway_node_0_constant_config.key, private_grpc, private_http, private_tcp, private_rpc, public_grpc, public_http, public_tcp, public_rpc, archway_service_config.password)
    else:
        chain_config = archway_node_service.get_service_config(archway_node_0_constant_config.chain_id, archway_node_0_constant_config.key, archway_private_ports.grpc, archway_private_ports.http, archway_private_ports.tcp, archway_private_ports.rpc, archway_node_0_constant_config.grpc, archway_node_0_constant_config.http, archway_node_0_constant_config.tcp, archway_node_0_constant_config.rpc, archway_service_config.password)

    return archway_node_service.start_cosmos_node(plan, chain_config)
