package utils

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"math/big"
	"os"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckSQLError(err error) {
	if err != nil {
		panic("SQL error")
	}
}

func Zeroing(buf []byte) {
	for i := range buf {
		buf[i] = 0
	}
}

func Config() map[string]string {
	yamlString, err := os.ReadFile("config.yaml")
	CheckError(err)

	var config map[string]interface{}
	err = yaml.Unmarshal([]byte(yamlString), &config)
	CheckError(err)
	return convertToStringMap(config)
}

func convertToStringMap(mapInterface map[string]interface{}) map[string]string {
	mapString := make(map[string]string)
	for key, value := range mapInterface {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)

		mapString[strKey] = strValue
	}
	return mapString
}
func StringToBigInt(val string) (big.Int, error) {
	n := new(big.Int)
	n, ok := n.SetString(val, 10)
	if !ok {
		return big.Int{}, errors.New("SetString: error")
	}
	return *n, nil
}
