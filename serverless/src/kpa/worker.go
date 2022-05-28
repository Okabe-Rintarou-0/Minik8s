package kpa

import (
	"context"
	"minik8s/apiObject"
	"minik8s/entity"
	"minik8s/util/logger"
	"time"
)

var logWorker = logger.Log("kpa worker")

type TriggerFn func(string, entity.FunctionData) (entity.FunctionData, error)

const numRetries = 5

type worker struct {
	Target     *apiObject.Workflow
	DAG        *dag
	Ctx        context.Context
	Cancel     context.CancelFunc
	TriggerFn  TriggerFn
	ResultChan chan *entity.FunctionTriggerResult
}

func (w *worker) Run() {
	w.doJob()
}

func (w *worker) finishTask(function string, data entity.FunctionData) (result entity.FunctionData, err error) {
	restRetries := numRetries
	for restRetries > 0 {
		result, err = w.TriggerFn(function, data)
		if err != nil {
			restRetries -= 1
			logWorker("Trigger function %s failed, wait for 5 secs and retry(rest retry num: %d)", function, restRetries)
			time.Sleep(time.Second * 5)
		} else {
			break
		}
	}
	return
}

func (w *worker) doJob() {
	if w.DAG == nil {
		logWorker("No available dag")
		return
	}
	curNode := w.DAG.Root
	var data = entity.FunctionData(w.DAG.EntryParams)
	var err error
	for curNode != nil {
		select {
		case <-w.Ctx.Done():
			logWorker("Canceled, worker stop working on workflow[ID = %s]", w.Target.UID())
			return
		default:
		}
		logWorker("last function's result is: %+v", data)
		switch curNode.Type {
		case taskType:
			// this is a task node, so it should call the function instance in w.finishTask
			logWorker("current node is a task node: %s", curNode.Function)
			data, err = w.finishTask(curNode.Function, data)
			logWorker("finished job and got result: %s", data)
			result := &entity.FunctionTriggerResult{
				Data:              data,
				Time:              time.Now(),
				WorkflowNamespace: w.Target.Namespace(),
				WorkflowName:      w.Target.Name(),
				FinishedAll:       false,
			}
			if err != nil {
				result.Error = err.Error()
				result.Status = entity.TriggerError
				logWorker("error occurs: %s, return", err.Error())
				return
			} else {
				result.Error = ""
				result.Status = entity.TriggerSucceed
			}
			logWorker("Task finished successfully, no error, result = %s", data)
			logWorker("send result %+v to result channel", result)
			w.ResultChan <- result
			curNode = curNode.Next
		case choiceType:
			logWorker("current node is a choice node:")
			// choose branch and go to this branch
			curNode = curNode.ChooseSatisfied(data)
		}
	}
	w.ResultChan <- &entity.FunctionTriggerResult{
		WorkflowNamespace: w.Target.Namespace(),
		WorkflowName:      w.Target.Name(),
		Data:              data,
		Time:              time.Now(),
		Status:            entity.TriggerSucceed,
		Error:             "",
		FinishedAll:       true,
	}
}

func NewWorker(ctx context.Context, target *apiObject.Workflow, cancel context.CancelFunc, triggerFn TriggerFn, resultChan chan *entity.FunctionTriggerResult) *worker {
	return &worker{
		Target:     target,
		DAG:        Workflow2DAG(target),
		Ctx:        ctx,
		Cancel:     cancel,
		TriggerFn:  triggerFn,
		ResultChan: resultChan,
	}
}
