# Import required modules and constants
constants = import_module("../../../../package_io/constants.star")
ibc_relay_config = constants.IBC_RELAYER_SERVICE

def start_cosmos_relay(plan, src_key, src_chain_id, dst_key, dst_chain_id, src_config, dst_config, links):
    """
    Start a Cosmos relay service with given source and destination chains configuration.

    Args:
        plan (plan): plan.
        src_key (str): The key for the source chain.
        src_chain_id (str): The ID of the source chain.
        dst_key (str): The key for the destination chain.
        dst_chain_id (str): The ID of the destination chain.
        src_config (dict): Configuration for the source chain.
        dst_config (dict): Configuration for the destination chain.
        links (dict): Additional arguments.

    Returns:
        struct: Configuration information for the relay service.
    """

    plan.print("Starting Cosmos relay")

    plan.upload_files(src = ibc_relay_config.run_file_path, name = "run")

    cosmos_config = read_file(ibc_relay_config.ibc_relay_config_file_template)

    cfg_template_data = {
        "KEY": src_key,
        "CHAINID": src_chain_id,
        "CHAIN": links["src"],
    }
    plan.render_templates(
        config = {
            "cosmos-%s.json" % src_chain_id: struct(
                template = cosmos_config,
                data = cfg_template_data,
            ),
        },
        name = "config-%s" % src_chain_id,
    )

    cfg_template_data = {
        "KEY": dst_key,
        "CHAINID": dst_chain_id,
        "CHAIN": links["dst"],
    }
    plan.render_templates(
        config = {
            "cosmos-%s.json" % dst_chain_id: struct(
                template = cosmos_config,
                data = cfg_template_data,
            ),
        },
        name = "config-%s" % dst_chain_id,
    )

    # Install 'jq' based on the type of chain (neutron or archway) for the source
    if links["src"] == "neutron":
        plan.exec(service_name = src_config["service_name"], recipe = ExecRecipe(command = ["/bin/sh", "-c", "apt install jq"]))
    elif links["src"] == "archway":
        plan.exec(service_name = src_config["service_name"], recipe = ExecRecipe(command = ["/bin/sh", "-c", "apk add jq"]))

    # Retrieve the seed for the source chain
    src_chain_seed = plan.exec(service_name = src_config["service_name"], recipe = ExecRecipe(command = ["/bin/sh", "-c", "jq -r '.mnemonic' ../../start-scripts/key_seed.json | tr -d '\n\r'"]))

    # Install 'jq' based on the type of chain (neutron or archway) for the destination
    if links["dst"] == "neutron":
        plan.exec(service_name = dst_config["service_name"], recipe = ExecRecipe(command = ["/bin/sh", "-c", "apt install jq"]))
    elif links["src"] == "archway":
        plan.exec(service_name = dst_config["service_name"], recipe = ExecRecipe(command = ["/bin/sh", "-c", "apk add jq"]))

    # Retrieve the seed for the destination chain
    dst_chain_seed = plan.exec(service_name = dst_config["service_name"], recipe = ExecRecipe(command = ["/bin/sh", "-c", "jq -r '.mnemonic' ../../start-scripts/key_seed.json | tr -d '\n\r'"]))

    # Configure the Cosmos relay service
    relay_service = ServiceConfig(
        image = ibc_relay_config.relay_service_image,
        files = {
            ibc_relay_config.relay_config_files_path + src_chain_id: "config-%s" % src_chain_id,
            ibc_relay_config.relay_config_files_path + dst_chain_id: "config-%s" % dst_chain_id,
            ibc_relay_config.relay_config_files_path: "run",
        },
        entrypoint = ["/bin/sh", "-c", "chmod +x ../script/run.sh && sh ../script/run.sh '%s' '%s' '%s' '%s' '%s' '%s' '%s' '%s'" % (src_chain_id, dst_chain_id, src_key, dst_key, src_config["endpoint"], dst_config["endpoint"], src_chain_seed["output"], dst_chain_seed["output"])],
    )

    plan.print(relay_service)

    plan.add_service(name = ibc_relay_config.relay_service_name, config = relay_service)

    return struct(
        service_name = ibc_relay_config.relay_service_name,
    )

