def deploy(plan, chain_id, contract_name, message,service_name,password):
    contract = "../contracts/%s.wasm" % contract_name

    passcode = password

    plan.exec(service_name = service_name, recipe = ExecRecipe(command = ["/bin/sh", "-c", "echo '%s' | archwayd tx wasm store  %s --from node1-account --chain-id %s --gas auto --gas-adjustment 1.3 -y --output json -b block | jq -r '.logs[0].events[-1].attributes[0].value' > code_id.json " % (passcode, contract, chain_id)]))

    # Getting the code id

    code_id = plan.exec(service_name = service_name, recipe = ExecRecipe(command = ["/bin/sh", "-c", "cat code_id.json | tr -d '\n\r' "]))

    # instantiation

    plan.print("Instantiating the contract")

    exec = ExecRecipe(command = ["/bin/sh", "-c", "echo '%s' |  archwayd tx wasm instantiate %s '%s' --from node1-account --label %s --chain-id %s --no-admin --gas auto --gas-adjustment 1.3 -y -b block | tr -d '\n\r' " % (passcode, code_id["output"],contract_name, message, chain_id)])
    plan.exec(service_name = service_name, recipe = exec)

    # Getting the contract address

    contract = plan.exec(service_name = service_name, recipe = ExecRecipe(command = ["/bin/sh", "-c", "echo %s | archwayd query wasm list-contract-by-code %s --output json | jq -r '.contracts[-1]' | tr -d '\n\r' " % (passcode, code_id["output"])]))

    return contract["output"]
