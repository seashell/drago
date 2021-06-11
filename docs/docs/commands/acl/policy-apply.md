# Command: acl policy apply

The `acl policy apply` command is used to create or update ACL policies.

## Usage

```
drago acl policy apply <name> [options]
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

## Apply Options

- `--description=<description>`: Sets the description of the ACL policy.
