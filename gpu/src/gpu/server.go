package gpu

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io/fs"
	"io/ioutil"
	"minik8s/apiObject/types"
	"minik8s/entity"
	"minik8s/gpu/src/ssh"
	"minik8s/util/logger"
	"minik8s/util/recoverutil"
	"minik8s/util/topicutil"
	"minik8s/util/uidutil"
	"minik8s/util/wait"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type JobArgs struct {
	JobName         string
	WorkDir         string
	Output          string
	Error           string
	NumProcess      int
	NumTasksPerNode int
	CpusPerTask     int
	GpuResources    string
	RunScripts      string
	CompileScripts  string
	Username        string
	Password        string
}

const pollPeriod = time.Minute * 1
const DefaultJobURL = "./usr/local/jobs"

var (
	client = redis.NewClient(&redis.Options{
		Addr:     "10.119.11.101:6379",
		Password: "",
		DB:       0,
	})
	ctx = context.Background()
)

type Server interface {
	Run()
}

type server struct {
	cli     ssh.Client
	args    JobArgs
	uid     types.UID
	jobsURL string
	jobID   string
}

func (s *server) recover() {
	if err := recover(); err != nil {
		fmt.Println(recoverutil.Trace(fmt.Sprintf("%v\n", err)))
		s.cli.Reconnect()
	}
}

func (s *server) uploadJobStatus(jobStatus *entity.GpuJobStatus) {
	msg, _ := json.Marshal(jobStatus)
	client.Publish(ctx, topicutil.GpuJobStatusTopic(), msg)
}

func (s *server) poll() bool {
	defer s.recover()
	fmt.Println("Poll")
	state, completed := s.cli.JobCompleted(s.jobID)
	jobName := s.args.JobName
	parts := strings.Split(jobName, "/")
	var namespace, name string
	if len(parts) == 2 {
		namespace, name = parts[0], parts[1]
	} else {
		namespace, name = "default", jobName
	}
	status := &entity.GpuJobStatus{
		Namespace:    namespace,
		Name:         name,
		State:        state,
		LastSyncTime: time.Now(),
	}
	s.uploadJobStatus(status)
	fmt.Printf("Upload status: %+v\n", status)
	return !completed
}

func (s *server) getCudaFiles() []string {
	var cudaFiles []string
	_ = filepath.WalkDir(s.jobsURL, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			fileName := d.Name()
			if strings.HasSuffix(fileName, ".cu") {
				cudaFiles = append(cudaFiles, fileName)
			}
		}
		return nil
	})
	fmt.Printf("cudaFiles: %v\n", cudaFiles)
	return cudaFiles
}

func (s *server) uploadSmallFiles(filenames []string) error {
	if resp, err := s.cli.Mkdir(s.args.WorkDir); err != nil {
		fmt.Println(resp)
		return err
	}
	for _, filename := range filenames {
		if file, err := os.Open(path.Join(s.jobsURL, filename)); err == nil {
			if content, err := ioutil.ReadAll(file); err == nil {
				_, _ = s.cli.WriteFile(path.Join(s.args.WorkDir, filename), string(content))
			}
		}
	}
	return nil
}

func (s *server) scriptPath() string {
	return path.Join(s.args.WorkDir, s.args.JobName+"-"+s.uid+".slurm")
}

func (s *server) createJobScript() error {
	template := `#!/bin/bash
#SBATCH --job-name=%s
#SBATCH --partition=dgx2
#SBATCH --output=%s
#SBATCH --error=%s
#SBATCH -N %d
#SBATCH --ntasks-per-node=%d
#SBATCH --cpus-per-task=%d
#SBATCH --gres=%s

%s
`
	script := fmt.Sprintf(
		template,
		s.args.JobName,
		s.args.Output,
		s.args.Error,
		s.args.NumProcess,
		s.args.NumTasksPerNode,
		s.args.CpusPerTask,
		s.args.GpuResources,
		strings.Replace(s.args.RunScripts, ";", "\n", -1),
	)
	_, err := s.cli.WriteFile(s.scriptPath(), script)
	return err
}

func (s *server) compile() error {
	_, err := s.cli.Compile(s.args.CompileScripts)
	return err
}

func (s *server) submitJob() (err error) {
	if s.jobID, err = s.cli.SubmitJob(s.scriptPath()); err == nil {
		fmt.Printf("submit succeed, got jod ID: %s\n", s.jobID)
	}
	return err
}

func (s *server) prepare() (err error) {
	cudaFiles := s.getCudaFiles()
	if len(cudaFiles) == 0 {
		return fmt.Errorf("no available cuda files")
	}
	if err = s.uploadSmallFiles(cudaFiles); err != nil {
		return err
	}
	fmt.Println("upload cuda files successfully")
	if err = s.compile(); err != nil {
		return err
	}
	fmt.Println("compile successfully")
	if err = s.createJobScript(); err != nil {
		return err
	}
	fmt.Println("create job script successfully")
	return nil
}

func (s *server) downloadResult() {
	outputFile := s.args.Output
	if content, err := s.cli.ReadFile(outputFile); err == nil {
		if file, err := os.Create(path.Join(s.jobsURL, outputFile)); err == nil {
			defer file.Close()
			_, _ = file.Write([]byte(content))
		}
	}

	errorFile := s.args.Error
	if content, err := s.cli.ReadFile(errorFile); err == nil {
		if file, err := os.Create(path.Join(s.jobsURL, errorFile)); err == nil {
			defer file.Close()
			_, _ = file.Write([]byte(content))
		}
	}
}

func (s *server) Run() {
	_, _ = s.cli.RmDir(s.args.WorkDir)
	if err := s.prepare(); err != nil {
		logger.Error("prepare: " + err.Error())
		return
	}
	if err := s.submitJob(); err != nil {
		logger.Error("submit: " + err.Error())
		return
	}
	wait.PeriodWithCondition(pollPeriod, pollPeriod, s.poll)
	fmt.Println("Job finished, now download the result")
	s.downloadResult()
	fmt.Println("Down load successfully, now hang forever.")
	wait.Forever()
}

func NewServer(args JobArgs, jobsURL string) Server {
	return &server{
		cli:     ssh.NewClient(args.Username, args.Password),
		args:    args,
		uid:     uidutil.New(),
		jobsURL: jobsURL,
	}
}
