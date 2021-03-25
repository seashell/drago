# Command: node status

The `node status` command is used to list the status of one or more registered client nodes.

## Usage

```
drago node status [options] [node]
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
 

## Info Options

- `-self`: Query the status of the local node.

- `-json`: Enable JSON output.