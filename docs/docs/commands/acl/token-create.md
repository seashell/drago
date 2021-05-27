# Command: acl token create

The `acl token create` command is used to issue new ACL tokens.

## Usage

```
drago acl token create [options]
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

## Create Options

- `--name=<name>`: Sets the name of the ACL token.

- `--type=<type>`: Sets the type of the ACL token. Must be either "client" or "management". If not provided, defaults to "client".

- `--policy=<policy>`: Specifies policies to associate with a client token. Can be specified multiple times.

- `--json`: Enable JSON output.
