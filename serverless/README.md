### 创建文件
`$ kubectl `



## 运行Hello World
### 采用http trigger方法
```
curl localhost:49164 -X POST -d '{"a":55,"b":88,"name":"Taylor"}' --header "Content-Type: application/json"

Hello world, Taylor,55+88=143

```
此处借鉴OpenWhisk的方法，返回一个JSON

# Reference
- https://zhuanlan.zhihu.com/p/141726465