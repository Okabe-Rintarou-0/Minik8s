package httputil

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
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

func PostForm(URL string, form map[string]string) string {
	values := url.Values{}
	for key, value := range form {
		values.Add(key, value)
	}

	var err error
	var resp *http.Response
	if resp, err = http.PostForm(URL, values); err == nil {
		defer resp.Body.Close()
		var body []byte
		if body, err = ioutil.ReadAll(resp.Body); err == nil {
			return string(body)
		}
	}
	return err.Error()
}
