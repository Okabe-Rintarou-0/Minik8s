package dns

import (
	"os"
	"strings"
)

type Manager interface {
	AddEntry(host, ip string) error
	DelIfExistEntry(host string) error
}

type dnsManager struct {
	filename string
}

func New(filename string) Manager {
	return &dnsManager{
		filename: filename,
	}
}

func (dm *dnsManager) getMapping() (map[string]string, error) {
	data, err := os.ReadFile(dm.filename)
	mp := make(map[string]string)
	if err != nil {
		if os.IsNotExist(err) {
			return mp, nil
		}
		return nil, err
	}

	entries := strings.Split(string(data), "\n")
	for _, entry := range entries {
		records := strings.Split(entry, " ")
		if len(records) != 2 {
			continue
		}
		mp[strings.Trim(records[1], " ")] = strings.Trim(records[0], " ")
	}

	return mp, nil
}

func (dm *dnsManager) writeBack(mp map[string]string) error {
	file, err := os.Create(dm.filename)
	if err != nil {
		return err
	}
	for key, value := range mp {
		entry := value + " " + key + "\n"
		_, err := file.WriteString(entry)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dm *dnsManager) AddEntry(host, ip string) error {
	mp, err := dm.getMapping()
	if err != nil {
		return err
	}
	mp[host] = ip
	return dm.writeBack(mp)
}

func (dm *dnsManager) DelIfExistEntry(host string) error {
	mp, err := dm.getMapping()
	if err != nil {
		return err
	}
	delete(mp, host)
	return dm.writeBack(mp)
}
