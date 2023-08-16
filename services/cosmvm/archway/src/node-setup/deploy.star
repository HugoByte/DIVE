constants = import_module("github.com/hugobyte/dive/package_io/constants.star")
cosmos_node_constants = constants.COSMOS_NODE_CLIENT

def deploy(plan, chain_id, contract_name, message):
    contract = "../contracts/%s.wasm" % contract_name

    passcode = cosmos_node_constants.password

    plan.exec(service_name = "cosmos", recipe = ExecRecipe(command = ["/bin/sh", "-c", "echo '%s' | archwayd tx wasm store  %s --from node1-account --chain-id %s --gas auto --gas-adjustment 1.3 -y --output json -b block | jq -r '.logs[0].events[-1].attributes[0].value' > code_id.json " % (passcode, contract, chain_id)]))

    # Getting the code id

    code_id = plan.exec(service_name = "cosmos", recipe = ExecRecipe(command = ["/bin/sh", "-c", "cat code_id.json | tr -d '\n\r' "]))

    # instantiation

    plan.print("Instantiating the contract")

    exec = ExecRecipe(command = ["/bin/sh", "-c", "echo '%s' |  archwayd tx wasm instantiate %s '%s' --from node1-account --label xcall --chain-id %s --no-admin --gas auto --gas-adjustment 1.3 -y -b block | tr -d '\n\r' " % (passcode, code_id["output"], message, chain_id)])
    plan.exec(service_name = "cosmos", recipe = exec)

    # Getting the contract address

    contract = plan.exec(service_name = "cosmos", recipe = ExecRecipe(command = ["/bin/sh", "-c", "echo %s | archwayd query wasm list-contract-by-code %s --output json | jq -r '.contracts[-1]' | tr -d '\n\r' " % (passcode, code_id["output"])]))

    return contract["output"]
