ICON_NODE_CLIENT = struct(
    node_image = "iconloop/goloop-icon:v1.3.9",
    config_files_directory = "/goloop/config/",
    contracts_directory = "/goloop/contracts/",
    keystore_directory = "/goloop/keystores/",
    config_files_path = "github.com/hugobyte/dive/services/jvm/icon/static-files/config/",
    contract_files_path = "github.com/hugobyte/dive/services/jvm/icon/static-files/contracts/",
    keystore_files_path = "github.com/hugobyte/dive/services/bridges/btp/static-files/keystores/keystore.json",
    rpc_endpoint_path = "api/v3/icon_dex",
    service_name = "icon-node-",
    genesis_file_path = "/goloop/genesis/",
)

HARDHAT_NODE_CLIENT = struct(
    node_image = "node:lts-alpine",
    port = 8545,
    config_files_path = "github.com/hugobyte/dive/services/evm/eth/static-files/hardhat.config.js",
    config_files_directory = "/config/",
    service_name = "hardhat-node",
    network = "0x539.hardhat",
    network_id = "0x539",
    keystore_path = "keystores/hardhat_keystore.json",
    keypassword = "hardhat",
)

CONTRACT_DEPLOYMENT_SERVICE_ETHEREUM = struct(
    node_image = "node:lts-alpine",
    static_file_path = "github.com/hugobyte/dive/services/evm/eth/static-files/",
    static_files_directory_path = "/static-files/",
    service_name = "eth-contract-deployer",
    template_file = "github.com/hugobyte/dive/services/evm/eth/static-files/hardhat.config.ts.tmpl",
    rendered_file_directory = "/static-files/rendered/",
)

ETH_NODE_CLIENT = struct(
    service_name = "el-1-geth-lighthouse",
    network_name = "eth",
    network = "0x301824.eth",
    nid = "0x301824",
    keystore_path = "keystores/eth_keystore.json",
    keypassword = "password",
)
ARCHWAY_SERVICE_CONFIG = struct(
    start_script = "github.com/hugobyte/dive/services/cosmvm/archway/static_files/start.sh",
    default_contract_path = "github.com/hugobyte/dive/services/cosmvm/archway/static_files/contracts",
    service_name = "node-service",
    image = "archwaynetwork/archwayd:v2.0.0",
    path = "/start-scripts/",
    contract_path = "/root/contracts/",
    config_files = "github.com/hugobyte/dive/services/cosmvm/archway/static_files/config/",
    password = "password",
)

NEUTRON_SERVICE_CONFIG = struct(
    service_name = "neutron-node",
    image = "hugobyte/neutron-node:v0.2",
    init_script = "github.com/hugobyte/dive/services/cosmvm/neutron/static_files/init.sh",
    start_script = "github.com/hugobyte/dive/services/cosmvm/neutron/static_files/start.sh",
    init_nutrond_script = "github.com/hugobyte/dive/services/cosmvm/neutron/static_files/init-neutrond.sh",
    path = "/start-scripts/",
)

IBC_RELAYER_SERVICE = struct(
    ibc_relay_config_file_template = "github.com/hugobyte/dive/services/bridges/ibc/static-files/config/cosmosjson.tpl",
    relay_service_name = "cosmos-ibc-relay",
    # updated the ibc relay image
    relay_service_image = "hugobyte/ibc-relay:v0.1",
    relay_config_files_path = "/script/",
    run_file_path = "github.com/hugobyte/dive/services/bridges/ibc/static-files/run.sh",
    relay_service_image_icon_to_cosmos = "hugobyte/icon-ibc-relay:v0.1",
    relay_service_name_icon_to_cosmos = "ibc-relayer",
    config_file_path = "github.com/hugobyte/dive/services/bridges/ibc/static-files/config",
    ibc_relay_wasm_file_template = "github.com/hugobyte/dive/services/bridges/ibc/static-files/config/archwayibc.json.tpl",
    ibc_relay_java_file_template = "github.com/hugobyte/dive/services/bridges/ibc/static-files/config/icon.json.tpl",
    icon_keystore_file = "github.com/hugobyte/dive/services/bridges/btp/static-files/keystores/keystore.json",
    relay_keystore_path = "/root/.relayer/keys/"
)

NETWORK_PORT_KEYS_AND_IP_ADDRESS = struct(
    grpc = "grpc",
    rpc = "rpc",
    http = "http",
    tcp = "tcp",
    public_ip_address = "127.0.0.1",
)

ARCHAY_NODE0_CONFIG = struct(
    chain_id = "constantine-3",
    grpc = 9090,
    http = 9091,
    tcp = 26658,
    rpc = 4564,
    key = "constantine-3-key",
    
)

ARCHAY_NODE1_CONFIG = struct(
    chain_id = "archway-node-1",
    grpc = 9080,
    http = 9092,
    tcp = 26659,
    rpc = 4566,
    key = "archway-node-1-key",
)

COMMON_ARCHWAY_PRIVATE_PORTS = struct(
    grpc = 9090,
    http = 9091,
    tcp = 26656,
    rpc = 26657,
)

NEUTRON_PRIVATE_PORTS = struct(
    http = 1317,
    rpc = 26657,
    tcp = 26656,
    grpc = 9090,
)

NEUTRON_NODE1_CONFIG = struct(
    http = 1317,
    rpc = 26669,
    tcp = 26656,
    grpc = 8090,
    chain_id = "test-chain1",
    key = "test-key",
    password = "clock post desk civil pottery foster expand merit dash seminar song memory figure uniform spice circle try happy obvious trash crime hybrid hood cushion",
)

NEUTRON_NODE2_CONFIG = struct(
    http = 1311,
    rpc = 26653,
    tcp = 26652,
    grpc = 8091,
    chain_id = "test-chain2",
    key = "test-key",
    password = "clock post desk civil pottery foster expand merit dash seminar song memory figure uniform spice circle try happy obvious trash crime hybrid hood cushion",
)
