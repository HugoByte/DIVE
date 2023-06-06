contract_deployment_service = import_module("github.com/hugobyte/chain-package/services/jvm/icon/src/contract_deploy.star")

def deploy_bmc(plan,args):
    plan.print("Deploying BMC Contract")

     contract_deployment_service.contract_deployer(plan,args)


def deploy_xcall(plan,args):

    plan.print("Deploying xCall Contract")


def deploy_bmv(plan,args):
    plan.print("Deploying BMV contract")