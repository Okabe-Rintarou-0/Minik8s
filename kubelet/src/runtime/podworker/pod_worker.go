package podworker

type podWorkType byte

const (
	containerStart podWorkType = iota
)

type podWorkArgs interface{}

type podWork struct {
	WorkType podWorkType
	Args     podWorkArgs
}

type podWorker struct {
}

func (w *podWorker) doWork(work podWork) {

}

func (w *podWorker) Run(workCh workChan) {
	for {
		select {
		case work := <-workCh:
			w.doWork(work)
		}
	}
}
