package weaveutil

import (
	"minik8s/util/logger"
	"os/exec"
)

func WeaveAttach(id, ip string) error {
	if out, err := exec.Command("weave", "attach", ip, id).Output(); err != nil {
		logger.Log("weave attach err")(err.Error())
		logger.Log("weave attach output")(string(out))
		return err
	} else {
		logger.Log("weave attach success")(string(out))
		return nil
	}
}
