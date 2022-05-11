package wait

import "runtime"

func Forever() {
	for {
		runtime.Gosched()
	}
}
