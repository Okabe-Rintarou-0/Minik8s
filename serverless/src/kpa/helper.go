package kpa

import (
	"context"
	"fmt"
	"minik8s/apiObject"
	"minik8s/apiserver/src/url"
	"minik8s/serverless/src/function"
	"minik8s/serverless/src/registry"
	"minik8s/util/apiutil"
	"minik8s/util/httputil"
	"minik8s/util/imageutil"
	"path"
	"strconv"
	"time"
)

func (c *controller) updateReplicaSetToApiServer(funcReplicaSet *functionReplicaSet) {
	URL := url.Prefix + path.Join(url.ReplicaSetURL, "function", funcReplicaSet.Function)
	resp := httputil.PutForm(URL, map[string]string{
		"replicas": strconv.Itoa(funcReplicaSet.NumReplicas),
	})
	logManager("update rs and get resp: %s", resp)
}

func (c *controller) scaleToHalf(funcReplicaSet *functionReplicaSet) {
	funcReplicaSet.NumReplicas /= 2
	logManager("Scaling, num replicas changed to %d", funcReplicaSet.NumReplicas)
	c.updateReplicaSetToApiServer(funcReplicaSet)
}

// createReplicaSet is coroutine-safe, should be called without lock
func (c *controller) createReplicaSet(apiFunc *apiObject.Function) {
	c.scaleLock.Lock()
	defer c.scaleLock.Unlock()
	imageName := registry.RegistryHost + "/" + imageutil.FormatImageName(apiFunc.Name)
	replicaSet := apiObject.ReplicaSet{
		Base: apiObject.Base{
			ApiVersion: "api/v1",
			Kind:       "ReplicaSet",
			Metadata: apiObject.Metadata{
				Name:      apiFunc.Name,
				Namespace: "function",
			},
		},
		Spec: apiObject.ReplicaSetSpec{
			Replicas: 2,
			Template: apiObject.PodTemplateSpec{
				Spec: apiObject.PodSpec{
					//NodeSelector: apiObject.Labels{
					//	"type": "master",
					//},
					Containers: []apiObject.Container{
						{
							Name:      apiFunc.Name,
							Image:     imageName,
							Resources: apiObject.ContainerResources{},
							Ports: []apiObject.ContainerPort{
								{
									//HostIP:        registry.RegistryHostIP,
									ContainerPort: "8080",
								},
							},
						},
					},
				},
			},
		},
	}

	logManager("Add rs to api-server now")
	URL := url.Prefix + url.ReplicaSetURL
	apiutil.ApplyApiObjectToApiServer(URL, replicaSet)

	c.functionReplicaSetMap[apiFunc.Name] = &functionReplicaSet{
		Function:        apiFunc.Name,
		NumRequest:      0,
		NumReplicas:     2,
		LastRequestTime: time.Now(),
	}
}

// createReplicaSet is coroutine-safe, should be called without lock
func (c *controller) removeReplicaSet(apiFunc *apiObject.Function) {
	URL := url.Prefix + path.Join(url.ReplicaSetURL, "function", apiFunc.Name)
	resp := httputil.DeleteWithoutBody(URL)
	logManager("remove replicaSet and got resp: %s", resp)
	c.scaleLock.Lock()
	defer c.scaleLock.Unlock()
	delete(c.functionReplicaSetMap, apiFunc.Name)
}

func (c *controller) createFunction(apiFunc *apiObject.Function) error {
	c.scaleLock.RLock()
	replicaSet := c.functionReplicaSetMap[apiFunc.Name]
	c.scaleLock.RUnlock()
	if replicaSet == nil {
		logManager("Now create replica set")
		if err := function.CreateFunctionImage(apiFunc.Name, apiFunc.Path); err != nil {
			return err
		} else {
			logManager("Now should create replicaSet")
			c.createReplicaSet(apiFunc)
		}
	}
	return nil
}

func (c *controller) deleteFunction(apiFunc *apiObject.Function) error {
	c.removeReplicaSet(apiFunc)
	return nil
}

func (c *controller) updateFunction(apiFunc *apiObject.Function) error {
	_ = c.deleteFunction(apiFunc)
	return c.createFunction(apiFunc)
}

func (c *controller) createWorkflowWorker(workflow *apiObject.Workflow) error {
	c.workerLock.Lock()
	defer c.workerLock.Unlock()
	fullName := path.Join(workflow.Namespace(), workflow.Name())
	if _, exists := c.workers[fullName]; exists {
		return fmt.Errorf("worker already exists")
	} else {
		ctx, cancel := context.WithCancel(bgCtx)
		w := NewWorker(ctx, workflow, cancel, c.TriggerFunc, c.resultChan)
		c.workers[fullName] = w
		go w.Run()
	}
	return nil
}

func (c *controller) removeWorkflowWorker(fullName string) error {
	c.workerLock.Lock()
	defer c.workerLock.Unlock()
	if w, exists := c.workers[fullName]; !exists {
		return fmt.Errorf("worker does not exist")
	} else {
		w.Cancel()
		delete(c.workers, fullName)
		logManager("remove worker %s successfully", fullName)
	}
	return nil
}
