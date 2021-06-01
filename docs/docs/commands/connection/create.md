# Command: connection create

The `connection create` command is used to create a connection.

## Usage

```
drago connection create <src_node_id> <dst_node_id> <network> [options]
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

## Info Options

- `--json`: Enable JSON output.

- `--allow-all`: Enables routing of all traffic in this connection.

- `--keepalive=<seconds>`: Time interval between persistent keepalive packets. Defaults to 0, which disables the feature.
