package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/DiegoAraujoJS/go-bot/utils"
	"github.com/bwmarrin/discordgo"
)

func getJiraTicket(ticket_prefix string, ticket_id string, config utils.ConfigStruct) (*http.Response, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://"+config.Jira_user+":"+config.Jira_token+"@lenox-test.atlassian.net/rest/api/2/issue/"+ticket_prefix+"-"+ticket_id, nil)

	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)

	if err != nil {
		defer func() {
			if _err := recover(); _err != nil {
				fmt.Println("recover", _err)
			}
		}()
		return response, err
	}

	return response, err
}

func getTicketPhoto(content string, config utils.ConfigStruct) *http.Response {

	content = strings.Split(content, "//")[1]

	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://"+config.Jira_user+":"+config.Jira_token+"@"+content, nil)

	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)

	if err != nil {
		panic(err.Error())
	}

	return response
}

type JiraResponse struct {
	Fields struct {
		Status struct {
			Name string `json:"name"`
		} `json:"status"`
		Summary     string `json:"summary"`
		Description string `json:"description"`
		Creator     struct {
			DisplayName string `json:"displayName"`
		} `json:"creator"`
		Attachment []struct {
			Content  string `json:"content"`
			MimeType string `json:"mimeType"`
		} `json:"attachment"`
	} `json:"fields"`
    Key string `json:"key"`
}

var jiraRegexp = regexp.MustCompile(`([A-Z]+-|ticket )\d+`)
var imageNameRegexp = regexp.MustCompile(`!.*!`)

func JiraExpandTicket(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == utils.BotUserId {
        return
    }

    match := jiraRegexp.Find([]byte(m.Content))

    if match != nil {

        split := strings.Split(string(match), "-")
        if len(split) == 1 {
            split = strings.Split(string(match), " ")
        }
        prefix, ticket_id := split[0], split[len(split)-1]
        if prefix == "ticket" {
            prefix = "LW"
        }

        response, err := getJiraTicket(prefix, ticket_id, utils.Config)
        if err != nil {
            return
        }
        defer response.Body.Close()

        if strings.Contains(response.Status, "404") {
            s.ChannelMessageSend(m.ChannelID, "No existe el ticket "+ticket_id)
            return
        }

        var json_body JiraResponse
        body, _ := ioutil.ReadAll(response.Body)
        json.Unmarshal(body, &json_body)

        description_no_image_name := imageNameRegexp.ReplaceAll([]byte(json_body.Fields.Description), []byte(""))

        message := discordgo.MessageEmbed{
            Author: &discordgo.MessageEmbedAuthor{
                Name: json_body.Fields.Creator.DisplayName,
            },
            Title:       json_body.Fields.Summary,
            Description: string(description_no_image_name),
            URL:         "https://lenox-test.atlassian.net/browse/" + prefix + "-" + ticket_id,
            Color:       16711680,
        }

        var discord_response = make([]*discordgo.MessageEmbed, len(json_body.Fields.Attachment)+1)
        discord_response[0] = &message

        wg := sync.WaitGroup{}

        for i, att := range json_body.Fields.Attachment {
            if !strings.Contains(att.MimeType, "image") {
                continue
            }

            wg.Add(1)

            go func(i int, content string) {
                photo := getTicketPhoto(content, utils.Config)
                image := discordgo.MessageEmbed{
                    Image: &discordgo.MessageEmbedImage{
                        URL: photo.Request.Response.Header["Location"][0],
                    },
                }

                discord_response[i+1] = &image
                wg.Done()
                }(i, att.Content)
        }

        wg.Wait()

        var discord_response_clean = make([]*discordgo.MessageEmbed, 0, len(discord_response))

        for _, v := range discord_response { if v != nil { discord_response_clean = append(discord_response_clean, v) } }

        defer func() {
            if _err := recover(); _err != nil {
                fmt.Print("Error -->", _err)
            }
        }()
        s.ChannelMessageSendEmbeds(m.ChannelID, discord_response_clean)
    }

}
