# Drago Commands (CLI)


Drago is controlled via a very easy to use command-line interface (CLI). Drago is only a single command-line application: `drago`, which takes subcommands such as `agent` or `status`. The complete list of subcommands is in the navigation to the left.

The Drago CLI is a well-behaved command line application. In erroneous cases, a non-zero exit status will be returned. It also responds to `-h` and `--help` as you would most likely expect.

To view a list of the available commands at any time, just run Drago with no arguments. To get help for any specific subcommand, run the subcommand with the `-h` argument.

Each command has been conveniently documented on this website. Links to each command can be found on the left.

### Remote usage

The Drago CLI can be used to interact with a remote Drago agent.

To do so, set the `DRAGO_ADDR` environment variable or use the `--address=<addr>` flag when running commands.

```
$ DRAGO_ADDR=https://<remote_addr>:8080 drago agent-info
$ drago agent-info --address=https://remote-address:4646
```

The provided address must be reachable from your local machine.
