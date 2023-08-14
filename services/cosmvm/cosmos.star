cosmvm_node = import_module("github.com/hugobyte/dive/services/cosmvm/src/node-setup/start_node.star")
constants = import_module("github.com/hugobyte/dive/package_io/constants.star")
node_constants = constants.COSMOS_NODE_CLIENT

# spins up the 2 comsos nodes
def start_node_service_cosmos_to_cosmos(plan):
    src_chain_config = cosmvm_node.get_service_config(constants.COSMOS_NODE_CLIENT.chain_id, constants.COSMOS_NODE_CLIENT.key, node_constants.private_port_grpc, node_constants.private_port_http, node_constants.private_port_tcp, node_constants.private_port_rpc, node_constants.public_port_grpc, node_constants.public_port_http, node_constants.public_port_tcp, node_constants.public_port_rpc, node_constants.password)
    dst_chain_config = cosmvm_node.get_service_config(constants.COSMOS_NODE_CLIENT.chain_id_1, constants.COSMOS_NODE_CLIENT.key1, node_constants.private_port_grpc, node_constants.private_port_http, node_constants.private_port_tcp, node_constants.private_port_rpc, node_constants.public_port_grpc_node, node_constants.public_port_http_node, node_constants.public_port_tcp_node, node_constants.public_port_rpc_node, node_constants.password)

    source_chain_response = cosmvm_node.start_cosmos_node(plan, src_chain_config)
    destination_chain_response = cosmvm_node.start_cosmos_node(plan, dst_chain_config)

    src_service_config = {
        "service_name": source_chain_response.service_name,
        "endpoint": source_chain_response.endpoint,
        "endpoint_public": source_chain_response.endpoint_public,
    }

    dst_service_config = {
        "service_name": destination_chain_response.service_name,
        "endpoint": destination_chain_response.endpoint,
        "endpoint_public": destination_chain_response.endpoint_public,
    }

    return struct(
        src_config = src_service_config,
        dst_config = dst_service_config,
    ), src_chain_config, dst_chain_config

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
        chain_config = cosmvm_node.get_service_config(constants.COSMOS_NODE_CLIENT.chain_id, constants.COSMOS_NODE_CLIENT.key, private_grpc, private_http, private_tcp, private_rpc, public_grpc, public_http, public_tcp, public_rpc, node_constants.password)
    else:
        chain_config = cosmvm_node.get_service_config(constants.COSMOS_NODE_CLIENT.chain_id, constants.COSMOS_NODE_CLIENT.key, node_constants.private_port_grpc, node_constants.private_port_http, node_constants.private_port_tcp, node_constants.private_port_rpc, node_constants.public_port_grpc, node_constants.public_port_http, node_constants.public_port_tcp, node_constants.public_port_rpc, node_constants.password)

    return cosmvm_node.start_cosmos_node(plan, chain_config)
