# NOTE: If you're a VSCode user, you might like our VSCode extension: https://marketplace.visualstudio.com/items?itemName=Kurtosis.kurtosis-extension
cosmvm = import_module("github.com/hugobyte/chain-package/services/cosmvm/start_node.star")
contract = import_module("github.com/hugobyte/chain-package/services/cosmvm/deploy.star")


# For more information on...
#  - the 'run' function:  https://docs.kurtosis.com/concepts-reference/packages#runnable-packages
#  - the 'plan' object:   https://docs.kurtosis.com/starlark-reference/plan
#  - the 'args' object:   https://docs.kurtosis.com/next/concepts-reference/args
def run(plan, args):
    
    cosmvm.run(plan,args)

    # service_name = args.get("service_name","cosmos")
    plan.exec(service_name="cosmos", recipe=ExecRecipe(command=["/bin/sh","-c", "apk add jq"]))
    message = args.get("message", "message")
    contract_name=args.get("contract_name","contract")

    result = contract.deploy(plan,args,contract_name, message)
    return result


    # Try out a plan.add_service here (https://docs.kurtosis.com/starlark-reference/plan#add_service)
