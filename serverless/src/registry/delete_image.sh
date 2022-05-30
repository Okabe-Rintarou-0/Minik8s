#!/bin/bash
# Inspired by: https://gist.github.com/jaytaylor/86d5efaddda926a25fa68c263830dac1

set -o errexit

if [ -z "$1" ]
  then
    echo "Error: The image name arg is mandatory"
    exit 1
fi


registry='0.0.0.0:5000'
name=$1

curl -v -sSL -X DELETE "http://${registry}/v2/${name}/manifests/$(
    curl -sSL -I \
        -H "Accept: application/vnd.docker.distribution.manifest.v2+json" \
            "http://${registry}/v2/${name}/manifests/$(
            curl -sSL "http://${registry}/v2/${name}/tags/list" | jq -r '.tags[0]'
        )" \
    | awk '$1 == "Docker-Content-Digest:" { print $2 }' \
    | tr -d $'\r' \
)"