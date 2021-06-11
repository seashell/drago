# Command: node leave

The `node leave` command is used to have the local client node leave a specific network.

## Usage

```
drago node leave <network> [options]
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
