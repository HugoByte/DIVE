node_service = import_module("github.com/hugobyte/chain-package/services/icon/node-service/src/main.star")
contract_deployment_service = import_module("github.com/hugobyte/chain-package/services/icon/contract-deploy-service/src/deploy.star")

def icon(plan,args):

    if args["deploy"] == "node":

        node_service.main(plan,args)

    elif args["deploy"] == "contract":
        plan.print("deploying Contract")

        contract_deployment_service.start_contract_deploy_service(plan,args)
