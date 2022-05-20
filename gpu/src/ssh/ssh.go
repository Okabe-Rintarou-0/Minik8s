package ssh

import (
	"fmt"
	"github.com/melbahja/goph"
	"github.com/spf13/cast"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// reference: https://docs.hpc.sjtu.edu.cn
// https://docs.hpc.sjtu.edu.cn/job/slurm.html

const (
	gpuUser      = "stu633"
	gpuPasswd    = "8uhlGet%"
	gpuLoginAddr = "login.hpc.sjtu.edu.cn"
	gpuDataAddr  = "data.hpc.sjtu.edu.cn"
	accountType  = "acct-stu"
)

type JobInfo struct {
	JobID     string
	JobName   string
	Partition string
	Account   string
	AllocCPUS int
	State     string
	ExitCode  string
}

type QueueInfo struct {
	Partition string
	Available string
	TimeLimit string
	Nodes     int
	State     string
	NodeList  string
}

type Client interface {
	Close()

	GetQueueInfoByPartition(partition string) []*QueueInfo
	GetAllQueueInfo() []*QueueInfo                  //Sinfo
	GetJobById(jobID string) *JobInfo               //Squeue
	SubmitScript(scriptPath string) (string, error) //Sbatch
	//Scancel() ([]byte, error)         //取消指定作业

	Compile(cmd string) (string, error)

	Scp(localPath, remotePath string) error
	Rsync(localPath, remotePath string) error

	Mkdir(dir string) (string, error)
	CreateFile(filename string) (string, error)
	WriteFile(filename, content string) (string, error)
	ReadFile(filename string) (string, error)
}

type client struct {
	username string
	password string
	sshCli   *goph.Client
}

func (cli *client) Compile(cmd string) (string, error) {
	if resp, err := cli.loadCuda(); err != nil {
		return resp, err
	}

	resp, err := cli.sshCli.Run(cmd)
	return string(resp), err
}

func (cli *client) Scp(localPath, remotePath string) error {
	if runtime.GOOS == "linux" {
		remoteAddr := fmt.Sprintf("%s@%s:%s", cli.username, gpuDataAddr, remotePath)
		cmd := exec.Command("scp", "-r", localPath, remoteAddr)
		return cmd.Run()
	}
	return fmt.Errorf("scp is not supported in your os")
}

func (cli *client) Rsync(localPath, remotePath string) error {
	if runtime.GOOS == "linux" {
		remoteAddr := fmt.Sprintf("%s@%s:%s", cli.username, gpuDataAddr, remotePath)
		cmd := exec.Command("rsync", "--archive", "--partial", "--progress", remoteAddr, localPath)
		return cmd.Run()
	}
	return fmt.Errorf("rsync is not supported in your os")
}

func (cli *client) loadCuda() (string, error) {
	cmd := "module load cuda/9.2.88-gcc-4.8.5"
	resp, err := cli.sshCli.Run(cmd)
	return string(resp), err
}

func (cli *client) SubmitScript(scriptPath string) (string, error) {
	cmd := fmt.Sprintf("sbatch %s", scriptPath)
	respRaw, err := cli.sshCli.Run(cmd)
	return string(respRaw), err
}

func (cli *client) Mkdir(dir string) (string, error) {
	cmd := fmt.Sprintf("mkdir %s", dir)
	resp, err := cli.sshCli.Run(cmd)
	return string(resp), err
}

func (cli *client) CreateFile(filename string) (string, error) {
	cmd := fmt.Sprintf("touch %s", filename)
	resp, err := cli.sshCli.Run(cmd)
	return string(resp), err
}

func (cli *client) WriteFile(filename, content string) (string, error) {
	cmd := fmt.Sprintf("echo \"%s\" > %s", strconv.Quote(content), filename)
	resp, err := cli.sshCli.Run(cmd)
	return string(resp), err
}

func (cli *client) ReadFile(filename string) (string, error) {
	cmd := fmt.Sprintf("cat %s", filename)
	resp, err := cli.sshCli.Run(cmd)
	return string(resp), err
}

func (cli *client) Close() {
	cli.sshCli.Close()
}

func (cli *client) GetJobById(jobID string) *JobInfo {
	cmd := fmt.Sprintf("sacct -j %s | tail -n +3 | awk '{print $1, $2, $3, $4, $5, $6, $7}'", jobID)
	if raw, err := cli.sshCli.Run(cmd); err == nil {
		resp := string(raw)
		rows := strings.Split(resp, "\n")
		if len(rows) > 0 {
			row := rows[0]
			cols := strings.Split(row, " ")
			return &JobInfo{
				JobID:     cols[0],
				JobName:   cols[1],
				Partition: cols[2],
				Account:   cols[3],
				AllocCPUS: cast.ToInt(cols[4]),
				State:     cols[5],
				ExitCode:  cols[6],
			}
		}
	}
	return nil
}

func parseQueueInfoTable(raw string) (infos []*QueueInfo) {
	rows := strings.Split(raw, "\n")
	for _, row := range rows {
		cols := strings.Split(row, " ")
		if len(cols) != 6 {
			continue
		}
		infos = append(infos, &QueueInfo{
			Partition: cols[0],
			Available: cols[1],
			TimeLimit: cols[2],
			Nodes:     cast.ToInt(cols[3]),
			State:     cols[4],
			NodeList:  cols[5],
		})
	}
	return
}

func (cli *client) GetAllQueueInfo() (infos []*QueueInfo) {
	cmd := "sinfo | tail -n +2 | awk '{print $1, $2, $3, $4, $5, $6}'"
	if raw, err := cli.sshCli.Run(cmd); err == nil {
		return parseQueueInfoTable(string(raw))
	}
	return nil
}

func (cli *client) GetQueueInfoByPartition(partition string) (infos []*QueueInfo) {
	cmd := fmt.Sprintf("sinfo --partition=%s | tail -n +2 | awk '{print $1, $2, $3, $4, $5, $6}'", partition)
	if raw, err := cli.sshCli.Run(cmd); err == nil {
		return parseQueueInfoTable(string(raw))
	}
	return nil
}

func (cli *client) Scancel() ([]byte, error) {
	return cli.sshCli.Run("scancel")
}

func (cli *client) Upload(localPath, remotePath string) error {
	return cli.sshCli.Upload(localPath, remotePath)
}

func newSSHClient(username, password string) *goph.Client {
	if cli, err := goph.NewUnknown(username, gpuLoginAddr, goph.Password(password)); err == nil {
		return cli
	}
	return nil
}

func NewClient(username, password string) Client {
	sshCli := newSSHClient(username, password)
	if sshCli == nil {
		return nil
	}
	return &client{
		username: username,
		password: password,
		sshCli:   sshCli,
	}
}
