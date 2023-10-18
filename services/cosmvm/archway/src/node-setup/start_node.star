# Import necessary modules
constants = import_module("../../../../../package_io/constants.star")
cosmos_node_constants = constants.ARCHWAY_SERVICE_CONFIG
network_port_keys_and_ip = constants.NETWORK_PORT_KEYS_AND_IP_ADDRESS

def start_cosmos_node(plan, chain_id, key, password, service_name, private_grpc, private_http, private_tcp, private_rpc, public_grpc, public_http, public_tcp, public_rpc):
    """
    Configure and start a Cosmos node for Archway.

    Args:
        plan (plan): Plan object for service deployment.
        chain_id (str): Chain ID.
        key (str): Key.
        password (str): Password.
        service_name (str): Name of the service.
        private_grpc (int): Private gRPC port.
        private_http (int): Private HTTP port.
        private_tcp (int): Private TCP port.
        private_rpc (int): Private RPC port.
        public_grpc (int): Public gRPC port.
        public_http (int): Public HTTP port.
        public_tcp (int): Public TCP port.
        public_rpc (int): Public RPC port.

    Returns:
        struct: Configuration information for the service.
    """
    plan.print("Launching " + service_name + " deployment service")

    start_script_file = "start-script-%s" % chain_id
    contract_files = "contract-%s" % chain_id
    plan.upload_files(src = cosmos_node_constants.start_script, name = start_script_file)
    plan.upload_files(src = cosmos_node_constants.default_contract_path, name = contract_files)

    cosmwasm_node_config = ServiceConfig(
        image = cosmos_node_constants.image,
        files = {
            cosmos_node_constants.path: start_script_file,
            cosmos_node_constants.contract_path: contract_files,
        },
        ports = {
            network_port_keys_and_ip.grpc: PortSpec(number = private_grpc, transport_protocol = network_port_keys_and_ip.tcp.upper(), application_protocol = network_port_keys_and_ip.http),
            network_port_keys_and_ip.http: PortSpec(number = private_http, transport_protocol = network_port_keys_and_ip.tcp.upper(), application_protocol = network_port_keys_and_ip.http),
            network_port_keys_and_ip.tcp: PortSpec(number = private_tcp, transport_protocol = network_port_keys_and_ip.tcp.upper(), application_protocol = network_port_keys_and_ip.http),
            network_port_keys_and_ip.rpc: PortSpec(number = private_rpc, transport_protocol = network_port_keys_and_ip.tcp.upper(), application_protocol = network_port_keys_and_ip.http),
        },
        public_ports = {
            network_port_keys_and_ip.grpc: PortSpec(number = public_grpc, transport_protocol = network_port_keys_and_ip.tcp.upper(), application_protocol = network_port_keys_and_ip.http),
            network_port_keys_and_ip.http: PortSpec(number = public_http, transport_protocol = network_port_keys_and_ip.tcp.upper(), application_protocol = network_port_keys_and_ip.http),
            network_port_keys_and_ip.tcp: PortSpec(number = public_tcp, transport_protocol = network_port_keys_and_ip.tcp.upper(), application_protocol = network_port_keys_and_ip.http),
            network_port_keys_and_ip.rpc: PortSpec(number = public_rpc, transport_protocol = network_port_keys_and_ip.tcp.upper(), application_protocol = network_port_keys_and_ip.http),
        },
        entrypoint = ["/bin/sh", "-c", "cd ../..%s && chmod +x start.sh && ./start.sh %s %s %s" % (cosmos_node_constants.path, chain_id, key, password)],
    )

    node_service_response = plan.add_service(name = service_name, config = cosmwasm_node_config)

    plan.print(node_service_response)

    public_url = get_service_url(network_port_keys_and_ip.public_ip_address, cosmwasm_node_config.public_ports)
    private_url = get_service_url(node_service_response.ip_address, node_service_response.ports)

    return struct(
        service_name = service_name,
        endpoint = private_url,
        endpoint_public = public_url,
        chain_id = chain_id,
        chain_key = key,
    )

# returns url
def get_service_url(ip_address, ports):
    port_id = ports["rpc"].number
    protocol = ports["rpc"].application_protocol
    url = "{0}://{1}:{2}".format(protocol, ip_address, port_id)
    return url
