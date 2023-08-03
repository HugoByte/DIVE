
cosmvm_cosmos = import_module("github.com/hugobyte/dive/services/cosmvm/cosmos.star")
ARCHWAY = "github.com/hugobyte/dive/services/cosmvm/static_files/config/"
RELAY_SERVICE_NAME = "cosmos-relay"
RELAY_SERVICE_IMAGE = "relay1"
RELAY_CONFIG_FILES_PATH = "/script/"

def start_cosmos_relay(plan, args , src, dst):

    plan.print("starting cosmos relay")

    plan.upload_files(src=ARCHWAY, name="archway_config")

    relay_service = ServiceConfig(
        image=RELAY_SERVICE_IMAGE,
        files={
            RELAY_CONFIG_FILES_PATH: "archway_config"
        },
    
        entrypoint=["/bin/sh"]

    )

    plan.add_service(name=RELAY_SERVICE_NAME,config=relay_service)

    return struct(
        service_name = RELAY_SERVICE_NAME
    )



