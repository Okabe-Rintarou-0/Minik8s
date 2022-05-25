# Minik8s

## Description

Group project of SE3356 Cloud Operating System Design and Practice, Spring 2022.

## Structure

### Kubelet

The structure of `kubelet` in minik8s is similar to k8s, but it's greatly simplified.

![Our kubelet](./readme-images/kubelet.svg)

### Pod Resources Monitor

The pod resources monitor is based on `cAdvisor`, `Prometheus` and `Grafana`.

We recommend you to use grafana dashboard with UID `11277` and `893`.

Here is a good reference: [Build up Prometheus + Grafana + cAdvisor](https://blog.51cto.com/jiachuanlin/2538983)

#### Hint

Because all these components are running in containers, so you can't access other running component by simply
using `localhost`(Even if they are running in `host` network mode).
Please use the ip instead.

### GPU
Users only need to specify the scripts needed to compile cuda files and run them, and also the work directory. 

The cuda files(ended with `.cu`) will be recognized and uploaded to the π2.0 platform. The slurm script will be created automatically according to given parameters.

The jobs should be independent of each other, so we adopt a sidecar structure. The `gpu-server` will upload cuda files, compile them, create slurm script and finally submit the job by using command `sbatch`.

Since we don't have a good idea to be aware of the completion of submitted jobs(π2.0 supports email alert, but it's not suitable for this situation). So we temporarily adopt the strategy of polling(every 5 minutes). Once the job has been completed(can be known by using command `sacct`. If the job returned is `COMPLETED` in its `State` field, then it is completed), the `gpu-server` will download the output file and error file(`xxx.out`, `xxx.err`, specified by users). Users can then browse and download the results of jobs using `nginx-fileserver`.

![](./readme-images/gpu-pod-struct.svg)

#### π2.0 GPU Support
See:
+ https://github.com/SJTU-HPC/docs.hpc.sjtu.edu.cn
+ https://docs.hpc.sjtu.edu.cn/index.html
+ https://docs.hpc.sjtu.edu.cn/job/slurm.html
+ https://studio.hpc.sjtu.edu.cn/

## Tools

### Container Management
For `windows`, we have `Docker Desktop` to monitor the stats of all containers.
But in linux, we don't have such convenience.

![Portainer](./readme-images/portainer.png)

Fortunately, `portainer` performs even better than `Docker Desktop`.
It can be deployed easily by using docker. You can type `./portainer-run.sh` to start the portainer.
Then you can access it at http://localhost:9000.

### Automatic deployment

![Jenkins](readme-images/Jenkins.png)

`Jenkins` is super convenient for our project.

#### Q & A
Q: Why `nohup` does not work?

A: Killed by `Jenkins`. Try to add `BUILD_ID=dontKillMe` to the shell script.

Q: Why `go: command not found`?

A: Please add environment variables it needs manually to `Jenkins`.

![](readme-images/Jenkins_Env.png)