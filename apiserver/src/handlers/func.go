package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"minik8s/apiObject"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/util/topicutil"
	"net/http"
	"path"
)

func HandleApplyFunc(c *gin.Context) {
	apiFunc := apiObject.Function{}
	if err := readAndUnmarshal(c.Request.Body, &apiFunc); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	etcdURL := path.Join(url.FuncURL, apiFunc.Name)
	if raw, err := etcd.Get(etcdURL); err == nil {
		oldFunc := apiObject.Function{}
		if err = json.Unmarshal([]byte(raw), &oldFunc); err == nil {
			c.String(http.StatusOK, fmt.Sprintf("function %s already exists", apiFunc.Name))
			return
		}
	}

	functionJson, _ := json.Marshal(apiFunc)
	if err := etcd.Put(etcdURL, string(functionJson)); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	topic := topicutil.FunctionUpdateTopic()
	updateMsg, _ := json.Marshal(entity.FunctionUpdate{
		Action: entity.CreateAction,
		Target: apiObject.Function{
			Name: apiFunc.Name,
			Path: apiFunc.Path,
		},
	})

	listwatch.Publish(topic, updateMsg)
	c.String(http.StatusOK, "ok")
}
