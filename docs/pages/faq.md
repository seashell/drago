# FAQ

## Why Drago?

WireGuard® is an extremely simple yet fast and modern VPN that utilizes state-of-the-art cryptography. It aims to be faster, simpler, leaner, and more useful than IPsec, while avoiding the massive headache. It intends to be considerably more performant than OpenVPN. WireGuard is designed as a general purpose VPN for running on embedded interfaces and super computers alike, fit for many different circumstances. Initially released for the Linux kernel, it is now cross-platform and widely deployable, being regarded as the most secure, easiest to use, and simplest VPN solution in the industry.

WireGuard presents several advantages over other VPN solutions, but it does not allow for the dynamic configuration of network parameters such as IP addresses and firewall rules. Drago builds on top of WireGuard, allowing users to dynamically manage the configuration of their VPN networks, providing a unified control plane for overlays spanning from containers to virtual machines to IoT devices.

## Where can I use Drago?

- Secure home automation, SSH access, etc
- Establish secure VPNs for your company
- Manage access to sensitive services deployed to private hosts
- Expose development servers for debugging and demonstration purposes
- Establish multi-cloud clusters with ease
- Build your own cloud with RaspberryPIs

## What are Drago's main features?

- Single-binary, lightweight
- Encrypted node-to-node communication
- Support for multiple storage backends
- Dynamic network configuration
- Automatic key rotation
- Extensible via REST API
- Slick management dashboard
- Automatic IP assignment (coming soon)
