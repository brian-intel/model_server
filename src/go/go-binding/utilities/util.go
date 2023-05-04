package utilities

import (
	"math"
	"os"
	"bufio"
	"strings"
	"fmt"
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

func GetModelNameFromConfig(configfilepath string) (string, error) {
	if configfilepath == "" {
		return "", fmt.Errorf("Invalid input configfilepath: %v", configfilepath)
	}

	file, err := os.Open(configfilepath)
	if err != nil {
		return "", fmt.Errorf("Error open file: %v", configfilepath)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	var modelName = ""
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			return modelName, fmt.Errorf("Error readline: %v", err)
		}
		if strings.Contains(string(line), "\"name\": ") {
			sections := strings.Split(string(line), ":")
			names := strings.Split(sections[1], "\"")
			modelName := names[1]
			return modelName, nil
		}
	}
	return "", fmt.Errorf("Can not find model name in configfilepath: %v", configfilepath)
}
