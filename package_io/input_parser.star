def generate_config_data(args):

    data = get_args_data(args)
    config_data = {
        "links": data.links,
        "chains" : {
            "%s" % data.src : {},
            "%s" % data.dst : {}
        },
        "contracts" : {
            "%s" %  data.src : {},
            "%s" %  data.dst : {}
        },
        "bridge" : data.bridge
    }

    return config_data

def get_args_data(args):

    links = args["links"]
    source_chain = links["src"]
    destination_chain = links["dst"]

    if destination_chain == "icon":
        destination_chain = "icon-1"
    
    if source_chain == "eth" or source_chain == "hardhat":
        if destination_chain == "icon":
            destination_chain = source_chain
            source_chain = "icon" 

    bridge = args["bridge"]

    return struct(
        links = links,
        src = source_chain,
        dst = destination_chain,
        bridge = bridge
    )