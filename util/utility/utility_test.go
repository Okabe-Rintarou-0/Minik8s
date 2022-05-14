package utility

import (
	"fmt"
	"testing"
	"time"
)

func TestUtility(t *testing.T) {
	start := time.Now()
	fmt.Println(GetCpuAndMemoryUsage())
	end := time.Now()
	fmt.Println("Takes ", end.Sub(start).Seconds(), "seconds")
}
