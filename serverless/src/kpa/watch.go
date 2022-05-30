package kpa

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"minik8s/entity"
	"minik8s/util/logger"
	"path"
)

func (c *controller) handleFunctionUpdate(msg *redis.Message) {
	functionUpdate := entity.FunctionUpdate{}
	if err := json.Unmarshal([]byte(msg.Payload), &functionUpdate); err != nil {
		logger.Error(err.Error())
		return
	}
	logManager("Receive %s function: %s(%s)", functionUpdate.Action.String(), functionUpdate.Target.Name, functionUpdate.Target.Path)
	apiFunc := functionUpdate.Target

	var err error
	switch functionUpdate.Action {
	case entity.CreateAction:
		err = c.createFunction(&apiFunc)
	case entity.DeleteAction:
		err = c.deleteFunction(&apiFunc)
	case entity.UpdateAction:
		err = c.updateFunction(&apiFunc)
	}


	if err != nil {
		logger.Error(err.Error())
	}
}

func (c *controller) handleWorkflowUpdate(msg *redis.Message) {
	workflowUpdate := entity.WorkflowUpdate{}
	if err := json.Unmarshal([]byte(msg.Payload), &workflowUpdate); err != nil {
		logger.Error(err.Error())
		return
	}
	logManager("Receive %s workflow: %s", workflowUpdate.Action.String(), workflowUpdate.Target.Name())

	wf := workflowUpdate.Target

	var err error
	switch workflowUpdate.Action {
	case entity.CreateAction:
		err = c.createWorkflowWorker(&wf)
	case entity.DeleteAction:
		err = c.removeWorkflowWorker(path.Join(wf.Namespace(), wf.Name()))
	}

	if err != nil {
		logger.Error(err.Error())
	}
}
