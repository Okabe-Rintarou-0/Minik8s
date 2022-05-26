mkdir -p /etc/kube/dns/
cp ./apiserver/src/dns/CoreFile /etc/kube/dns/Corefile
cp ./apiserver/src/dns/hosts /etc/kube/dns/hosts
docker stop coredns 2>/dev/null >/dev/null
docker rm coredns 2>/dev/null >/dev/null
docker pull coredns/coredns:latest 2>/dev/null >/dev/null
docker run -d --name coredns -v /etc/kube/dns/:/data/ -v /etc/kube/dns/Corefile:/Corefile coredns/coredns:latest