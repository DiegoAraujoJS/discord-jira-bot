package commands

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func getJiraTicket(ticket_id string, config ConfigStruct) *http.Response {

	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://"+config.Jira_user+":"+config.Jira_token+"@lenox-test.atlassian.net/rest/api/2/issue/LW-"+ticket_id, nil)

	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	return response
}

type jiraResponse struct {
	Fields struct {
		Summary     string `json:"summary"`
		Description string `json:"description"`
		Creator     struct {
			DisplayName string `json:"displayName"`
		} `json:"creator"`
	} `json:"fields"`
}

var jiraRegexp = regexp.MustCompile(`(?i)(LW-|ticket |LW )\d+`)

func JiraExpandTicket(BotId string, config ConfigStruct) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {

		if m.Author.ID == BotId {
			return
		}

		match := jiraRegexp.Find([]byte(m.Content))

		if match != nil {

			split := strings.Split(string(match), "-")
			if len(split) == 1 {
				split = strings.Split(string(match), " ")
			}
			ticket_id := split[len(split)-1]

			response := getJiraTicket(ticket_id, config)

			var json_body jiraResponse
			body, _ := ioutil.ReadAll(response.Body)
			json.Unmarshal(body, &json_body)

			s.ChannelMessageSend(m.ChannelID, "https://lenox-test.atlassian.net/browse/LW-"+ticket_id+"\n"+json_body.Fields.Summary+"\n\n"+json_body.Fields.Description+"\n\n"+json_body.Fields.Creator.DisplayName)
		}
	}
}
