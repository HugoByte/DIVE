archway_node_service = import_module("./archway/archway.star")
neutron_node_service = import_module("./neutron/neutron.star")
parser = import_module("../../package_io/input_parser.star")

def start_cosmvm_chains(plan, node_name, chain_id = None, key = None, password = None, public_grpc = None, public_http = None, public_tcp = None, public_rpc = None):
    if node_name == "archway":
        return archway_node_service.start_node_service(plan, chain_id, key, password, public_grpc, public_http, public_tcp, public_rpc)

    elif node_name == "neutron":
        return neutron_node_service.start_node_service(plan, chain_id, key, password, public_grpc, public_http, public_tcp, public_rpc)

def start_ibc_between_cosmvm_chains(plan, chain_a, chain_b):
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
