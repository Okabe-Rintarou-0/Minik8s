package knative

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"minik8s/serverless/src/kpa"
	"minik8s/util/recoverutil"
)

type Knative struct {
	kpaController kpa.Controller
	httpServer    *gin.Engine
}

func NewKnative() *Knative {
	return &Knative{
		kpaController: kpa.NewController(),
		httpServer:    gin.Default(),
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
	kn.httpServer.POST("/:function", kn.kpaController.HandleTriggerFunc)
	log.Fatal(kn.httpServer.Run(":8081"))
}
