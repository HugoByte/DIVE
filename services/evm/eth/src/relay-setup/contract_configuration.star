eth_contract_deployer_service = import_module("../node-setup/contract-deployer.star")

def deploy_bmc(plan, chain_name, network, network_name):
    """
    Deploy BMC on the ETH network.

    Args:
        plan (Plan):  plan.
        chain_name (str): The name of the blockchain network.
        network (str): The network identifier.
        network_name (str): The name of the network.

    Returns:
        struct: A struct containing addresses of deployed BMC-related contracts:
            - bmcm: Address of the BMC Manager contract.
            - bmcs: Address of the BMC Storage contract.
            - bmc: Address of the BMC Contract.
    """
    plan.print("Deploying BMC Contract on %s" % network)

    eth_contract_deployer_service.deploy_contract(plan, "bmc", '{"link":"%s","chainNetwork":"%s"}' % (network_name, network), "localnet")

    bmc_address = eth_contract_deployer_service.get_contract_address(plan, "bmc", chain_name)

    bmcm_address = eth_contract_deployer_service.get_contract_address(plan, "bmcm", chain_name)

    bmcs_address = eth_contract_deployer_service.get_contract_address(plan, "bmcs", chain_name)

    return struct(
        bmcm=bmcm_address,
        bmcs=bmcs_address,
        bmc=bmc_address
    )


def deploy_xcall(plan, chain_name, network, network_name):
    """
    Deploy an xCall Contract on the ETH network.

    Args:
        plan (Plan):  plan.
        chain_name (str): The name of the blockchain network.
        network (str): The network identifier.
        network_name (str): The name of the xCall contract.

    Returns:
        str: The address of the deployed xCall Contract.
    """
    plan.print("Deploying xCall Contract on %s" % network)

    eth_contract_deployer_service.deploy_contract(plan, "xcall", '{"name":"%s"}' % network_name, "localnet")

    xcall_address = eth_contract_deployer_service.get_contract_address(plan, "xcall", chain_name)

    return xcall_address


def deploy_dapp(plan, chain_name, network, network_name):
    """
    Deploy a Dapp Contract on the ETH network.

    Args:
        plan (Plan): The deployment plan.
        chain_name (str): The name of the blockchain network.
        network (str): The network identifier.
        network_name (str): The name of the Dapp contract.

    Returns:
        str: The address of the deployed Dapp Contract.
    """
    plan.print("Deploying Dapp Contract on %s" % network)

    eth_contract_deployer_service.deploy_contract(plan, "dapp", '{"name":"%s"}' % network_name, "localnet")

    dapp_address = eth_contract_deployer_service.get_contract_address(plan, "dapp", chain_name)

    return dapp_address



def deploy_bmv_bridge(plan, lastblock_height, src_bmc_address, src_chain_network, chain_name, network, network_name):
    """
    Deploy a BmvBridge Contract on the ETH network and return its address.

    Args:
        plan (Plan): The deployment plan.
        lastblock_height (str): The last block height on the source chain.
        src_bmc_address (str): The address of the source BMC (Blockchain Management Contract).
        src_chain_network (str): The network of the source chain.
        chain_name (str): The name of the blockchain network.
        network (str): The network identifier.
        network_name (str): The name of the network.

    Returns:
        str: The address of the deployed BmvBridge Contract.
    """
    plan.print("Deploying Bmv-Bridge Contract on %s" % network)

    params = '{"current_chain":{"name":"%s"},"src":{"lastBlockHeight":"%s","bmc":"%s","network":"%s"}}' % (network_name, lastblock_height, src_bmc_address, src_chain_network)

    eth_contract_deployer_service.deploy_contract(plan, "bmv_bridge", params, "localnet")

    bmvb = eth_contract_deployer_service.get_contract_address(plan, "bmvb", chain_name)

    return bmvb


def deploy_bmv(plan, src_first_block_header, src_bmc_address, src_chain_network, src_chain_network_type_id, chain_name, network, network_name):
    """
    Deploy a Bmv Contract on the specified network and return its address.

    Args:
        plan (Plan): plan.
        src_first_block_header (str): The first block header on the source chain.
        src_bmc_address (str): The address of the source BMC (Blockchain Management Contract).
        src_chain_network (str): The network of the source chain.
        src_chain_network_type_id (str): The network type ID of the source chain.
        chain_name (str): The name of the blockchain network.
        network (str): The network identifier.
        network_name (str): The name of the network.

    Returns:
        str: The address of the deployed Bmv Contract.
    """
    plan.print("Deploying Bmv Contract on %s" % network)

    params = '{"current_chain":{"name":"%s"},"src":{"firstBlockHeader":"%s","bmc":"%s","network":"%s","networkTypeId":"%s"}}' % (network_name, src_first_block_header, src_bmc_address, src_chain_network, src_chain_network_type_id)

    eth_contract_deployer_service.deploy_contract(plan, "bmv", params, "localnet")

    bmv = eth_contract_deployer_service.get_contract_address(plan, "bmv", chain_name)

    return bmv



