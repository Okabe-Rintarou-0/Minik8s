package gpu

import (
	"testing"
)

func TestJobScript(t *testing.T) {
	//	template := `#!/bin/bash
	//#SBATCH --job-name=%s
	//#SBATCH --partition=dgx2
	//#SBATCH --output=%s
	//#SBATCH --error=%s
	//#SBATCH -N %d
	//#SBATCH --ntasks-per-node=%d
	//#SBATCH --cpus-per-task=%d
	//#SBATCH --gres=%s
	//
	//%s
	//`
	args := JobArgs{
		JobName:         "123",
		Output:          "123",
		Error:           "123.out",
		NumProcess:      1,
		NumTasksPerNode: 1,
		CpusPerTask:     1,
		GpuResources:    "gpu:1",
		CompileScripts:  "module load cuda/9.2.88-gcc-4.8.5;nvcc gpu-job/cublashello.cu -o gpu-job/cublashello -lcublas",
		RunScripts:      "module load cuda/9.2.88-gcc-4.8.5;./gpu-job/cublashello",
		Username:        "stu633",
		Password:        "8uhlGet%",
		WorkDir:         "gpu-job",
	}
	//script := fmt.Sprintf(
	//	template,
	//	args.JobName,
	//	args.Output,
	//	args.Error,
	//	args.NumProcess,
	//	args.NumTasksPerNode,
	//	args.CpusPerTask,
	//	args.GpuResources,
	//	strings.Replace(args.RunScripts, ";", "\n", -1),
	//)

	server := NewServer(args, "D:/gpu")
	server.Run()
	//fmt.Println(strings.Split(script, "\n"))
}
