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
     
    exec = ExecRecipe(command=["/bin/sh", "-c", "echo '%s' |  archwayd tx wasm instantiate %s '%s' --from node1-account --label xcall --chain-id my-chain --no-admin --gas auto --gas-adjustment 1.3 -y -b block" % (PASSCODE, CODE_ID["output"], message) ])
    plan.exec(service_name="cosmos", recipe=exec)

    # Getting the contract address

    CONTRACT = plan.exec(service_name="cosmos", recipe=ExecRecipe(command=["/bin/sh", "-c", "echo %s | archwayd query wasm list-contract-by-code %s --output json | jq -r '.contracts[-1]'" % (PASSCODE, CODE_ID["output"])]))
    
    plan.print(CONTRACT)
