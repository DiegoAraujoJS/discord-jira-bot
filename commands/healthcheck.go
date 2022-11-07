package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"net/http"
	"strconv"
	"sync"
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

			var wg = sync.WaitGroup{}

			for k, v := range servers {
				wg.Add(1)
				then := time.Now()
				go func(k string, v string, then time.Time) {
					_, err := client.Get(v)
					time_elapsed := time.Since(then)
					if err != nil {
						response += "status " + k + " " + v + ": down " + strconv.Itoa(int(time_elapsed.Milliseconds())) + " ms" + "\n"
					} else {
						response += "status " + k + " " + v + ": ok " + strconv.Itoa(int(time_elapsed.Milliseconds())) + " ms" + "\n"
					}
					wg.Done()

				}(k, v, then)
			}

			wg.Wait()

			_, _ = s.ChannelMessageSend(m.ChannelID, response)
			return
		}

	}
}
