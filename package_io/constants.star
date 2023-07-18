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
          service_name = "el-client-0",
          network_name= "eth",
          network = "0x301824.eth",
          nid = "0x301824",
		  keystore_path = "keystores/eth_keystore.json",
		  keypassword = "password"
)


 