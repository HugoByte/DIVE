def generate_config_data(args):
    data = get_args_data(args)
    config_data = {
        "links": data.links,
        "chains": {
            "%s" % data.src: {},
            "%s" % data.dst: {},
        },
        "contracts": {
            "%s" % data.src: {},
            "%s" % data.dst: {},
        },
        "bridge": data.bridge,
    }

    return config_data

def get_args_data(args):
    links = args["links"]
    source_chain = links["src"]
    destination_chain = links["dst"]

    if source_chain == "eth" or source_chain == "hardhat":
        if destination_chain == "icon":
            destination_chain = source_chain
            source_chain = "icon"

    if destination_chain == "cosmwasm" and source_chain == "cosmwasm":
        destination_chain = "cosmwasm1"

    bridge = args["bridge"]

    return struct(
        links = links,
        src = source_chain,
        dst = destination_chain,
        bridge = bridge,
    )

def generate_new_config_data(links, srcchain_service_name, dst_chain_service_name, bridge):
    config_data = "" 
    if bridge == "":
        config_data = {
        "links": links,
        "chains": {
            "%s" % srcchain_service_name: {},
            "%s" % dst_chain_service_name: {},
        },
        "contracts": {
            "%s" % srcchain_service_name: {},
            "%s" % dst_chain_service_name: {},
        },
        }
    else:

        config_data = {
         "links": links,
        "chains": {
            "%s" % srcchain_service_name: {},
            "%s" % dst_chain_service_name: {},
        },
        "contracts": {
            "%s" % srcchain_service_name: {},
            "%s" % dst_chain_service_name: {},
        },
        "bridge": bridge,
        }

    return config_data

def generate_new_config_data_cosmvm_cosmvm(links, srcchain_service_name, dst_chain_service_name):
    config_data = {
        "links": links,
        "chains": {
            "%s" % srcchain_service_name: {},
            "%s" % dst_chain_service_name: {},
        },
    }

    return config_data


def struct_to_dict(s):
    fields = dir(s)
    return {field: getattr(s, field) for field in fields if not field.startswith("_")}