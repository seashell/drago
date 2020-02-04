# drago

TODO

## Build

System requirements:
- Golang 1.14+
- Node v___
-  yarn v___ (for building the Web UI)


```
go generate
go build
```

## Usage

```
drago agent
```

## Overview

Drago connects nodes in a cluster by providing an encrypted layer 3 network that can span across data centers and public clouds. By allowing pools of nodes in different locations to communicate securely, Drago enables the operation of multi-cloud clusters as well as the connection of edge and IoT devices.

Drago's design allows clients to VPN to a cluster in order to securely access services running on the cluster.

## How it works

Drago uses [WireGuard](https://www.wireguard.com/), a performant and secure VPN, to connect the different nodes in the cluster.

The Drago agent runs on every node in the cluster, setting up the public and private keys for the VPN as well as the necessary rules to route packets between locations.

Drago can operate both as a complete, independent networking provider as well as an add-on complimenting the cluster-networking solution currently installed on a cluster.