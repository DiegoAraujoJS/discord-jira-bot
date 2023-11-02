package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/DiegoAraujoJS/go-bot/utils"
	"github.com/bwmarrin/discordgo"
)

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

func getJiraTicket(ticket_prefix string, ticket_id string) (*http.Response, error) {
	client := &http.Client{}

    url := utils.Endpoint + ".atlassian.net/rest/api/2/issue/"+ticket_prefix+"-"+ticket_id
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

var jiraRegexp = regexp.MustCompile(`([A-Z]+-|ticket )\d+`)
var imageNameRegexp = regexp.MustCompile(`!.*!`)

func JiraExpandTicket(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.Bot { return }

    defer func() {
        if err := recover(); err != nil {
            fmt.Println("Recover from panic; Error ->", err)
        }
    }()

    utils.ExposeUsageDetails("jira-ticket", func(record string, records map[string]int) {
        fmt.Printf("%v\tFetch jira ticket.\tUsage: %v", time.Now().Add(-3 * time.Hour).Format("2022-12-12 15:32:12") + "\tGMT -3", records[record])
    })

    match := jiraRegexp.Find([]byte(m.Content))
    if match == nil { return }

    go func() {
        defer func() {
            if err := recover(); err != nil {
                fmt.Println("Recover from panic while sending reaction. Error ->", err)
            }
        }()
        s.MessageReactionAdd(m.ChannelID, m.ID, "👀")
    }()

    split := strings.Split(string(match), "-")
    if len(split) == 1 {
        split = strings.Split(string(match), " ")
    }
    prefix, ticket_id := split[0], split[len(split)-1]
    if prefix == "ticket" {
        prefix = "LW"
    }

    response, err := getJiraTicket(prefix, ticket_id)
    if err != nil {
        return
    }
    defer response.Body.Close()

    if strings.Contains(response.Status, "404") {
        s.ChannelMessageSend(m.ChannelID, "No existe el ticket " + prefix + "-" + ticket_id)
        return
    }

    var json_body JiraResponse
    body, _ := ioutil.ReadAll(response.Body)
    json.Unmarshal(body, &json_body)

    s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
        Author: &discordgo.MessageEmbedAuthor{
            Name: json_body.Fields.Creator.DisplayName,
        },
        Title:       json_body.Fields.Summary,
        Description: string(imageNameRegexp.ReplaceAll([]byte(json_body.Fields.Description), []byte(""))),
        URL:         "https://" + utils.Url + ".atlassian.net/browse/" + prefix + "-" + ticket_id,
        Color:       16711680,
    })
}
