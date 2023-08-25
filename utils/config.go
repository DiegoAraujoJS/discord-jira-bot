package utils

import (
	"encoding/json"
	"os"
)

type guildConfiguration struct {
    JiraEndpoint string
    JiraUser string
    JiraToken string
    JiraUrl string
    BotPrefix string
}

var guildConfigurations = map[string]guildConfiguration{}

type configJson struct {
    BotUserId   string
    BotToken    string
    GuildConfigurations guildConfiguration
}

var config *configJson

type publicGuildConfiguration struct {
    JiraEndpoint string
    BotPrefix string
}

func getConfig() *configJson {
    if config != nil {
        var configData, err = os.ReadFile("./config.json")
        if err != nil {
            panic("No config.json found at root")
        }
        err = json.Unmarshal(configData, config)
        if err != nil {
            panic("Cannot store config.json data")
        }
        return config
    }
    return config
}

func GetGuildConfig(guildId string) *publicGuildConfiguration {
    return &publicGuildConfiguration{
        JiraEndpoint: "https://"+guildConfigurations[guildId].JiraUser+":"+guildConfigurations[guildId].JiraToken+"@" + guildConfigurations[guildId].JiraUrl,
        BotPrefix: guildConfigurations[guildId].BotPrefix,
    }
}

func GetBotToken() string {
    return config.BotToken
}

func GetBotUserId() string {
    return config.BotUserId
}

func SetBotUserId(botUserId string) {
    config.BotUserId = botUserId
}
