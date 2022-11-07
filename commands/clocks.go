package commands

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var clockRegexp = regexp.MustCompile(`SN=[^ ]+`)

func clockHealth(s *discordgo.Session, m *discordgo.MessageCreate) *http.Response {

	match := clockRegexp.Find([]byte(m.Content))

	var response *http.Response

	if match != nil {
		split := strings.Split(string(match), " ")

		sn := split[len(split)-1]

		client := &http.Client{}

		req, _ := http.NewRequest("GET", "http://cloud.sistemaslenox.com:48971/iclock/ping?SN="+sn, nil)

		then := time.Now()
		_, err := client.Do(req)
		time_elapsed := time.Since(then)

		if err != nil {
			log.Fatal(err)
		}
		s.ChannelMessageSend(m.ChannelID, "http://cloud.sistemaslenox.com:48971/iclock/ping?SN="+sn+"\t"+strconv.Itoa(int(time_elapsed.Milliseconds()))+" ms")

	}
	return response

}
