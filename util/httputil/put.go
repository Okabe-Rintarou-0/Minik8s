package httputil

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func PutForm(URL string, form map[string]string) string {
	cli := http.Client{}
	formJson, _ := json.Marshal(form)
	r := bytes.NewReader(formJson)
	req, _ := http.NewRequest(http.MethodPut, URL, r)

	resp, _ := cli.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func PutJson(URL string, v interface{}) (*http.Response, error) {
	cli := http.Client{}
	vJson, _ := json.Marshal(v)
	r := bytes.NewReader(vJson)
	req, _ := http.NewRequest(http.MethodPut, URL, r)

	return cli.Do(req)
}
