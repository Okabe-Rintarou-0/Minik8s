package nginx

import (
	"minik8s/apiserver/src/url"
	"minik8s/util/logger"
	"os"
	"os/exec"
	"path"
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
	ApplyLoadBalance(servers []Server) error
	Start() error
	Shutdown() error
	GetName() string
	Reload() error
}

type nginxManager struct {
	dirPath  string
	filename string
	UID      string
}

func New(UID string) Manager {
	return &nginxManager{
		dirPath:  url.NginxDirPath,
		filename: path.Join(url.NginxDirPath, url.NginxFileName),
		UID:      UID,
	}
}

func (nm *nginxManager) Start() error {
	cmd := exec.Command("docker", "run", "--name", nm.GetName(),
		"-v", nm.dirPath+":/etc/nginx/", "-d", "nginx")
	logger.Log("nginx start")(cmd.String())
	if out, err := cmd.Output(); err != nil {
		logger.Log("nginx start")(string(out))
		return err
	} else {
		logger.Log("nginx start")(string(out))
	}
	return nil
}

func (nm *nginxManager) Reload() error {
	if _, err := exec.Command("docker", "exec", nm.GetName(),
		"nginx", "-s", "reload").Output(); err != nil {
		return err
	}
	return nil
}

func (nm *nginxManager) Shutdown() error {
	if _, err := exec.Command("docker", "stop", nm.GetName()).Output(); err != nil {
		return err
	}
	if _, err := exec.Command("docker", "rm", nm.GetName()).Output(); err != nil {
		return err
	}
	return nil
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
			_, err = file.WriteString("\n\t\tlocation = " + location.Addr + " {\n\t\t\tproxy_pass http://" + location.Dest + "/;\n\t\t}\n")
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
	return nil
}

func (nm *nginxManager) ApplyLoadBalance(servers []Server) error {
	file, err := os.Create(nm.filename)
	if err != nil {
		return err
	}
	if _, err = file.WriteString("events {\n\tworker_connections 1024;\n}\n\nhttp {\n"); err != nil {
		return err
	}
	for i, server := range servers {
		if len(server.Locations) != 0 {
			_, err = file.WriteString("\tupstream backend" + strconv.Itoa(i) + " {\n")
			for _, location := range server.Locations {
				_, err = file.WriteString("\t\tserver " + location.Dest + ";\n")
			}
			_, err = file.WriteString("\t}\n")
		}
	}

	for i, server := range servers {
		if len(server.Locations) != 0 {
			port := strconv.Itoa(server.Port)
			if _, err = file.WriteString("\n\tserver {\n\t\tlisten " + port + ";\n\t\tserver_name localhost;\n"); err != nil {
				return err
			}
			if _, err = file.WriteString("\t\tlocation / {\n\t\t\tproxy_pass http://backend" + strconv.Itoa(i) + ";\n\t\t}\n"); err != nil {
				return err
			}
			if _, err = file.WriteString("\t}\n\n"); err != nil {
				return err
			}
		}
	}

	if _, err = file.WriteString("}"); err != nil {
		return err
	}
	return nil
}

func (nm *nginxManager) GetName() string {
	return "nginx-" + nm.UID
}
