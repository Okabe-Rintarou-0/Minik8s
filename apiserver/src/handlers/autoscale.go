package handlers

import (
	"github.com/gin-gonic/gin"
	"minik8s/apiObject"
	"minik8s/util/uidutil"
	"net/http"
	"strconv"
)

func HandleAutoscale(c *gin.Context) {
	name := c.Param("name")
	target := c.PostForm("target")
	cpu, _ := strconv.ParseFloat(c.PostForm("cpu"), 64)
	mem, _ := strconv.ParseFloat(c.PostForm("mem"), 64)
	min, _ := strconv.Atoi(c.PostForm("min"))
	max, _ := strconv.Atoi(c.PostForm("max"))
	log("Read hpa params: cpu = %v, mem = %v, min = %v, max = %v", cpu, mem, min, max)

	if name == "" {
		name = "hpa-" + uidutil.New()
	}

	hpa := &apiObject.HorizontalPodAutoscaler{
		Base: apiObject.Base{
			ApiVersion: "autoscaling/v1",
			Kind:       "HorizontalPodAutoscaler",
			Metadata: apiObject.Metadata{
				Name:      name,
				Namespace: "default",
				UID:       uidutil.New(),
			},
		},
		Spec: apiObject.HPASpec{
			MinReplicas: min,
			MaxReplicas: max,
			ScaleTargetRef: apiObject.ScaleTargetRef{
				ApiVersion: "v1",
				Kind:       "ReplicaSet",
				Metadata: apiObject.Metadata{
					Name:      target,
					Namespace: "default",
				},
			},
			Metrics: &apiObject.Metrics{
				CPUUtilizationPercentage: cpu,
				MemUtilizationPercentage: mem,
			},
		},
	}

	if err := addHPA(hpa); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	c.String(http.StatusOK, "Apply successfully!")
}
