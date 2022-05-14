package handlers

import (
	"github.com/gin-gonic/gin"
	"minik8s/entity"
	"strconv"
)

func HandleSetNodeStatus(c *gin.Context) {
	hostname := c.Param("name")
	lifecycleInt64, _ := strconv.Atoi(c.PostForm("lifecycle"))
	lifecycle := entity.NodeLifecycle(lifecycleInt64)
	log("Received lifecycle %v from %v", lifecycle.String(), hostname)
}
