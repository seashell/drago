# Command: acl token update

The `acl token update` command is used to update an existing ACL token.

## Usage

```
drago acl token update [options] <id>
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

## Update Options

- `-name`: Sets the name of the ACL token.


- `-type`: Sets the type of the ACL token. Must be either "client" or "management". If not provided, defaults to "client".


- `-policy`: Specifies policies to associate with a client token. Can be specified multiple times.

- `-json`: Enable JSON output.