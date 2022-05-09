DIR=/opt/grafana-storage
if [ ! -d "$DIR" ]; then
	mkdir "$DIR"
fi
chmod 777 -R "$DIR"

docker run -d \
-p 3000:3000 \
--name=grafana \
-v /opt/grafana-storage:/var/lib/grafana \
grafana/grafana