package kpa

import (
	"context"
	"minik8s/apiObject"
	"minik8s/entity"
	"minik8s/util/logger"
)

var logWorker = logger.Log("kpa worker")

type worker struct {
	Target  *apiObject.Workflow
	DAG     *dag
	Ctx     context.Context
	MsgChan chan *entity.FunctionMsg
	Cancel  context.CancelFunc
}

func (w *worker) Run() {
	w.doJob()
}

func (w *worker) finishTask(data entity.FunctionData) {

}

func (w *worker) fetchFuncData() entity.FunctionData {
	return ""
}

func (w *worker) handleFunctionMsg(msg *entity.FunctionMsg) bool {
	return true
}

func (w *worker) doJob() {
	if w.DAG == nil {
		logWorker("No available dag")
		return
	}
	curNode := w.DAG.Root
	for curNode != nil {
		select {
		case msg := <-w.MsgChan:
			if !w.handleFunctionMsg(msg) {
				logWorker("Handle msg finished, Worker should stop working on workflow[ID = %s]", w.Target.UID())
				return
			}
		case <-w.Ctx.Done():
			logWorker("Canceled, worker stop working on workflow[ID = %s]", w.Target.UID())
			return
		}

		// data is the params and results of a function
		// it is in a form of json and its type is map[string]interface{}
		// interface{} means it could be any type, mapping string -> any type of value
		// (Because all type can be converted to interface{} type)
		// You can convert interface{} to type T, by using v.(type) expression
		// But for numeric, interface can only be converted to float64
		// If you want an integer, try int(v.(float64))
		data := w.fetchFuncData()
		logWorker("last function's result is: %+v", data)
		switch curNode.Type {
		case taskType:
			// this is a task node, so it should call the function instance in w.finishTask
			logWorker("current node is a task node: %s", curNode.Function)
			w.finishTask(data)
			curNode = curNode.Next
		case choiceType:
			logWorker("current node is a choice node:")
			// choose branch and go to this branch
			curNode = curNode.ChooseSatisfied(data)
		}
	}
}

func NewWorker(ctx context.Context, target *apiObject.Workflow, cancel context.CancelFunc) *worker {
	return &worker{
		Target:  target,
		DAG:     Workflow2DAG(target),
		Ctx:     ctx,
		MsgChan: make(chan *entity.FunctionMsg),
		Cancel:  cancel,
	}
}
