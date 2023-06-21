ETH_CONTRACT_DEPLOYMENT_SERVICE_IMAGE = "node:lts-alpine"
CONTRACT_DEPLOYMENT_STATIC_FILES= "github.com/hugobyte/chain-package/services/evm/eth/static-files/"
CONTRACT_DEPLOYEMNT_STATIC_FILES_DIR_PATH = "/static-files/"
SERVICE_NAME = "eth-contract-deployer"
HARDHAT_CONFIG = "github.com/hugobyte/chain-package/services/evm/eth/static-files/hardhat.config.ts.tmpl"
HARDHAT_CONFIG_DIR = "/static-files/rendered/"

def start_deploy_service(plan,args):

    plan.print("Starting Contract Deploy Service")

    endpoint = args["endpoint"]

    plan.upload_files(src=CONTRACT_DEPLOYMENT_STATIC_FILES,name="static-files")

    hardhat_config = read_file(HARDHAT_CONFIG)
    cfg_template_data = {
        "URL": endpoint
    }
    node_cfg = plan.render_templates(
        config= {
            "hardhat.config.ts": struct(
                template = hardhat_config,
                data = cfg_template_data,
            ),
        },
        name="config"
    )


    service_config = ServiceConfig(
        image=ETH_CONTRACT_DEPLOYMENT_SERVICE_IMAGE,
        files={
            CONTRACT_DEPLOYEMNT_STATIC_FILES_DIR_PATH : "static-files",
            HARDHAT_CONFIG_DIR : "config"


        },
        entrypoint= ["/bin/sh","-c","mv /static-files/rendered/hardhat.config.ts /static-files/ &&  apk add jq &&  sleep 9999999999"]

    )

    service_response = plan.add_service(name=SERVICE_NAME,config=service_config)

    plan.exec(service_name=SERVICE_NAME,recipe=ExecRecipe(command=["/bin/sh","-c","cd static-files && npm install && npm install hardhat && npx hardhat compile"]))

    return service_response

def get_latest_block(plan,current_chain,network_name):

    file_name = "get_block_number.ts"

    params = '{"current_chain":"%s"}' % current_chain

    

    exec_command = ["/bin/sh","-c","cd static-files &&  params='{0}' npx hardhat --network {1} run scripts/{2}".format(params,network_name,file_name)]

    plan.exec(service_name="eth-contract-deployer",recipe=ExecRecipe(exec_command))


    exec_command = ["/bin/sh","-c","cd static-files && cat deployments.json | jq -r .eth.blockNum | tr -d '\n\r'"]

    response = plan.exec(service_name="eth-contract-deployer",recipe=ExecRecipe(exec_command))

    return response["output"]
