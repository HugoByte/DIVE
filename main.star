icon = import_module("github.com/hugobyte/chain-package/services/icon/icon.star")


def run(plan,args):

    plan.print("Starting Kurtosis Package")
    
    icon.icon(plan,args)
   
        

