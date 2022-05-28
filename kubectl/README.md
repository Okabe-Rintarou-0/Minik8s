# Command Guide and Reference Manual

## kubectl apply

+ `kubectl apply -f [filename]`:

  If you want to create api object in a declarative way, you can use this command. 
  
  For example, if you have a pod template yaml file called `pod.yaml` in current directory, then you can type `kubectl apply -f ./pod.yaml` to create a pod according to your specified template.

## kubectl get

+ `kubectl get pod [pod name]`

  This command will show the status of the given pod in a table.
  
  For example, if you have a pod called `example` in the default namespace, then you can check its status by `kubectl get pod example`(since it's in the default namespace, so the namespace can be omitted. Otherwise, you must type `namespace/name` to specify a pod).

  <details>
  <summary>Example</summary>
  <img src="readme-images/kubectl_get_pod.png">
  </details>

+ `kubectl get pods`
  
  This command will show the status of all pods(in all namespaces) in a table.

  <details>
  <summary>Example</summary>
  <img src="readme-images/kubectl_get_pods.png">
  </details>