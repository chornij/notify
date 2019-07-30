package slack

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Client struct {
	Username   string
	Icon       string
	WebhookURL string
}

func (c *Client) Notify(message Message, fields []map[string]interface{}) bool {
	if c.WebhookURL != "" {
		data := getPayload(c.Username, c.Icon, message, fields)

		resp, err := http.PostForm(c.WebhookURL, url.Values{
			"payload": {toJson(data)},
		})

		if err != nil {
			log.Println(err)
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
			}

			if string(body) == "ok" {
				return true
			}
		}
	}

	return false
}

// https://api.slack.com/docs/message-attachments
func getPayload(username, icon string, message Message, fields []map[string]interface{}) map[string]interface{} {
	attachment := map[string]interface{}{
		"text":  message.Text,
		"color": message.Color,
	}

	if pretext := os.Getenv("NOTIFY_SLACK_PRETEXT"); pretext != "" {
		attachment["pretext"] = pretext
	}

	if fields != nil {
		attachment["fields"] = fields
	}

	data := map[string]interface{}{
		"username":   username,
		"icon_emoji": icon,
		"attachments": []map[string]interface{}{
			attachment,
		},
	}

	return data
}

func toJson(data map[string]interface{}) string {
	buf, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}

	return string(buf)
}
