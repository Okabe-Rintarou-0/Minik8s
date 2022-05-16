module minik8s

go 1.16

require (
	github.com/coreos/go-iptables v0.6.0
	github.com/docker/docker v20.10.14+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/fatih/color v1.7.0
	github.com/gin-gonic/gin v1.7.7
	github.com/go-redis/redis/v8 v8.11.5
	github.com/rodaine/table v1.0.1
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra v1.4.0
	github.com/stretchr/testify v1.7.1
	go.etcd.io/etcd v3.3.27+incompatible
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/containerd/containerd v1.6.2 // indirect
	github.com/coreos/etcd v3.3.27+incompatible // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/tklauser/go-sysconf v0.3.10 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	golang.org/x/net v0.0.0-20220513224357-95641704303c // indirect
	golang.org/x/sys v0.0.0-20220513210249-45d2b4557a2a // indirect
	google.golang.org/genproto v0.0.0-20220505152158-f39f71e6c8f3 // indirect
	google.golang.org/grpc v1.46.2 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
