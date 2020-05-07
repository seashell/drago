# dist

The dist folder contains configuration samples for various platforms.

## Conventions

On unix systems we will place agent configs under `/etc/drago` and store data under `/var/lib/drago/`. We assume that these directories exist, and that Drago is installed to `/usr/bin/drago`.

## Agent configuration

The following configuration files are provided as examples:

    * `server.yml`
    * `client.yml`

Place one of these under `/etc/drago.d` depending on the host's role. You should use `server.yml` to configure a host as a server or `client.yml` to configure a host as a client.

## systemd

On systems using systemd, the basic systemd unit file under `systemd/drago.service` starts and stops the drago agent. Place it under `/etc/systemd/system/drago.service`.

You can control Drago with `systemctl start|stop|restart drago`.