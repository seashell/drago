# Command: interface update

The `interface update` command is used to update an existing interface.

## Usage

```
drago interface update <interface_id> [options]
```

## General Options

- `--address=<addr>`
    The address of the Drago server.
    Overrides the `DRAGO_ADDR` environment variable if set.
    Defaults to `http://127.0.0.1:8080`.

- `--token=<token>`
    The token used to authenticate with the Drago server.
    Overrides the `DRAGO_TOKEN` environment variable if set.
    Defaults to `""`.

## Update Options

- `--address`: Interface IP address in CIDR notation
