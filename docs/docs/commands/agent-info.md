# Command: agent-info

The `agent-info` command is used to display configurations and status from the agent to which the CLI is connected.

## Usage

```
drago agent-info [options]
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

## Agent Info Options

- `--json`: Enable JSON output.
