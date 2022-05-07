package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"minik8s/entity"
	"minik8s/kubectl/src/util"
	"strings"
	"time"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
	}
}

var filePath string

func init() {
	applyCmd.Flags().StringVarP(&filePath, "filePath", "f", "", "filePath of api object yaml file")
	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(describeCmd)
}

var rootCmd = &cobra.Command{
	Use:   "kubectl",
	Short: "Kubectl is for better control of minik8s",
	Long: `By using kubectl, you can create api object in minik8s, or know details of them by using kubectl describe command.
For example: kubectl apply -f ./example.yaml; kubectl describe pod examplePod`,
	Run: runRoot,
}

func runRoot(cmd *cobra.Command, args []string) {
	// Reach here if no args
	fmt.Printf("execute %s args:%v \n", cmd.Name(), args)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Kubectl apply is used to create api object in a declarative way",
	Long:  "Kubectl apply is used to create api object in a declarative way",
	Run:   apply,
}

func apply(cmd *cobra.Command, args []string) {
	content, err := util.LoadContent(filePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tp, err := util.ParseApiObjectType(content)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	switch tp {
	case util.Pod:
		pod, err := util.ParsePod(content)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("Get pod %v\n", *pod)
	}
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Kubectl get is used to get brief information of the api object with given unique name",
	Long:  "Kubectl get is used to get brief information of the api object with given unique name",
	// exactly two args, one is the type of api object, another is the unique name
	// for example, kubectl get pod example, where pod is the type and example is the unique name
	Args: cobra.MinimumNArgs(1),
	Run:  get,
}

func podStatusForTest() *entity.PodStatus {
	return &entity.PodStatus{
		ID:        uuid.NewV4().String(),
		Name:      "example",
		Labels:    nil,
		Namespace: "default",
		Status:    entity.ContainerCreating,
		SyncTime:  time.Now(),
	}
}

func getPodFromApiServer(name string) *entity.PodStatus {
	//TODO just for test now, replace it with api-server
	return podStatusForTest()
}

func getPodsFromApiServer() []*entity.PodStatus {
	//TODO just for test now, replace it with api-server
	pod1 := podStatusForTest()
	pod2 := podStatusForTest()
	pod2.Status = entity.Error
	pod2.Error = "Image Pull Error"
	return []*entity.PodStatus{pod1, pod2}
}

func podDescriptionForTest(name string) *entity.PodDescription {
	//TODO just for test now, replace it with api-server

	var logs []entity.PodStatusLogEntry

	podStatus := podStatusForTest()
	logs = append(logs, entity.PodStatusLogEntry{
		Status: podStatus.Status,
		Time:   podStatus.SyncTime,
		Error:  podStatus.Error,
	})

	logs = append(logs, entity.PodStatusLogEntry{
		Status: entity.Running,
		Time:   time.Now().Add(time.Minute * 30),
		Error:  "",
	})

	return &entity.PodDescription{
		CurrentStatus: *podStatusForTest(),
		Logs:          logs,
	}
}

func getPodDescriptionFromApiServer(name string) *entity.PodDescription {
	//TODO just for test now, replace it with api-server
	podDesc := podDescriptionForTest(name)
	return podDesc
}

func podStatusTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Name", "UID", "Status", "Last Sync Time", "Error")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}

func podStatusLogTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Time", "Status", "Error")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}

func printSpecifiedPodStatus(name string) error {
	podStatus := getPodFromApiServer(name)
	if podStatus == nil {
		return fmt.Errorf("no such pod")
	}

	tbl := podStatusTbl()
	tbl.AddRow(podStatus.Name, podStatus.ID, podStatus.Status.String(), podStatus.SyncTime.Format(time.RFC3339), podStatus.Error)
	tbl.Print()
	return nil
}

func printSpecifiedPodDescription(name string) error {
	podDesc := getPodDescriptionFromApiServer(name)
	if podDesc == nil {
		return fmt.Errorf("no such pod")
	}

	logs := podDesc.Logs
	tbl := podStatusLogTbl()
	fmt.Println("History log:")
	for _, log := range logs {
		tbl.AddRow(log.Time.Format(time.RFC3339), log.Status.String(), log.Error)
	}
	tbl.Print()

	fmt.Println("Current status:")
	podStatus := podDesc.CurrentStatus
	tbl = podStatusTbl()
	tbl.AddRow(podStatus.Name, podStatus.ID, podStatus.Status.String(), podStatus.SyncTime.Format(time.RFC3339), podStatus.Error)
	tbl.Print()

	return nil
}

func printPodStatuses() error {
	podStatuses := getPodsFromApiServer()

	tbl := podStatusTbl()
	for _, podStatus := range podStatuses {
		tbl.AddRow(podStatus.Name, podStatus.ID, podStatus.Status.String(), podStatus.SyncTime.Format(time.RFC3339), podStatus.Error)
	}
	tbl.Print()
	return nil
}

func get(cmd *cobra.Command, args []string) {
	apiObjectType := args[0]
	var name string
	if len(args) > 1 {
		name = args[1]
	}
	apiObjectType = strings.ToLower(apiObjectType)
	if !util.IsValidApiObjectType(apiObjectType) {
		fmt.Printf("invalid api object type \"%s\", acceptable api object type is pod, service, etc.", apiObjectType)
		return
	}

	var err error
	switch apiObjectType {
	case "pod":
		err = printSpecifiedPodStatus(name)
	case "pods":
		err = printPodStatuses()
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	//fmt.Printf("Get information of a %s named %s\n", apiObjectType, name)
}

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Kubectl describe is used to get detailed information of the api object with given unique name",
	Long:  "Kubectl describe is used to get detailed information of the api object with given unique name",
	// exactly two args, one is the type of api object, another is the unique name
	// for example, kubectl describe pod example, where pod is the type and example is the unique name
	Args: cobra.ExactValidArgs(2),
	Run:  describe,
}

func describe(cmd *cobra.Command, args []string) {
	apiObjectType := args[0]
	name := args[1]

	var err error
	switch apiObjectType {
	case "pod":
		err = printSpecifiedPodDescription(name)
	default:
		fmt.Println("Invalid api object type!")
		return
	}

	if err != nil {
		fmt.Println(err.Error())
	}
	//fmt.Printf("Describe information of a %s named %s\n", apiObjectType, name)
}
