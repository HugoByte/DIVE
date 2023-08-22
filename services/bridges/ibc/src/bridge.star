constants = import_module("github.com/hugobyte/dive/package_io/constants.star")
ibc_relay_config = constants.IBC_RELAYER_SERVICE

def start_cosmos_relay(plan, src_key, src_chain_id, dst_key, dst_chain_id, src_config, dst_config):
    plan.print("starting cosmos relay")

    plan.upload_files(src = ibc_relay_config.run_file_path, name = "run")
    comos_config = read_file(ibc_relay_config.ibc_relay_config_file_template)
    cfg_template_data = {
        "KEY": src_key,
        "CHAINID": src_chain_id,
    }
    plan.render_templates(
        config = {
            "cosmos-%s.json" % src_chain_id: struct(
                template = comos_config,
                data = cfg_template_data,
            ),
        },
        name = "config-%s" % src_chain_id,
    )

    cfg_template_data = {
        "KEY": dst_key,
        "CHAINID": dst_chain_id,
    }
    plan.render_templates(
        config = {
            "cosmos-%s.json" % dst_chain_id: struct(
                template = comos_config,
                data = cfg_template_data,
            ),
        },
        name = "config-%s" % dst_chain_id,
    )

    plan.exec(service_name = src_config["service_name"], recipe = ExecRecipe(command = ["/bin/sh", "-c", "apk add jq"]))

    src_chain_seed = plan.exec(service_name = src_config["service_name"], recipe = ExecRecipe(command = ["/bin/sh", "-c", "jq -r '.mnemonic' ../../start-scripts/key_seed.json | tr -d '\n\r'"]))

    plan.exec(service_name = dst_config["service_name"], recipe = ExecRecipe(command = ["/bin/sh", "-c", "apk add jq"]))

    dst_chain_seed = plan.exec(service_name = dst_config["service_name"], recipe = ExecRecipe(command = ["/bin/sh", "-c", "jq -r '.mnemonic' ../../start-scripts/key_seed.json | tr -d '\n\r'"]))

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

def start_cosmos_relay_for_icon_to_cosmos(plan,src_chain_config,dst_chain_config):

    plan.print("starting the cosmos relay for icon to cosmos")

    plan.upload_files(src=ibc_relay_config.config_file_path, name="archway_config")
    plan.upload_files(src=ibc_relay_config.icon_keystore_file,name="icon-keystore")

    wasm_config = read_file(ibc_relay_config.ibc_relay_wasm_file_template)
    java_config = read_file(ibc_relay_config.ibc_relay_java_file_template)

    cfg_template_data_wasm = {
        "KEY": dst_chain_config["key"],
        "CHAINID": dst_chain_config["chain_id"],
        "RPCADDRESS":dst_chain_config["rpc_address"],
        "IBCADDRESS":dst_chain_config["ibc_address"],
    }

    cfg_template_data_java = {
        "CHAINID": src_chain_config["chain_id"],
        "RPCADDRESS":src_chain_config["rpc_address"],
        "IBCADDRESS":src_chain_config["ibc_address"],
    }

    plan.render_templates(
        config = {
            "ibc-icon.json": struct(
                template = java_config,
                data = cfg_template_data_java,
            ),
        },
        name = "config-icon",
    )

    plan.render_templates(
        config = {
            "ibc-cosmos.json": struct(
                template = wasm_config,
                data = cfg_template_data_wasm,
            ),
        },
        name = "config-wasm",
    )



    relay_service = ServiceConfig(
        image= ibc_relay_config.relay_service_image_icon_to_cosmos,
        files= {
            ibc_relay_config.relay_config_files_path + "icon": "config-icon",
            ibc_relay_config.relay_config_files_path + "wasm": "config-wasm",
            ibc_relay_config.relay_keystore_path + src_chain_config["chain_id"] : "icon-keystore"
        },
        entrypoint=["/bin/sh"]
    )

    plan.print(relay_service)

    plan.add_service(name = ibc_relay_config.relay_service_name_icon_to_cosmos, config = relay_service)

    

    return struct(
        service_name = ibc_relay_config.relay_service_name_icon_to_cosmos,
    )

def setup_relay(plan,src_chain_config,dst_chain_config):

    src_chain_id = src_chain_config["chain_id"]
    src_password = src_chain_config["password"]

    dst_chain_id = dst_chain_config["chain_id"]
    dst_chain_key = dst_chain_config["key"]
    dst_chain_service_name = dst_chain_config["service_name"]

    plan.exec(service_name= ibc_relay_config.relay_service_name_icon_to_cosmos, recipe=ExecRecipe(command=["/bin/sh", "-c", "apk add yq"]))

    seed = plan.exec(service_name=dst_chain_service_name, recipe=ExecRecipe(command=["/bin/sh", "-c", "jq -r '.mnemonic' ../../start-scripts/key_seed.json | tr -d '\n\r'"]))
    plan.print("starting the relay")

    plan.exec(service_name= ibc_relay_config.relay_service_name_icon_to_cosmos, recipe=ExecRecipe(command=["/bin/sh", "-c", "rly config init"]))

    plan.print("Adding the chain1")

    plan.exec(service_name=ibc_relay_config.relay_service_name_icon_to_cosmos, recipe=ExecRecipe(command=["/bin/sh", "-c", "rly chains add --file ../script/wasm/ibc-cosmos.json %s" % dst_chain_id]))

    plan.print("Adding the chain2")

    plan.exec(service_name=ibc_relay_config.relay_service_name_icon_to_cosmos, recipe=ExecRecipe(command=["/bin/sh", "-c", "rly chains add --file ../script/icon/ibc-icon.json %s" % src_chain_id]))


    plan.print("Adding the keys")

    plan.exec(service_name=ibc_relay_config.relay_service_name_icon_to_cosmos, recipe=ExecRecipe(command=["/bin/sh", "-c", "rly keys restore %s %s '%s' " % (dst_chain_id,dst_chain_key,seed["output"])]))

    
    # plan.exec(service_name=ibc_relay_config.relay_service_name_icon_to_cosmos, recipe=ExecRecipe(command=["/bin/sh", "-c", "rly keys add %s keystore --password %s" % (src_chain_id,src_password)]))


    plan.print("Adding the paths")

    plan.exec(service_name=ibc_relay_config.relay_service_name_icon_to_cosmos, recipe=ExecRecipe(command=["/bin/sh", "-c", "rly paths new %s %s icon-cosmos" % (src_chain_id,dst_chain_id)]))
   
    plan.exec(service_name=ibc_relay_config.relay_service_name_icon_to_cosmos, recipe=ExecRecipe(command=["/bin/sh", "-c", "rly tx clients icon-cosmos --client-tp 2000000m -d"]))

    plan.exec(service_name=ibc_relay_config.relay_service_name_icon_to_cosmos, recipe=ExecRecipe(command=["/bin/sh", "-c", "rly tx connection icon-cosmos"]))





