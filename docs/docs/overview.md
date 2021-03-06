# Overview

> [!WARNING]
> This section is still a work-in-progress. If you think you can contribute, please see our [contribution guidelines](docs/../../../contributing.md).


Drago is a flexible configuration manager for WireGuard networks that is designed to make it simple to configure network overlays spanning heterogeneous nodes distributed across different clouds and physical locations.

Drago is meant to be simple and provide a solid foundation for higher-level functionality. Need automatic IP assignment, dynamic firewall rules, or some kind of telemetry? You are free to implement it on top of the already existing APIs.

## Use-cases

- Secure home automation, SSH access, etc
- Establish secure VPNs for your company
- Manage access to sensitive services deployed to private hosts
- Expose development servers for debugging and demonstration purposes
- Establish multi-cloud clusters with ease
- Build your own cloud with RaspberryPIs

## Main features

- Single-binary, lightweight
- Encrypted node-to-node communication
- Support for different WireGuard implementations
- Slick management dashboard
- Extensible via REST API