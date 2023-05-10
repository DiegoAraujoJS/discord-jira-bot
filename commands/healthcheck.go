package commands

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/DiegoAraujoJS/go-bot/utils"
	"github.com/bwmarrin/discordgo"
)

func getReponse (k string, v string, client *http.Client, wg *sync.WaitGroup, response *string, ok *bool) {
    then := time.Now()
    _, err := client.Get(v)
    time_elapsed := time.Since(then)
    if err != nil {
        *response += "status " + k + " " + v + ": 🔴 " + strconv.Itoa(int(time_elapsed.Milliseconds())) + " ms" + "\n"
        *ok = false
    } else {
        *response += "status " + k + " " + v + ": ✅ " + strconv.Itoa(int(time_elapsed.Milliseconds())) + " ms" + "\n"
    }
    wg.Done()
}

func HealthCheck(s *discordgo.Session, m *discordgo.MessageCreate) {

    if m.Author.ID == utils.BotUserId {
        return
    }

    client := &http.Client{
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return http.ErrUseLastResponse
        },
        Timeout: 5 * time.Second,
    }

    if m.Content == utils.BotPrefix + "upstatus" {
        s.MessageReactionAdd(m.ChannelID, m.ID, "📡")
        var response = ""
        var ok = true

        var wg = sync.WaitGroup{}

        for k, v := range utils.Servers {
            wg.Add(1)
            go getReponse(k, v, client, &wg, &response, &ok)
        };
        wg.Wait()

        message, _ := s.ChannelMessageSend(m.ChannelID, response)

        if ok {
            go s.MessageReactionAdd(m.ChannelID, message.ID, "😎")
            return
        }
        go s.MessageReactionAdd(m.ChannelID, message.ID, "😱")
    }
}
