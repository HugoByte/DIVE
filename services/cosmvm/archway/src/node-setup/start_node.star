constants = import_module("github.com/hugobyte/dive/package_io/constants.star")
cosmos_node_constants = constants.ARCHWAY_SERVICE_CONFIG
network_port_keys_and_ip = constants.NETWORK_PORT_KEYS_AND_IP_ADDRESS

def start_cosmos_node(plan, args):
    chain_id = args["cid"]
    key = args["key"]
    password = args["password"]
    service_name = args["service_name"]

    plan.print("Launching " + service_name + " deployment service")

    start_script_file = "start-script-%s" % chain_id
    plan.upload_files(src = cosmos_node_constants.start_script, name = start_script_file)

    cosmwasm_node_config = ServiceConfig(
        image = cosmos_node_constants.image,
        files = {
            cosmos_node_constants.path: start_script_file,
        },
        ports = {
           network_port_keys_and_ip.grpc: PortSpec(number = args["private_grpc"], transport_protocol =network_port_keys_and_ip.tcp.upper(), application_protocol =network_port_keys_and_ip.http),
           network_port_keys_and_ip.http: PortSpec(number = args["private_http"], transport_protocol =network_port_keys_and_ip.tcp.upper(), application_protocol =network_port_keys_and_ip.http),
           network_port_keys_and_ip.tcp: PortSpec(number = args["private_tcp"], transport_protocol =network_port_keys_and_ip.tcp.upper(), application_protocol =network_port_keys_and_ip.http),
           network_port_keys_and_ip.rpc: PortSpec(number = args["private_rpc"], transport_protocol =network_port_keys_and_ip.tcp.upper(), application_protocol =network_port_keys_and_ip.http),
        },
        public_ports = {
           network_port_keys_and_ip.grpc: PortSpec(number = args["public_grpc"], transport_protocol =network_port_keys_and_ip.tcp.upper(), application_protocol =network_port_keys_and_ip.http),
           network_port_keys_and_ip.http: PortSpec(number = args["public_http"], transport_protocol =network_port_keys_and_ip.tcp.upper(), application_protocol =network_port_keys_and_ip.http),
           network_port_keys_and_ip.tcp: PortSpec(number = args["public_tcp"], transport_protocol =network_port_keys_and_ip.tcp.upper(), application_protocol =network_port_keys_and_ip.http),
           network_port_keys_and_ip.rpc: PortSpec(number = args["public_rpc"], transport_protocol =network_port_keys_and_ip.tcp.upper(), application_protocol =network_port_keys_and_ip.http),
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
        chain_key = key
    )

# returns url
def get_service_url(ip_address, ports):
    port_id = ports["rpc"].number
    protocol = ports["rpc"].application_protocol
    url = "{0}://{1}:{2}".format(protocol, ip_address, port_id)
    return url

def get_service_config(cid, key, private_grpc, private_http, private_tcp, private_rpc, public_grpc, public_http, public_tcp, public_rpc, password):
    service_name = "{0}-{1}".format(cosmos_node_constants.service_name,cid)
    config = {
        "public_grpc" : public_grpc,
        "public_http" : public_http,
        "public_tcp" : public_tcp,
        "public_rpc" : public_rpc,
        "private_http" : private_http,
        "private_tcp" : private_tcp,
        "private_rpc" : private_rpc,
        "private_grpc" : private_grpc,
        "service_name" : service_name,
        "cid" : cid,
        "key" : key,
        "password" : password,
    } 

    return config
