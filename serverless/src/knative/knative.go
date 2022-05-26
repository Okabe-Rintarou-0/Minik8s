package knative

import (
	"fmt"
	"minik8s/serverless/src/kpa"
	"minik8s/util/recoverutil"
	"minik8s/util/wait"
)

type Knative struct {
	kpaController kpa.Controller
}

func NewKnative() *Knative {
	return &Knative{
		kpaController: kpa.NewController(),
	}
}

func (kn *Knative) recover() {
	if err := recover(); err != nil {
		fmt.Println(recoverutil.Trace(fmt.Sprintf("%v\n", err)))
	}
}

func (kn *Knative) Run() {
	defer kn.recover()
	kn.kpaController.Run()
	wait.Forever()
}
