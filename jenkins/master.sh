#!/bin/bash

BUILD_ID=dontKillMe
pwd
echo "Clear and init"
source ~/.bashrc

go version

go env -w GOPATH="/root/go/SDK"
go env -w GOPROXY="https://goproxy.cn,direct"
go env -w GOROOT="/usr/local/lib/go"

log_dir=$WORKSPACE/logs

rm -f /root/gpu/matrix_op*

kill -9 $(ps -e | grep api-server | awk '{print $1}')
kill -9 $(ps -e | grep scheduler | awk '{print $1}')
kill -9 $(ps -e | grep controller-mana | awk '{print $1}')
kill -9 $(ps -e | grep kubelet | awk '{print $1}')
kill -9 $(ps -e | grep kubeproxy | awk '{print $1}')
kill -9 $(ps -e | grep knative | awk '{print $1}')

docker rm -f $(docker ps -a | grep k8s)

docker start redis
echo "Start redis"

echo "Create log dir: $log_dir"
if [ ! -d $log_dir ]; then
	mkdir $log_dir
fi

echo "Clear and init finished"

echo "Start weave"
./script/weave-start-master.sh
echo "Weave is ready"

echo "Start dns"
./script/dns-restart.sh
echo "Dns is ready"

echo "Start building up minik8s control plane now"
go build -o apiserver/run/api-server apiserver/run/main.go
nohup apiserver/run/api-server > $log_dir/api-server-log.txt &
sleep 20s
echo "Api server is ready"

go build -o scheduler/run/scheduler scheduler/run/main.go
nohup scheduler/run/scheduler > $log_dir/scheduler-log.txt &
echo "Scheduler is ready"

go build -o controller/run/controller-manager controller/run/main.go
nohup controller/run/controller-manager > $log_dir/controller-manager-log.txt &
echo "Controller Manager is ready"

go build -o kubectl/run/kubectl kubectl/run/main.go
cp kubectl/run/kubectl ~/kubectl
echo "Kubectl is ready"

echo "Start cadvisor"
docker rm -f cadvisor
script/cadvisor-run.sh

go build -o kubelet/run/kubelet kubelet/run/main.go
nohup kubelet/run/kubelet --ip=10.119.11.101 > $log_dir/kubelet-log.txt &
echo "Kubelet is ready"

go build -o proxy/run/kubeproxy proxy/run/main.go
nohup proxy/run/kubeproxy > $log_dir/proxy-log.txt &
echo "Proxy is ready"

go build -o serverless/run/knative serverless/run/main.go
nohup serverless/run/knative > $log_dir/knative-log.txt &2>1  &
echo "Knative is ready"

docker rm -f prometheus
go build -o script/prometheus_utils/create-yaml script/prometheus_utils/create_prometheus_yaml.go

script/prometheus_utils/create-yaml

kubectl/run/kubectl label nodes cloudos1 type=master

cd script

./prometheus-server-run.sh
docker rm -f grafana
./grafana-run.sh

echo "All done!"