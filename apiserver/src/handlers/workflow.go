package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/util/httputil"
	"net/http"
	"path"
)

func HandlePutWorkflowResult(c *gin.Context) {
	result := entity.FunctionTriggerResult{}
	if err := httputil.ReadAndUnmarshal(c.Request.Body, &result); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	} else {
		etcdURL := path.Join(url.WorkflowURL, "result", result.WorkflowNamespace, result.WorkflowName)
		log("Receive workflow result: %+v", result)
		resultJson, _ := json.Marshal(result)
		if err = etcd.Put(etcdURL, string(resultJson)); err != nil {
			c.String(http.StatusOK, err.Error())
			return
		} else {
			c.String(http.StatusOK, "ok")
		}
	}
}
