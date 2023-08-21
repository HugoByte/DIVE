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

def start_cosmos_relay_for_icon_to_cosmos(plan,args):

    plan.print("starting the cosmos relay for icon to cosmos")

    plan.upload_files(src=ibc_relay_config.config_file_path, name="archway_config")

    relay_service = ServiceConfig(
        image= ibc_relay_config.relay_service_image_icon_to_cosmos,
        files= {
            ibc_relay_config.relay_config_files_path: "archway_config"
        },
        entrypoint=["/bin/sh"]
    )

    plan.print(relay_service)

    plan.add_service(name = ibc_relay_config.relay_service_name_icon_to_cosmos, config = relay_service)

    return struct(
        service_name = ibc_relay_config.relay_service_name_icon_to_cosmos,
    )

def setup_relay(plan,args,SEED0):

    plan.print("starting the relay")

    plan.exec(service_name= ibc_relay_config.relay_service_name_icon_to_cosmos, recipe=ExecRecipe(command=["/bin/sh", "-c", "rly config init"]))

    plan.print("Adding the chain1")

    plan.exec(service_name=ibc_relay_config.relay_service_name_icon_to_cosmos, recipe=ExecRecipe(command=["/bin/sh", "-c", "rly chains add --file ../script/archway1.json my-chain"]))

    plan.print("Adding the chain2")

    plan.exec(service_name=ibc_relay_config.relay_service_name_icon_to_cosmos, recipe=ExecRecipe(command=["/bin/sh", "-c", "rly chains add --file ../script/icon.json 0xacbc4e "]))


    plan.print("Adding the keys")

    plan.exec(service_name=ibc_relay_config.relay_service_name_icon_to_cosmos, recipe=ExecRecipe(command=["/bin/sh", "-c", "rly keys restore my-chain fd '%s' " % (SEED0["output"])]))

    # plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "rly keys restore chain-2 default '%s' " % (SEED1["output"])]))
    # plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "rly keys add 0xacbc4e keystore --password gochain"]))

    # plan.exec(service_name="cosmos-relay", recipe=ExecRecipe(command=["/bin/sh", "-c", "rly keys add 0xacbc4e keystore --password password"]))

    plan.print("Adding the paths")

    plan.exec(service_name=ibc_relay_config.relay_service_name_icon_to_cosmos, recipe=ExecRecipe(command=["/bin/sh", "-c", "rly paths new 0xacbc4e my-chain demo"]))
   
    plan.exec(service_name=ibc_relay_config.relay_service_name_icon_to_cosmos, recipe=ExecRecipe(command=["/bin/sh", "-c", "rly tx clients demo --client-tp 10000000m "]))

    plan.exec(service_name=ibc_relay_config.relay_service_name_icon_to_cosmos, recipe=ExecRecipe(command=["/bin/sh", "-c", "rly tx connection demo"]))

    # plan.exec(service_name="cosmos-relay",recipe=ExecRecipe(command=["/bin/sh", "-c", "rly start &"]))

    # plan.exec(service_name="cosmos-relay",recipe=ExecRecipe(command=["/bin/sh", "-c", "rly tx channel demo --dst-port our-port --order ordered"]))





