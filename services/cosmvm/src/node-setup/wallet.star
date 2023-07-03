
PASSWORD = "password"

def add_balance(plan,args):
    exec = ExecRecipe(command=["/bin/sh","-c", "archwayd keys show fd -a --keyring-backend=test | tr -d '\n\r'" ])
    address = plan.exec(service_name="cosmos", recipe=exec)
    fd =address["output"]

    exec = ExecRecipe(command=["/bin/sh", "-c", "echo '%s' | archwayd keys show test-account -a | tr -d '\n\r'" % (PASSWORD)])
    test_address = plan.exec(service_name="cosmos", recipe=exec)
    test_address = test_address["output"]

    exec_command = ExecRecipe(command=["/bin/sh", "-c", "archwayd tx bank send '%s' '%s' 10000000stake --keyring-backend test --chain-id my-chain -y" % (fd, test_address) ])
    result = plan.exec(service_name="cosmos", recipe=exec_command)
    plan.print(result)



