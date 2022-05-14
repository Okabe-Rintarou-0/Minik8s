package parseutil

import (
	"encoding/json"
)

func MarshalAny(v interface{}) []byte {
	result, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return result
}
