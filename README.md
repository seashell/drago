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

Drago is a flexible configuration manager for WireGuard networks which is designed to make it simple to configure network overlays spanning heterogeneous hosts distributed across different clouds and physical locations.

<p align="center"> 
<img src="../assets/demo.gif"/>
</p>

## Overview

[WireGuard®](https://www.wireguard.com/) is an extremely simple yet fast and modern VPN that utilizes state-of-the-art cryptography. It aims to be faster, simpler, leaner, and more useful than IPsec, while avoiding the massive headache. It intends to be considerably more performant than OpenVPN. WireGuard is designed as a general purpose VPN for running on embedded interfaces and super computers alike, fit for many different circumstances. Initially released for the Linux kernel, it is now cross-platform and widely deployable, being regarded as the most secure, easiest to use, and simplest VPN solution in the industry. 

WireGuard presents several advantages over other VPN solutions, but it does not allow for the dynamic configuration of network parameters such as IP addresses and firewall rules.

Drago builds on top of WireGuard, allowing users to dynamically manage the configuration of their VPN networks, providing a unified control plane for overlays spanning from containers to virtual machines to IoT devices.

## How it works

Drago follows a client-server paradigm, in which a centralized server provides multiple clients running alongside WireGuard with their desired state. The desired state is retrieved from the server in a periodic basis and applied to WireGuard running on each host. In other words, the Drago server works as a gateway for accessing network configurations safely stored in a database. 

<h1 align="center"><br>
    <a"><img src="../assets/architecture.png" width="400px"></a>
<br></h1>

It exposes a comprehensive API through which these configurations can be retrieved and modified, implements authentication mechanisms to prevent unauthorized access, and serves a slick web UI to facilitate the process of managing and visualizing the state of the managed networks.

The Drago client, which runs on every host in the network, is responsible for directly interacting with the server through the API, and for retrieving the most up-to-date configurations. Through a simple reconciliation process, the Drago client then guarantees that the WireGuard configurations on each host match the desired state stored in the database. When running in client mode, Drago also takes care of automatically generating key pairs for WireGuard, and sharing the public key so that hosts can always connect to each other.

The only assumption made by Drago is that each host running the client is also running WireGuard, and that the host in which the configuration server is located is reachable through the network.

Drago does not enforce any specific network topology. Its sole responsibility is to distribute the desired configurations, and guarantee that they are correctly applied to WireGuard on every single registered host. This means that it is up to you to define how your hosts are connected to each other and how your network should look like.

Drago is meant to be simple, and provide a solid foundation for higher-level functionality. Need automatic IP assignment, dynamic firewall rules, or some kind of telemetry? You are free to implement on top of the already existing API.

## Build

System requirements:
- Golang 1.14+
- Node 10.17.0+
- yarn 1.12.3+

```
$ go generate
$ go build
```

Alternatively, you can build with `make`, for example:
```
$ make dev
...
```

To see help on building with make, run:
 ```
 $ make help
 ...
 ```
 
## Usage

```
drago --help
drago agent --config=<config_file>
```

## Development

Build the project:
```
go generate
go build
```

Start the Drago server with `air`:

```
./air.sh
```

Once Drago server is up and running, start a dev server for the web UI:

```
cd ui
yarn start
```

Both the backend server and the UI server have hot-reloading capabilities.

Build and run the Drago client:

```
sudo ./drago --config="./dist/client.yml"
```

## Roadmap
- [ ] Website
- [ ] Code coverage
- [ ] Backend API for issuing volatile tokens
- [ ] Integration with Hashicorp Vault
- [ ] Integration with userspace WireGuard implementations such as `cloudflare/boringtun`
- [ ] Integration with firewall tools for more sophisticated networking rules
