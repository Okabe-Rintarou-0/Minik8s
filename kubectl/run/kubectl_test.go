package main

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"minik8s/apiserver/src/url"
	"net/http"
	"os"
	"path"
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
	fileURL := url.HttpScheme + path.Join("localhost:8000/files/", "matrix_add.cu")
	if resp, err := http.Get(fileURL); err == nil {
		if content, err := ioutil.ReadAll(resp.Body); err == nil {
			//fmt.Println(string(content))
			defer resp.Body.Close()
			URL := path.Join("D:/", "matrix_add.cu")
			if file, err := os.Create(URL); err == nil {
				w := bufio.NewWriter(file)
				_, err = w.Write(content)
				_ = w.Flush()
			}
		}
	}
	//for _, file := range files {
	//	if strings.HasSuffix(file, "/") {
	//		blue := color.New(color.FgBlue)
	//		_, _ = blue.Print(file[0:len(file)-1], " ")
	//	} else {
	//		fmt.Printf("%s ", file)
	//	}
	//}
}
