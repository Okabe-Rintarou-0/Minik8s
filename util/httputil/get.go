package httputil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetAndUnmarshal(URL string, target interface{}) error {
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	var content []byte
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, target)
}
