package ipgen

import (
	"errors"
	"fmt"
	"math/big"
	"minik8s/apiserver/src/etcd"
	"net"
	"strconv"
	"strings"
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
	Clear() error
	ClearIfInit() error
}

type ipGenerator struct {
	url  string
	base int64
	mask int
}

func New(url string, net string) (Generator, error) {
	sp := strings.Split(net, "/")
	if len(sp) != 2 {
		return nil, errors.New("invalid subnet")
	}
	ip := sp[0]
	mask, err := strconv.Atoi(sp[1])
	if err != nil {
		return nil, err
	}
	return &ipGenerator{
		url:  url,
		base: inetAtoN(ip),
		mask: mask,
	}, nil
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
	return inetNtoA(ig.base + int64(num)), nil
}

func (ig *ipGenerator) GetNext() (string, error) {
	ret, err := etcd.Get(ig.url)
	if err != nil {
		return ret, err
	}
	num, err := strconv.Atoi(ret)
	newNum := (num + 1) & ((1 << (32 - ig.mask)) - 1)
	if err != nil {
		return ret, err
	}
	if err := etcd.Put(ig.url, strconv.Itoa(newNum)); err != nil {
		return ret, err
	}
	return inetNtoA(ig.base + int64(newNum)), nil
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

func (ig *ipGenerator) Clear() error {
	return etcd.Put(ig.url, "1")
}

func (ig *ipGenerator) ClearIfInit() error {
	if ret, err := etcd.Get(ig.url); err != nil {
		return err
	} else {
		if ret == "" {
			return etcd.Put(ig.url, "1")
		}
	}
	return nil
}
