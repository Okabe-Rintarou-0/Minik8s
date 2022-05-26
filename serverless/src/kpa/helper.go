package kpa

import (
	"context"
	"fmt"
	"minik8s/apiObject"
	"minik8s/apiserver/src/url"
	"minik8s/serverless/src/function"
	"minik8s/util/httputil"
	"path"
	"strconv"
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
	c.updateReplicaSetToApiServer(funcReplicaSet)
}

func (c *controller) createFunction(apiFunc *apiObject.Function) error {
	if _, exists := c.functionReplicaSetMap[apiFunc.Name]; !exists {
		if err := function.InitFunction(apiFunc.Name, apiFunc.Name); err != nil {
			return err
		}
	}
	return nil
}

func (c *controller) createWorkflowWorker(workflow *apiObject.Workflow) error {
	c.workerLock.Lock()
	defer c.workerLock.Unlock()
	if _, exists := c.workers[workflow.UID()]; exists {
		return fmt.Errorf("worker already exists")
	} else {
		ctx, cancel := context.WithCancel(bgCtx)
		w := NewWorker(ctx, workflow, cancel)
		c.workers[workflow.UID()] = w
		go w.Run()
	}
	return nil
}

func (c *controller) removeWorkflowWorker(workflow *apiObject.Workflow) error {
	c.workerLock.Lock()
	defer c.workerLock.Unlock()
	if w, exists := c.workers[workflow.UID()]; !exists {
		return fmt.Errorf("worker does not exist")
	} else {
		w.Cancel()
		delete(c.workers, workflow.UID())
	}
	return nil
}
