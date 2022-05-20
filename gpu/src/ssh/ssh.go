package ssh

import (
	"fmt"
	"github.com/melbahja/goph"
	"github.com/spf13/cast"
	"strconv"
	"strings"
)

// reference: https://docs.hpc.sjtu.edu.cn
// https://docs.hpc.sjtu.edu.cn/job/slurm.html

const (
	gpuUser   = "stu633"
	gpuPasswd = "8uhlGet%"
	gpuAddr   = "login.hpc.sjtu.edu.cn"
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
	GetAllQueueInfo() []*QueueInfo    //查看队列状态和信息
	GetJobById(jobID string) *JobInfo //显示用户作业历史
	//Sbatch() ([]byte, error)          //提交作业
	//Scancel() ([]byte, error)         //取消指定作业
	//Upload(localPath, remotePath string) error

	Mkdir(dir string) (string, error)
	CreateFile(filename string) (string, error)
	WriteFile(filename, content string) (string, error)
	ReadFile(filename string) (string, error)
}
type client struct {
	sshCli *goph.Client
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
	if cli, err := goph.NewUnknown(username, gpuAddr, goph.Password(password)); err == nil {
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
		sshCli: sshCli,
	}
}
