package knative

import (
	"minik8s/serverless/src/kpa"
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

func (kn *Knative) Run() {
	kn.kpaController.Run()
	wait.Forever()
}
