eth_contract_deployer_service = import_module("github.com/hugobyte/chain-package/services/evm/eth/src/node-setup/contract-deployer.star")


def deploy_bmc(plan,args):

    eth_config_data = args["chains"]["eth"]

    network = eth_config_data["network"]
    network_name = eth_config_data["network_name"]

    plan.print("Deploying BMC Contract on %s" % network)

    eth_contract_deployer_service.deploy_contract(plan,"bmc",'{\"link\":\"%s\",\"chainNetwork\":\"%s\"}' % (network_name,network),"localnet")

    bmc_address = eth_contract_deployer_service.get_contract_address(plan,"bmc")

    bmcm_address = eth_contract_deployer_service.get_contract_address(plan,"bmcm")

    bmcs_address = eth_contract_deployer_service.get_contract_address(plan,"bmcs")

    return struct(
        bmcm = bmcm_address,
        bmcs = bmcs_address,
        bmc = bmc_address
    )

def deploy_xcall(plan,args):

    eth_config_data = args["chains"]["eth"]
    network = eth_config_data["network"]
    network_name = eth_config_data["network_name"]

    plan.print("Deploying xCall Contract on %s" % network)

    eth_contract_deployer_service.deploy_contract(plan,"xcall",'{"name":"%s"}' % network_name,"localnet")

    xcall_address = eth_contract_deployer_service.get_contract_address(plan,"xcall")

    return xcall_address


def deploy_dapp(plan,args):

    eth_config_data = args["chains"]["eth"]

    network = eth_config_data["network"]
    network_name = eth_config_data["network_name"]

    plan.print("Deploying dapp Contract on %s" % network)

    eth_contract_deployer_service.deploy_contract(plan,"dapp",'{"name":"%s"}' % network_name,"localnet")

    dapp_address = eth_contract_deployer_service.get_contract_address(plan,"dapp")

    return dapp_address


def deploy_bmv_bridge(plan,args,lastblock_height,src_bmc_address,srcchain_network):

    eth_config_data = args["chains"]["eth"]

    network = eth_config_data["network"]
    network_name = eth_config_data["network_name"]

    plan.print("Deploying Bmv-Bridge Contract on %s" % network)

    params = '{"current_chain":{"name":"%s"},"src":{"lastBlockHeight":"%s","bmc":"%s","network":"%s"}}' % (network_name,lastblock_height,src_bmc_address,srcchain_network)

    eth_contract_deployer_service.deploy_contract(plan,"bmv_bridge",params,"localnet")

    bmvb = eth_contract_deployer_service.get_contract_address(plan,"bmvb")

    return bmvb


def deploy_bmv(plan,args,src_first_block_header,src_bmc_address,srcchain_network,srcchain_network_type_id):

    eth_config_data = args["chains"]["eth"]

    network = eth_config_data["network"]
    network_name = eth_config_data["network_name"]

    plan.print("Deploying Bmv Contract on %s" % network)

    params = '{"current_chain":{"name":"%s"},"src":{"firstBlockHeader":"%s","bmc":"%s","network":"%s","networkTypeId":"%s"}}' % (network_name,src_first_block_header,src_bmc_address,srcchain_network,srcchain_network_type_id)

    eth_contract_deployer_service.deploy_contract(plan,"bmv",params,"localnet")

    bmv = eth_contract_deployer_service.get_contract_address(plan,"bmv")

    return bmv






