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
    fmt.Println("Bot token: ", utils.BotToken)
	if err != nil {
		fmt.Println(err.Error())
		return
	}


	user, err := goBot.User("@me")
    fmt.Println("Bot user: ", user.ID)
    utils.BotUserId = user.ID


	if err != nil {
		fmt.Println(err.Error())
		return
	}

    goBot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
        fmt.Printf("Author: %s, Message: %s, Channell: %s\n", m.Author.Username, m.Content, m.ChannelID)
    })
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

    halt := make(chan struct{})
    <-halt
}
