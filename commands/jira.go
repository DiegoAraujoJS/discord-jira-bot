package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
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

func getTicketPhoto(content string, config ConfigStruct) *http.Response {

	content = strings.Split(content, "//")[1]

	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://"+config.Jira_user+":"+config.Jira_token+"@"+content, nil)

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
		Attachment []struct {
			Content string `json:"content"`
		} `json:"attachment"`
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

			photo := getTicketPhoto(json_body.Fields.Attachment[0].Content, config)

			fmt.Println(photo.Request.Response.Header["Location"][0])

			message := discordgo.MessageEmbed{
				Title:       json_body.Fields.Summary,
				Description: "https://lenox-test.atlassian.net/browse/LW-" + ticket_id + "\n\n" + json_body.Fields.Description + "\n\n" + json_body.Fields.Creator.DisplayName,
			}

			var images discordgo.MessageEmbed

			for _, v := range photo.Request.Response.Header {
				fmt.Println(images, v)
			}

			image := discordgo.MessageEmbed{
				Image: &discordgo.MessageEmbedImage{
					URL: photo.Request.Response.Header["Location"][0],
				},
			}

			s.ChannelMessageSendEmbeds(m.ChannelID, []*discordgo.MessageEmbed{&message, &image})

		}
	}
}
