# serverless模块使用指南

## Function
函数支持Python，默认Python3.10版本。
### Function创建
当Knative服务启动时，会在master节点创建一个私有的Docker Registry。worker节点接收到创建函数请求后，会将函数与本机的模板一起打包成一个镜像，上传至Docker Registry。之后在进行扩容时，worker节点从Docker Registry拉取镜像，并进行容器的创建。


### Function使用
法一：采用http trigger方法，发送POST请求，参数放在JSON内自动解析（参考OpenWhisk）
```
curl localhost:49164 -X POST -d '{"a":55,"b":88,"name":"Taylor"}' --header "Content-Type: application/json"

Hello world, Taylor, 55+88=143
```

法二：采用命令行解析
- invoke指令：
- update指令：
（参考Knative）


## Workflow
我们将workflow看成是一个状态机，workflow上每一个节点对应一个State（参考aws step function）
### Workflow定义
根据要求，我们首先定义了Workflow所支持的功能。




### Workflow创建

### Workflow使用
Workflow使用的

## 扩/缩容机制


# Reference
- https://zhuanlan.zhihu.com/p/141726465
- https://blog.csdn.net/u013595878/article/details/106292026
- https://docs.aws.amazon.com/step-functions/latest