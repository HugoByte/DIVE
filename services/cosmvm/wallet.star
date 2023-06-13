def wallet(plan,service_name, wallet_name):
    plan.print("creating the wallet")

    wallet_file = "{0}.json".format(wallet_name)

    execute_cmd = ExecRecipe(command=["archway", "accounts", "--add", "wallet-name", wallet_name],)
    plan.exec(service_name="service_name",recipe=execute_cmd)

    return wallet_file