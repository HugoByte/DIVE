"""
Creates Wallet with Given name and password
'wallet_name' - Naming the genereated keystore
'wallet_password' - Pasword used for sigining and Decrypting the Generated Keystore
"""
def create_wallet(plan,service_name,wallet_name,wallet_password):
    plan.print("Creating Wallet")

    wallet_file = "{0}.json".format(wallet_name)
   
    execute_cmd = ExecRecipe(command=["./bin/goloop","ks","gen","-p",wallet_password,"-o",wallet_file])
    result = plan.exec(service_name=service_name,recipe=execute_cmd)

    return wallet_file


# Returns Network Wallet Address
def get_network_wallet_address(plan,service_name):

    execute_cmd = ExecRecipe(command=["jq",".address","config/keystore.json"])
    result = plan.exec(service_name=service_name,recipe=execute_cmd)

    execute_cmd = ExecRecipe(command=["/bin/sh", "-c","echo \"%s\" | tr -d '\n\r'" % result["output"] ])
    result = plan.exec(service_name=service_name,recipe=execute_cmd)

    return result["output"]

# Returns Network Wallet Public Key
def get_network_wallet_public_key(plan,service_name):
    execute_cmd = ExecRecipe(command=["./bin/goloop","ks","pubkey","-k","config/keystore.json","-p","gochain"])
    result = plan.exec(service_name=service_name,recipe=execute_cmd)

    execute_cmd = ExecRecipe(command=["/bin/sh", "-c","echo \"%s\" | tr -d '\n\r'" % result["output"] ])
    result = plan.exec(service_name=service_name,recipe=execute_cmd)

    return result["output"]
