package ipgen

import (
	"fmt"
	"math/big"
	"minik8s/apiserver/src/etcd"
	"net"
	"strconv"
)

func inetNtoA(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

func inetAtoN(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

type Generator interface {
	GetCurrent() (string, error)
	GetNext() (string, error)
	GetCurrentWithMask() (string, error)
	GetNextWithMask() (string, error)
	Clear(ip string) error
	ClearIfInit(ip string) error
}

type ipGenerator struct {
	url  string
	mask int
}

func New(url string, mask int) Generator {
	return &ipGenerator{
		url:  url,
		mask: mask,
	}
}

func (ig *ipGenerator) GetCurrent() (string, error) {
	ret, err := etcd.Get(ig.url)
	if err != nil {
		return ret, err
	}
	num, err := strconv.Atoi(ret)
	if err != nil {
		return ret, err
	}
	return inetNtoA(int64(num)), nil
}

func (ig *ipGenerator) GetNext() (string, error) {
	ret, err := etcd.Get(ig.url)
	if err != nil {
		return ret, err
	}
	num, err := strconv.Atoi(ret)
	if err != nil {
		return ret, err
	}
	if err = etcd.Put(ig.url, strconv.Itoa(num+1)); err != nil {
		return ret, err
	}
	return inetNtoA(int64(num + 1)), nil
}

func (ig *ipGenerator) GetCurrentWithMask() (string, error) {
	ret, err := ig.GetCurrent()
	if err != nil {
		return ret, err
	}
	return ret + "/" + strconv.Itoa(ig.mask), nil
}

func (ig *ipGenerator) GetNextWithMask() (string, error) {
	ret, err := ig.GetNext()
	if err != nil {
		return ret, err
	}
	return ret + "/" + strconv.Itoa(ig.mask), nil
}

func (ig *ipGenerator) Clear(ip string) error {
	return etcd.Put(ig.url, strconv.Itoa(int(inetAtoN(ip))))
}

func (ig *ipGenerator) ClearIfInit(ip string) error {
	if ret, err := etcd.Get(ig.url); err != nil {
		return err
	} else {
		if ret == "" {
			return etcd.Put(ig.url, strconv.Itoa(20+int(inetAtoN(ip))))
		}
	}
	return nil
}
