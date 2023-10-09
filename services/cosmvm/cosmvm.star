archway_node_service = import_module("./archway/archway.star")
neutron_node_service = import_module("./neutron/neutron.star")
parser = import_module("../../package_io/input_parser.star")

def start_cosmvm_chains(plan, node_name, chain_id = None, key = None, password = None, public_grpc = None, public_http = None, public_tcp = None, public_rpc = None):
    """_summary_

    Args:
        plan: Plan
        node_name: Cosmos Supported Chain Name
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

    if node_name == "archway":
        return archway_node_service.start_node_service(plan, chain_id, key, password, public_grpc, public_http, public_tcp, public_rpc)

    elif node_name == "neutron":
        return neutron_node_service.start_node_service(plan, chain_id, key, password, public_grpc, public_http, public_tcp, public_rpc)

def start_ibc_between_cosmvm_chains(plan, chain_a, chain_b):
    """
    To Start IBC Connection between Chain A and Chain B.

    Args:
        plan: plan object.
        chain_a: Sepcify Chain A.
        chain_b:Sepcify Chain B.

    Returns:
        _type_: _description_
    """

    if chain_a == "archway" and chain_b == "archway":
        return archway_node_service.start_nodes_services_archway(plan)

    elif chain_a == "neutron" and chain_b == "neutron":
        return neutron_node_service.start_node_services(plan)

    elif chain_a == "neutron" and chain_b == "archway":
        chain_a_service = parser.struct_to_dict(neutron_node_service.start_node_service(plan))
        chain_b_service = parser.struct_to_dict(archway_node_service.start_node_service(plan))
        return struct(
            src_config = chain_a_service,
            dst_config = chain_b_service,
        )

    elif chain_a == "archway" and chain_b == "neutron":
        chain_a_service = parser.struct_to_dict(archway_node_service.start_node_service(plan))
        chain_b_service = parser.struct_to_dict(neutron_node_service.start_node_service(plan))
        return struct(
            src_config = chain_a_service,
            dst_config = chain_b_service,
        )
