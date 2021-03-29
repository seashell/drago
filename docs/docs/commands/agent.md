# Command: agent

The `agent` command is likely Drago's most important command. It is used to start client and/or server agent.

## Command-line Options

The `agent` command accepts the following arguments (note that some of them can be overriden by CLI flags):

- `-client`: Enable client mode on the local agent.

- `-config=<path>`: Indicate the path to a configuration file containing configurations to be used by the Drago agent. Can be speficied multiple times.

- `-data-dir=<path>`: 


- `-dev`: Enable development mode on the local agent. This means the agent will execute both as a client and as a server, with sane configurations that are ideal for development.

- `-node=<name>`: 

- `-server`: Enable server mode on the local agent.

- `-servers=<host:port>`: 

- `-wireguard-path`: Path to a userspace WireGuard implementation. Both `wireguard-go` and `cloudflare/boringtun` are supported. If not provided, the WireGuard kernel module will be used.


## Example

```
$ drago agent --server
```
