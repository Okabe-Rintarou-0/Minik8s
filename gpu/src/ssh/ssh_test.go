package ssh

import (
	"fmt"
	"github.com/melbahja/goph"
	"log"
	"testing"
)

func TestSSH(t *testing.T) {
	cli, err := goph.NewUnknown(gpuUser, gpuLoginAddr, goph.Password(gpuPasswd))
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
	job := cli.GetJobById("13277555")
	fmt.Println(job)
	//allQueues := cli.GetAllQueueInfo()
	//fmt.Println(allQueues)
	//smallQueue := cli.GetQueueInfoByPartition("small")
	//fmt.Println(smallQueue)
	//fmt.Println(cli.WriteFile("test.txt", "hello world!"))
	//fmt.Println(cli.ReadFile("test.txt"))
	//fmt.Println(cli.CreateFile("test2.txt"))
	//fmt.Println(cli.Mkdir("./test223"))
}
