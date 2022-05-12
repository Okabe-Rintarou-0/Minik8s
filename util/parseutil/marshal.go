package parseutil

import (
	"encoding/json"
	"minik8s/entity"
)

func MarshalPodStatus(status *entity.PodStatus) []byte {
	result, err := json.Marshal(status)
	if err != nil {
		return nil
	}
	return result
}
