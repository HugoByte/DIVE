archway_node_service = import_module("github.com/hugobyte/dive/services/cosmvm/archway/archway.star")
neutron_node_service = import_module("github.com/hugobyte/dive/services/cosmvm/neutron/neutron.star")
parser = import_module("github.com/hugobyte/dive/package_io/input_parser.star")


def start_cosmvm_chains(plan,node_name,args):
    if node_name == "archway":
        return archway_node_service.start_node_service(plan,args)

    elif node_name == "neutron":
        return neutron_node_service.start_node_service(plan,args)
        

def start_ibc_between_cosmvm_chains(plan, chain_a, chain_b, args):
    if chain_a == "archway" and chain_b == "archway":
        return archway_node_service.start_nodes_services_archway(plan)

    elif chain_a == "neutron" and chain_b == "neutron":
        return neutron_node_service.start_node_services(plan, args)
    
    elif chain_a == "neutron" and chain_b == "archway":
        chain_a_service = parser.struct_to_dict(neutron_node_service.start_node_service(plan, args["src_chain"]))
        chain_b_service = parser.struct_to_dict(archway_node_service.start_node_service(plan, args["dst_chain"]))
        return struct(
            src_config = chain_a_service,
            dst_config = chain_b_service
        )

    elif chain_a == "archway" and chain_b == "neutron":
        chain_a_service = parser.struct_to_dict(archway_node_service.start_node_service(plan, args["src_chain"]))
        chain_b_service = parser.struct_to_dict(neutron_node_service.start_node_service(plan, args["dst_chain"]))
        return struct(
            src_config = chain_a_service,
            dst_config = chain_b_service
        )
        
