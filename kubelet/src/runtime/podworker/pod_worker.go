package podworker

import (
	"fmt"
	"minik8s/apiObject"
)

const (
	workChanSize = 5
)

type podWorker struct {
	podCreateFn                  PodCreateFn
	podDeleteFn                  PodDeleteFn
	podContainerStartFn          PodContainerStartFn
	podContainerRestartFn        PodContainerRestartFn
	podContainerCreateAndStartFn PodContainerCreateAndStartFn
	podContainerRemoveFn         PodContainerRemoveFn
	workCh                       chan podWork
	currentWork                  podWork
}

func (w *podWorker) AddWork(work podWork) {
	if w.needDo(&work) {
		w.workCh <- work
	}
}

func (w *podWorker) Done() {
	close(w.workCh)
}

// needDo decide whether the worker should do this job.
func (w *podWorker) needDo(work *podWork) bool {
	lastWorkType := w.currentWork.WorkType
	workType := work.WorkType

	if lastWorkType == podCreate {
		// If the worker is creating the pod, it should reject container start, container create jobs
		// But it should reject pod delete jobs.
		return workType == podDelete
	} else if lastWorkType == podDelete {
		// similar to podCreate
		return workType == podCreate
	}
	return true
}

func (w *podWorker) handleError(err error, errPod *apiObject.Pod) {
	if err != nil && errPod != nil {
		fmt.Println(err.Error())
		w.publishPodStatus(w.errorPodStatus(errPod))
	}
}

func (w *podWorker) doWork(work podWork) {
	var err error
	var errPod *apiObject.Pod
	defer w.handleError(err, errPod)

	switch work.WorkType {
	case podCreate:
		arg := work.Arg.(podCreateFnArg)
		fmt.Printf("pod worker received pod create job %s\n", arg.pod.UID())
		w.publishPodStatus(w.containerCreatingPodStatus(arg.pod))
		if err = w.podCreateFn(arg.pod); err == nil {
			w.publishPodStatus(w.runningPodStatus(arg.pod))
		} else {
			errPod = arg.pod
		}
	case podDelete:
		arg := work.Arg.(podDeleteFnArg)
		fmt.Printf("pod worker received pod delete job %s\n", arg.pod.UID())
		if err = w.podDeleteFn(arg.pod); err == nil {
			w.publishPodStatus(w.deletedPodStatus(arg.pod))
		} else {
			errPod = arg.pod
		}
	case podContainerCreateAndStart:
		arg := work.Arg.(podContainerCreateAndStartFnArg)
		fmt.Printf("pod worker received pod create and start job %s\n", arg.pod.UID())
		err = w.podContainerCreateAndStartFn(arg.pod, arg.target)
	case podContainerRemove:
		arg := work.Arg.(podContainerRemoveFnArg)
		err = w.podContainerRemoveFn(arg.podUID, arg.ID)
	case podContainerStart:
		arg := work.Arg.(podContainerStartFnArg)
		err = w.podContainerStartFn(arg.podUID, arg.ID)
	case podContainerRestart:
		fmt.Println("pod worker received restart job")
		arg := work.Arg.(podContainerRestartFnArg)
		err = w.podContainerRestartFn(arg.pod, arg.ID, arg.fullName)
	}

	w.currentWork = noWork
}

func (w *podWorker) Run() {
	for {
		select {
		case work, open := <-w.workCh:
			if !open {
				fmt.Println("Work channel has been closed!")
				return
			}
			fmt.Println("Worker received job:", work)
			w.currentWork = work
			w.doWork(work)
		}
	}
}

func newWorker(podCreateFn PodCreateFn, podDeleteFn PodDeleteFn, podContainerCreateAndStartFn PodContainerCreateAndStartFn,
	podContainerStartFn PodContainerStartFn, podContainerRemoveFn PodContainerRemoveFn, podContainerRestartFn PodContainerRestartFn) *podWorker {
	return &podWorker{
		podCreateFn:                  podCreateFn,
		podDeleteFn:                  podDeleteFn,
		podContainerStartFn:          podContainerStartFn,
		podContainerRestartFn:        podContainerRestartFn,
		podContainerCreateAndStartFn: podContainerCreateAndStartFn,
		podContainerRemoveFn:         podContainerRemoveFn,
		workCh:                       make(chan podWork, workChanSize),
		currentWork:                  noWork,
	}
}
