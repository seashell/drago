
<h1 align="center"><br>
    <a href="https://perun.network/"><img src="./misc/logo.png" alt="Drago" width="96"></a>
    <br>
    Drago
<br></h1>

<h5 align="center">
An easy-to-use, generic control plane for Wireguard network overlays
</h5>

------------------

<p align="center">
  <a href="https://goreportcard.com/report/github.com/seashell/drago"><img src="https://goreportcard.com/badge/github.com/seashell/drago" alt="Go report: A+"></a>
  <img alt="GitHub" src="https://img.shields.io/github/license/seashell/drago">
</p>

Drago is a cloud-agnostic configuration manager for WireGuard designed to make it simple to configure network overlays spanning cloud VMs and edge devices in a hybrid-cloud fashion.

<p align="center"> 
<img src="misc/drago-demo.gif"/>
</p>

## Overview

[WireGuard®](https://www.wireguard.com/) is an extremely simple yet fast and modern VPN that utilizes state-of-the-art cryptography. It aims to be faster, simpler, leaner, and more useful than IPsec, while avoiding the massive headache. It intends to be considerably more performant than OpenVPN. WireGuard is designed as a general purpose VPN for running on embedded interfaces and super computers alike, fit for many different circumstances. Initially released for the Linux kernel, it is now cross-platform and widely deployable, being regarded as the most secure, easiest to use, and simplest VPN solution in the industry. 

Although Wireguard presents several advantages when compared to other VPN solutions, it does not allow for the dynamic configuration of network parameters such as IP addresses and firewall rules. This can pose some difficulties in scenarios such as edge deployments, in which the VPN is the only way to communicate with remote devices.


## How it works

Drago is built on top of WireGuard, a performant and secure VPN that allows for the creation of secure network overlays across heterogeneous hosts, from containers to virtual machines, to IoT devices. It extends Wireguard's functionality by providing users with a unified control plane for the dynamic configuration of the underlying network.

Drago follows a client-server paradigm, in which a centralized server provides multiple clients with their desired state, which is then used to derive changes that must be applied to each host.

This means that the Drago server works as a gateway for accessing network configurations safely stored in a database. It exposes a comprehensive API through which these configurations can be retrieved and modified, and also implements authentication mechanisms to prevent unauthorized access. For the sake of convenience, the Drago server also provides an easy-to-use web UI to facilitate the process of managing and visualizing network parameters and topology.

The Drago client, which runs on every host in the network, is responsible for directly interacting with the server through the API, and for retrieving the most up-to-date configurations. Through a simple reconciliation process, the Drago client then guarantees that the Wireguard configurations on each host match the desired state stored in the database. When running in client mode, Drago also takes care of automatically generating key pairs for Wireguard, and sharing the public key so that hosts can always connect to each other.

No assumption is made by Drago other than that each host has Wireguard installed, and that the configuration server is reachable through the network.

## Build

System requirements:
- Golang 1.14+
- Node 10.17.0+
- yarn 1.12.3+

```
go generate
go build
```

## Usage

```
drago --help
drago agent --config=<config_file>
```

## Development

For development purposes, we suggest that you first start the Drago server:

```
go run main.go agent --config="./demo/server.yml"
```

Once the server is up and running, you can run a dev server for the web UI:

```
cd ui
yarn start
```

Finally, you can build and run the Drago client:

```
go build
sudo ./drago --config="./dist/client.yml"
```

## TODO
- Simple auth for the management API
- Certificate-based client authentication
- Automatic generation of client tokens (in addition to CLI command)
- Integration with a production-grade DB such as Postgres
- Collect metrics from links (e.g., upstream/downstream traffic, last handshake, etc)
- Collection of host metrics (e.g., last seen)
- Input validation (backend + frontend)
- Allow for the management of multiple overlay networks
- Improvements in the overall architecture
- Refactoring (project layout, variables name)
- Persistent connections (e.g., using Websockets) for enchanced responsiveness
- Filtering + Pagination
- Topology graph improvements (labels)
- Implement other storage backends (BoltDB, file, etc)
- Import / Export network topology from e.g., JSON file
- DB query optimization
