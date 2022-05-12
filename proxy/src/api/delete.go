package api

//
//import (
//	"context"
//	"fmt"
//	clientv3 "github.com/coreos/etcd/clientv3"
//	"github.com/spf13/cobra"
//)
//
//var BeforeMinutes int
//var Completed bool
//
//func init() {
//	deleteCmd.Flags().BoolVarP(&Completed, "completed", "c", true, "for delete, if set, will delete pods/jobs that have been Completed.")
//	deleteCmd.Flags().IntVarP(&BeforeMinutes, "minutes", "m", 0, "for delete, default 0. If > 0, will delete pods/jobs that started before [minutes] minutes ago. If completed and minutes flags are both set, will delete pods/jobs that met both conditions.")
//	RootCmd.AddCommand(deleteCmd)
//}
//
//var deleteCmd = &cobra.Command{
//	Use:   "delete [job or pod]",
//	Short: "delete job/pod that started before [minutes] ago or has been completed.",
//	Long:  "delete job/pod that started before [minutes] ago or has been completed.",
//	Run: func(cmd *cobra.Command, args []string) {
//		fmt.Println(BeforeMinutes)
//		fmt.Println(Completed)
//		if len(args) != 1 || (args[0] != "pod" && args[0] != "job") || ClusterId == "" || Namespace == "" || EtcdServer == "" ||
//			(BeforeMinutes == 0 && Completed == false) {
//			cmd.Help()
//			return
//		}
//		doDelete(args[0])
//	},
//}
//
//func deleteFromEtcd(cli *clientv3.Client, key string) error {
//	if _, err := cli.Delete(context.TODO(), key, clientv3.WithPrefix()); err != nil {
//		return fmt.Errorf("Delete %s failed: %v", err)
//	}
//	return nil
//}
//
//func doDelete(typ string) {
//	cli, err := clientv3.New(clientv3.Config{
//		Endpoints: []string{EtcdServer},
//	})
//	if err != nil {
//		fmt.Printf("new etcd cli failed:%v", err)
//		return
//	}
//	defer cli.Close()
//
//	var etcdNode string
//	if typ == "pod" {
//		etcdNode = "/" + ClusterId + "/" + "pods/" + Namespace
//	} else {
//		etcdNode = "/" + ClusterId + "/" + "jobs/" + Namespace
//	}
//
//	res, err := cli.Get(context.TODO(), etcdNode, clientv3.WithPrefix())
//	if err != nil {
//		fmt.Printf("Get from etcd error:%v", err)
//		return
//	}
//
//	for _, kv := range res.Kvs {
//		if typ == "pod" {
//			pod, err := decodePod(kv.Value)
//			if err != nil {
//				fmt.Println(err)
//				return
//			}
//			if Completed {
//				if !isPodCompleted(pod) {
//					continue
//				}
//			}
//			if BeforeMinutes > 0 {
//				if !podStartedAtMinutesAgo(pod, BeforeMinutes) {
//					continue
//				}
//			}
//		} else {
//			job, err := decodeJob(kv.Value)
//			if err != nil {
//				fmt.Println(err)
//				return
//			}
//			if Completed {
//				if !isJobCompleted(job) {
//					continue
//				}
//			}
//			if BeforeMinutes > 0 {
//				if !jobStartedAtMinutesAgo(job, BeforeMinutes) {
//					continue
//				}
//			}
//		}
//
//		if err := deleteFromEtcd(cli, string(kv.Key)); err != nil {
//			fmt.Printf("Delete %s failed: %v\n", string(kv.Key), err)
//			return
//		} else {
//			fmt.Printf("Delete %s succeded.\n", string(kv.Key))
//		}
//	}
//
//}
