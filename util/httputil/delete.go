package httputil

import (
	"io/ioutil"
	"net/http"
)

func DeleteWithoutBody(URL string) string {
	req, err := http.NewRequest(http.MethodDelete, URL, nil)
	if err != nil {
		return err.Error()
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}
	return string(body)
}
