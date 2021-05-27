# Architecture

Drago follows a client-server paradigm, in which a centralized server provides multiple clients running alongside WireGuard with their desired state. The desired state is periodically retrieved from the server and applied to the WireGuard interfaces on each node.

Drago exposes two APIs through which configurations can be retrieved and modified. The RPC API is primarily for client agents to synchronize their configurations, whereas the HTTP API is meant for management purposes.

Drago implements authentication mechanisms to prevent unauthorized access, and serves a slick web UI to facilitate the process of updating and visualizing the state of the managed networks. Everything is nicely bundled within the same binary.

The Drago client, expected to run on every node in the network, is responsible for directly interacting with the server through the API, and for retrieving the most up-to-date configurations. Through a simple reconciliation process, the Drago client then guarantees that the WireGuard configurations on each node match the desired state stored in the database. When running in client mode, Drago also takes care of automatically generating key pairs for WireGuard, and sharing the public key so that nodes can always connect to each other.

The only assumption made by Drago is that each node running the client is also running WireGuard and that the node in which the configuration server is located is reachable through the network.

Drago does not enforce any specific network topology. Its sole responsibility is to distribute the desired configurations, and guarantee that they are correctly applied to WireGuard on every single registered node. This means that it is up to you to define how your nodes are connected to each other and how your network should look like.
