package utils

import (
	"encoding/json"

	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/model"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/service"
)

// convert node info
func ConvertNodeInfo(a interface{}) map[string]int {
	mapInt := make(map[string]int)
	mapInterface := a.(map[string]interface{})

	for key, value := range mapInterface {
		mapInt[key] = int(value.(float64))
	}

	return mapInt
}

func ConvertAppKey(a []interface{}) (appKey, appSecret string) {
	keyVar := a[0].([]interface{})
	secretVar := a[1].([]interface{})
	var keyStr, secretStr []string

	for _, value := range keyVar {
		keyStr = append(keyStr, value.(string))
	}
	for _, value := range secretVar {
		secretStr = append(secretStr, value.(string))
	}

	appKey = keyStr[2]
	appSecret = secretStr[2]
	return
}

func ConvertApiSet(a model.JSON) service.ApiSet {
	values, _ := a.Value()
	var apiSet service.ApiSet
	json.Unmarshal(values.([]byte), &apiSet)
	return apiSet
}

func ConvertChannels(a model.JSON) []int32 {
	values, _ := a.Value()
	channels := make([]int32, 0)
	json.Unmarshal(values.([]byte), &channels)
	return channels
}

func ConvertNodeMap(a model.JSON) service.ApiSet {
	values, _ := a.Value()
	var apiSet service.ApiSet
	json.Unmarshal(values.([]byte), &apiSet)
	return apiSet
}

func StringSliceIncludes(a string, b []string) bool {
	if (a == ``) != (b == nil) {
		return false
	}

	for i := range b {
		if a == b[i] {
			return true
		}
	}

	return false
}

func ConvertUpstreamMap(a model.JSON) []service.Upstream {
	values, _ := a.Value()
	var upstreamMap []service.Upstream
	json.Unmarshal(values.([]byte), &upstreamMap)
	return upstreamMap
}
