# NOTE: If you're a VSCode user, you might like our VSCode extension: https://marketplace.visualstudio.com/items?itemName=Kurtosis.kurtosis-extension
cosmvm = import_module("github.com/hugobyte/chain-package/services/cosmvm/start_node.star")
contract = import_module("github.com/hugobyte/chain-package/services/cosmvm/deploy.star")

def run(plan, args):
    
    cosmvm.run(plan,args)

    # service_name = args.get("service_name","cosmos")
    plan.exec(service_name="cosmos", recipe=ExecRecipe(command=["/bin/sh","-c", "apk add jq"]))
    message = args.get("message", "message")
    contract_name = args.get("contract_name","contract")

    # result = contract.deploy(plan,args,contract_name, message)
    ibc_core = contract.deploy_core(plan,args)
    return ibc_core

    xcall = contract.deploy_xcall(plan,args,timeout_height, ibc_host)
    return xcall

    light_client = contract.deploy_light_client(plan,args)
    return light_client