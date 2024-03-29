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
    fmt.Println("Starting bot")
	goBot, err := discordgo.New("Bot " + utils.BotToken)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	goBot.AddHandler(commands.HealthCheck)
	goBot.AddHandler(commands.JiraExpandTicket)
	goBot.AddHandler(commands.GetTickets)
	goBot.AddHandler(commands.ClockHealth)
    goBot.AddHandler(commands.GetCommands)
    goBot.AddHandler(commands.ReactionsHandler)
    goBot.AddHandler(commands.GetGuildData)

	err = goBot.Open()

	if err != nil {
		fmt.Println("start", err.Error())
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
