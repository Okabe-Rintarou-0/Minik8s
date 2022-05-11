package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"minik8s/entity"
	"time"
)

func podStatusTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Name", "UID", "PodLifecycle", "Last Sync Time", "PodError")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}

func podStatusLogTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Time", "PodLifecycle", "PodError")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}

func getPodFromApiServer(name string) *entity.PodStatus {
	//TODO just for test now, replace it with api-server
	return podStatusForTest()
}

func getPodsFromApiServer() []*entity.PodStatus {
	//TODO just for test now, replace it with api-server
	pod1 := podStatusForTest()
	pod2 := podStatusForTest()
	pod2.Lifecycle = entity.PodError
	pod2.Error = "Image Pull PodError"
	return []*entity.PodStatus{pod1, pod2}
}

func printSpecifiedPodStatus(name string) error {
	podStatus := getPodFromApiServer(name)
	if podStatus == nil {
		return fmt.Errorf("no such pod")
	}

	tbl := podStatusTbl()
	tbl.AddRow(podStatus.Name, podStatus.ID, podStatus.Lifecycle.String(), podStatus.SyncTime.Format(time.RFC3339), podStatus.Error)
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
	tbl.AddRow(podStatus.Name, podStatus.ID, podStatus.Lifecycle.String(), podStatus.SyncTime.Format(time.RFC3339), podStatus.Error)
	tbl.Print()

	return nil
}

func printPodStatuses() error {
	podStatuses := getPodsFromApiServer()

	tbl := podStatusTbl()
	for _, podStatus := range podStatuses {
		tbl.AddRow(podStatus.Name, podStatus.ID, podStatus.Lifecycle.String(), podStatus.SyncTime.Format(time.RFC3339), podStatus.Error)
	}
	tbl.Print()
	return nil
}
