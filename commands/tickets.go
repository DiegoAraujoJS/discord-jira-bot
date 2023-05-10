package commands

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/DiegoAraujoJS/go-bot/utils"
	"github.com/bwmarrin/discordgo"
)

type MultipleJiraResponse struct {
	Issues []JiraResponse `json:"issues"`
}

var askTicketsRegex = regexp.MustCompile(`!tickets-\w+-(\w+|"[\w ]+")`)
var quotedStateRegex = regexp.MustCompile(`"[\w ]+"`)

func getTicketUrl(ticket_key string) string {
    ticket_prefix := strings.Split(ticket_key, "-")[0]
    ticket_id := strings.Split(ticket_key, "-")[1]
    return utils.Endpoint + ".atlassian.net/browse/" + ticket_prefix + "-" + ticket_id
}

func GetTickets(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == utils.BotUserId {
        return
    }
    match := askTicketsRegex.Find([]byte(m.Content))
    if match == nil {
        return
    }
    go s.MessageReactionAdd(m.ChannelID, m.ID, "📝")
    status := strings.Split(string(match), "-")[2]
    match_quotes := quotedStateRegex.Find([]byte(status))
    if match_quotes != nil {
        quoted_string := string(match_quotes)
        status = quoted_string[1 : len(quoted_string)-1]
    }
    url := utils.Endpoint + ".atlassian.net/rest/api/3/search"
    headers := map[string]string{
        "Accept":       "application/json",
        "Content-Type": "application/json",
    }
    project := strings.Split(string(match), "-")[1]

    _payload := map[string]interface{}{
        "expand": []string{
            "names",
            "schema",
            "operations",
        },
        "jql": "project = " + project + " AND status = \"" + status + "\" ORDER BY created DESC",
        "fieldsByKeys": false,
        "fields": []string{
            "status",
            "summary",
        },
        "startAt": 0,
    }

    payload, _ := json.Marshal(_payload)

    client := &http.Client{}
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))

    for h, v := range headers { req.Header.Set(h, v) }

    resp, err := client.Do(req)
    if err != nil {
        panic(err.Error())
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    var jira_response MultipleJiraResponse
    json.Unmarshal(body, &jira_response)
    var message_body string

    var count int

    for _, v := range jira_response.Issues {
        if count == 13 {break}
        count += 1
        message_body += "\n\n" + "["+v.Key+"]"+"("+getTicketUrl(v.Key)+")" + "\t" + v.Fields.Summary
    }

    if message_body == "" { message_body = "No se encontraron tickets de " + project + " en " + status  }
    s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
        Description: message_body,
    })
}

