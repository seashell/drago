# Dist

The dist folder contains sample configs for various platforms.

## Conventions

On unixes we will place agent configs under `/etc/drago` and store data under `/var/lib/drago/`. You will need to create both of these directories. We assume that drago is installed to `/usr/bin/drago`.

## Agent Configs

The following example configuration files are provided:

    * `server.yml`
    * `client.yml`

Place one of these under `/etc/drago.d` depending on the host's role. You should use `server.yml` to configure a host as a server or `client.yml` to configure a host as a client.

## Systemd

On systems using systemd the basic systemd unit file under `systemd/drago.service` starts and stops the drago agent. Place it under `/etc/systemd/system/drago.service`.

You can control drago with `systemctl start|stop|restart drago`.