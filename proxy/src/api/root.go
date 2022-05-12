package api

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var ClusterId string
var Namespace string
var EtcdServer string

func init() {
	RootCmd.PersistentFlags().StringVarP(&ClusterId, "cluster-id", "i", "", "cluster instance id.")
	RootCmd.PersistentFlags().StringVarP(&Namespace, "namespace", "n", "", "namespace name.")
	RootCmd.PersistentFlags().StringVarP(&EtcdServer, "etcd-server", "e", "", "etcd server.")
}

var RootCmd = &cobra.Command{
	Use:   "kube-apiserver get/delete pod/job",
	Short: "get all the job/pod in the namesapce, or delete job/pod that started before [minutes] ago or has been completed.",
	Long:  "get all the job/pod in the namespace, or delete job/pod that started before [minutes] ago or has been completed.",
	Example: `kube-apiserver get pod -i cls-xxx -n default -e http://1.2.3.4:2379 (list pod in the default namespace)
etcd-tool delete job -i cls-xxx -n default -e http://1.2.3.4:2379 -c -m 600 (delete jobs that have been completed and started before 10h ago)`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
