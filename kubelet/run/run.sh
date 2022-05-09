#!/bin/bash
DIR="$(dirname "$0")"
"$DIR"/cadvisor-run.sh

"$DIR"/kubelet