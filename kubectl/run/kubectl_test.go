package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"minik8s/util/colorwrapper"
	"net/http"
	"path"
	"strings"
	"testing"
)

func TestNginxFileServer(t *testing.T) {
	fmt.Println(path.Join("123", "123"))
	URL := "http://localhost:8000/files/"
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println(err)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	var files []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		files = append(files, s.Text())
	})

	for _, file := range files {
		if strings.HasSuffix(file, "/") {
			fmt.Printf("%s ", colorwrapper.Green(file[0:len(file)-1]))
		} else {
			fmt.Printf("%s ", file)
		}
	}
}
