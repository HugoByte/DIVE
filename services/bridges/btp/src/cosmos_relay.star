constants = import_module("github.com/hugobyte/dive/package_io/constants.star")

def start_cosmos_relay(plan, args , src, dst):

    plan.print("starting cosmos relay")

    cosmos_node_constants = constants.COSMOS_NODE_CLIENT

    plan.upload_files(src=cosmos_node_constants.config_files, name="archway_config")

    relay_service = ServiceConfig(
        image=cosmos_node_constants.relay_service_image,
        files={
            cosmos_node_constants.relay_config_files_path: "archway_config"
        },
    
        entrypoint=["/bin/sh"]
    )

    plan.add_service(name=cosmos_node_constants.relay_service_name,config=relay_service)

    return struct(
        service_name = cosmos_node_constants.relay_service_name
    )
