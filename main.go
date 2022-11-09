package main

import (
	"fmt"
	"github.com/DiegoAraujoJS/go-bot/commands"
	"github.com/DiegoAraujoJS/go-bot/utils"
	"github.com/bwmarrin/discordgo"
)

var config commands.ConfigStruct

var servers map[string]string

var (
	BotId string
	goBot *discordgo.Session
)

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	user, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	BotId = user.ID
	goBot.AddHandler(commands.HealthCheck(BotId, config, servers))
	goBot.AddHandler(commands.JiraExpandTicket(BotId, config))
	goBot.AddHandler(commands.ClockHealth(BotId))

	err = goBot.Open()

	if err != nil {
		fmt.Println("start ", err.Error())
		return
	}
	fmt.Println("bot is running")
}

// Mover estos valores a un json y agregarlo al gitignore.

func main() {
	err := utils.ReadConfig("./config.json", &config)

	if err != nil {
		return
	}

	err = utils.ReadConfig("./healthcheck_routes.json", &servers)

	if err != nil {
		return
	}

	Start()

	<-make(chan struct{})
}
