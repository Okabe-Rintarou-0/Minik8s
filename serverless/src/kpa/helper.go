package kpa

import (
	"minik8s/apiserver/src/url"
	"minik8s/util/httputil"
	"path"
	"strconv"
)

func (c *controller) updateReplicaSetToApiServer(funcReplicaSet *functionReplicaSet) {
	URL := url.Prefix + path.Join(url.ReplicaSetURL, "function", funcReplicaSet.Function)
	resp := httputil.PutForm(URL, map[string]string{
		"replicas": strconv.Itoa(funcReplicaSet.NumReplicas),
	})
	logManager("update rs and get resp: %s", resp)
}

func (c *controller) scaleToHalf(funcReplicaSet *functionReplicaSet) {
	funcReplicaSet.NumReplicas /= 2
	c.updateReplicaSetToApiServer(funcReplicaSet)
}
