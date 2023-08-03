cosmvm_node = import_module("github.com/hugobyte/dive/services/cosmvm/src/node-setup/start_node.star")

def deploy(plan,args,contract_name, message):

    contract = "../contracts/%s.wasm" % contract_name

    passcode="password"
   
    res = plan.exec(service_name="cosmos", recipe=ExecRecipe(command=["/bin/sh", "-c" ,"echo '%s' | archwayd tx wasm store  %s --from node1-account --chain-id my-chain --gas auto --gas-adjustment 1.3 -y --output json -b block | jq -r '.logs[0].events[-1].attributes[0].value' > code_id.json " % (passcode,contract)]) )
   
    # Getting the code id

    code_id = plan.exec(service_name="cosmos", recipe=ExecRecipe(command=["/bin/sh", "-c", "cat code_id.json | tr -d '\n\r' "]))
    
    plan.print(code_id)
    
    # instantiation

    plan.print("Instantiating the contract")
     
    exec = ExecRecipe(command=["/bin/sh", "-c", "echo '%s' |  archwayd tx wasm instantiate %s '%s' --from node1-account --label xcall --chain-id my-chain --no-admin --gas auto --gas-adjustment 1.3 -y -b block | tr -d '\n\r' " % (passcode, code_id["output"], message) ])
    plan.exec(service_name="cosmos", recipe=exec)

    # Getting the contract address

    contract = plan.exec(service_name="cosmos", recipe=ExecRecipe(command=["/bin/sh", "-c", "echo %s | archwayd query wasm list-contract-by-code %s --output json | jq -r '.contracts[-1]' | tr -d '\n\r' " % (passcode, code_id["output"])]))
    
    return contract["output"]


