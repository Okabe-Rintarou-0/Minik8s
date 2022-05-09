#!/bin/bash

docker run \
-itd \
--name=minik8s-redis \
-p 6379:6379 \
--net=host \
redis
