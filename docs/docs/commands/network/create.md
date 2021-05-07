# Command: network create

The `network create` command is used to create a new network.

## Usage

```
drago network create [options]
```

## General Options

- `-address=<addr>`
    The address of the Drago server.
    Overrides the DRAGO_ADDR environment variable if set.
    Defaults to `http://127.0.0.1:8080`


- `-token=<token>`
    The token used to authenticate with the Drago server.
    Overrides the `DRAGO_TOKEN` environment variable if set.
    Defaults to `""`
 

## Create Options

- `-name`: Network name.

- `-ip_range`: Network IP address range in CIDR notation.

- `-json`: Enable JSON output.