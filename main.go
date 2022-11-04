package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
	"time"
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

var servers = map[string]string{
	"cloud":            "https://cloud.sistemaslenox.com",
	"ADMS cloud":       "http://cloud.sistemaslenox.com:48971/iclock/ping?SN=COVS220960036",
	"PWA cloud":        "https://cloud.sistemaslenox.com/mobile",
	"PWA server cloud": "https://cloud.sistemaslenox.com:48970/api",
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == BotId {
		return
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 5 * time.Second,
	}

	if m.Content == config.BotPrefix+"upstatus" {
		var response = ""

		for k, v := range servers {
			_, err := client.Get(v)
			if err != nil {
				response += "status " + k + " " + v + ": down\n"
			} else {
				response += "status " + k + " " + v + ": ok\n"
			}
		}

		_, _ = s.ChannelMessageSend(m.ChannelID, response)
		return
	}
}

func main() {
	err := ReadConfig()

	if err != nil {
		return
	}

	fmt.Println(config)

	Start()

	<-make(chan struct{})
}
