package replicaSet

import (
	"minik8s/apiserver/src/url"
	"minik8s/util/httputil"
	"path"
)

func deletePodToApiServer(namespace, name string) {
	logWorker("Pod to delete is Pod[%s/%s]", namespace, name)
	URL := url.Prefix + path.Join(url.PodURL, namespace, name)
	resp := httputil.DeleteWithoutBody(URL)
	logWorker("Delete pod %s/%s and get resp: %s", namespace, name, resp)
}
