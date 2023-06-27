eth_network_module = import_module("github.com/kurtosis-tech/eth-network-package/main.star")
constants = import_module("github.com/hugobyte/dive/package_io/constants.star")

# Spins Up the ETH Node
def start_eth_node(plan,args):

    eth_contstants = constants.ETH_NODE_CLIENT
    eth_network_participants, cl_genesis_timestamp = eth_network_module.run(plan, args)
    network_address = get_network_address(eth_network_participants[0].el_client_context.ip_addr,eth_network_participants[0].el_client_context.rpc_port_num)
    return struct(
          service_name = eth_contstants.service_name,
          network_name= eth_contstants.network_name,
          network = eth_contstants.network,
          nid = eth_contstants.nid,
          endpoint = network_address,
		  endpoint_public = "",
		  keystore_path = eth_contstants.keystore_path,
		  keypassword = eth_contstants.keypassword
		  )


# Returns Network Address
def get_network_address(ip_addr,rpc_port):
     return '{0}:{1}'.format(ip_addr,rpc_port)

def start_node_service(plan,args,node_type):

	if node_type == "eth":
		return start_eth_node(plan,args)
	
	else:
		return start_hardhat_node(plan)

# Spins up Hardhat Node
def start_hardhat_node(plan):

	plan.print("Starting Hardhat Node")

	hardhat_constants = constants.HARDHAT_NODE_CLIENT


	plan.upload_files(src=hardhat_constants.config_files_path,name="hardhat-config")

	service_config = ServiceConfig(
		image=hardhat_constants.node_image,
		ports={
			hardhat_constants.port_key : PortSpec(number=hardhat_constants.port,transport_protocol="TCP",application_protocol="http")
		},
		public_ports = {
            hardhat_constants.port_key : PortSpec(number=hardhat_constants.port,transport_protocol="TCP",application_protocol="http")
        },
		files={
			hardhat_constants.config_files_directory : "hardhat-config"
		},
		entrypoint=["/bin/sh","-c","mkdir -p /app && cd app && npm install hardhat && /app/node_modules/.bin/hardhat --config ../config/hardhat.config.js node 2>&1 | tee /app/logs/hardhat.log"]
	)

	response = plan.add_service(name=hardhat_constants.service_name,config=service_config)

	private_url = get_network_address(response.ip_address,hardhat_constants.port)
	public_url = get_network_address("127.0.0.1",hardhat_constants.port)
	return struct(
          service_name = hardhat_constants.service_name,
          network_name= "hardhat",
          network = hardhat_constants.network,
          nid = hardhat_constants.network_id,
          endpoint = private_url,
		  endpoint_public = public_url,
		  keystore_path = "config/hardhat_keystore.json",
		  keypassword = "hardhat"
     )