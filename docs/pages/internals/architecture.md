# Architecture

Drago is a flexible configuration manager for WireGuard networks which is designed to make it simple to configure network overlays spanning heterogeneous nodes distributed across different clouds and physical locations.

Drago follows a client-server paradigm, in which a centralized server provides multiple clients running alongside WireGuard with their desired state. The desired state is retrieved from the server in a periodic basis and applied to WireGuard running on each node. In other words, the Drago server works as a gateway for accessing network configurations safely stored in a database.

It exposes a comprehensive API through which these configurations can be retrieved and modified, implements authentication mechanisms to prevent unauthorized access, and serves a slick web UI to facilitate the process of managing and visualizing the state of the managed networks.

The Drago client, which runs on every node in the network, is responsible for directly interacting with the server through the API, and for retrieving the most up-to-date configurations. Through a simple reconciliation process, the Drago client then guarantees that the WireGuard configurations on each node match the desired state stored in the database. When running in client mode, Drago also takes care of automatically generating key pairs for WireGuard, and sharing the public key so that nodes can always connect to each other.

The only assumption made by Drago is that each node running the client is also running WireGuard, and that the node in which the configuration server is located is reachable through the network.

Drago does not enforce any specific network topology. Its sole responsibility is to distribute the desired configurations, and guarantee that they are correctly applied to WireGuard on every single registered node. This means that it is up to you to define how your nodes are connected to each other and how your network should look like.

Drago is meant to be simple, and provide a solid foundation for higher-level functionality. Need automatic IP assignment, dynamic firewall rules, or some kind of telemetry? You are free to implement on top of the already existing API.
