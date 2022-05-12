package httputil

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func PostJson(URL string, content interface{}) (*http.Response, error) {
	cli := http.Client{}
	b, _ := json.Marshal(content)
	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	return cli.Do(req)
}
