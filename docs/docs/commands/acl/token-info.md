# Command: acl token info

The `acl token info` command is used to display detailed information on an existing ACL token.

## Usage

```
drago acl token info [options] <id>
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