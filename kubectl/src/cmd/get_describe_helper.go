package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/util/httputil"
	"time"
)

func podStatusTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Name", "UID", "Status", "Last Sync Time", "Error")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}

func nodeStatusTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Hostname", "Status", "Ipv4", "Cpu", "Memory", "Pods", "Last Sync Time", "Error")
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

func getPodFromApiServer(name string) (pod *entity.PodStatus, err error) {
	err = httputil.GetAndUnmarshal(url.Prefix+url.PodURL+name, &pod)
	return
}

func getPodsFromApiServer() (pods []*entity.PodStatus, err error) {
	err = httputil.GetAndUnmarshal(url.Prefix+url.PodURL, &pods)
	return
}

func getNodesFromApiServer() (nodes []*entity.NodeStatus, err error) {
	err = httputil.GetAndUnmarshal(url.Prefix+url.NodeURL, &nodes)
	return
}

func getNodeFromApiServer(name string) (node *entity.NodeStatus, err error) {
	err = httputil.GetAndUnmarshal(url.Prefix+url.NodeURL+name, &node)
	return
}

func printSpecifiedPodStatus(name string) error {
	podStatus, err := getPodFromApiServer(name)
	if err != nil {
		return err
	}
	if podStatus == nil {
		return fmt.Errorf("no such pod")
	}

	tbl := podStatusTbl()
	tbl.AddRow(podStatus.Name, podStatus.ID, podStatus.Lifecycle.String(), podStatus.SyncTime.Format(time.RFC3339), podStatus.Error)
	tbl.Print()
	return nil
}

func printSpecifiedPodDescription(name string) error {
	podDesc, err := getPodDescriptionFromApiServer(name)
	if err != nil {
		return err
	}
	if podDesc == nil {
		return fmt.Errorf("no such pod")
	}

	logs := podDesc.Logs
	tbl := podStatusLogTbl()
	fmt.Println("History logger:")
	for _, log := range logs {
		tbl.AddRow(log.Time.Format(time.RFC3339), log.Status.String(), log.Error)
	}
	tbl.Print()

	fmt.Println("Current status:")
	podStatus := podDesc.CurrentStatus
	tbl = podStatusTbl()
	tbl.AddRow(podStatus.Name, podStatus.ID, podStatus.Lifecycle.String(), podStatus.SyncTime.Format(time.RFC3339), podStatus.Error)
	tbl.Print()

	return nil
}

func printSpecifiedNodeStatus(name string) error {
	nodeStatus, err := getNodeFromApiServer(name)
	if err != nil {
		return err
	}
	if nodeStatus == nil {
		return fmt.Errorf("no such node")
	}

	tbl := nodeStatusTbl()
	tbl.AddRow(
		nodeStatus.Hostname,
		nodeStatus.Lifecycle.String(),
		nodeStatus.Ip,
		nodeStatus.CpuPercent,
		nodeStatus.MemPercent,
		nodeStatus.NumPods,
		nodeStatus.SyncTime.Format(time.RFC3339),
		nodeStatus.Error,
	)
	tbl.Print()

	return nil
}

func printPodStatuses() error {
	podStatuses, err := getPodsFromApiServer()
	if err != nil {
		return err
	}

	tbl := podStatusTbl()
	for _, podStatus := range podStatuses {
		tbl.AddRow(podStatus.Name, podStatus.ID, podStatus.Lifecycle.String(), podStatus.SyncTime.Format(time.RFC3339), podStatus.Error)
	}
	tbl.Print()
	return nil
}

func printNodeStatuses() error {
	nodeStatuses, err := getNodesFromApiServer()
	if err != nil {
		return err
	}
	tbl := nodeStatusTbl()
	for _, nodeStatus := range nodeStatuses {
		tbl.AddRow(
			nodeStatus.Hostname,
			nodeStatus.Lifecycle.String(),
			nodeStatus.Ip,
			nodeStatus.CpuPercent,
			nodeStatus.MemPercent,
			nodeStatus.NumPods,
			nodeStatus.SyncTime.Format(time.RFC3339),
			nodeStatus.Error,
		)
	}
	tbl.Print()
	return nil
}
