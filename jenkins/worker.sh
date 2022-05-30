#!/bin/bash

BUILD_ID=dontKillMe
pwd
echo "Clear and init"
source /etc/profile

go version

go env -w GOPATH="/data/go"
go env -w GOPROXY="https://goproxy.cn,direct"
go env -w GOROOT="/usr/local/go"

log_dir=$WORKSPACE/logs

kill -9 $(ps -e | grep kubelet | awk '{print $1}')
kill -9 $(ps -e | grep kubeproxy | awk '{print $1}')

docker rm -f $(docker ps -a | grep k8s)

echo "Create log dir: $log_dir"
if [ ! -d $log_dir ]; then
	mkdir $log_dir
fi

echo "Clear and init finished"

weave reset
weave launch 10.119.11.101
weave expose 10.44.0.2/16
echo "Weave Net is ready"


echo "nameserver 10.44.0.9" > /etc/resolv.conf
echo "nameserver 114.114.114.114" >> /etc/resolv.conf
echo "DNS is ready"

go build -o kubectl/run/kubectl kubectl/run/main.go
cp kubectl/run/kubectl ~/kubectl
echo "Kubectl is ready"

echo "Start cadvisor"
docker rm -f cadvisor
script/cadvisor-run.sh

go build -o kubelet/run/kubelet kubelet/run/main.go
nohup kubelet/run/kubelet --ip=10.119.11.93 > $log_dir/kubelet-log.txt &
echo "Kubelet is ready"


kubectl/run/kubectl label nodes cloudos2 type=worker

echo "All done!"