def start_cosmos_relay_for_icon_to_cosmos(plan, src_chain_config, dst_chain_config, args):
    plan.print("starting the cosmos relay for icon to cosmos")

    source_chain = args["links"]["src"]
    destination_chain = args["links"]["dst"]

    if destination_chain == "archway":
        plan.upload_files(src = ibc_relay_config.config_file_path, name = "archway_config")
    elif destination_chain == "neutron":
        plan.upload_files(src = ibc_relay_config.config_file_path, name = "neutron_config")

    plan.upload_files(src = ibc_relay_config.icon_keystore_file, name = "icon-keystore")


    if destination_chain == "archway":
        wasm_config = read_file(ibc_relay_config.ibc_relay_wasm_file_template)
    elif destination_chain == "neutron":
        wasm_config = read_file(ibc_relay_config.ibc_relay_neutron_wasm_file_template)

    java_config = read_file(ibc_relay_config.ibc_relay_java_file_template)

    cfg_template_data_wasm = {
        "KEY": dst_chain_config["key"],
        "CHAINID": dst_chain_config["chain_id"],
        "RPCADDRESS": dst_chain_config["rpc_address"],
        "IBCADDRESS": dst_chain_config["ibc_address"],
    }

    cfg_template_data_java = {
        "CHAINID": src_chain_config["chain_id"],
        "RPCADDRESS": src_chain_config["rpc_address"],
        "IBCADDRESS": src_chain_config["ibc_address"],
    }

    plan.render_templates(
        config = {
            "%s.json" % src_chain_config["chain_id"]: struct(
                template = java_config,
                data = cfg_template_data_java,
            ),
        },
        name = "config-%s" % src_chain_config["chain_id"],
    )

    plan.render_templates(
        config = {
            "%s.json" % dst_chain_config["chain_id"]: struct(
                template = wasm_config,
                data = cfg_template_data_wasm,
            ),
        },
        name = "config-%s" % dst_chain_config["chain_id"],
    )

    relay_service = ServiceConfig(
        image = ibc_relay_config.relay_service_image_icon_to_cosmos,
        files = {
            ibc_relay_config.relay_config_files_path + "java": "config-%s" % src_chain_config["chain_id"],
            ibc_relay_config.relay_config_files_path + "wasm": "config-%s" % dst_chain_config["chain_id"],
            ibc_relay_config.relay_keystore_path + src_chain_config["chain_id"]: "icon-keystore",
        },
        entrypoint = ["/bin/sh"],
    )

    plan.add_service(name = ibc_relay_config.relay_service_name_icon_to_cosmos, config = relay_service)

    return struct(
        service_name = ibc_relay_config.relay_service_name_icon_to_cosmos,
    )

def setup_relay(plan, src_chain_config, dst_chain_config):
    src_chain_id = src_chain_config["chain_id"]
    src_password = src_chain_config["password"]

    dst_chain_id = dst_chain_config["chain_id"]
    dst_chain_key = dst_chain_config["key"]
    dst_chain_service_name = dst_chain_config["service_name"]

    plan.exec(service_name = ibc_relay_config.relay_service_name_icon_to_cosmos, recipe = ExecRecipe(command = ["/bin/sh", "-c", "apk add yq"]))

    seed = plan.exec(service_name = dst_chain_service_name, recipe = ExecRecipe(command = ["/bin/sh", "-c", "jq -r '.mnemonic' ../../start-scripts/key_seed.json | tr -d '\n\r'"]))
    plan.print("starting the relay")

    plan.exec(service_name = ibc_relay_config.relay_service_name_icon_to_cosmos, recipe = ExecRecipe(command = ["/bin/sh", "-c", "rly config init"]))

    plan.print("Adding the chain1")

    plan.exec(service_name = ibc_relay_config.relay_service_name_icon_to_cosmos, recipe = ExecRecipe(command = ["/bin/sh", "-c", "rly chains add --file ../script/wasm/%s.json %s" % (dst_chain_id, dst_chain_id)]))

    plan.print("Adding the chain2")

    plan.exec(service_name = ibc_relay_config.relay_service_name_icon_to_cosmos, recipe = ExecRecipe(command = ["/bin/sh", "-c", "rly chains add --file ../script/java/%s.json %s" % (src_chain_id, src_chain_id)]))

    plan.print("Adding the keys")

    plan.exec(service_name = ibc_relay_config.relay_service_name_icon_to_cosmos, recipe = ExecRecipe(command = ["/bin/sh", "-c", "rly keys restore %s %s '%s' " % (dst_chain_id, dst_chain_key, seed["output"])]))

    plan.print("Adding the paths")

    path_name = "{0}-{1}".format(src_chain_id, dst_chain_id)

    plan.exec(service_name = ibc_relay_config.relay_service_name_icon_to_cosmos, recipe = ExecRecipe(command = ["/bin/sh", "-c", "rly paths new %s %s %s" % (src_chain_id, dst_chain_id, path_name)]))

    plan.exec(service_name = ibc_relay_config.relay_service_name_icon_to_cosmos, recipe = ExecRecipe(command = ["/bin/sh", "-c", "rly tx clients %s --client-tp 2000000m -d" % path_name]))

    plan.exec(service_name = ibc_relay_config.relay_service_name_icon_to_cosmos, recipe = ExecRecipe(command = ["/bin/sh", "-c", "rly tx connection %s" % path_name]))

    return path_name

