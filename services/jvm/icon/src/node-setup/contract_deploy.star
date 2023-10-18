DEFAULT_STEP_LIMIT = "500000000000"


def deploy_contract(plan, contract_name, init_message, service_name, uri, keystore_path, keystore_password, nid):
    """
    Deploy a smart contract to the icon.

    Args:
        plan (Plan): The kurtosis  plan.
        contract_name (str): The name of the contract.
        init_message (str): The initialization message for the contract.
        service_name (str): The name of the service to execute the deployment.
        uri (str): The URI for connecting to the Goloop blockchain.
        keystore_path (str): The path to the keystore for signing the deployment.
        keystore_password (str): The password for the keystore.
        nid (str): The network ID of the Goloop blockchain.

    Returns:
        str: The output of the deployment.

    Note:
        This function deploys a smart contract to the Goloop blockchain by sending a transaction.

    """
    contract = contract_name + ".jar"

    execute_command = [
        "./bin/goloop",
        "rpc",
        "sendtx",
        "deploy",
        "contracts/" + contract,
        "--content_type",
        "application/java",
        "--params",
        init_message,
        "--key_store",
        keystore_path,
        "--key_password",
        keystore_password,
        "--step_limit",
        DEFAULT_STEP_LIMIT,
        "--uri",
        uri,
        "--nid",
        nid,
    ]

    plan.print(execute_command)
    result = plan.exec(service_name=service_name, recipe=ExecRecipe(command=execute_command))

    return result["output"]


def get_score_address(plan, service_name, tx_hash):
    """
    Get the score address for a given transaction hash.

    Args:
        plan (Plan): The kurtosis plan.
        service_name (str): The name of the service to execute the request.
        tx_hash (str): The transaction hash for which to retrieve the score address.

    Returns:
        str: The score address extracted from the transaction result.

    Note:
        This function sends a POST request to the Icon blockchain to retrieve the score address
        associated with a given transaction hash.

    """
    post_request = PostHttpRequestRecipe(
        port_id="rpc",
        endpoint="/api/v3/icon_dex",
        content_type="application/json",
        body='{"jsonrpc": "2.0", "method": "icx_getTransactionResult", "id": 1, "params": {"txHash": %s }}' % tx_hash,
        extract={
            "score_address": ".result.scoreAddress"
        }
    )
   
    result = plan.wait(service_name=service_name, recipe=post_request, field="code", assertion="==", target_value=200)

    return result["extract.score_address"]
