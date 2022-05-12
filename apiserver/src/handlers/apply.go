package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"minik8s/apiObject"
)

func HandleApplyPod(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	pod := &apiObject.Pod{}
	_ = json.Unmarshal(body, pod)
	fmt.Printf("receive pod[ID = %v]: %v\n", pod.UID(), pod)
}
