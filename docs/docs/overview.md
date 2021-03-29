# Overview

> [!WARNING]
> This section is still a work-in-progress. If you think you can contribute, please see our [contribution guidelines](docs/../../../contributing.md).


Drago is a flexible configuration manager for WireGuard networks which is designed to make it simple to configure network overlays spanning heterogeneous nodes distributed across different clouds and physical locations.

Drago is meant to be simple, and provide a solid foundation for higher-level functionality. Need automatic IP assignment, dynamic firewall rules, or some kind of telemetry? You are free to implement on top of the already existing APIs.

## Use-cases

- Secure home automation, SSH access, etc
- Establish secure VPNs for your company
- Manage access to sensitive services deployed to private hosts
- Expose development servers for debugging and demonstration purposes
- Establish multi-cloud clusters with ease
- Build your own cloud with RaspberryPIs

## Main features

- Single-binary, lightweight
- Slick management dashboard
- Encrypted node-to-node communication
- Support for multiple storage backends
- Support for userspace WireGuard implementations
- Dynamic network configuration
- Automatic key rotation
- Extensible via REST API
- Automatic IP assignment