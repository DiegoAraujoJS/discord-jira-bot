package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
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

	if m.Author.ID == BotId {
		return
	}

	if m.Content == config.BotPrefix+"upstatus" {
		cloudResponse, _ := http.Get("https://cloud.sistemaslenox.com.ar")
		body, _ := ioutil.ReadAll(cloudResponse.Body)
		testResponse, _ := http.Get("https://test.sistemaslenox.com.ar")
		testBody, _ := ioutil.ReadAll(testResponse.Body)
        _, _ = s.ChannelMessageSend(m.ChannelID, string(body) + "\n" + string(testBody))
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
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
