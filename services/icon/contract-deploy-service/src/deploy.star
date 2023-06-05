CONTRACT_DEPLOYER_SERVICE = "contract_deployer_service"
CONTRACT_DEPLOYER_SERVICE_IMAGE = "alpine"
DEPLOYER_SERVICE_PRIVATE_PORT = 80
DEPLOYER_SERVICE_PORT_KEY = "deployer_rpc"
CONTRACT_DIR_PATH = "/contracts"
DEFAULT_CONTRACT_DIR_PATH = "github.com/hugobyte/chain-package/services/icon/contract-deploy-service/static-files/"


def start_contract_deploy_service(plan,args):

    plan.print("Starting "+CONTRACT_DEPLOYER_SERVICE)

    contract_file_path = args["contracts_path"]

    plan.upload_files(src=contract_file_path,name="contracts")


    service_config = ServiceConfig(
        image=CONTRACT_DEPLOYER_SERVICE_IMAGE,  
        files={
            CONTRACT_DIR_PATH : "contracts"

        }      
    )

    plan.add_service(name=CONTRACT_DEPLOYER_SERVICE,config=service_config)

    install_dependecy_command = ["/bin/sh","-c","apk add jq go python3 git make"]

    plan.exec(service_name=CONTRACT_DEPLOYER_SERVICE,recipe=ExecRecipe(command=install_dependecy_command))

    # Build goloop binary for contract deployment and execution 

    goloop_build_command = ["/bin/sh","-c","git clone https://github.com/icon-project/goloop.git && cd goloop && GOBUILD_TAGS= make goloop && mv bin/goloop /usr/bin"]

    plan.exec(service_name=CONTRACT_DEPLOYER_SERVICE,recipe=ExecRecipe(command=goloop_build_command))

    wallet_info = create_wallet(plan,CONTRACT_DEPLOYER_SERVICE,"contract_name","password")


def create_wallet(plan,service_name,wallet_name,wallet_password):
    plan.print("Creating Wallet")

    wallet_file = "{0}.json".format(wallet_name)
   
    execute_cmd = ExecRecipe(command=["goloop","ks","gen","-p",wallet_password,"-o",wallet_file])
    result = plan.exec(service_name=service_name,recipe=execute_cmd)

    return wallet_file


def deploy_contract(plan,service_name, contract_name,args):

    plan.print("Deploying %s Contract" % contract_name)




    wallet_info = create_wallet(plan,"contract_name","password")








