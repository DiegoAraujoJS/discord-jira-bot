package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
)

var config configStruct

type configStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
}

func ReadConfig() error {
	file, err := ioutil.ReadFile("./config.json")

	json.Unmarshal(file, &config)

	if err != nil {
		fmt.Println("read config ", err.Error())
		return err
	}

	fmt.Println(config)
	return err
}

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

	goBot.AddHandler(messageHandler)

	err = goBot.Open()

	if err != nil {
		fmt.Println("start ", err.Error())
		return
	}
	fmt.Println("bot is running")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "" {
		chanMessages, _ := s.ChannelMessages(m.ChannelID, 1, "", "", m.ID)
		m.Content = chanMessages[0].Content
	}
	fmt.Println("message:", m.Content)

	if m.Author.ID == BotId {
		return
	}
	if true {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
	}
}

func main() {
	err := ReadConfig()

	if err != nil {
		fmt.Println("main ", err.Error())
		return
	}

	fmt.Println(config)

	Start()

	<-make(chan struct{})
}
