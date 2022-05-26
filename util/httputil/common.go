package httputil

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func ReadAndUnmarshal(body io.ReadCloser, target interface{}) error {
	content, err := ioutil.ReadAll(body)
	defer body.Close()
	if err != nil {
		return err
	}
	if err = json.Unmarshal(content, target); err != nil {
		return err
	}
	return nil
}
