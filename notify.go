package notify

import (
	"github.com/getsentry/sentry-go"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/chornij/notify/slack"
)

var (
	slackClient *slack.Client
)

type (
	Scope struct {
		device  *Device
		tags    []Tag
		request *Request
	}

	Device struct {
		ID string
	}

	Tag struct {
		Title string
		Value string
		Href  string
	}

	Request struct {
		URL    string
		Method string
		IP     string
	}
)

// Notify error
func Error(err error, action string) {
	if err == nil {
		return
	}

	scope := &Scope{
		tags: []Tag{
			{
				Title: "Action",
				Value: action,
			},
		},
	}

	stdout(err.Error(), scope)

	slackClient := getSlackClient()
	go slackClient.Notify(slack.Error(err.Error()), getSlackTags(scope))

	sentry.WithScope(func(s *sentry.Scope) {
		s.SetTags(getSentryTags(scope))
		sentry.CaptureException(err)
	})
}

// Notify error request
func ErrorRequest(err error, r *http.Request) {
	if err == nil {
		return
	}

	scope := &Scope{}

	if r != nil {
		url := ""

		if r.URL != nil {
			url = r.URL.Path
		}

		request := Request{
			Method: r.Method,
			URL:    url,
			IP:     r.RemoteAddr,
		}
		scope.request = &request
	}

	stdout(err.Error(), scope)

	slackClient := getSlackClient()
	go slackClient.Notify(slack.Error(err.Error()), getSlackTags(scope))

	sentry.ConfigureScope(func(s *sentry.Scope) {
		s.SetTags(getSentryTags(scope))
		sentry.CaptureException(err)
	})
}

// Logging message with tags
func Log(text string, tags []Tag) {
	scope := &Scope{
		tags: tags,
	}

	stdout(text, scope)

	slackClient := getSlackClient()
	go slackClient.Notify(slack.Info(text), getSlackTags(scope))

	sentry.WithScope(func(s *sentry.Scope) {
		s.SetTags(getSentryTags(scope))
		sentry.CaptureMessage(text)
	})
}

// Slack and stdout info message
func LogInfo(text string, scope *Scope) {
	stdout(text, scope)

	slackClient := getSlackClient()
	go slackClient.Notify(slack.Info(text), getSlackTags(scope))
}

// Slack and stdout warning message
func LogWarning(text string, scope *Scope) {
	stdout(text, scope)

	slackClient := getSlackClient()
	go slackClient.Notify(slack.Warning(text), getSlackTags(scope))
}

// Slack and stdout error message
func LogError(err error, scope *Scope) {
	if err == nil {
		return
	}

	stdout(err.Error(), scope)

	slackClient := getSlackClient()
	go slackClient.Notify(slack.Error(err.Error()), getSlackTags(scope))
}

func stdout(text string, scope *Scope) {
	var v []string
	v = append(v, text)

	if scope != nil {
		if scope.device != nil {
			v = append(v, "Device ID: "+scope.device.ID)
		} else if scope.request != nil {
			v = append(v, "IP: "+scope.request.IP)
		}
	}

	log.Println(strings.Join(v, " "))
}

func getSlackClient() slack.Client {
	if slackClient == nil {
		slackClient = &slack.Client{
			Username:   os.Getenv("NOTIFY_SLACK_USERNAME"),
			Icon:       os.Getenv("NOTIFY_SLACK_ICON"),
			WebhookURL: os.Getenv("NOTIFY_SLACK_TOKEN"),
		}
	}

	return *slackClient
}

func getSlackTags(scope *Scope) (tags []map[string]interface{}) {
	if scope != nil {
		if scope.device != nil {
			value := scope.device.ID
			if url := os.Getenv("NOTIFY_SLACK_DEVICE_URL"); url != "" {
				value = "<" + strings.Replace(url, "{value}", scope.device.ID, -1) + ">"
			}

			tags = append(tags, map[string]interface{}{
				"title": "Device",
				"value": value,
				"short": true,
			})
		}

		for _, t := range scope.tags {
			v := t.Value

			if t.Href != "" {
				v = "<" + strings.Replace(t.Href, "{value}", t.Value, -1) + ">"
			}

			tags = append(tags, map[string]interface{}{
				"title": t.Title,
				"value": v,
				"short": true,
			})
		}

		if scope.request != nil {
			tags = append(tags, map[string]interface{}{
				"title": "Request",
				"value": scope.request.Method + " " + scope.request.URL,
				"short": true,
			})

			if scope.request.IP != "" {
				if logIp := os.Getenv("NOTIFY_SLACK_LOG_IP"); logIp == "true" {
					tags = append(tags, map[string]interface{}{
						"title": "Request IP",
						"value": scope.request.IP,
						"short": true,
					})
				}
			}
		}
	}

	return
}

func getSentryTags(scope *Scope) (tags map[string]string) {
	if scope.device != nil && scope.request != nil {
		sentry.ConfigureScope(func(s *sentry.Scope) {
			s.SetUser(sentry.User{
				ID:        scope.device.ID,
				IPAddress: scope.request.IP,
			})
		})
	}

	tags = make(map[string]string)

	if scope != nil {
		if scope.device != nil {
			value := scope.device.ID
			if url := os.Getenv("NOTIFY_SLACK_DEVICE_URL"); url != "" {
				value = strings.Replace(url, "{value}", scope.device.ID, -1)
			}

			tags["Device"] = value
		}

		for _, t := range scope.tags {
			v := t.Value

			if t.Href != "" {
				v = strings.Replace(t.Href, "{value}", t.Value, -1)
			}

			tags[t.Title] = v
		}

		if scope.request != nil {
			tags["Request"] = scope.request.Method + " " + scope.request.URL

			if scope.request.IP != "" {
				if logIp := os.Getenv("NOTIFY_SLACK_LOG_IP"); logIp == "true" {
					tags["Request IP"] = scope.request.IP
				}
			}
		}
	}

	return
}
