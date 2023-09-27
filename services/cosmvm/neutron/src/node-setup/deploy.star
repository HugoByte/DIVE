def deploy(plan, chain_id, chain_key, contract_name, message, service_name):
    """
    Deploy a contract on a Neutron node.

    Args:
        plan (plan): The execution plan.
        chain_id (str): The chain ID.
        chain_key (str): The chain key.
        contract_name (str): The name of the contract to deploy.
        message (str): The message to pass during contract deployment.
        service_name (str): The name of the Neutron node service.
        
    Returns:
        str: The contract address.
    """

    # Define the contract file path
    contract = "../../root/contracts/%s.wasm" % contract_name

    # Execute a command to store the contract on the chain and retrieve the code ID
    plan.exec(service_name = service_name, recipe = ExecRecipe(command = ["/bin/sh", "-c", "neutrond tx wasm store  %s --from  %s --home ./data/%s --keyring-backend test --chain-id %s --gas auto --gas-adjustment 1.3 -y --output json -b block | jq -r '.logs[0].events[-1].attributes[1].value' > code_id.json " % (contract, chain_key, chain_id, chain_id)]))

    # Getting the code id
    code_id = plan.exec(service_name = service_name, recipe = ExecRecipe(command = ["/bin/sh", "-c", "cat code_id.json | tr -d '\n\r' "]))

    # Instantiate the contract
    plan.print("Instantiating the contract")
    exec = ExecRecipe(command = ["/bin/sh", "-c", "neutrond tx wasm instantiate %s '%s' --from %s --home ./data/%s --keyring-backend test --label %s --chain-id %s --no-admin --gas auto --gas-adjustment 1.3 -y -b block | tr -d '\n\r' " % (code_id["output"], message, chain_key, chain_id, contract_name, chain_id)])
    plan.exec(service_name = service_name, recipe = exec)

    # Getting the contract address
    contract = plan.exec(service_name = service_name, recipe = ExecRecipe(command = ["/bin/sh", "-c", "neutrond query wasm list-contract-by-code %s --output json | jq -r '.contracts[-1]' | tr -d '\n\r' " % (code_id["output"])]))

    return contract["output"]
