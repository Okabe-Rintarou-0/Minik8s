package cmd

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
	"io/ioutil"
	"minik8s/apiObject"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/util/colorwrapper"
	"minik8s/util/httputil"
	"net/http"
	"os"
	"path"
	"strings"
)

var downloadFile string

var directory string

var gpuCmd = &cobra.Command{
	Use:   "gpu",
	Short: "Kubectl gpu is used to get gpu task results",
	Long:  "Kubectl gpu is used to get gpu task results",
	Args:  cobra.ExactArgs(1),
	Run:   handleGpu,
}

func listFiles(URL string) {
	resp, err := http.Get(url.HttpScheme + URL)
	if err != nil {
		fmt.Println(err)
		return
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
	fmt.Println()
}

func downloadFiles(baseURL string) {
	fileURL := url.HttpScheme + path.Join(baseURL, downloadFile)
	if resp, err := http.Get(fileURL); err == nil {
		if content, err := ioutil.ReadAll(resp.Body); err == nil {
			defer resp.Body.Close()
			URL := path.Join(directory, downloadFile)
			if file, err := os.Create(URL); err == nil {
				w := bufio.NewWriter(file)
				_, _ = w.Write(content)
			}
		}
	}
}

func handleGpu(cmd *cobra.Command, args []string) {
	fullName := args[0]
	namespace, name := parseName(fullName)
	gpuJob := apiObject.GpuJob{}
	gpuURL := url.Prefix + path.Join(url.GpuURL, namespace, name)
	var err error
	if err = httputil.GetAndUnmarshal(gpuURL, &gpuJob); err == nil {
		pod := entity.PodStatus{}
		podURL := url.Prefix + path.Join(url.PodURL, "status", gpuJob.Namespace(), gpuJob.Name())
		if err = httputil.GetAndUnmarshal(podURL, &pod); err == nil {
			// get node ip
			scheduledNode := pod.Node
			nodeNamespace, nodeName := parseName(scheduledNode)
			nodeURL := url.Prefix + path.Join(url.NodeURL, "status", nodeNamespace, nodeName)

			// get node ip
			node := entity.NodeStatus{}
			ip := "127.0.0.1"
			if err = httputil.GetAndUnmarshal(nodeURL, &node); err == nil {
				ip = node.Ip
			}

			port := "80"
			for podPort, portBinding := range pod.PortBindings {
				if podPort.Port() == "80" {
					port = portBinding[0].HostPort
				}
			}
			URL := path.Join(ip+":"+port, "files") + "/"
			if downloadFile == "" {
				listFiles(URL)
			} else {
				downloadFiles(URL)
			}
		}
	}
	if err != nil {
		fmt.Println(err.Error())
	}
}
