# Minik8s

## Description

Group project of SE3356 Cloud Operating System Design and Practice, Spring 2022.

## Structure

### Project overview

<img src="https://img.shields.io/badge/license-Apache-blue" alt="">
<img src="https://img.shields.io/badge/go-1.16-blue" alt="">

<details>
<summary>Dependencies</summary>
<ul>
<li>gin: <a href="https://github.com/gin-gonic/gin">https://github.com/gin-gonic/gin</a></li>
<li>cobra: <a href="https://github.com/spf13/cobra">https://github.com/spf13/cobra</a></li>
<li>cast: <a href="https://github.com/spf13/cast">https://github.com/spf13/cast</a></li>
<li>gjson: <a href="https://github.com/tidwall/gjson">https://github.com/tidwall/gjson</a></li>
<li>go.uuid: <a href="https://github.com/satori/go.uuid">https://github.com/satori/go.uuid</a></li>
<li>etcd: <a href="https://github.com/etcd-io/etcd">https://github.com/etcd-io/etcd</a></li>
<li>table: <a href="https://github.com/rodaine/table">https://github.com/rodaine/table</a></li>
<li>goph: <a href="https://github.com/melbahja/goph">https://github.com/melbahja/goph</a></li>
<li>redis: <a href="https://github.com/go-redis/redis">https://github.com/go-redis/redis</a></li>
<li>goquery: <a href="https://github.com/PuerkitoBio/goquery">https://github.com/PuerkitoBio/goquery</a></li>
<li>yaml: <a href="https://github.com/go-yaml/yaml">https://github.com/go-yaml/yaml</a></li>
<li>testify: <a href="https://github.com/stretchr/testify">https://github.com/stretchr/testify</a></li>
<li>go-iptables: <a href="https://github.com/coreos/go-iptables">https://github.com/coreos/go-iptables</a></li>
</ul>
</details>

### Kubectl

`kubectl` is a command line tool that helps user controller `minik8s`. It's similar to `kubectl` in `Kubenetes`, but
it's simplified and different in some commands. It is also based on `cobra`.

<img src="./readme-images/cobra.png" alt="cobra 5 " style="zoom:50%;" />

We support basic command like `kubectl get pods`, `kubectl apply -f xxx.yaml`. For more info,
see [kubectl README](/kubectl/README.md).

### Kubelet

The structure of `kubelet` in `minik8s` is similar to k8s, but it's greatly simplified.

![Our kubelet](./readme-images/kubelet.svg)

#### Core: How to create a pod

Start an infra container first(default image is `registry.aliyuncs.com/google_containers/pause:3.6`). The infra
container provides network namespace and volumes for all the other containers. So they can communicate with each other
through `localhost` and share same volumes.

#### How to allocate unique IP for pods

