# Appendix 1: Command Guide and Reference Manual

## kubectl apply

+ `kubectl apply -f [filename]`:

  If you want to create api object in a declarative way, you can use this command.

  For example, if you have a pod template yaml file called `pod.yaml` in the current directory, then you can type `kubectl apply -f ./pod.yaml` to create a pod according to your specified template.

## kubectl get

+ `kubectl get pod [pod name]`

  This command will show the status of the given pod in a table.

  For example, if you have a pod called `example` in the default namespace, then you can check its status by `kubectl get pod example` (since it's in the default namespace, so the namespace can be omitted. Otherwise, you must type `namespace/name` to specify a pod).
  
  <img src="../kubectl/readme-images/kubectl_get_pod.png" alt="">
  
+ `kubectl get pods`

  This command will show the status of all pods(in all namespaces) in a table.

  <img src="../kubectl/readme-images/kubectl_get_pods.png" alt="">
  
+ `kubectl get node [node name]`

  This command will show the status of the given node in a table.
  
  <img src="../kubectl/readme-images/kubectl_get_node.png" alt="">
  
+ `kubectl get nodes`

  This command will show the status of all nodes(in all namespaces) in a table.
  
  <img src="../kubectl/readme-images/kubectl_get_nodes.png" alt="">
  
+ `kubectl get rs [replicaSet name]`

  This command will show the status of the given replicaSet in a table.
  
  <img src="../kubectl/readme-images/kubectl_get_rs.png" alt="">
  
+ `kubectl get rss`

  This command will show the status of all replicaSets(in all namespaces) in a table.
  
  <img src="../kubectl/readme-images/kubectl_get_rss.png" alt="">
  
+ `kubectl get hpa [hpa name]`

  This command will show the status of the given horizontal pod autoscaler in a table.
  
  <img src="../kubectl/readme-images/kubectl_get_hpa.png" alt="">
  
+ `kubectl get hpas`

  This command will show the status of all horizontal pod autoscalers(in all namespaces) in a table.
  
  <img src="../kubectl/readme-images/kubectl_get_hpas.png" alt="">
  
+ `kubectl get wf [workflow name]`

  This command will show the status of the given workflow in a table.
  
  <img src="../kubectl/readme-images/kubectl_get_wf.png" alt="" style="zoom:55%;">
  
+ `kubectl get wfs`

  This command will show the status of all workflows(in all namespaces) in a table.
  
  <img src="../kubectl/readme-images/kubectl_get_wfs.png" alt="" style="zoom:55%;">
  
+ `kubectl get func [function name]`
  
  This command will show the status of the given function in a table.
  
  <img src="../kubectl/readme-images/kubectl_get_func.png" alt="" style="zoom:55%;">
  
+ `kubectl get funcs`

  This command will show the status of all functions in a table.
  
  <img src="../kubectl/readme-images/kubectl_get_funcs.png" alt="" style="zoom:55%;">
  
+ `kubectl get gpu [gpu job name]`

  This command will show the status of the given gpu job in a table.
  
  <img src="../kubectl/readme-images/kubectl_get_gpu.png" alt="" style="zoom:55%;">
  
+ `kubectl get gpus`

  This command will show the status of all gpu jobs in a table.
  
  <img src="../kubectl/readme-images/kubectl_get_gpus.png" alt="" style="zoom:55%;">
- `kubectl get service [service name]`

  This command will show the status of the given service in a table.

- `kubectl get services`

  This command will show the status of all services in a table.

- `kubectl get dns [dns name]`

  This command will show the status of the given dns in a table.

- `kubectl get dnses`

  This command will show the status of all dnses in a table.

## kubectl delete

+ `kubectl delete [api object type] [name]`
  
  <img src="../kubectl/readme-images/kubectl_delete.png" alt="">

## kubectl autoscale

+ `kubectl autoscale [hpa name]`

| Parameter | Description                                                  | Example      |
| --------- | ------------------------------------------------------------ | ------------ |
| target    | Required. Specify the target of hpa                          | --target=rs1 |
| min       | Optional, default value is 1. Specify the minimum number of replicas | --min=1      |
| max       | Optional, default value is 1. Specify the maximum number of replicas | --max=4      |
| c         | Optional. Choose CPU utilization as metrics, and the given value as its benchmark | -c 0.5       |
| m         | Optional. Choose memory utilization as metrics, and the given value as its benchmark | -m 0.5       |
| i         | Optional, the default value is 15. Specify the interval of scaling. | -i 30        |

This command will create a hpa according to the given parameters. Notice that if you specify both CPU and memory metrics,
the hpa controller will choose CPU(CPU has higher priority).

You can also create a hpa in a declarative way: `kubectl apply -f hpa.yaml`.

We recommend you to use `kubectl autoscale` command, for it is more obvious.

<img src="../kubectl/readme-images/kubectl_autoscale.png" alt="" style="zoom:55%;">


## kubectl reset

This command is for test only(of course you can also feel free to use it). It will remove all the K-V pairs stored in `etcd`, thus resetting the status of the whole system.

## kubectl gpu

+ `kubectl gpu [gpu job name] -d [directory] -f [file to download]`

  This command will list or download the files throw the `nginx-fileserver`. Files in blue colors are directories.

  `-d` flag stands for the directory to hold the download files(only works when flag `-f` is used).

  `-f` flag stands for the file you want to download. It will be downloaded through `GET` http request into the directory you have specified in flag `-d`.
  
  If you don't use any flags, this command will list the files, just like `ls` command.
  
  - List files example
  
    <img src="../kubectl/readme-images/kubectl_gpu_ls.png" alt="" style="zoom:55%;">
  
    - Download files example
  
      You can see that the file `matrix-op.out` has been successfully downloaded to the current directory.
  
      <img src="../kubectl/readme-images/kubectl_gpu_download.png" alt="" style="zoom:55%;">

## kubectl func

+ `kubectl func add -f [function name] -p [function code path]`

  This command will register a function to the `knative`.

  `-f` flag stands for the name of the function.

  `-p` flag stands for the file path of the code.

  <img src="../kubectl/readme-images/kubectl_func_add.png" alt="" style="zoom:55%;">
  
+ `kubectl func rm -f [function name]`

  This command will remove the function and its instances(pods and replicaSet).
  
  <img src="../kubectl/readme-images/kubectl_func_rm.png" alt="" style="zoom:55%;">
  
+ `kubectl func update -f [function name] -p [function code path]`

  This command will update a function to the `Knative`.

## kubectl trigger

+ `kubectl trigger [function name] -d [function params json]`

  This command will call a serverless function through `http trigger`.

  `-d` flag stands for the JSON form of parameters the function needs.

  <img src="../kubectl/readme-images/kubectl_trigger.png" alt="" style="zoom:55%;">

## kubectl wf

+ `kubectl wf apply -f [workflow.json filepath]`

  This command will create a workflow in a declarative way.

  `-f` flag stands for the `${workflow_name}.json` filepath.

  <img src="../kubectl/readme-images/kubectl_trigger.png" alt="" style="zoom:55%;">
  
+ `kubectl wf rm [workflow name]`

  This command will remove the specified workflow.
  
  <img src="../kubectl/readme-images/kubectl_wf_rm.png" alt="" style="zoom:55%;">

## kubectl label

+ `kubectl label nodes [node name] [label]`

  This command can label a node with given labels.
  
  <img src="../kubectl/readme-images/kubectl_label.png" alt="" style="zoom:55%;">

## kubectl cfg

+ `kubectl cfg sched=[schedule strategy]`

  This command can dynamically change the schedule strategy.

| strategy | Description                                              |
| -------- | -------------------------------------------------------- |
| min-pods | Schedule pod to the node with the minimum number of pods |
| max-pods | Schedule pod to the node with the maximum number of pods |
| min-cpu  | Schedule pod to the node with minimum CPU utilization    |
| min-mem  | Schedule pod to the node with minimum memory utilization |
| random   | Schedule pod to a random node                            |