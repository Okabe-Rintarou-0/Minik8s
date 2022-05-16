package nginx

import (
	"os"
	"os/exec"
	"strconv"
)

type Location struct {
	Dest string
	Addr string
}

type Server struct {
	Locations []Location
	Port      int
}

type Manager interface {
	Apply(servers []Server) error
}

type nginxManager struct {
	filename string
}

func New(filename string) Manager {
	return &nginxManager{
		filename: filename,
	}
}

func (nm *nginxManager) Apply(servers []Server) error {
	file, err := os.Create(nm.filename)
	if err != nil {
		return err
	}
	_, err = file.WriteString("events {\n\tworker_connections 1024;\n}\n\nhttp {\n")
	if err != nil {
		return err
	}
	for _, server := range servers {
		port := strconv.Itoa(server.Port)
		_, err = file.WriteString("\tserver {\n\t\tlisten " + port + ";\n\t\tserver_name localhost;\n")
		if err != nil {
			return err
		}
		for _, location := range server.Locations {
			_, err = file.WriteString("\n\t\tlocation = " + location.Addr + " {\n\t\t\tproxy_pass " + location.Dest + ";\n\t\t}\n")
			if err != nil {
				return err
			}
		}
		_, err = file.WriteString("\t}\n\n")
		if err != nil {
			return err
		}
	}
	_, err = file.WriteString("}")
	if err != nil {
		return err
	}
	cmd := exec.Command("nginx", "-s", "reload")
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
