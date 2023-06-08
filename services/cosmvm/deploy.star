wallet_config = import_module("github.com/hugobyte/chain-package/services/cosmvm/wallet.star")

def deploy(plan,service_name,artifacts_path, node_uri, init_message):
    plan.print("Deploying the contract")

    execute_cmd = ExecRecipe(command=["archwayd", "tx", "wasm", "store", artifacts_path, "--from", wallet_config, "--node", node_uri, "--chain-id", "constantine-2", "--gas-prices", "0.25aconst", "--gas", "auto", "--gas-adjustment", "1.3", "-y", "--output","json"],)
    plan.exec(service_name=service_name, recipe=execute_cmd)
    RES = plan.print(execute_cmd)

    # Getting the code id

    execute_cmd1 = ExecRecipe(command=["/bin/sh", "-c", "CODE_ID=$(echo $RES | jq -r '.logs[0].events[] | select(.type=="store_code") | .attributes[] | select(.key=="code_id") | .value')"],)
    plan.exec(service_name=service_name, recipe=execute_cmd1)
    CODE_ID = plan.print(execute_cmd1)
    
    # instantiation

    plan.print("Instantiating the contract")
    exec = ExecRecipe(command=["archwayd", "tx", "wasm", "instantiate", CODE_ID, "--from", wallet_config, "--node", node_uri, "--chain-id", "constantine-2", "--gas-prices", "0.25aconst", "--gas auto", "--gas-adjustment", "1.3", "--no-admin" ],)
    plan.exec(service_name="service_name", recipe=exec)

    # Getting the contract address

    execute = ExecRecipe(command=["archwayd", "query", "wasm", "list-contract-by-code", CODE_ID, "--node", node_uri, "--output", "json"],)
    plan.exec(service_name="service_name", recipe=execute)

    plan.print(execute)

    contract = ExecRecipe(command=[execute, "jq", "-r", ".contracts[-1]"],)
    plan.exec(service_name="service_name", recipe=contract)

    # checking the contract details

    contract_details = ExecRecipe(command=["archwayd", "query", "wasm", "contract", contract, "--node", node_uri],)
    plan.exec(service_name="service_name", recipe=contract_details)

    # checking the balances
    plan.print("The total balances is")
    balance = ExecRecipe(command=["archwayd", "query", "bank", "balances", contract, "--node", node_uri],)
    plan.exec(service_name="service_name", recipe=balance)

    # Querying the entire contract state

    query = ExecRecipe(command=["archwayd", "query", "wasm", "contract-state", "all", contract, "--node" ,node_uri ],)



    


   
        




