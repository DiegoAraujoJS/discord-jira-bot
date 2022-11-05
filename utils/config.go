package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ReadConfig(json_name string, memory_store interface{}) error {

	file, err := ioutil.ReadFile(json_name)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	json.Unmarshal(file, memory_store)

	return err
}

var Exported string
