package etcd

import (
	"encoding/json"
	"fmt"
)

func strToPtr(s string) *string {
	return &s
}

func resourceKey(resourceType, resourceID string) string {
	key := fmt.Sprintf("%s/%s/%s/%s", defaultPrefix, resourceType, defaultNamespace, resourceID)
	return key
}

func encodeValue(in interface{}) string {
	encoded, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	return string(encoded)
}

func decodeValue(data []byte, out interface{}) error {
	return json.Unmarshal(data, out)
}
