RELAY_SERVICE_IMAGE = 'hugobyte/btp-relay'
RELAY_SERVICE_NAME = "btp-relay"
RELAY_CONFIG_FILES_PATH = "/relay/config/"

# Starts BTP relayer
def start_relayer(plan,src_chain,dst_chain,args,src_btp_address,dst_btp_address,bridge):

    plan.print("Starting Relay Service")

    src_config = args["chains"][src_chain]
    src_service_name = src_config["service_name"]
    src_endpoint = src_config["endpoint"]
    src_keystore = src_config["keystore_path"]
    src_keypassword =src_config["keypassword"]

    dst_config = args["chains"][dst_chain]
    dst_service_name = dst_config["service_name"]
    dst_endpoint = dst_config["endpoint"]
    dst_keystore = dst_config["keystore_path"]
    dst_keypassword =dst_config["keypassword"]

    relay_service = ServiceConfig(
        image=RELAY_SERVICE_IMAGE,
        files={
            RELAY_CONFIG_FILES_PATH: "config-files-0"
        },
        cmd=["/bin/sh","-c","./bin/relay --direction both --log_writer.filename log/relay.log --src.address %s --src.endpoint %s --src.key_store %s --src.key_password %s  --src.bridge_mode=%s --dst.address %s --dst.endpoint %s --dst.key_store %s --dst.key_password %s start " % (src_btp_address,src_endpoint,src_keystore,src_keypassword,bridge, dst_btp_address, dst_endpoint, dst_keystore, dst_keypassword)]

    )

    plan.add_service(name=RELAY_SERVICE_NAME,config=relay_service)


