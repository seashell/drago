# Command: acl token self

The `acl token self` command is used to display detailed information on the currently set ACL token.

## Usage

```
drago acl token self [options]
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

## Info Options

- `-json`: Enable JSON output.