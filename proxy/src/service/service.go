package service

import (
	"github.com/coreos/go-iptables/iptables"
	"minik8s/proxy/src/constant"
	"strconv"
	"strings"
)

type EndPoint struct {
	Name string
	Ip   string
	Port string
}

func (ep *EndPoint) getTarget() string {
	return ep.Ip + ":" + ep.Port
}

type Manager interface {
	Init() error
	CreateService(service, ip string, port string) error
	DeleteService(service, ip string, port string) error
	// CreateEndpoints Endpoints should be deleted before re-create endpoints
	CreateEndpoints(service string, eps []EndPoint) error
	DeleteEndPoints(service string, eps []EndPoint) error
}

// serviceManager nat table only
type serviceManager struct {
	tab *iptables.IPTables
}

func New() (Manager, error) {
	tab, err := iptables.New()
	if err != nil {
		return nil, err
	}
	sm := &serviceManager{
		tab: tab,
	}
	return sm, err
}

func (sm *serviceManager) newChain(chain string) error {
	if exist, err := sm.tab.ChainExists("nat", chain); err != nil {
		return err
	} else if exist {
		return nil
	}
	return sm.tab.NewChain("nat", chain)
}

func (sm *serviceManager) delChain(chain string) error {
	return sm.tab.DeleteChain("nat", chain)
}

func (sm *serviceManager) appendEndpoints(eps []EndPoint) error {
	for _, ep := range eps {
		err := sm.tab.AppendUnique("nat", ep.Name, "-p", "tcp", "-m", "tcp",
			"-j", "DNAT", "--to-destination", ep.getTarget())
		if err != nil {
			return err
		}
	}
	return nil
}

func (sm *serviceManager) deleteEndpoints(eps []EndPoint) error {
	for _, ep := range eps {
		err := sm.tab.Delete("nat", ep.Name, "-p", "tcp", "-m", "tcp",
			"-j", "DNAT", "--to-destination", ep.getTarget())
		if err != nil {
			return err
		}
	}
	return nil
}

func (sm *serviceManager) appendChainToChains(from string, tos []EndPoint) error {
	num := len(tos)
	for i, to := range tos {
		if i+1 == num {
			err := sm.tab.AppendUnique("nat", from, "-j", to.Name)
			if err != nil {
				return err
			}
		} else {
			prob := 1.0 / float64(num-i)
			err := sm.tab.AppendUnique("nat", from, "-m", "statistic", "--mode", "random",
				"--probability", strconv.FormatFloat(prob, 'f', 8, 32), "-j", to.Name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (sm *serviceManager) deleteChainToChains(from string, tos []EndPoint) error {
	num := len(tos)
	for i, to := range tos {
		if i+1 == num {
			err := sm.tab.Delete("nat", from, "-j", to.Name)
			if err != nil {
				return err
			}
		} else {
			prob := 1.0 / float64(num-i)
			err := sm.tab.Delete("nat", from, "-m", "statistic", "--mode", "random",
				"--probability", strconv.FormatFloat(prob, 'f', 8, 32), "-j", to.Name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (sm *serviceManager) appendServiceToChain(from, to, ip string, port string) error {
	return sm.tab.AppendUnique("nat", from, "-d", ip, "-p", "tcp", "-m", "tcp",
		"--dport", port, "-j", to)
}

func (sm *serviceManager) deleteServiceToChain(from, to, ip string, port string) error {
	return sm.tab.Delete("nat", from, "-d", ip, "-p", "tcp", "-m", "tcp",
		"--dport", port, "-j", to)
}

func (sm *serviceManager) Init() (err error) {
	if err = sm.tab.ClearChain("nat", constant.KubeServices); err != nil {
		return err
	}
	if chains, err := sm.tab.ListChains("nat"); err != nil {
		return err
	} else {
		for _, chain := range chains {
			if strings.Contains(chain, constant.KubeEndpoint) {
				if err = sm.tab.ClearChain("nat", chain); err != nil {
					return err
				}
			} else if strings.Contains(chain, constant.KubeService) {
				if err = sm.tab.ClearChain("nat", chain); err != nil {
					return err
				}
			}
		}
		for _, chain := range chains {
			if strings.Contains(chain, constant.KubeEndpoint) {
				if err = sm.tab.DeleteChain("nat", chain); err != nil {
					return err
				}
			} else if strings.Contains(chain, constant.KubeService) {
				if err = sm.tab.DeleteChain("nat", chain); err != nil {
					return err
				}
			}
		}
	}

	if err = sm.newChain(constant.KubeServices); err != nil {
		return err
	}
	if err = sm.tab.AppendUnique("nat", constant.PreRouting, "-j", constant.KubeServices); err != nil {
		return err
	}
	if err = sm.tab.AppendUnique("nat", constant.Output, "-j", constant.KubeServices); err != nil {
		return err
	}
	return nil
}

func (sm *serviceManager) CreateService(service, ip string, port string) error {
	if err := sm.newChain(service); err != nil {
		return err
	}
	return sm.appendServiceToChain(constant.KubeServices, service, ip, port)
}

func (sm *serviceManager) DeleteService(service, ip string, port string) error {
	if err := sm.deleteServiceToChain(constant.KubeServices, service, ip, port); err != nil {
		return err
	}
	return sm.delChain(service)
}

func (sm *serviceManager) CreateEndpoints(service string, eps []EndPoint) error {
	for _, ep := range eps {
		if err := sm.newChain(ep.Name); err != nil {
			return err
		}
	}
	if err := sm.appendEndpoints(eps); err != nil {
		return err
	}
	return sm.appendChainToChains(service, eps)
}

func (sm *serviceManager) DeleteEndPoints(service string, eps []EndPoint) error {
	if err := sm.deleteEndpoints(eps); err != nil {
		return err
	}
	//if err := sm.deleteChainToChains(service, eps); err != nil {
	//	return err
	//}
	if err := sm.tab.ClearChain("nat", service); err != nil {
		return err
	}
	for _, ep := range eps {
		if err := sm.delChain(ep.Name); err != nil {
			return err
		}
	}
	return nil
}
