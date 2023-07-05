<<<<<<< HEAD
ETH_CONTRACT_DEPLOYMENT_SERVICE_IMAGE = "node:lts-alpine"
CONTRACT_DEPLOYMENT_STATIC_FILES= "github.com/hugobyte/dive/services/evm/eth/static-files/"
CONTRACT_DEPLOYMENT_STATIC_FILES_DIR_PATH = "/static-files/"
SERVICE_NAME = "eth-contract-deployer"
HARDHAT_CONFIG = "github.com/hugobyte/dive/services/evm/eth/static-files/hardhat.config.ts.tmpl"
HARDHAT_CONFIG_DIR = "/static-files/rendered/"
=======
constants = import_module("github.com/hugobyte/dive/package_io/constants.star")
>>>>>>> main

# Starts the eth deploy service
def start_deploy_service(plan,args):

    deployer_constants = constants.CONTRACT_DEPLOYMENT_SERVICE_ETHEREUM

    plan.print("Starting Contract Deploy Service")
    endpoint = args["endpoint"]

    plan.upload_files(src=deployer_constants.static_file_path,name="static-files")

    hardhat_config = read_file(deployer_constants.template_file)
    cfg_template_data = {
        "URL": endpoint
    }
    plan.render_templates(
        config= {
            "hardhat.config.ts": struct(
                template = hardhat_config,
                data = cfg_template_data,
            ),
        },
        name="config"
    )


    service_config = ServiceConfig(
        image=deployer_constants.node_image,
        files={
            deployer_constants.static_files_directory_path : "static-files",
            deployer_constants.rendered_file_directory : "config"


        },
        entrypoint = ["/bin/sh","-c","mv /static-files/rendered/hardhat.config.ts /static-files/ &&  apk add jq &&  sleep 9999999999"]

    )

    service_response = plan.add_service(name=deployer_constants.service_name,config=service_config)
    plan.exec(service_name=deployer_constants.service_name,recipe=ExecRecipe(command=["/bin/sh","-c","cd static-files && npm install && npm install hardhat && npx hardhat compile"]))

    return service_response

# Returns Latest block 
def get_latest_block(plan,current_chain,network_name):

    file_name = "get_block_number.ts"
    params = '{"current_chain":"%s"}' % current_chain


    exec_command = ["/bin/sh","-c","cd static-files &&  params='{0}' npx hardhat --network {1} run scripts/{2}".format(params,network_name,file_name)]
    plan.exec(service_name=constants.CONTRACT_DEPLOYMENT_SERVICE_ETHEREUM.service_name,recipe=ExecRecipe(exec_command))

    exec_command = ["/bin/sh","-c","cd static-files && cat deployments.json | jq -r .%s.blockNum | tr -d '\n\r'" % current_chain]
    response = plan.exec(service_name=constants.CONTRACT_DEPLOYMENT_SERVICE_ETHEREUM.service_name,recipe=ExecRecipe(exec_command))

    return response["output"]
