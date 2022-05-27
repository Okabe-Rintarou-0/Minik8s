mkdir -p /etc/kube/dns/
cp ./apiserver/src/dns/CoreFile /etc/kube/dns/Corefile
cp ./apiserver/src/dns/hosts /etc/kube/dns/hosts
docker stop coredns
docker rm coredns
docker run -d --name coredns -v /etc/kube/dns/:/data/ -v /etc/kube/dns/Corefile:/Corefile coredns/coredns:latest
weave attach 10.44.0.9/16 coredns
echo "nameserver 10.44.0.9" > /etc/resolv.conf
echo "nameserver 114.114.114.114" >> /etc/resolv.conf