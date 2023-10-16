def create_wallet(plan, service_name, wallet_name, wallet_password):
    """
    Create a wallet with the specified name and password.

    Args:
        plan (Plan): kurtosis plan.
        service_name (str): The name of the service.
        wallet_name (str): The name of the wallet to be created.
        wallet_password (str): The password for the wallet.

    Returns:
        str: The name of the wallet file that was created (e.g., "wallet_name.json").
    """
    plan.print("Creating Wallet")

    wallet_file = "{0}.json".format(wallet_name)

    execute_cmd = ExecRecipe(command=["./bin/goloop", "ks", "gen", "-p", wallet_password, "-o", wallet_file])
    result = plan.exec(service_name=service_name, recipe=execute_cmd)

    return wallet_file

def get_network_wallet_address(plan, service_name):
    """
    Get the address of the network wallet from the keystore configuration.

    Args:
        plan (Plan): The kurtosis plan.
        service_name (str): The name of the service.

    Returns:
        str: The network wallet address extracted from the keystore configuration.
    """
    execute_cmd = ExecRecipe(command=["jq", ".address", "config/keystore.json"])
    result = plan.exec(service_name=service_name, recipe=execute_cmd)

    execute_cmd = ExecRecipe(command=["/bin/sh", "-c", "echo \"%s\" | tr -d '\n\r'" % result["output"]])
    result = plan.exec(service_name=service_name, recipe=execute_cmd)

    return result["output"]

def get_network_wallet_public_key(plan, service_name):
    """
    Get the public key of the network wallet from the keystore configuration.

    Args:
        plan (Plan): The kurtosis plan.
        service_name (str): The name of the service.

    Returns:
        str: The public key of the network wallet extracted from the keystore configuration.
    """
    execute_cmd = ExecRecipe(command=["./bin/goloop", "ks", "pubkey", "-k", "config/keystore.json", "-p", "gochain"])
    result = plan.exec(service_name=service_name, recipe=execute_cmd)

    execute_cmd = ExecRecipe(command=["/bin/sh", "-c", "echo \"%s\" | tr -d '\n\r'" % result["output"]])
    result = plan.exec(service_name=service_name, recipe=execute_cmd)

    return result["output"]
