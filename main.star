
cosmvm = import_module("github.com/hugobyte/chain-package/services/cosmvm/start_node.star")
contract = import_module("github.com/hugobyte/chain-package/services/cosmvm/deploy.star")

def run(plan, args):
    
    cosmvm.run(plan,args)

    plan.exec(service_name="cosmos", recipe=ExecRecipe(command=["/bin/sh","-c", "apk add jq"]))
    message = args.get("message", "message")
    contract_name = args.get("contract_name","contract")

    ibc_core = contract.deploy_core(plan,args)
    plan.print(ibc_core)
    
    xcall = contract.deploy_xcall(plan,args, ibc_core)
   
    light_client = contract.deploy_light_client(plan,args)
   
