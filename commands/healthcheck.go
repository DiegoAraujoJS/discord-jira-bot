package commands

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/DiegoAraujoJS/go-bot/utils"
	"github.com/bwmarrin/discordgo"
)


func HealthCheck(s *discordgo.Session, m *discordgo.MessageCreate) {

    if m.Author.ID == utils.BotUserId {
        return
    }

    fmt.Println("content -->", m.Content)

    client := &http.Client{
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return http.ErrUseLastResponse
        },
        Timeout: 5 * time.Second,
    }

    if m.Content == utils.Config.BotPrefix+"upstatus" {
        var response = ""

        var wg = sync.WaitGroup{}

        for k, v := range utils.Servers {
            wg.Add(1)
            go func(k string, v string) {
                then := time.Now()
                _, err := client.Get(v)
                time_elapsed := time.Since(then)
                if err != nil {
                    response += "status " + k + " " + v + ": down " + strconv.Itoa(int(time_elapsed.Milliseconds())) + " ms" + "\n"
                } else {
                    response += "status " + k + " " + v + ": ok " + strconv.Itoa(int(time_elapsed.Milliseconds())) + " ms" + "\n"
                }
                wg.Done()

            }(k, v)
        }

        wg.Wait()

        _, _ = s.ChannelMessageSend(m.ChannelID, response)
        return
    }
}
