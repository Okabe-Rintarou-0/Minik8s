package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"minik8s/apiObject"
	"net/http"
)

func HandleLabelNode(c *gin.Context) {
	name := c.Param("name")
	body := c.Request.Body

	defer body.Close()
	content, err := ioutil.ReadAll(body)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	labels := &apiObject.Labels{}
	err = json.Unmarshal(content, labels)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	log("Add labels %v to node[hostname = %v]", labels, name)
}
