package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"minik8s/apiObject"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"net/http"
	"path"
	"strconv"
)

func HandleLabelNode(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	body := c.Request.Body
	overwrite, _ := strconv.ParseBool(c.Query("overwrite"))
	if overwrite {
		log("Add labels with overwrite")
	}

	labels := apiObject.Labels{}
	if err := readAndUnmarshal(body, &labels); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	log("Add labels %v to node[hostname = %v]", labels, name)

	node := apiObject.Node{}
	etcdNodeURL := path.Join(url.NodeURL, namespace, name)
	var err error
	var raw string
	if raw, err = etcd.Get(etcdNodeURL); err == nil {
		log("got %s from etcd", raw)
		if err = json.Unmarshal([]byte(raw), &node); err == nil {
			nodeLabels := node.Labels()
			if nodeLabels == nil {
				nodeLabels = make(apiObject.Labels)
			}
			for key, value := range labels {
				if _, exists := nodeLabels[key]; !exists || overwrite {
					nodeLabels[key] = value
				}
			}
			nodeJson, _ := json.Marshal(node)
			if err = etcd.Put(etcdNodeURL, string(nodeJson)); err == nil {
				c.String(http.StatusOK, "ok")

				etcdNodeStatusURL := path.Join(url.NodeURL, "status", node.Namespace(), node.Name())
				if nodeStatusStr, err := etcd.Get(etcdNodeStatusURL); err == nil {
					nodeStatus := entity.NodeStatus{}
					if err = json.Unmarshal([]byte(nodeStatusStr), &nodeStatus); err == nil {
						nodeStatus.Labels = nodeLabels.DeepCopy()
						nodeStatusJson, _ := json.Marshal(nodeStatus)
						_ = etcd.Put(etcdNodeStatusURL, string(nodeStatusJson))
					}
				}
				return
			}
		}
	}
	fmt.Printf("no such node %s/%s\n", node.Namespace(), node.Name())
	c.String(http.StatusOK, err.Error())
}
