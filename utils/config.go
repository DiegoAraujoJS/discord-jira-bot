package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ConfigStruct struct {
	Token      string `json:"Token"`
	BotPrefix  string `json:"BotPrefix"`
	Jira_token string `json:"Jira_token"`
	Jira_user  string `json:"Jira_user"`
    Url        string `json:"Url"`
}

var Config ConfigStruct
var BotUserId string
var Servers map[string]string

func ReadConfig(json_name string, memory_store interface{}) error {

	file, err := ioutil.ReadFile(json_name)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	json.Unmarshal(file, memory_store)

	return err
}
