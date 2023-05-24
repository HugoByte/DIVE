icon_node_launcher = import_module("github.com/hugobyte/chain-package/services/icon/icon.star")
def run(plan, args):
   
    plan.print("Starting Deployment Tool")

    if args["chain"] == "ICON":
        ip = icon_node_launcher.launch_icon_node(plan,args)
        plan.print(ip)
        response = plan.exec(service_name="icon",recipe=ExecRecipe(command=["../bin/goloop","rpc","lastblock","--uri","http://"+ip+"/api/v3"]),)
        plan.print(response)
    else:
        plan.print("Not Configured")

