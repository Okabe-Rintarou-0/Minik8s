package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"minik8s/apiObject"
	"minik8s/apiserver/src/url"
	"minik8s/util/httputil"
	"minik8s/util/logger"
	"strconv"
	"strings"
)

var overwrite bool

var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "Kubectl label is used to change labels of a given api object",
	Long:  "Kubectl label is used to change labels of a given api object",
	Args:  cobra.MinimumNArgs(3),
	Run:   label,
}

func labelSpecifiedNode(fullName string, labels apiObject.Labels) {
	namespace, name := parseName(fullName)
	specified := strings.Replace(url.NodeLabelsURLWithSpecifiedName, ":namespace", namespace, -1)
	specified = strings.Replace(specified, ":name", name, -1)
	URL := url.Prefix + specified + "?overwrite=" + strconv.FormatBool(overwrite)
	resp, err := httputil.PostJson(URL, labels)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	body := resp.Body
	defer body.Close()
	content, err := ioutil.ReadAll(body)
	if err != nil {
		logger.Error(err.Error())
	}
	fmt.Printf("Label node %s with %v and get resp: %s\n", name, labels, string(content))
}

func parseLabels(args []string) apiObject.Labels {
	labels := make(apiObject.Labels)
	for _, arg := range args {
		keyAndValue := strings.Split(arg, "=")
		if len(keyAndValue) == 2 {
			key, value := keyAndValue[0], keyAndValue[1]
			labels[key] = value
		}
	}
	return labels
}

func label(cmd *cobra.Command, args []string) {
	labelTarget := strings.ToLower(args[0])
	name := args[1]
	switch labelTarget {
	// we only support nodes now
	case "nodes":
		labels := parseLabels(args[2:])
		labelSpecifiedNode(name, labels)
		if overwrite {
			fmt.Println("Should overwrite!")
		}
		return
	default:
		fmt.Println("Unsupported api object!")
		return
	}
}