[Weave Net](https://www.weave.works/) can be used as a Docker plugin. A Docker network named `weave` is created
by `weave launch`, which can be visible in the whole cluster. Under the Weave Net, containers can be allocated
its `ClusterIP` in the cluster.

After the `pause` container has been created, `kubelet` runs the command `weave attach <ip> <pause_container_id>` to
attach `ClusterIP` to the pod.

To make the Weave Net visible from host, run the command `weave expose <ip>` to join the Weave Net.

#### Support & References

+ Docker http client: [Moby](https://pkg.go.dev/github.com/docker/docker/client)
+ Docker api document: [Docker Engine API (v1.41)](https://docs.docker.com/engine/api/v1.41/#)
+ A good article on pod
  creation: [2.2 从 Pause 容器理解 Pod 的本质](https://k8s.iswbm.com/c02/p02_learn-kubernetes-pod-via-pause-container.html)

### Scheduler

The function of `scheduler` is very clear and simple, scheduling the newly created pods.

The process of scheduling can be broken down into several steps:

+ Get the list of available nodes from `api-server` through REST apis.
+ If there is a `nodeSelector` in the specification of the pod, then filter the list according to given labels.
+ Select one node according to a certain strategy from the filtered nodes.

![scheduler](readme-images/scheduler.svg)

Here are some strategies we support now:

| strategy                   | Description                                     |
|----------------------------|-------------------------------------------------|
| random                     | Select a node randomly                          |
| minimum pods               | Select the node with minimum number of pods     |
| maximum pods               | Select the node with maximum number of pods     |
| minimum cpu utilization    | Select the node with minimum cpu utilization    |
| minimum memory utilization | Select the node with minimum memory utilization |

The strategy is by default `random`. You can also dynamically change the strategy by using `kubectl cfg` command. For
more info about this command, see the README.md of `kubectl`.

### Api-server

<img alt="gin" align="right" height="150" src="./readme-images/gin.png"/>

`Api-server` is the center of `minik8s`. It should expose REST apis for other components of the control plane. For fast
development, we adopted a mature framework: `gin`

### Proxy

`Proxy` is responsible for allocating virtual service IP, which is unique and visible in the whole cluster. `Proxy` will
prepare an `nginx` container for each service. The `nginx` container will attain its service IP in Weave Net and `proxy`
will configure the `nginx.conf`. Users and other pods in the Weave Net can visit services just through these service IP.

### DNS

[CoreDNS](https://coredns.io/manual/toc/)  is a DNS server, which can be configured through its `Corefile`.
When `api-server` receives an `apply` request of `DNS`,  `api-server` will first start an `nginx` container to deal with
the path-service mapping issue and then add the IP-name mapping to CoreDNS.

### Autoscaler

#### Structure

`Kubelet` in each node will collect runtime status through `docker status`, including cpu and memory utilization. All
these statuses will be published to a certain topic, on which both `api-server` and `controller-manager` are watching.

Here is a shared cache in the `controller-manager`. It can receive the statuses published by `kubelet` and
do `incremental synchronization`. Meanwhile, `api-server` will persist these statuses into `etcd`, a distributed K-V
store system.
`etcd` is the one who truly indicates the status of the whole system. So, the cache in the `controller-manager` has to
periodically do full synchronization with `api-server`, in order to stay consistent with `etcd`.

![Autoscaler](./readme-images/autoscaler_structure.svg)

#### Visualization

The pod resources monitor is based on `cAdvisor`, `Prometheus` and `Grafana`.

<img src="./readme-images/autoscaler_visualization.svg" alt="Autoscaler 5" style="zoom:50%;" />

We recommend you to use grafana dashboard with UID `11277` and `893`.

Here is a good reference: [Build up Prometheus + Grafana + cAdvisor](https://blog.51cto.com/jiachuanlin/2538983)

#### Hint

Because all these components are running in containers, so you can't access other running component by simply
using `localhost`(Even if they are running in `host` network mode). Please use the ip instead.

### GPU

Users only need to specify the scripts needed to compile cuda files and run them, and also the work directory.

The cuda files(ended with `.cu`) will be recognized and uploaded to the π2.0 platform. The slurm script will be created
automatically according to given parameters.

The jobs should be independent of each other, so we adopt a sidecar structure. The `gpu-server` will upload cuda files,
compile them, create slurm script and finally submit the job by using command `sbatch`.

Since we don't have a good idea to be aware of the completion of submitted jobs(π2.0 supports email alert, but it's not
suitable for this situation). So we temporarily adopt the strategy of polling(every 5 minutes). Once the job has been
completed(can be known by using command `sacct`. If the job returned is `COMPLETED` in its `State` field, then it is
completed), the `gpu-server` will download the output file and error file(`xxx.out`, `xxx.err`, specified by users).
Users can then browse and download the results of jobs using `nginx-fileserver`.

![](./readme-images/gpu-pod-struct.svg)

#### Cuda

Each block in `cuda grid` is corresponding to an area in a matrix. We can map a cell in the block to an element in a
matrix.

```c
int i = blockIdx.x * blockDim.x + threadIdx.x;
int j = blockIdx.y * blockDim.y + threadIdx.y;
```

![cuda grid](readme-images/cuda.png)

`blockIdx` stands for the coordinate of a block. For example, the block in the upper left corner has `blockIdx(0,0)`

`blockDim` stands for the dimension of a block. A block is two-dimensional, so `blockDim.x` stands for the width
while `blockDim.y` stands for the height.

`threadIdx` stands for the coordinate of a thread inside a block. Similar to `blockIdx`.

Because `gpu` is a `device` and it does not share memory with `cpu`. Special functions should be used to do memory
operations in `gpu`. 

Keyword `__global__` can be used to define a function that will be called in `gpu`. If you have defined a function f:
```c
__global__ void f() {}
```

Then you can call it by:
```
f <<<blockDim, threadDim>>>();
```
The types of both `blockDim` and `threadDim` are `Dim3`. You can define a variable of type `Dim3` by:
```
Dim3 var(x, y)
```
You only need specify the first two dimension, for the third dimension is always 1.

For matrix addition:
```
__global__ void matrix_add(int **A, int **B, int **C) {
    int i = blockIdx.x * blockDim.x + threadIdx.x;
    int j = blockIdx.y * blockDim.y + threadIdx.y;
    C[i][j] = A[i][j] + B[i][j];
}
```

For matrix multiplication:
```c
__global__ void matrix_multiply(int **A, int **B, int **C) {
    int i = blockIdx.x * blockDim.x + threadIdx.x;
    int j = blockIdx.y * blockDim.y + threadIdx.y;
    C[i][j] += A[i][j] * B[j][i];
}
```
#### π2.0 GPU Support

See:

+ https://github.com/SJTU-HPC/docs.hpc.sjtu.edu.cn
+ https://docs.hpc.sjtu.edu.cn/index.html
+ https://docs.hpc.sjtu.edu.cn/job/slurm.html
+ https://studio.hpc.sjtu.edu.cn/

### Serverless

#### Structure

The structure of our serverless system draws lessons from `Knative` but is quite simplified. Users can register
functions to `api-server`. `KPA controller` will create corresponding function image and push it into docker registry.
It will also create a replicaSet through `api-server` apis.
** NOTE: If you want to delete function image from docker image, you should set the environment
variable `REGISTRY_STORAGE_DELETE_ENABLED=true`, otherwise deleting process will fail. **

The `ReplicaSet Controller` can then create pods on nodes. Notice that there is a http server running on master node
(port `8081`), and you can call a function by http trigger.

![Knative](./readme-images/knative.svg)

#### Function Registration

User can register a function (we only support `python` now) to the `api-server`. Here is an example of function:

```python
def main(params):
    x = params["x"]
    x = x + 5
    result = {
        "x": x
    }
    return result
```

This function needs a parameter `x` and `x` is passed in the form of `json`, and will add 5 to `x` and return a
dictionary/json.

In our system, all parameters and results can be transferred in the form of `json` , and there is no need to pass the
parameter type in `json`, since interpreter itself can infer variable type automatically.

Once a function is registered, a corresponding image will be pushed to the registry and a replicaSet will be created,
which will create pods(function instances) on worker nodes.

#### Http Trigger

We support a convenient way to call a function by http trigger. You can type `kubectl trigger [funcname] -d [data]` to
send http trigger to the specified function instances.

Because the function instances are maintained by a `replicaSet`, so the http server in `Knative` will randomly choose
one pod in the replicaSet and call it.

Take `addFive` for example, you can type `kubectl trigger addFive -d '{"x": 100}'`, and you will get a
response: `'{"x": 105}'`

All pods have their own unique ip, so they can be called by `POST` http request to `${pod_ip}:8080`.

#### Workflow

A workflow is equivalent to a DAG of functions. It can be defined in the form of `json`,
see [workflow](apiObject/examples/workflow) for examples.

Our implementation draws lessons from AWS. We also support `Choice` and `Task`.

<details>
<summary>Workflow example(Graph)</summary>
<img src="readme-images/workflow.svg"/>
</details>

<details>
<summary>Workflow example(Json)</summary>
<pre><code>{
  "apiVersion": "/api/v1",
  "kind": "Workflow",
  "metadata": {
    "namespace": "default",
    "name": "print"
  },
  "startAt": "addFive",
  "params": {
    "x": 5
  },
  "nodes": {
    "addFive": {
      "type": "Task",
      "next": "judge"
    },
    "judge": {
      "type": "Choice",
      "choices": [
        {
          "variable": "x",
          "numericEquals": 10,
          "next": "printEquals"
        },
        {
          "variable": "x",
          "numericNotEquals": 10,
          "next": "printNotEquals"
        }
      ]
    },
    "printEquals": {
      "type": "Task"
    },
    "printNotEquals": {
      "type": "Task"
    }
  }
}
</code></pre>
</details>

##### Reference

+ [创建无服务器工作流](https://aws.amazon.com/cn/getting-started/hands-on/create-a-serverless-workflow-step-functions-lambda/)

## Tools

### Container Management

For `windows`, we have `Docker Desktop` to monitor the stats of all containers. But in linux, we don't have such
convenience.

![Portainer](./readme-images/portainer.png)

Fortunately, `portainer` performs even better than `Docker Desktop`. It can be deployed easily by using docker. You can
type `./portainer-run.sh` to start the portainer. Then you can access it at http://localhost:9000.

### Automatic deployment

![Jenkins](readme-images/Jenkins.png)

`Jenkins` is super convenient for our project.

#### Q & A

Q: Why `nohup` does not work?

A: Killed by `Jenkins`. Try to add `BUILD_ID=dontKillMe` to the shell script.

Q: Why `go: command not found`?

A: Please add environment variables it needs manually to `Jenkins`.
