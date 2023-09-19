constants = import_module("github.com/hugobyte/dive/package_io/constants.star")
participant_network = import_module("github.com/kurtosis-tech/eth-network-package/src/participant_network.star")
input_parser = import_module("github.com/kurtosis-tech/eth-network-package/package_io/input_parser.star")
static_files = import_module("github.com/kurtosis-tech/eth-network-package/static_files/static_files.star")
genesis_constants = import_module("github.com/kurtosis-tech/eth-network-package/src/prelaunch_data_generator/genesis_constants/genesis_constants.star")

network_keys_and_public_address = constants.NETWORK_PORT_KEYS_AND_IP_ADDRESS

# Spins Up the ETH Node
def start_eth_node(plan,args):
 	eth_contstants = constants.ETH_NODE_CLIENT
 	args_with_right_defaults = input_parser.get_args_with_default_values(args)
 	num_participants = len(args_with_right_defaults.participants)
 	network_params = args_with_right_defaults.network_params

 	all_participants, cl_genesis_timestamp, _ = participant_network.launch_participant_network(plan, args_with_right_defaults.participants, network_params, args_with_right_defaults.global_client_log_level)
	
 	network_address = get_network_address(all_participants[0].el_client_context.ip_addr,all_participants[0].el_client_context.rpc_port_num)

	return struct(
          service_name = all_participants[0].el_client_context.service_name,
          network_name= eth_contstants.network_name,
          network = eth_contstants.network,
          nid = eth_contstants.nid,
          endpoint = "http://%s" % network_address,
		  endpoint_public = "http://",
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
			network_keys_and_public_address.rpc : PortSpec(number=hardhat_constants.port,transport_protocol=network_keys_and_public_address.tcp.upper(),application_protocol=network_keys_and_public_address.http)
		},
		public_ports = {
            network_keys_and_public_address.rpc : PortSpec(number=hardhat_constants.port,transport_protocol=network_keys_and_public_address.tcp.upper(),application_protocol=network_keys_and_public_address.http)
        },
		files={
			hardhat_constants.config_files_directory : "hardhat-config"
		},
		entrypoint=["/bin/sh","-c","mkdir -p /app && cd app && npm install hardhat && /app/node_modules/.bin/hardhat --config ../config/hardhat.config.js node 2>&1 | tee /app/logs/hardhat.log"]
	)

	response = plan.add_service(name=hardhat_constants.service_name,config=service_config)

	private_url = get_network_address(response.ip_address,hardhat_constants.port)
	public_url = get_network_address(network_keys_and_public_address.public_ip_address,hardhat_constants.port)
	return struct(
          service_name = hardhat_constants.service_name,
          network_name= "hardhat",
          network = hardhat_constants.network,
          nid = hardhat_constants.network_id,
          endpoint = "http://%s" % private_url,
		  endpoint_public = "http://%s" % public_url,
		  keystore_path = hardhat_constants.keystore_path,
		  keypassword = hardhat_constants.keypassword
     )
