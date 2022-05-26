package trigger

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/util/httputil"
	"net/http"
	"strings"
)

func getFuncPods(funcName string) []*entity.PodStatus {
	URL := url.Prefix + strings.Replace(url.FuncPodsURLWithSpecifiedName, ":name", funcName, -1)

	var pods []*entity.PodStatus
	err := httputil.GetAndUnmarshal(URL, pods)
	if err != nil {
		return nil
	}
	fmt.Printf("Got pods: %v\n", pods)
	return pods
}

func Trigger(function string, data entity.FunctionData) (result entity.FunctionData, err error) {
	pods := getFuncPods(function)
	n := len(pods)
	if n == 0 {
		return "", fmt.Errorf("no available function instances")
	}

	randomIdx := rand.Intn(n)
	pod := pods[randomIdx]
	ip := pod.Ip
	URL := fmt.Sprintf("http://%s:8080", ip)

	var resp *http.Response
	if resp, err = httputil.PostString(URL, string(data)); err == nil {
		var content []byte
		if content, err = ioutil.ReadAll(resp.Body); err == nil {
			defer resp.Body.Close()
			return entity.FunctionData(content), nil
		}
	}
	return "", err
}
