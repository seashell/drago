# Command: acl policy upsert

The `acl policy upsert` command is used to create or update ACL policies.

## Usage

```
drago acl policy upsert [options] <name>
```

## General Options

- `-address=<addr>`
    The address of the Drago server.
    Overrides the DRAGO_ADDR environment variable if set.
    Defaults to `http://127.0.0.1:8080.`


- `-token=<token>`
    The token used to authenticate with the Drago server.
    Overrides the `DRAGO_TOKEN` environment variable if set.
    Defaults to `""`

## Create Options

- `-description`: Sets the description of the ACL policy.