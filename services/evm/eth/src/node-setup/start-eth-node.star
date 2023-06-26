eth_network_module = import_module("github.com/kurtosis-tech/eth-network-package/main.star")
NODE_SERVICE_IMAGE = "node:lts-alpine"
SERVICE_RPC_PORT_KEY = "hardhar-rpc"
SERVICE_RPC_PRIVATE_PORT = 8545
SERVICE_RPC_PUBLIC_PORT = 8554
HARDHAT_CONFIG_PATH = "github.com/hugobyte/dive/services/evm/eth/static-files/hardhat.config.js"
HARDHAT_CONFIG_DIR = "/config/"
HARDHAT_SERVICE_NAME = "hardhat-client"
HARDHAT_NETWORK = "0x539.hardhat"
HARDHAT_NID = "0x539"

# Spins Up the ETH Node
def start_eth_node(plan,args):

     eth_network_participants, cl_genesis_timestamp = eth_network_module.run(plan, args)
     network_address = get_network_address(eth_network_participants[0].el_client_context.ip_addr,eth_network_participants[0].el_client_context.rpc_port_num)

     return struct(
          service_name = "el-client-0",
          network_name= "eth",
          network = "0x301824.eth",
          nid = "0x301824",
          endpoint = network_address,
		  endpoint_public = "",
		  keystore_path = "config/eth_keystore.json",
		  keypassword = "password"

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

	plan.print("Starting Hardhat")

	plan.upload_files(src=HARDHAT_CONFIG_PATH,name="hardhat-config")

	service_config = ServiceConfig(
		image=NODE_SERVICE_IMAGE,
		ports={
			SERVICE_RPC_PORT_KEY : PortSpec(number=SERVICE_RPC_PRIVATE_PORT,transport_protocol="TCP",application_protocol="http")
		},
		public_ports = {
            SERVICE_RPC_PORT_KEY : PortSpec(number=SERVICE_RPC_PUBLIC_PORT,transport_protocol="TCP",application_protocol="http")
        },
		files={
			HARDHAT_CONFIG_DIR : "hardhat-config"
		},
		entrypoint=["/bin/sh","-c","mkdir -p /app && cd app && npm install hardhat && /app/node_modules/.bin/hardhat --config ../config/hardhat.config.js node 2>&1 | tee /app/logs/hardhat.log"]
	)

	response = plan.add_service(name=HARDHAT_SERVICE_NAME,config=service_config)

	private_url = get_network_address(response.ip_address,SERVICE_RPC_PRIVATE_PORT)
	public_url = get_network_address("127.0.0.1",SERVICE_RPC_PUBLIC_PORT)
	return struct(
          service_name = HARDHAT_SERVICE_NAME,
          network_name= "hardhat",
          network = HARDHAT_NETWORK,
          nid = HARDHAT_NID,
          endpoint = private_url,
		  endpoint_public = public_url,
		  keystore_path = "config/hardhat_keystore.json",
		  keypassword = "hardhat"
     )