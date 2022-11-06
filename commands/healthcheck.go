package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"net/http"
	"time"
)

type ConfigStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
}

func HealthCheck(BotId string, config ConfigStruct, servers map[string]string) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {

		if m.Author.ID == BotId {
			return
		}

		fmt.Println("content -->", m.Content)

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
}
