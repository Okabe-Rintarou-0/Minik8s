docker run -d --net="host" --name node_exporter --restart=unless-stopped -p 9100:9100 \
-v "\proc:/host/proc:ro" \
-v "\sys:/host/sys:ro" \
-v "\:/rootfs:ro" \
prom/node-exporter
