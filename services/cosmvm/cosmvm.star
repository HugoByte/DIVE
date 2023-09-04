archway_node_service = import_module("github.com/hugobyte/dive/services/cosmvm/archway/archway.star")
neutron_node_service = import_module("github.com/hugobyte/dive/services/cosmvm/neutron/neutron.star")

def start_cosmvm_chains(plan,node_name,args):
    if node_name == "archway":
        return archway_node_service.start_node_service(plan,args)

    elif node_name == "neutron":
        return neutron_node_service.start_node_service(plan,args)
        

def start_ibc_between_cosmvm_chains(plan,chain_a,chain_b):
    if chain_a == "archway" and chain_b == "archway":
        return archway_node_service.start_nodes_services_archway(plan)