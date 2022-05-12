package api

//
//import (
//	"context"
//	"fmt"
//	clientv3 "github.com/coreos/etcd/clientv3"
//	"github.com/spf13/cobra"
//)
//
//func init() {
//	RootCmd.AddCommand(getCmd)
//}
//
//var getCmd = &cobra.Command{
//	Use:   "get [job or pod]",
//	Short: "get all the job/pod belong to the namespace.",
//	Long:  "get all the job/pod belong to the namespace.",
//	Run: func(cmd *cobra.Command, args []string) {
//		if len(args) != 1 || (args[0] != "pod" && args[0] != "job") || ClusterId == "" || Namespace == "" || EtcdServer == "" {
//			cmd.Help()
//			return
//		}
//		doGet(args[0], ClusterId, Namespace)
//	},
//}
//
//func doGetJob(cli *clientv3.Client, cluster string, namespace string) {
//
//	jobNode := "/" + cluster + "/" + "jobs/" + namespace
//
//	res, err := cli.Get(context.TODO(), jobNode, clientv3.WithPrefix())
//	if err != nil {
//		fmt.Printf("Get from etcd error:%v", err)
//		return
//	}
//
//	fmt.Println("JOB NAME|JOB CREATION TIME|JOB START TIME|JOB COMPLETION TIME|IS JOB COMPLETION")
//	for _, kv := range res.Kvs {
//		job, err := decodeJob(kv.Value)
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//
//		isCompleted := "false"
//		if isJobCompleted(job) {
//			isCompleted = "true"
//		}
//		fmt.Printf("%s|%s|%s|%s|%s\n", job.Name, StringTime(&job.CreationTimestamp), StringTime(job.Status.StartTime), StringTime(job.Status.CompletionTime), isCompleted)
//	}
//
//	return
//}
//
//func doGetPod(cli *clientv3.Client, cluster string, namespace string) {
//
//	podNode := "/" + cluster + "/" + "pods/" + namespace
//
//	res, err := cli.Get(context.TODO(), podNode, clientv3.WithPrefix())
//	if err != nil {
//		fmt.Printf("Get from etcd error:%v", err)
//		return
//	}
//
//	fmt.Println("POD NAME|POD CREATION TIME|POD START TIME|IS POD COMPLETION")
//	for _, kv := range res.Kvs {
//		pod, err := decodePod(kv.Value)
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		isCompleted := "false"
//		if isPodCompleted(pod) {
//			isCompleted = "true"
//		}
//		fmt.Printf("%s|%s|%s|%s\n", pod.Name, StringTime(&pod.CreationTimestamp), StringTime(pod.Status.StartTime), isCompleted)
//	}
//
//	return
//}
//
//func doGet(typ, cluster, namespace string) {
//	cli, err := clientv3.New(clientv3.Config{
//		Endpoints: []string{EtcdServer},
//	})
//	if err != nil {
//		fmt.Printf("new etcd cli failed:%v", err)
//		return
//	}
//	defer cli.Close()
//
//	if typ == "pod" {
//		doGetPod(cli, cluster, namespace)
//	} else {
//		doGetJob(cli, cluster, namespace)
//	}
//}
