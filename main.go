package main

import (
	"fmt"

	"github.com/DiegoAraujoJS/go-bot/commands"
	"github.com/DiegoAraujoJS/go-bot/utils"
	"github.com/bwmarrin/discordgo"
)



var (
	goBot *discordgo.Session
)

func Start() {
	goBot, err := discordgo.New("Bot " + utils.BotToken)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	user, err := goBot.User("@me")
    utils.BotUserId = user.ID

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	goBot.AddHandler(commands.HealthCheck)
	goBot.AddHandler(commands.JiraExpandTicket)
	goBot.AddHandler(commands.GetTickets)
	goBot.AddHandler(commands.ClockHealth)
    goBot.AddHandler(commands.GetCommands)

	err = goBot.Open()

	if err != nil {
		fmt.Println("start ", err.Error())
		return
	}
	fmt.Println("bot is running")
}

// Mover estos valores a un json y agregarlo al gitignore.

func main() {
    err := utils.ReadConfig()

	if err != nil {
		return
	}

	Start()

	<-make(chan struct{})
}
