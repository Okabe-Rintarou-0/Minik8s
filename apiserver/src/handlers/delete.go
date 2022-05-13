package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func deleteSpecifiedPod(name string) {
	log("Pod to delete is %s", name)
}

func HandleDeletePod(c *gin.Context) {
	name := c.Param("name")
	deleteSpecifiedPod(name)
	c.String(http.StatusOK, "Delete successfully")
}
