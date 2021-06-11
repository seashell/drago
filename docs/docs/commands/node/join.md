# Command: node join

The `node join` command is used to have the local client node join an existing network.

## Usage

```
drago node join <network> [options]
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
