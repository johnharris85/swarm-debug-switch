# Description
During debugging it can be useful to switch all Docker nodes in a swarm into / out of `debug` mode. This tool can do that :)

# Usage
Due to `docker service` not supporting some required functionality (`--privileged` and `--pid=host`) we need to start a 'jump' container. The command below should be run on a manager node.

    docker service create \
        --name=debug-switcher \
        --restart-condition=on-failure \
        --restart-max-attempts=3 \
        --mode=global \
        --mount type=bind,source=/var/run/docker.sock,destination=/var/run/docker.sock \
        --mount type=bind,source=/etc/docker,destination=/etc/docker \
        docker:18 \
            run \
            --rm \
            -v /etc/docker:/etc/docker \
            --privileged \
            --pid=host \
            johnharris85/swarm-debug-switch:0.1

## Disclaimer
This container must run as privileged so be aware of the security implications of that.