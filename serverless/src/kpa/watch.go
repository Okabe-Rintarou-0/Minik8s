package kpa

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"minik8s/entity"
	"minik8s/util/logger"
)

type functionExecResult struct {
}

func (c *controller) handleFunctionUpdate(msg *redis.Message) {
	functionUpdate := entity.FunctionUpdate{}
	_ = json.Unmarshal([]byte(msg.Payload), &functionUpdate)

	apiFunc := functionUpdate.Target

	var err error
	switch functionUpdate.Action {
	case entity.CreateAction:
		err = c.createFunction(&apiFunc)
	case entity.DeleteAction:

	}

	if err != nil {
		logger.Error(err.Error())
	}
}

func (c *controller) handleWorkflowUpdate(msg *redis.Message) {
	workflowUpdate := entity.WorkflowUpdate{}
	_ = json.Unmarshal([]byte(msg.Payload), &workflowUpdate)

	wf := workflowUpdate.Target

	var err error
	switch workflowUpdate.Action {
	case entity.CreateAction:
		err = c.createWorkflowWorker(&wf)
	case entity.DeleteAction:
		err = c.removeWorkflowWorker(&wf)
	}

	if err != nil {
		logger.Error(err.Error())
	}
}
