package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"minik8s/apiObject"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/util/httputil"
	"path"
	"strconv"
	"time"
)

func podStatusTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Name", "UID", "Status", "Node", "Cpu", "Memory", "Last Sync Time", "Error")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}

func replicaSetStatusTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Name", "UID", "Status", "Replicas", "Cpu", "Memory", "Last Sync Time", "Error")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}

func hpaStatusTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Name", "UID", "Status", "Min Replicas", "Max Replicas", "Current", "Metrics", "Benchmark", "Last Sync Time", "Error")
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

func getPodFromApiServer(fullName string) (pod *entity.PodStatus, err error) {
	namespace, name := parseName(fullName)
	URL := url.Prefix + path.Join(url.PodURL, "status", namespace, name)
	err = httputil.GetAndUnmarshal(URL, &pod)
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

func getNodeFromApiServer(fullName string) (node *entity.NodeStatus, err error) {
	namespace, name := parseName(fullName)
	URL := url.Prefix + path.Join(url.NodeURL, "status", namespace, name)
	err = httputil.GetAndUnmarshal(URL, &node)
	return
}

func getReplicaSetsFromApiServer() (replicaSets []*entity.ReplicaSetStatus, err error) {
	err = httputil.GetAndUnmarshal(url.Prefix+url.ReplicaSetURL, &replicaSets)
	return
}

func getReplicaSetFromApiServer(fullName string) (replicaSet *entity.ReplicaSetStatus, err error) {
	namespace, name := parseName(fullName)
	URL := url.Prefix + path.Join(url.ReplicaSetURL, "status", namespace, name)
	err = httputil.GetAndUnmarshal(URL, &replicaSet)
	return
}

func getHPAsFromApiServer() (hpas []*entity.HPAStatus, err error) {
	err = httputil.GetAndUnmarshal(url.Prefix+url.HPAURL, &hpas)
	return
}

func getHPAFromApiServer(fullName string) (hpa *entity.HPAStatus, err error) {
	namespace, name := parseName(fullName)
	URL := url.Prefix + path.Join(url.HPAURL, "status", namespace, name)
	err = httputil.GetAndUnmarshal(URL, &hpa)
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
	fullName := path.Join(podStatus.Namespace, podStatus.Name)
	tbl.AddRow(
		fullName,
		podStatus.ID,
		podStatus.Lifecycle.String(),
		podStatus.Node,
		podStatus.CpuPercent,
		podStatus.MemPercent,
		podStatus.SyncTime.Format(time.RFC3339),
		podStatus.Error,
	)
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
	fullName := path.Join(podStatus.Namespace, podStatus.Name)
	tbl.AddRow(fullName, podStatus.ID, podStatus.Lifecycle.String(), podStatus.SyncTime.Format(time.RFC3339), podStatus.Error)
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
		fullName := path.Join(podStatus.Namespace, podStatus.Name)
		tbl.AddRow(
			fullName,
			podStatus.ID,
			podStatus.Lifecycle.String(),
			podStatus.Node,
			podStatus.CpuPercent,
			podStatus.MemPercent,
			podStatus.SyncTime.Format(time.RFC3339),
			podStatus.Error,
		)
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

func printReplicaSetStatuses() error {
	replicaSetStatuses, err := getReplicaSetsFromApiServer()
	if err != nil {
		return err
	}

	tbl := replicaSetStatusTbl()
	for _, replicaSetStatus := range replicaSetStatuses {
		fullName := path.Join(replicaSetStatus.Namespace, replicaSetStatus.Name)
		replicas := strconv.Itoa(replicaSetStatus.NumReady) + "/" + strconv.Itoa(replicaSetStatus.NumReplicas)
		tbl.AddRow(
			fullName,
			replicaSetStatus.ID,
			replicaSetStatus.Lifecycle.String(),
			replicas,
			replicaSetStatus.CpuPercent,
			replicaSetStatus.MemPercent,
			replicaSetStatus.SyncTime.Format(time.RFC3339),
			replicaSetStatus.Error,
		)
	}
	tbl.Print()
	return nil
}

func printSpecifiedReplicaSetStatus(name string) error {
	replicaSetStatus, err := getReplicaSetFromApiServer(name)
	if err != nil {
		return err
	}
	if replicaSetStatus == nil {
		return fmt.Errorf("no such replicaSet")
	}

	tbl := replicaSetStatusTbl()
	fullName := path.Join(replicaSetStatus.Namespace, replicaSetStatus.Name)
	replicas := strconv.Itoa(replicaSetStatus.NumReady) + "/" + strconv.Itoa(replicaSetStatus.NumReplicas)
	tbl.AddRow(
		fullName,
		replicaSetStatus.ID,
		replicaSetStatus.Lifecycle.String(),
		replicas,
		replicaSetStatus.CpuPercent,
		replicaSetStatus.MemPercent,
		replicaSetStatus.SyncTime.Format(time.RFC3339),
		replicaSetStatus.Error,
	)
	tbl.Print()
	return nil
}

func printHPAStatuses() error {
	hpaStatuses, err := getHPAsFromApiServer()
	if err != nil {
		return err
	}

	tbl := hpaStatusTbl()
	for _, hpaStatus := range hpaStatuses {
		fullName := path.Join(hpaStatus.Namespace, hpaStatus.Name)
		replicas := strconv.Itoa(hpaStatus.NumReady) + "/" + strconv.Itoa(hpaStatus.NumTarget)
		tbl.AddRow(
			fullName,
			hpaStatus.ID,
			hpaStatus.Lifecycle.String(),
			hpaStatus.MinReplicas,
			hpaStatus.MaxReplicas,
			replicas,
			hpaStatus.Metrics,
			hpaStatus.Benchmark,
			hpaStatus.SyncTime.Format(time.RFC3339),
			hpaStatus.Error,
		)
	}
	tbl.Print()
	return nil
}

func printSpecifiedHPAStatus(name string) error {
	hpaStatus, err := getHPAFromApiServer(name)
	if err != nil {
		return err
	}
	if hpaStatus == nil {
		return fmt.Errorf("no such hpa")
	}

	tbl := hpaStatusTbl()
	fullName := path.Join(hpaStatus.Namespace, hpaStatus.Name)
	replicas := strconv.Itoa(hpaStatus.NumReady) + "/" + strconv.Itoa(hpaStatus.NumTarget)
	tbl.AddRow(
		fullName,
		hpaStatus.ID,
		hpaStatus.Lifecycle.String(),
		hpaStatus.MinReplicas,
		hpaStatus.MaxReplicas,
		replicas,
		hpaStatus.Metrics,
		hpaStatus.Benchmark,
		hpaStatus.SyncTime.Format(time.RFC3339),
		hpaStatus.Error,
	)
	tbl.Print()
	return nil
}

func ServiceTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Name", "UID", "ClusterIp", "Ports")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}

func DNSTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Name", "UID", "Host", "Paths")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}

func getServiceFromApiServer(fullName string) (service *apiObject.Service, err error) {
	namespace, name := parseName(fullName)
	URL := url.Prefix + path.Join(url.ServiceURL, namespace, name)
	err = httputil.GetAndUnmarshal(URL, &service)
	return
}

func getServicesFromApiServer() (services []apiObject.Service, err error) {
	URL := url.Prefix + url.ServiceURL
	err = httputil.GetAndUnmarshal(URL, &services)
	return
}

func getDnsFromApiServer(fullName string) (dns *apiObject.Dns, err error) {
	namespace, name := parseName(fullName)
	URL := url.Prefix + path.Join(url.DNSURL, namespace, name)
	err = httputil.GetAndUnmarshal(URL, &dns)
	return
}

func getDnsesFromApiServer() (dnses []apiObject.Dns, err error) {
	URL := url.Prefix + url.DNSURL
	err = httputil.GetAndUnmarshal(URL, &dnses)
	return
}

func printSpecifiedService(name string) error {
	service, err := getServiceFromApiServer(name)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("no such service")
	}

	tbl := ServiceTbl()
	fullName := path.Join(service.Metadata.Namespace, service.Metadata.Name)
	tbl.AddRow(
		fullName,
		service.Metadata.UID,
		service.Spec.ClusterIP,
		service.Spec.Ports,
	)
	tbl.Print()
	return nil
}

func printServices() error {
	services, err := getServicesFromApiServer()
	if err != nil {
		return err
	}
	if services == nil {
		return fmt.Errorf("no such pod")
	}

	tbl := ServiceTbl()
	for _, service := range services {
		fullName := path.Join(service.Metadata.Namespace, service.Metadata.Name)
		tbl.AddRow(
			fullName,
			service.Metadata.UID,
			service.Spec.ClusterIP,
			service.Spec.Ports,
		)
	}
	tbl.Print()
	return nil
}

func printSpecifiedDns(name string) error {
	dns, err := getDnsFromApiServer(name)
	if err != nil {
		return err
	}
	if dns == nil {
		return fmt.Errorf("no such service")
	}

	tbl := DNSTbl()
	fullName := path.Join(dns.Metadata.Namespace, dns.Metadata.Name)
	tbl.AddRow(
		fullName,
		dns.Metadata.UID,
		dns.Spec.Host,
		dns.Spec.Paths,
	)
	tbl.Print()
	return nil
}

func printDnses() error {
	dnses, err := getDnsesFromApiServer()
	if err != nil {
		return err
	}
	if dnses == nil {
		return fmt.Errorf("no such pod")
	}

	tbl := DNSTbl()
	for _, dns := range dnses {
		fullName := path.Join(dns.Metadata.Namespace, dns.Metadata.Name)
		tbl.AddRow(
			fullName,
			dns.Metadata.UID,
			dns.Spec.Host,
			dns.Spec.Paths,
		)
	}
	tbl.Print()
	return nil
}
