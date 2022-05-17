package handlers

import (
	"github.com/gin-gonic/gin"
	"minik8s/apiObject"
	"minik8s/util/uidutil"
	"net/http"
	"strconv"
)

func HandleAutoscale(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	target := c.PostForm("target")
	targetNamespace, targetName := parseTargetName(target)
	cpu, _ := strconv.ParseFloat(c.PostForm("cpu"), 64)
	mem, _ := strconv.ParseFloat(c.PostForm("mem"), 64)
	min, _ := strconv.Atoi(c.PostForm("min"))
	max, _ := strconv.Atoi(c.PostForm("max"))
	interval, _ := strconv.Atoi(c.PostForm("interval"))
	log("Read hpa %s/%s params: cpu = %v, mem = %v, min = %v, max = %v", namespace, name, cpu, mem, min, max)

	uid := uidutil.New()
	if name == "" {
		name = "hpa-" + uid
	}

	hpa := &apiObject.HorizontalPodAutoscaler{
		Base: apiObject.Base{
			ApiVersion: "autoscaling/v1",
			Kind:       "HorizontalPodAutoscaler",
			Metadata: apiObject.Metadata{
				Name:      name,
				Namespace: namespace,
				UID:       uid,
			},
		},
		Spec: apiObject.HPASpec{
			MinReplicas: min,
			MaxReplicas: max,
			ScaleTargetRef: apiObject.ScaleTargetRef{
				ApiVersion: "v1",
				Kind:       "ReplicaSet",
				Metadata: apiObject.Metadata{
					Name:      targetName,
					Namespace: targetNamespace,
				},
			},
			Metrics: apiObject.Metrics{
				CPUUtilizationPercentage: cpu,
				MemUtilizationPercentage: mem,
			},
			ScaleInterval: interval,
		},
	}

	if err := addHPA(hpa); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	c.String(http.StatusOK, "Apply successfully!")
}
