ICON_NODE_CLIENT = struct(
    node_image = "iconloop/goloop-icon:v1.3.8",
    config_files_directory = "/goloop/config/",
    contracts_directory = "/goloop/contracts/",
    keystore_directory = "/goloop/keystores/",
    config_files_path = "github.com/hugobyte/dive/services/jvm/icon/static-files/config/",
    contract_files_path = "github.com/hugobyte/dive/services/jvm/icon/static-files/contracts/",
    keystore_files_path = "github.com/hugobyte/dive/services/bridges/btp/static-files/keystores/keystore.json",
    port_key = "rpc",
    public_ip_address = "127.0.0.1",
    rpc_endpoint_path = "api/v3/icon_dex",
    service_name = "icon-node-",
    genesis_file_path = "/goloop/genesis/"
)

HARDHAT_NODE_CLIENT = struct(
    node_image = "node:lts-alpine",
    port_key = "rpc",
    port = 8545,
    config_files_path = "github.com/hugobyte/dive/services/evm/eth/static-files/hardhat.config.js",
    config_files_directory = "/config/",
    service_name = "hardhat-node",
    network = "0x539.hardhat",
    network_id = "0x539",
    keystore_path = "keystores/hardhat_keystore.json",
	keypassword = "hardhat"
)

CONTRACT_DEPLOYMENT_SERVICE_ETHEREUM = struct(
     node_image = "node:lts-alpine",
     static_file_path = "github.com/hugobyte/dive/services/evm/eth/static-files/",
     static_files_directory_path = "/static-files/",
     service_name = "eth-contract-deployer",
     template_file = "github.com/hugobyte/dive/services/evm/eth/static-files/hardhat.config.ts.tmpl",
     rendered_file_directory = "/static-files/rendered/"
)

ETH_NODE_CLIENT = struct(
          service_name = "el-1-geth-lighthouse",
          network_name= "eth",
          network = "0x301824.eth",
          nid = "0x301824",
		  keystore_path = "keystores/eth_keystore.json",
		  keypassword = "password"
)

COSMOS_NODE_CLIENT = struct(
    start_cosmos = "github.com/hugobyte/dive/services/cosmvm/start-cosmos-0.sh",
    start_cosmos_1 = "github.com/hugobyte/dive/services/cosmvm/start-cosmos-1.sh",
    default_contract_path = "github.com/hugobyte/dive/services/cosmvm/static_files/contracts",
    service_name = "cosmos",
    service_name_1 = "cosmos1",
    image = "archwaynetwork/archwayd:constantine",
    path = "/start-scripts/",
    contract_path = "/root/contracts/",
    chain_id = "my-chain",
    chain_id_1 = "chain-1",
    public_ip_address = "127.0.0.1",
    cosmos_grpc_port_key = "grpc",
    cosmos_rpc_port_key = "rpc",
    cosmos_http_port_key = "http",
    cosmos_tcp_port_key = "tcp",
    private_port_grpc = 9090,
    private_port_http = 9091,
    private_port_tcp = 26656,
    private_port_rpc = 26657,
    public_port_grpc_node_1 = 9090,
    public_port_http_node_1 = 9091,
    public_port_tcp_node_1 = 26656,
    public_port_rpc_node_1 = 4564,
    public_port_grpc_node_2 = 9080,
    public_port_http_node_2 = 9092,
    public_port_tcp_node_2 = 26658,
    public_port_rpc_node_2 = 4566,
    config_files = "github.com/hugobyte/dive/services/cosmvm/static_files/config/",
    relay_service_name = "cosmos-relay",
    relay_service_image = "relay1",
    relay_config_files_path = "/script/",
)


 