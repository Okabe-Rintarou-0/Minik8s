package wait

import "time"

type handler func()

type conditionalHandler func() bool

func Forever() {
	<-make(chan struct{})
}

func Period(delay time.Duration, period time.Duration, handler handler) {
	<-time.After(delay)
	tick := time.Tick(period)
	for {
		handler()
		<-tick
	}
}

func PeriodWithCondition(delay time.Duration, period time.Duration, handler conditionalHandler) {
	<-time.After(delay)
	tick := time.Tick(period)
	for handler() {
		<-tick
	}
}

func After(d time.Duration, handler handler) {
	<-time.After(d)
	handler()
}

func Until(triggerTime time.Time, handler handler) {
	delta := triggerTime.Sub(time.Now())
	if delta < 0 {
		goto handle
	}
	<-time.After(delta)
handle:
	handler()
}
