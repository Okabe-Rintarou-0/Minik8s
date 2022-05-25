package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"minik8s/apiObject"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/url"
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

	node := apiObject.Node{}
	etcdURL := path.Join(url.NodeURL, namespace, name)
	var err error
	var raw string
	if raw, err = etcd.Get(etcdURL); err == nil {
		if err = json.Unmarshal([]byte(raw), &node); err == nil {
			nodeLabels := node.Labels()
			for key, value := range labels {
				if _, exists := nodeLabels[key]; !exists || overwrite {
					nodeLabels[key] = value
				}
			}
			nodeJson, _ := json.Marshal(node)
			if err = etcd.Put(etcdURL, string(nodeJson)); err == nil {
				c.String(http.StatusOK, "ok")
				log("Add labels %v to node[hostname = %v]", labels, name)
				return
			}
		}
	}
	fmt.Printf("no such node %s/%s\n", node.Namespace(), node.Name())
	c.String(http.StatusOK, err.Error())
}
