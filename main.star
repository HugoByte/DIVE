icon = import_module("github.com/hugobyte/chain-package/services/jvm/icon/main.star")


def run(plan,args):

    plan.print("Starting Kurtosis Package")
    
    icon.icon_service(plan,args)
   
        

