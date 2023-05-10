package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var Endpoint string
var BotPrefix string
var Jira_user string
var Jira_token string
var Url string

var BotUserId string
var BotToken string
var Servers map[string]string

var Secondary map[string]string

func ReadConfig() error {

	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

    var config struct{
        Token      string `json:"Token"`
        BotPrefix  string `json:"BotPrefix"`
        Jira_token string `json:"Jira_token"`
        Jira_user  string `json:"Jira_user"`
        Url        string `json:"Url"`
        Healthcheck map[string]string `json:"Healthcheck"`
    }

	json.Unmarshal(file, &config)

    Endpoint = "https://"+config.Jira_user+":"+config.Jira_token+"@" + config.Url
    BotPrefix = config.BotPrefix
    Jira_user = config.Jira_user
    Jira_token = config.Jira_token
    Url = config.Url
    BotToken = config.Token
    Servers = config.Healthcheck

	return err
}
