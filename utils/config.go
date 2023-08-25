package utils

import (
)

type GuildConfiguration struct {
    JiraEndpoint string
    JiraUser string
    JiraToken string
    JiraUrl string
    BotPrefix string
}

var GuildConfigurations = map[string]GuildConfiguration{}

type publicGuildConfiguration struct {
    JiraEndpoint string
    BotPrefix string
}

func GetGuildConfig(guildId string) *publicGuildConfiguration {
    return &publicGuildConfiguration{
        JiraEndpoint: "https://"+GuildConfigurations[guildId].JiraUser+":"+GuildConfigurations[guildId].JiraToken+"@" + GuildConfigurations[guildId].JiraUrl,
        BotPrefix: GuildConfigurations[guildId].BotPrefix,
    }
}


