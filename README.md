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
  <a href="https://gitter.im/seashell/drago"><img src="https://img.shields.io/badge/gitter-join%20chat-green?color=00cc99" alt="Gitter"></a>
</p>

Drago is a flexible configuration manager for WireGuard designed to make it simple to configure secure network overlays spanning heterogeneous nodes distributed across different clouds and physical locations. 

Drago is in active development, and we welcome contributions from the open-source community. Feedback and feature requests are particularly appreciated.

<p align="center"> 
<img src="../assets/demo.gif"/>
</p>

## Features

- Single-binary, lightweight
- Encrypted node-to-node communication
- Support for multiple storage backends
- Support for multiple WireGuard implementations
- Dynamic network configuration
- Automatic key rotation
- Extensible via REST API
- Slick management dashboard
- Automatic IP assignment
  
## Use cases
- Securely connect IoT devices
- Build your own cloud with Raspberry Pi's
- Connect services running on multiple cloud providers 
- Manage access to sensitive services deployed to private hosts
- Expose development servers for debugging and demonstration purposes
- Secure home automation, SSH access, etc
- Establish secure VPNs for your company

## Overview

[WireGuard®](https://www.wireguard.com/) is an extremely simple yet fast and modern VPN that utilizes state-of-the-art cryptography. It aims to be faster, simpler, leaner, and more useful than IPsec. It also intends to be considerably more performant than OpenVPN. WireGuard is designed as a general purpose VPN for running on embedded interfaces and super computers alike, fit for many different circumstances. Initially released for the Linux kernel, it is now cross-platform and widely deployable, being regarded as the most secure, easiest to use, and simplest VPN solution in the industry.

WireGuard presents several advantages over other VPN solutions, but it does not allow for the dynamic configuration of network parameters such as IP addresses and firewall rules. Drago builds on top of WireGuard, allowing users to dynamically manage the configuration of their VPN networks, providing a unified control plane for overlays spanning containers, virtual machines, and IoT devices.

## How it works

Drago follows a client-server paradigm, in which a centralized server provides multiple clients running alongside WireGuard with their desired state. The desired state is periodically retrieved from the server and applied to WireGuard running on each node. In other words, the Drago server works as a gateway for accessing network configurations safely stored in a database. 

<h1 align="center"><br>
    <a"><img src="../assets/architecture.png" width="400px"></a>
<br></h1>

The Drago server exposes a comprehensive API through which these configurations can be retrieved and modified, implements authentication mechanisms to prevent unauthorized access, and provides a slick web UI to facilitate the process of managing and visualizing the state of the user-defined networks.

The Drago client, in turn, runs on every node in the network, and is responsible for  retrieving the most up-to-date configurations from the server through the API. Thanks to a simple reconciliation process, the Drago client then guarantees that the WireGuard configurations on each node always match the desired state stored in the database. When running in client mode, Drago also takes care of automatically generating key pairs for WireGuard, and sharing the public key so that nodes can always connect to each other.

The only assumptions made by Drago is (i) that each node running a client has WireGuard available either as a kernel module or userspace application, and (ii) that the Drag server is reachable through the network.

Drago does not enforce any specific network topology. Its sole responsibility is to distribute the desired configurations, and guarantee that they are correctly applied to WireGuard on every single client node. This means that it is up to the user to define how nodes are connected to each other and how the network should look like.

Drago is meant to be simple, and provide a solid foundation for higher-level functionality. Need automatic IP assignment, dynamic firewall rules, or some kind of telemetry? Feel free to implement it on top of the already existing API.

## Usage
```
Usage: drago [--version] [--help] <command> [options]

Available commands:
    acl           Interact with ACL policies and tokens
    agent         Run a Drago agent
    agent-info    Display status information about the local agent
    connection    Interact with connections
    interface     Interact with interfaces
    network       Interact with networks
    node          Interact with nodes
    ui            Open the Drago web UI
    version       Print the Drago version
```

## Quickstart

A Docker container is provided for those interested in building and running Drago without having to install anything in their systems. 

In order to perform a containerized build of Drago's Docker image, run:

```
$ make container DOCKER=1
```

This will build a minimal Docker image containing the Drago binary. The `DOCKER` flag ensures that the build takes place within a Docker container, thus removing the entry barrier for potential users.

Once the build process finishes, start the Drago agent in development mode with:

```
$ docker run -ti -p 8080:8080 drago agent --dev
```

You can now interact with the system through the Web UI, available at `http://localhost:8080/`. Alternatively, you can also interact with Drago through the command-line interface:

```
$ docker run --network host drago
```

## Development

In order to develop Drago, your environment should meet the following requirements:
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

Note that the `--dev` flag also configures the server to use the in-memory storage backend instead of `etcd`. Therefore, any state will be destroyed whenever the agent is stopped.

To apply independent customizations to client and server, start them with:

```
$ go run main.go agent --client --config=<path-to-config-file>
```
and
```
$ go run main.go agent --server --config=<path-to-config-file>
```

We also provide the `air.sh` script, which makes use of `comstrek/air` to perform hot-reloading of the Drago agent. When running Drago through the `air.sh` script, the binary will be rebuilt and restarted whenever a change is detected in the codebase.

#### Web UI

The Drago web UI, can be launched independently from the binary so that developers can benefit from the covenient features offered by React's development server e.g., hot-reloading. To start Drago's UI in development mode and independently from the binary, `cd` into the `/ui/` directory and run:

```
$ yarn && yarn start
```

This will download all required dependencies, and launch React's development server. The UI can then be accessed at `http://localhost:3000/ui/`.

While running the Web UI independently from the Drago binary is possible, this subproject still lacks a more sophisticated data mocking mechanism to allow for its independent testing. For those interested in contributing to the UI, we suggest that they first start the Drago agent, and then run the Web UI according to the instructions above.


## Build

To build the Drago binary, run:
```
$ go generate
$ go build
```

Alternatively, you can build Drago with `make`. To do so, run:
```
$ make dev
```

In order to get a comprehensive list of build options, run:
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

- [ ] Project: 
  - [x] Website
  - [x] Documentation
  - [ ] Code coverage
  - [ ] E2E testing

- [ ] Features:
  - [x] RPC API
  - [ ] Input validation
  - [ ] Node pre-registration
  - [ ] Drago server clustering
  - [x] Fine-grained authorization
  - [x] Etcd storage backend
  - [x] Inmem storage backend
  - [x] Backend API for issuing volatile tokens
  - [ ] Integration of a plugin system
  - [x] Integration with userspace WireGuard implementations
    - [x]  `WireGuard/wireguard-go`
    - [x]  `cloudflare/boringtun`
  - [ ]  `miragejs` mocks for testing the UI
  - [ ]  Client-side input validation

- [ ] Improvements:
  - [ ] Repository transactions
  - [x] CLI improvements

- [ ] Plugins:
  - [ ]  Meshing
  - [ ]  Leasing
  - [ ]  Admission
  - [ ]  Notification

- [ ] Others:
  - [ ] Vault plugin
  - [ ] Terraform provider
  - [ ] [go-discover](https://github.com/hashicorp/go-discover) provider


## License
Drago is released under the Apache 2.0 license. See [LICENSE](https://github.com/seashell/drago/blob/master/LICENSE).
