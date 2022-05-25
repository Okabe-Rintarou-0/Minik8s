#!/bin/bash
docker run \
--name=cadvisor \
-p 8000:8080 \
--volume=/:/rootfs:ro \
--volume=/var/run:/var/run:rw \
--volume=/sys:/sys:ro \
--volume=/var/lib/docker/:/var/lib/docker:ro \
--detach=true \
google/cadvisor \
-storage_driver=influxdb \
-storage_driver_host=influxsrv:8086 \
