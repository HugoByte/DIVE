constants = import_module("github.com/hugobyte/dive/package_io/constants.star")

def start_cosmos_node(plan,args):

    cosmos_node_constants = constants.COSMOS_NODE_CLIENT

    plan.print("Launching " +cosmos_node_constants.service_name+  " deployment service")

    plan.upload_files(src=cosmos_node_constants.start_cosmos, name="start-script")
    plan.upload_files(src=cosmos_node_constants.default_contract_path, name="contract")

    cosmwasm_node_config = ServiceConfig(
        image=cosmos_node_constants.image,
        files={ 
            cosmos_node_constants.path: "start-script",
            cosmos_node_constants.contract_path: "contract",
        },
        ports={
            cosmos_node_constants.cosmos_grpc_port_key : PortSpec(number=cosmos_node_constants.private_port_grpc,transport_protocol="TCP",application_protocol="http"),
            cosmos_node_constants.cosmos_http_port_key : PortSpec(number=cosmos_node_constants.private_port_http,transport_protocol="TCP",application_protocol="http"),
            cosmos_node_constants.cosmos_tcp_port_key : PortSpec(number=cosmos_node_constants.private_port_tcp,transport_protocol="TCP",application_protocol="http"),
            cosmos_node_constants.cosmos_rpc_port_key : PortSpec(number=cosmos_node_constants.private_port_rpc,transport_protocol="TCP",application_protocol="http"),
        },
        public_ports={
            cosmos_node_constants.cosmos_grpc_port_key: PortSpec(number=cosmos_node_constants.public_port_grpc_node_1,transport_protocol="TCP",application_protocol="http"),
            cosmos_node_constants.cosmos_http_port_key : PortSpec(number=cosmos_node_constants.public_port_http_node_1,transport_protocol="TCP",application_protocol="http"),
            cosmos_node_constants.cosmos_tcp_port_key : PortSpec(number=cosmos_node_constants.public_port_tcp_node_1,transport_protocol="TCP",application_protocol="http"),
            cosmos_node_constants.cosmos_rpc_port_key : PortSpec(number=cosmos_node_constants.public_port_rpc_node_1,transport_protocol="TCP",application_protocol="http"),
           
        },
        
        entrypoint=["/bin/sh","-c","cd ../../start-scripts && chmod +x start-cosmos-0.sh && ./start-cosmos-0.sh"]
    )

    node_service_response = plan.add_service(name=cosmos_node_constants.service_name, config= cosmwasm_node_config)

    plan.print(node_service_response)

    public_url = get_service_url(cosmos_node_constants.public_ip_address,cosmwasm_node_config.public_ports)
    private_url = get_service_url(node_service_response.ip_address,node_service_response.ports)

    return struct(
        service_name = cosmos_node_constants.service_name,
        endpoint = private_url,
        endpoint_public = public_url
    )

# returns url
def get_service_url(ip_address,ports):
    port_id = ports["rpc"].number
    url = "{0}:{1}".format(ip_address,port_id)
    return url

def get_service_config(service_name, cid):

    return struct(
        service_name = service_name,
        cid = cid
    )

def start_cosmos_node_1(plan,args):

    cosmos_node_constant = constants.COSMOS_NODE_CLIENT

    plan.upload_files(src=cosmos_node_constant.start_cosmos_1, name="start-script1")

    cosmwasm_node_config_1 = ServiceConfig(
        image=constants.COSMOS_NODE_CLIENT.image,
        files={ 
            constants.COSMOS_NODE_CLIENT.path: "start-script1",
            constants.COSMOS_NODE_CLIENT.contract_path: "contract",
        },
        ports={
            cosmos_node_constant.cosmos_grpc_port_key : PortSpec(number=cosmos_node_constant.private_port_grpc,transport_protocol="TCP",application_protocol="http"),
            cosmos_node_constant.cosmos_http_port_key : PortSpec(number=cosmos_node_constant.private_port_http,transport_protocol="TCP",application_protocol="http"),
            cosmos_node_constant.cosmos_tcp_port_key : PortSpec(number=cosmos_node_constant.private_port_tcp,transport_protocol="TCP",application_protocol="http"),
            cosmos_node_constant.cosmos_rpc_port_key : PortSpec(number=cosmos_node_constant.private_port_rpc,transport_protocol="TCP",application_protocol="http"),
        },
        public_ports={
            cosmos_node_constant.cosmos_grpc_port_key : PortSpec(number=cosmos_node_constant.public_port_grpc_node_2,transport_protocol="TCP",application_protocol="http"),
            cosmos_node_constant.cosmos_http_port_key : PortSpec(number=cosmos_node_constant.public_port_http_node_2,transport_protocol="TCP",application_protocol="http"),
            cosmos_node_constant.cosmos_tcp_port_key  : PortSpec(number=cosmos_node_constant.public_port_tcp_node_2,transport_protocol="TCP",application_protocol="http"),
            cosmos_node_constant.cosmos_rpc_port_key  : PortSpec(number=cosmos_node_constant.public_port_rpc_node_2,transport_protocol="TCP",application_protocol="http"),
        },
        
        entrypoint=["/bin/sh","-c","cd ../../start-scripts && chmod +x start-cosmos-1.sh && ./start-cosmos-1.sh "]
    )

    node_service_response = plan.add_service(name=cosmos_node_constant.service_name_1, config= cosmwasm_node_config_1)

    public_url = get_service_url_1(cosmos_node_constant.public_ip_address,cosmwasm_node_config_1.public_ports)
    private_url = get_service_url_1(node_service_response.ip_address,node_service_response.ports)

    return struct(
        service_name = cosmos_node_constant.service_name_1,
        endpoint = private_url,
        endpoint_public = public_url
    )

# returns url
def get_service_url_1(ip_address,ports):
    port_id = ports["rpc"].number
    url = "{0}:{1}".format(ip_address,port_id)
    return url

def get_service_config_1(service_name, cid):

    return struct(
        service_name = service_name,
        cid = cid
    )
