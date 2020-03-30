
Simple demo network with 4 hosts:

- host_0 (container name server) is a bounce-server
- host_1 (container name `client_1`) and (container name `client_2`) are linked to the bounce-server
- host_3 (container name `client_3`) has no links

Instructions:

1. From project root directory, build drago binary and its container :
    ```
    $ STATIC=1 make dev && make container
    ```
2.  Go to demo directory and bring up the demo hosts:
    ```
    $ cd demo/docker-compose/ && \
    docker-compose up --force-recreate
    ```

    To test the network, run: 
    ```
    $ docker exec -ti <SRC_CONTAINER_NAME> ping <DST_CONTAINER_WG_IP>
    ```

    For example, to ping `client_2` from `client_1`:
    ```
    $ docker exec -ti client_1 ping 192.168.2.3
    ```
