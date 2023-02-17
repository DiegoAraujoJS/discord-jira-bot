package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
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

var askTicketsRegex = regexp.MustCompile(`!tickets=(\w+|"[\w ]+")`)
var quotedStateRegex = regexp.MustCompile(`"[\w ]+"`)

func GetTickets(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == utils.BotUserId {
        return
    }
    match := askTicketsRegex.Find([]byte(m.Content))
    if match == nil {
        return
    }
    status := strings.Split(string(match), "=")[1]
    match_quotes := quotedStateRegex.Find([]byte(status))
    if match_quotes != nil {
        quoted_string := string(match_quotes)
        status = quoted_string[1 : len(quoted_string)-1]
    }
    url := "https://" + utils.Config.Jira_user + ":" + utils.Config.Jira_token + "@lenox-test.atlassian.net/rest/api/3/search"
    headers := map[string]string{
        "Accept":       "application/json",
        "Content-Type": "application/json",
    }

    _payload := map[string]interface{}{
        "expand": []string{
            "names",
            "schema",
            "operations",
        },
        "jql": "project = LW AND status = \"" + status + "\" ORDER BY created DESC",
        // "maxResults":   3,
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

    for h, v := range headers {
        req.Header.Set(h, v)
    }
    resp, err := client.Do(req)
    if err != nil {
        panic(err.Error())
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)

    var jira_response MultipleJiraResponse

    json.Unmarshal(body, &jira_response)

    fmt.Println(len(jira_response.Issues))

    var message_body string

    for _, v := range jira_response.Issues {
        message_body += "\n" + v.Key + "\t" + v.Fields.Summary
    }

    s.ChannelMessageSend(m.ChannelID, message_body)
}

