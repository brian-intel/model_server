package utilities

import (
	"math"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

func ExpPow(val float32) float32 {
	val_float64 := float64(val)
	results := math.Exp(val_float64)
	return float32(results)
}

func Sigmoid(x float32) float32 {
	return 1.0 / (1.0 + ExpPow(-x))
}

func Min(value1 float32, value2 float32) float32 {
	if value1 < value2 {
		return value1
	}
	return value2
}

func Max(value1 float32, value2 float32) float32 {
	if value1 > value2 {
		return value1
	}
	return value2
}

func Clamp(val float32, min float32, max float32) float32 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func GetModelNameFromConfig(configfilepath string) ([]string, error) {
	if configfilepath == "" {
		return nil, fmt.Errorf("Invalid input configfilepath: %v", configfilepath)
	}

	type ModelConfigDetail struct {
		Name string `json:"name"`
	}

	type ModelConfig struct {
		Config ModelConfigDetail `json:"config"`
	}

	type ModelConfigList struct {
		ModelConfigList []ModelConfig `json:"model_config_list"`
	}

	content, err := ioutil.ReadFile(configfilepath)
	if err != nil {
		return nil, err
	}

	var modelConfigList ModelConfigList
	err = json.Unmarshal(content, &modelConfigList)
	if err != nil {
		return nil, err
	}

	var modelNameList []string
	for _, modelconfig := range modelConfigList.ModelConfigList {
		name := modelconfig.Config.Name
		modelNameList = append(modelNameList, name)
	}

	return modelNameList, nil
}
