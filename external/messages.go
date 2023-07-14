package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/DiegoAraujoJS/go-bot/utils"
)

func sendWhatsAppMessage(message string, recipient string) error {

    payload := map[string]interface{}{
        "recipient_type": "individual",
        "to": recipient,
        "type": "template",
        "template": map[string]interface{}{
            "name": "hello_world",
            "language": map[string]string{
                "code": "en_US",
            },
        },
        "messaging_product": "whatsapp",
    }

    jsonValue, _ := json.Marshal(payload)

    req, err := http.NewRequest("POST", utils.Whatsapp_endpoint, bytes.NewBuffer(jsonValue))
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", utils.Whatsapp_bot_token))

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    return nil
}

func sendTelegramMessage(message string, chat_id string) error {

    apiUrl := "https://api.telegram.org"
    resource := fmt.Sprintf("/bot%v/sendMessage", utils.Telegram_bot_token)

    u, _ := url.ParseRequestURI(apiUrl)
    u.Path = resource
    urlStr := u.String() // "https://api.telegram.org/bot{Your_Bot_Token}/sendMessage"

    data := url.Values{}
    data.Set("chat_id", chat_id)
    data.Set("text", message)

    client := &http.Client{}
    r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode())) 
    r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    resp, err := client.Do(r)
    if err != nil {
        return err
    }
    fmt.Println(resp.Status)
    return nil
}
