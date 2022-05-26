package apiutil

import (
	"fmt"
	"io/ioutil"
	"minik8s/util/httputil"
)

func ApplyApiObjectToApiServer(URL string, object interface{}) {
	resp, err := httputil.PostJson(URL, object)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(string(respBody))
}
