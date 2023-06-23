def deploy(plan,args,contract_name, message):
    plan.print("Deploying the contract")

    contract = "../contracts/%s.wasm" % contract_name

    PASSCODE="password"
   
    RES = plan.exec(service_name="cosmos", recipe=ExecRecipe(command=["/bin/sh", "-c" ,"echo '%s' | archwayd tx wasm store  %s --from node1-account --chain-id my-chain --gas auto --gas-adjustment 1.3 -y --output json -b block | jq -r '.logs[0].events[-1].attributes[0].value' > code_id.json " % (PASSCODE,contract)]) )
   
    # Getting the code id

    CODE_ID = plan.exec(service_name="cosmos", recipe=ExecRecipe(command=["/bin/sh", "-c", "cat code_id.json | tr -d '\n\r' "]))
    
    plan.print(CODE_ID)
    
    # instantiation

    plan.print("Instantiating the contract")
     
    exec = ExecRecipe(command=["/bin/sh", "-c", "echo '%s' |  archwayd tx wasm instantiate %s '%s' --from node1-account --label xcall --chain-id my-chain --no-admin --gas auto --gas-adjustment 1.3 -y -b block | tr -d '\n\r' " % (PASSCODE, CODE_ID["output"], message) ])
    plan.exec(service_name="cosmos", recipe=exec)

    # Getting the contract address

    CONTRACT = plan.exec(service_name="cosmos", recipe=ExecRecipe(command=["/bin/sh", "-c", "echo %s | archwayd query wasm list-contract-by-code %s --output json | jq -r '.contracts[-1]' | tr -d '\n\r' " % (PASSCODE, CODE_ID["output"])]))
    
    return CONTRACT["output"]

def deploy_core(plan,args):

    plan.print("Deploying core")

    message = '{}'

    contract_address = deploy(plan,args, "cw_ibc_core", message)

    return contract_address

def deploy_xcall(plan,args, contract_address):

    plan.print("Deploying xcall contract")

    message = '{"timeout_height":10 , "ibc_host":"%s"}' % contract_address 
    plan.print(message)

    tx_hash1 = deploy(plan,args, "cw_xcall", message)

    return tx_hash1

def deploy_light_client(plan,args):

    plan.print("Deploying the light client")

    message = '{}' 

    tx_hash = deploy(plan,args,"cw_icon_light_client", message)

    return tx_hash
