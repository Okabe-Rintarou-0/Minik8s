package handlers

import (
	"github.com/gin-gonic/gin"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/url"
	"net/http"
	"path"
)

func deleteSpecifiedNode(namespace, name string) (err error) {
	log("Node to delete is %s/%s", namespace, name)
	etcdNodeURL := path.Join(url.NodeURL, namespace, name)
	if err = etcd.Delete(etcdNodeURL); err == nil {
		etcdNodeStatusURL := path.Join(url.NodeURL, namespace, "status", name)
		err = etcd.Delete(etcdNodeStatusURL)
	}
	return
}

func deleteSpecifiedPod(name string) {
	log("Pod to delete is %s", name)
}

func HandleDeleteNode(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if err := deleteSpecifiedNode(namespace, name); err != nil {
		c.String(http.StatusOK, err.Error())
	}
	c.String(http.StatusOK, "Delete successfully")
}

func HandleDeletePod(c *gin.Context) {
	name := c.Param("name")
	deleteSpecifiedPod(name)
	c.String(http.StatusOK, "Delete successfully")
}
