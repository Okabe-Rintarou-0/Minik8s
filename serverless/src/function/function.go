package function

import (
	"minik8s/serverless/src/utils"
)

const (
	pythonImage = "python:3.10-slim"
)

func InitFunction() {
	utils.PullImage(pythonImage)
}
