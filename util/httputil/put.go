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
