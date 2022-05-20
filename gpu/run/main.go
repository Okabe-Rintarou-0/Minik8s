package main

import (
	"flag"
	"fmt"
)

var (
	jobName         string
	output          string
	err             string
	n               int
	numTasksPerNode int
	cpusPerTask     int
	gpuResources    string
)

func main() {
	flag.StringVar(&jobName, "job-name", "gpu-job", "gpu job name")
	flag.StringVar(&output, "output", "output", "output filename")
	flag.StringVar(&err, "error", "error", "err filename")
	flag.IntVar(&n, "N", 1, "number of processes(cpus)")
	flag.IntVar(&numTasksPerNode, "ntasks-per-node", 1, "number of tasks per node")
	flag.IntVar(&cpusPerTask, "cpus-per-task", 1, "number of cpus per task")
	flag.StringVar(&gpuResources, "gres", "gpu:1", "gpu resources")
	flag.Parse()

	fmt.Printf("Read args: jobName = %v, output = %v, error = %v, n = %v\n", jobName, output, err, n)
	fmt.Printf("Read args: numTasksPerNode = %v, cpusPerTask = %v, gpuResources = %v\n", numTasksPerNode, cpusPerTask, gpuResources)
}
