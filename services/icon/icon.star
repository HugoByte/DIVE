ICON_SERVICE_NAME = "icon"

ICON_NODE_IMAGE = "hemz1012/goloop"
ICON_RPC_PORT_ID = 9080
EXECUTABLE_PATH = "/bin/goloop"

def launch_icon_node(plan,args):
    plan.print("Launching The Icon Node")

    icon_node_service_config = ServiceConfig(
        image=ICON_NODE_IMAGE,
        ports={"http":PortSpec(number=ICON_RPC_PORT_ID,transport_protocol="TCP")}
    )

    icon_node_service = plan.add_service(name=ICON_SERVICE_NAME,config=icon_node_service_config)
    plan.print(icon_node_service)

    response = plan.wait(
            service_name=icon_node_service.name,
            recipe=PostHttpRequestRecipe(port_id="http",endpoint="/api/v3",content_type="application/json",body="{ \"jsonrpc\":\"2.0\", \"id\" :1, \"method\" :\"icx_getLastBlock\"}",extract={
                "height":".result.height"
            }),
            field="code",
            assertion="==",
            target_value=200,
            timeout="1m",
        )


    return "{0}:{1}".format(icon_node_service.ip_address,ICON_RPC_PORT_ID)