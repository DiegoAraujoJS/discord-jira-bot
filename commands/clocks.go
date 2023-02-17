package commands

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/DiegoAraujoJS/go-bot/utils"
	"github.com/bwmarrin/discordgo"
)

var clockRegexp = regexp.MustCompile(`SN=[^ ]+`)

func ClockHealth(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == utils.BotUserId {
        return
    }

    match := clockRegexp.Find([]byte(m.Content))

    if match == nil {return}

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
