RELAY_SERVICE_IMAGE = 'hugobyte/btp-relay'
RELAY_SERVICE_NAME = "btp-bridge"
RELAY_KEYSTORE_FILES_PATH = "/relay/keystores/"
RELAY_KEYSTORE_PATH = "../static-files/keystores/"

def start_relayer(plan, src_endpoint, src_keystore, src_keypassword, src_btp_address, dst_endpoint, dst_keystore, dst_keypassword, dst_btp_address, bridge):
    """
    Start the BTP Relay Service to relay data from a source to a destination using the specified parameters.

    Args:
        plan (Plan): The plan object used to manage services and configurations.
        src_endpoint (str): The source endpoint for relaying data.
        src_keystore (str): The path to the source keystore file.
        src_keypassword (str): The password for the source keystore.
        src_btp_address (str): The BTP address for the source.
        dst_endpoint (str): The destination endpoint for relaying data.
        dst_keystore (str): The path to the destination keystore file.
        dst_keypassword (str): The password for the destination keystore.
        dst_btp_address (str): The BTP address for the destination.
        bridge (str): The bridge mode for the relay.

    Returns:
        None

    """
    plan.print("Starting BTP Relay Service")

    plan.upload_files(src=RELAY_KEYSTORE_PATH, name="keystores")
    
    relay_service = ServiceConfig(
        image=RELAY_SERVICE_IMAGE,
        files={
            RELAY_KEYSTORE_FILES_PATH: "keystores"
        },
        cmd=[
            "/bin/sh",
            "-c",
            "./bin/relay --direction both --log_writer.filename log/relay.log --src.address %s --src.endpoint %s --src.key_store %s --src.key_password %s --src.bridge_mode=%s --dst.address %s --dst.endpoint %s --dst.key_store %s --dst.key_password %s start " % (src_btp_address, src_endpoint, src_keystore, src_keypassword, bridge, dst_btp_address, dst_endpoint, dst_keystore, dst_keypassword)
        ]
    )

    plan.add_service(name=RELAY_SERVICE_NAME, config=relay_service)
