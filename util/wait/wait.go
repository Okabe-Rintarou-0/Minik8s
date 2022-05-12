package wait

func Forever() {
	<-make(chan struct{})
}
