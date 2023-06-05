node_service = import_module("github.com/hugobyte/chain-package/services/jvm/icon/src/main.star")

def icon(plan,args):

    if args["service_method"] == "deploy_node":

      response =  node_service.node_service(plan,args)

      plan.print("Private_URL " + response.private_url)

      plan.print("Public_URL  "+response.public_url)

      

    elif args["service_method"] == "deploy_contract":

        response =  node_service.contract_deployer(plan,args)

        plan.print("ScoreAddress %s" % response)




