def deploy(plan, service_name, contract_path, wallet_name, node_url):
    exec_recipe = ExecRecipe(command=[
        contract_path
        "--from",
        wallet_name
        "--chain_id constantine-2",
        "--node ",
        node_url
        "--fees 3397uconst",
        "--gas auto",
        "--output json"
    ],)
deploy = plan.exec(
    service_name=service_name,
    recipe=exec_recipe
)

