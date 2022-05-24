package netutil

import (
	"fmt"
	"testing"
)

func TestUtil(t *testing.T) {
	//GetHostIp()

	fmt.Println(LookUpIp("www.baidu.com"))
}
