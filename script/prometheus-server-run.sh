DIR=/etc/prometheus
TEST_NODE_IP=192.168.1.103

if [ ! -d "$DIR" ]; then
	mkdir "$DIR"
fi

./prometheus_utils/create-yaml "$TEST_NODE_IP" > /etc/prometheus/prometheus.yml
chmod 777 /etc/prometheus/prometheus.yml
docker run --name=prometheus -d -p 9090:9090 -v /etc/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus

