package iptables

import (
	"github.com/coreos/go-iptables/iptables"
	"strconv"
)

const (
	prerouting = "PREROUTING"
	output     = "OUTPUT"
)

// Manager nat table only
type Manager interface {
	NewService(service string) error
	NewChain(chain string) error
	AppendChainToEndpoint(from, endpoint string) error
	AppendChainToChains(from string, tos []string) error
	AppendServiceToChain(from, to, ip string, port int32) error
}

type iptablesManager struct {
	tab *iptables.IPTables
}

func New() (Manager, error) {
	tab, err := iptables.New()
	if err != nil {
		return nil, nil
	}
	im := &iptablesManager{
		tab: tab,
	}
	return im, err
}

func (im *iptablesManager) NewService(service string) error {
	var err error
	if err = im.NewChain(service); err != nil {
		return err
	}
	err = im.tab.AppendUnique("nat", prerouting, "-j", service)
	if err != nil {
		return err
	}
	err = im.tab.AppendUnique("nat", output, "-j", service)
	if err != nil {
		return err
	}
	return nil
}

func (im *iptablesManager) NewChain(chain string) error {
	return im.tab.NewChain("nat", chain)
}

func (im *iptablesManager) AppendChainToEndpoint(from, endpoint string) error {
	return im.tab.AppendUnique("nat", from, "-p", "tcp", "-m", "tcp",
		"-j", "DNAT", "--to-destination", endpoint)
}

func (im *iptablesManager) AppendChainToChains(from string, tos []string) error {
	num := len(tos)
	for i, to := range tos {
		if i+1 == num {
			err := im.tab.AppendUnique("nat", from, "-j", to)
			if err != nil {
				return err
			}
		} else {
			prob := 1.0 / float64(num-i)
			err := im.tab.AppendUnique("nat", from, "-m", "statistic", "--mode", "random",
				"--probability", strconv.FormatFloat(prob, 'f', 8, 32), "-j", to)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (im *iptablesManager) AppendServiceToChain(from, to, ip string, port int32) error {
	return im.tab.AppendUnique("nat", from, "-d", ip, "-p", "tcp", "-m", "tcp",
		"--dport", strconv.Itoa(int(port)), "-j", to)
}