def get_relay_path_data(plan, service_name, path_name):
    src_chain_id_response = plan.exec(service_name = service_name, recipe = ExecRecipe(command = ["/bin/sh", "-c", "cat /root/.relayer/config/config.yaml | yq '.paths.%s.src.chain-id'" % path_name]))
    execute_cmd = ExecRecipe(command = ["/bin/sh", "-c", "echo \"%s\" | tr -d '\n\r'" % src_chain_id_response["output"]])
    src_chain_id = plan.exec(service_name = service_name, recipe = execute_cmd)

    src_client_id_response = plan.exec(service_name = service_name, recipe = ExecRecipe(command = ["/bin/sh", "-c", "cat /root/.relayer/config/config.yaml | yq '.paths.%s.src.client-id'" % path_name]))
    execute_cmd = ExecRecipe(command = ["/bin/sh", "-c", "echo \"%s\" | tr -d '\n\r'" % src_client_id_response["output"]])
    src_client_id = plan.exec(service_name = service_name, recipe = execute_cmd)

    src_connection_id_response = plan.exec(service_name = service_name, recipe = ExecRecipe(command = ["/bin/sh", "-c", "cat /root/.relayer/config/config.yaml | yq '.paths.%s.src.connection-id'" % path_name]))
    execute_cmd = ExecRecipe(command = ["/bin/sh", "-c", "echo \"%s\" | tr -d '\n\r'" % src_connection_id_response["output"]])
    src_connection_id = plan.exec(service_name = service_name, recipe = execute_cmd)

    dst_chain_id_response = plan.exec(service_name = service_name, recipe = ExecRecipe(command = ["/bin/sh", "-c", "cat /root/.relayer/config/config.yaml | yq '.paths.%s.dst.chain-id'" % path_name]))
    execute_cmd = ExecRecipe(command = ["/bin/sh", "-c", "echo \"%s\" | tr -d '\n\r'" % dst_chain_id_response["output"]])
    dst_chain_id = plan.exec(service_name = service_name, recipe = execute_cmd)

    dst_client_id_response = plan.exec(service_name = service_name, recipe = ExecRecipe(command = ["/bin/sh", "-c", "cat /root/.relayer/config/config.yaml | yq '.paths.%s.dst.client-id'" % path_name]))
    execute_cmd = ExecRecipe(command = ["/bin/sh", "-c", "echo \"%s\" | tr -d '\n\r'" % dst_client_id_response["output"]])
    dst_client_id = plan.exec(service_name = service_name, recipe = execute_cmd)

    dst_connection_id_response = plan.exec(service_name = service_name, recipe = ExecRecipe(command = ["/bin/sh", "-c", "cat /root/.relayer/config/config.yaml | yq '.paths.%s.dst.connection-id'" % path_name]))
    execute_cmd = ExecRecipe(command = ["/bin/sh", "-c", "echo \"%s\" | tr -d '\n\r'" % dst_connection_id_response["output"]])
    dst_connection_id = plan.exec(service_name = service_name, recipe = execute_cmd)

    config = struct(
        src_chain_id = src_chain_id["output"],
        src_client_id = src_client_id["output"],
        src_connection_id = src_connection_id["output"],
        dst_chain_id = dst_chain_id["output"],
        dst_client_id = dst_client_id["output"],
        dst_connection_id = dst_connection_id["output"],
    )

    return config

def start_channel(plan, service_name, path_name, src_port, dst_port):
    plan.print("Starting Channel")

    exec_cmd = ["/bin/sh", "-c", "rly tx chan %s --src-port=%s --dst-port=%s" % (path_name, src_port, dst_port)]

    plan.exec(service_name = service_name, recipe = ExecRecipe(command = exec_cmd))

def start_relay(plan, service_name):
    plan.print("Starting Relay")
    plan.exec(service_name = service_name, recipe = ExecRecipe(command = ["/bin/sh", "-c", "ln -sf /proc/1/fd/1 /root/.relayer/relay.log && rly start 1>&2"]))
