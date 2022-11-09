package commands

import (
	"github.com/bwmarrin/discordgo"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var clockRegexp = regexp.MustCompile(`SN=[^ ]+`)

func ClockHealth(BotId string) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == BotId {
			return
		}

		match := clockRegexp.Find([]byte(m.Content))

		if match != nil {

            client := &http.Client{
                CheckRedirect: func(req *http.Request, via []*http.Request) error {
                    return http.ErrUseLastResponse
                },
                Timeout: 5 * time.Second,
            }

            split := strings.Split(string(match), "=")

            sn := split[len(split)-1]

            req, _ := http.NewRequest("GET", "http://cloud.sistemaslenox.com:48971/iclock/ping?SN="+sn, nil)

            then := time.Now()
            res, _ := client.Do(req)
            if res == nil {
                sn = sn + "\tno response\t"
            }
            time_elapsed := time.Since(then)

            _, _ = s.ChannelMessageSend(m.ChannelID, "http://cloud.sistemaslenox.com:48971/iclock/ping?SN="+sn+"\t"+strconv.Itoa(int(time_elapsed.Milliseconds()))+" ms")

		}

	}
}
