eth_network_module = import_module("github.com/kurtosis-tech/eth-network-package/main.star")
ETH_DEPLOY_SERVICE_IMAGE = "node:lts-alpine"

# Spins Up the ETH Node
def start_eth_node(plan,args):

     eth_network_participants, cl_genesis_timestamp = eth_network_module.run(plan, args)
     network_address = get_network_address(eth_network_participants[0].el_client_context.ip_addr,eth_network_participants[0].el_client_context.rpc_port_num)

     return struct(
          service_name = "el-client-0",
          network_name= "eth",
          network = "0x301824.eth",
          nid = "0x301824",
          endpoint = network_address
     )

# Returns Network Address
def get_network_address(ip_addr,rpc_port):
     return '{0}:{1}'.format(ip_addr,rpc_port)