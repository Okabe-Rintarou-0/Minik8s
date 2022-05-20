package ssh

import (
	"fmt"
	"github.com/melbahja/goph"
	"log"
	"testing"
)

func TestSSH(t *testing.T) {
	cli, err := goph.NewUnknown(gpuUser, gpuAddr, goph.Password(gpuPasswd))
	defer cli.Close()
	if err != nil {
		fmt.Println("err!")
		log.Fatal(err)
	}

	resp, err := cli.Run("sinfo")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Got resp: %s\n", string(resp))
}

func TestGpuSSH(t *testing.T) {
	cli := NewClient(gpuUser, gpuPasswd)
	defer cli.Close()
	job := cli.GetJobById("13256987")
	fmt.Println(job)
	allQueues := cli.GetAllQueueInfo()
	fmt.Println(allQueues)
	smallQueue := cli.GetQueueInfoByPartition("small")
	fmt.Println(smallQueue)
	//resp, _ = cli.Squeue()
	//fmt.Println(string(resp))
	//resp, _ = cli.Sbatch()
	//fmt.Println(string(resp))
	//resp, _ = cli.Sinfo()
	//fmt.Println(string(resp))
}
