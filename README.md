<h1 align="center"><br>
    <a"><img src="../assets/dragopher.png" width="180"></a>
    <br>
    Drago
<br></h1>

<h5 align="center">
A flexible configuration manager for WireGuard networks
</h5>

------------------

<p align="center">
  <a href="https://goreportcard.com/report/github.com/seashell/drago"><img src="https://goreportcard.com/badge/github.com/seashell/drago" alt="Go report: A+"></a>
  <img alt="GitHub" src="https://img.shields.io/github/license/seashell/drago">
</p>

Drago is a flexible configuration manager for WireGuard which is designed to make it simple to configure secure network overlays spanning heterogeneous nodes distributed across different clouds and physical locations.

<p align="center"> 
<img src="../assets/demo.gif"/>
</p>

## Features

- Single-binary, lightweight
- Slick management dashboard
- Encrypted node-to-node communication
- Support for multiple storage backends
- Dynamic network configuration
- Automatic key rotation
- Extensible via REST API
- Automatic IP assignment
  
## Use cases
- Secure home automation, SSH access, etc
- Establish secure VPNs for your company
- Manage access to sensitive services deployed to private nodes
- Expose development servers for debugging and demonstration purposes
- Establish multi-cloud clusters with ease
- Build your own cloud with RaspberryPIs

## Overview

[WireGuard®](https://www.wireguard.com/) is an extremely simple yet fast and modern VPN that utilizes state-of-the-art cryptography. It aims to be faster, simpler, leaner, and more useful than IPsec. It also intends to be considerably more performant than OpenVPN. WireGuard is designed as a general purpose VPN for running on embedded interfaces and super computers alike, fit for many different circumstances. Initially released for the Linux kernel, it is now cross-platform and widely deployable, being regarded as the most secure, easiest to use, and simplest VPN solution in the industry. 

WireGuard presents several advantages over other VPN solutions, but it does not allow for the dynamic configuration of network parameters such as IP addresses and firewall rules.Drago builds on top of WireGuard, allowing users to dynamically manage the configuration of their VPN networks, providing a unified control plane for overlays spanning containers, virtual machines, and IoT devices.

## How it works

Drago follows a client-server paradigm, in which a centralized server provides multiple clients running alongside WireGuard with their desired state. The desired state is periodically retrieved from the server and applied to WireGuard running on each node. In other words, the Drago server works as a gateway for accessing network configurations safely stored in a database. 

<h1 align="center"><br>
    <a"><img src="../assets/architecture.png" width="400px"></a>
<br></h1>

The Drago server exposes a comprehensive API through which these configurations can be retrieved and modified, implements authentication mechanisms to prevent unauthorized access, and provides a slick web UI to facilitate the process of managing and visualizing the state of the user-defined networks.

The Drago client, in turn, runs on every node in the network, and is responsible for  retrieving the most up-to-date configurations from the server through the API. Thanks to a simple reconciliation process, the Drago client then guarantees that the WireGuard configurations on each node always match the desired state stored in the database. When running in client mode, Drago also takes care of automatically generating key pairs for WireGuard, and sharing the public key so that nodes can always connect to each other.

The only assumptions made by Drago is that each node running a client has WireGuard available either in the kernel or as a userspace binary, and that the server is reachable through the network.

Drago does not enforce any specific network topology. Its sole responsibility is to distribute the desired configurations, and guarantee that they are correctly applied to WireGuard on every single client node. This means that it is up to you to define how your nodes are connected to each other and how your network should look like.

Drago is meant to be simple, and provide a solid foundation for higher-level functionality. Need automatic IP assignment, dynamic firewall rules, or some kind of telemetry? You are free to implement on top of the already existing API.

## Usage
```
Usage: drago [--version] [--help] <command> [<args>]

Available commands:
    acl           Interact with ACL policies and tokens
    agent         Run a Drago agent
    connection    Interact with connections
    interface     Interact with interfaces
    network       Interact with networks
    node          Interact with nodes
    ui            Open the Drago web UI
    version       Prints the Drago version
    
```

## Development

Requirements:
- Golang 1.16+
- Node 10.17.0+
- yarn 1.12.3+

For the sake of convenience, the Drago agent can be initialized in development mode, meaning that it will execute both the client and the server logic. To start the Drago agent in development mode, run:

```
$ go run main.go agent --dev
```

While the agent is running, Drago's UI will be accessible at `127.0.0.1:8080`. If you see a message instead of the UI, it means that it hasn't been built properly. You can easily build the UI project with:

```
$ go generate
```

To apply independent customizations to client and server, start them with:

```
$ go run main.go agent --client --config=<path-to-config-file>
```
and
```
$ go run main.go agent --server --config=<path-to-config-file>
```

We also provide the `air.sh` script, which makes use of `comstrek/air` to perform hot-reloading of the Drago agent. When running Drago through the `air.sh` script, the binary will be rebuilt and restarted whenever a change is detected in the codebase.

### Web UI

The Drago web UI, can be launched independently from the binary. To do this, simply `cd` into the `/ui/` directory containing the UI codebase and run:

```
$ yarn start
```
This will launch React's development server, which comes with hot-reloading out-of-the-box.

While we still don't have any kind of data mocking (e.g., with `miragejs`), if you want to develop the UI you first need to make sure the Drago server is up and running (see instructions above).


## Build

To build the Drago binary, run:
```
$ go generate
$ go build
```

Alternatively, you can build Drago with `make`, for example:
```
$ make dev
```

Run the following to get a comprehensive list of build options:
 ```
 $ make help

 ```

## Contributing

- Fork the project on GitHub
- Clone your fork: `git clone https://github.com/your_username/drago`
- Create a new branch: `git checkout -b my-new-feature`
- Make changes and stage them: `git add .`
- Commit your changes: `git commit -m 'Add some feature'`
- Push to the branch: `git push origin my-new-feature`
- Create a new pull request

## Roadmap

- [x] Website
- [ ] Code coverage
- [x] Backend API for issuing volatile tokens
- [ ] Integration with Hashicorp Vault
- [ ] Integration with userspace WireGuard implementations
  - [x]  `WireGuard/wireguard-go`
  - [ ]  `cloudflare/boringtun`
- [ ] Integration with firewall tools for more sophisticated networking rules
- [ ] Auto-join and auto-meshing modes
- [x] Automatic IP assignment
- [ ] Automatic discovery
- [ ] Etcd storage backend
- [x] RPC API
- [x] Fine-grained authorization
- [x] CLI improvements 

## License
Drago is released under the Apache 2.0 license. See LICENSE.txt
