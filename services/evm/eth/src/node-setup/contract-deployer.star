# Deploy Contract to Eth Network
def deploy_contract(plan,contract_name,params,network_name):

    plan.print("Deploying Contract {0}".format(contract_name))


    contract_deployment_file_name = "deploy_{0}.ts".format(contract_name)

    exec_command = ["/bin/sh","-c","cd static-files &&  params='{0}' npx hardhat --network {1} run scripts/{2}".format(params,network_name,contract_deployment_file_name)]

    response = plan.exec(service_name="eth-contract-deployer",recipe=ExecRecipe(exec_command))

    plan.print(response)

# Returns Address of Deployed Contract
def get_contract_address(plan,contract_name,chain_type):

    exec_command = ["/bin/sh","-c","cd static-files && cat deployments.json | jq -r .%s.contracts.%s | tr -d '\n\r'" % (chain_type,contract_name)]

    response = plan.exec(service_name="eth-contract-deployer",recipe=ExecRecipe(exec_command))

    return response["output"]



