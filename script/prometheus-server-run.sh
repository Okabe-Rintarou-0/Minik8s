DIR=/etc/prometheus

if [ ! -d "$DIR" ]; then
	mkdir "$DIR"
fi

./prometheus_utils/create-yaml > /etc/prometheus/prometheus.yml
chmod 777 /etc/prometheus/prometheus.yml
docker run --name=prometheus -d -p 9090:9090 -v /etc/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus

