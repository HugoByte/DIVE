# Import modules and constants
neutron_node_service = import_module("github.com/hugobyte/dive/services/cosmvm/neutron/src/node-setup/start_node.star")
constants = import_module("github.com/hugobyte/dive/package_io/constants.star")
neutron_service_config = constants.NEUTRON_SERVICE_CONFIG
neutron_private_ports = constants.NEUTRON_PRIVATE_PORTS
neutron_public_ports = constants.NEUTRON_PUBLIC_PORTS

def start_node_service(plan, args):
    """
    Start a Neutron node service with the provided configuration.

    Args:
        plan (Plan): The deployment plan.
        args (dict): Arguments containing data for configuring the service.

    Returns:
        Any: The response from starting the Neutron node service.
    """

    data = args["data"]
    chain_config = ""

    if len(data) != 0:
        # Configure the service based on provided data
        private_grpc = data["private_grpc"]
        private_tcp = data["private_tcp"]
        private_http = data["private_http"]
        private_rpc = data["private_rpc"]

        public_grpc = data["public_grpc"]
        public_tcp = data["public_tcp"]
        public_http = data["public_http"]
        public_rpc = data["public_rpc"]

        chain_config = neutron_node_service.get_service_config(
            private_grpc, private_http, private_tcp, private_rpc,
            public_grpc, public_http, public_tcp, public_rpc
        )
    else:
        # Use predefined port values for configuration
        chain_config = neutron_node_service.get_service_config(
            neutron_private_ports.grpc, neutron_private_ports.http,
            neutron_private_ports.tcp, neutron_private_ports.rpc,
            neutron_public_ports.grpc, neutron_public_ports.http,
            neutron_public_ports.tcp, neutron_public_ports.rpc
        )

    # Start the Neutron node service and return the response
    return neutron_node_service.start_neutron_node(plan, chain_config)
