eth_contract_deployer_service = import_module("../node-setup/contract-deployer.star")

# Deploy Bmc contract on ETH and Returns it's address
def deploy_bmc(plan, network, network_name, chain_name):

    plan.print("Deploying BMC Contract on %s" % network)

    eth_contract_deployer_service.deploy_contract(plan,"bmc",'{\"link\":\"%s\",\"chainNetwork\":\"%s\"}' % (network_name,network),"localnet")

    bmc_address = eth_contract_deployer_service.get_contract_address(plan,"bmc",chain_name)

    bmcm_address = eth_contract_deployer_service.get_contract_address(plan,"bmcm",chain_name)

    bmcs_address = eth_contract_deployer_service.get_contract_address(plan,"bmcs",chain_name)

    return struct(
        bmcm = bmcm_address,
        bmcs = bmcs_address,
        bmc = bmc_address
    )

# Deploy xCall Contract and returns it's address
def deploy_xcall(plan,network, network_name,chain_name,service_name):

    plan.print("Deploying xCall Contract on %s" % network)

    eth_contract_deployer_service.deploy_contract(plan,"xcall",'{"name":"%s"}' % network_name,"localnet")

    xcall_address = eth_contract_deployer_service.get_contract_address(plan,"xcall",chain_name)

    return xcall_address

# Deploy dapp Contract and returns it's address
def deploy_dapp(plan,network, network_name,chain_name,service_name):

    plan.print("Deploying dapp Contract on %s" % network)

    eth_contract_deployer_service.deploy_contract(plan,"dapp",'{"name":"%s"}' % network_name,"localnet")

    dapp_address = eth_contract_deployer_service.get_contract_address(plan,"dapp",chain_name)

    return dapp_address

# Deploy BmvBridge Contract and returns it's address
def deploy_bmv_bridge(plan,network, network_name ,lastblock_height,src_bmc_address,srcchain_network,chain_name,service_name):


    plan.print("Deploying Bmv-Bridge Contract on %s" % network)

    params = '{"current_chain":{"name":"%s"},"src":{"lastBlockHeight":"%s","bmc":"%s","network":"%s"}}' % (network_name,lastblock_height,src_bmc_address,srcchain_network)

    eth_contract_deployer_service.deploy_contract(plan,"bmv_bridge",params,"localnet")

    bmvb = eth_contract_deployer_service.get_contract_address(plan,"bmvb",chain_name)

    return bmvb

# Deploy Bmv contract and returns it's address
def deploy_bmv(plan,network, network_name, src_first_block_header,src_bmc_address,srcchain_network,srcchain_network_type_id,chain_name):

    plan.print("Deploying Bmv Contract on %s" % network)

    params = '{"current_chain":{"name":"%s"},"src":{"firstBlockHeader":"%s","bmc":"%s","network":"%s","networkTypeId":"%s"}}' % (network_name,src_first_block_header,src_bmc_address,srcchain_network,srcchain_network_type_id)

    eth_contract_deployer_service.deploy_contract(plan,"bmv",params,"localnet")

    bmv = eth_contract_deployer_service.get_contract_address(plan,"bmv",chain_name)

    return bmv






