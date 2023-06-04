CONTRACT_DEPLOYER_SERVICE = "contract_deployer_service"
CONTRACT_DEPLOYER_SERVICE_IMAGE = "golang"
DEPLOYER_SERVICE_PRIVATE_PORT = 80
DEPLOYER_SERVICE_PORT_KEY = "deployer_rpc"
# CONTRACT_DIR_PATH = "contracts"

def start_contract_deploy_service(plan):

    plan.print("Starting "+CONTRACT_DEPLOYER_SERVICE)
    

    service_config = ServiceConfig(
        image=CONTRACT_DEPLOYER_SERVICE_IMAGE,
        
        
    )

    plan.add_service(name=CONTRACT_DEPLOYER_SERVICE,config=service_config)


