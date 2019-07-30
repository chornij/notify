package notify

import (
	"errors"
	"github.com/chornij/env"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestNotifyError(t *testing.T) {
	env.TryLoad("./test.env")

	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	})
	assert.NotEmpty(t, os.Getenv("SENTRY_DSN"))
	assert.NoError(t, err)

	Error(errors.New("some error"), "testing")
}

func TestNotifyErrorRequest(t *testing.T) {
	env.TryLoad("./test.env")

	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	})
	assert.NotEmpty(t, os.Getenv("SENTRY_DSN"))
	assert.NoError(t, err)

	r := http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "HTTPS",
			Host:   "example.com",
			Path:   "/users",
		},
		RemoteAddr: "192.168.0.1",
	}

	ErrorRequest(errors.New("some error"), &r)
}

func TestLoggingMessageWithTags(t *testing.T) {
	env.TryLoad("./test.env")

	Log("Some message", []Tag{
		{
			Title: "Big Mac",
			Value: "big-mac",
			Href:  "https://www.mcdonalds.com/us/en-us/product/{value}.html",
		},
		{
			Title: "Some title",
		},
		{
			Value: "Vvv",
		},
	})
}

func TestLoggingInfoMessage(t *testing.T) {
	env.TryLoad("./test.env")

	LogInfo("Some message", nil)
}

func TestLoggingWarningMessage(t *testing.T) {
	env.TryLoad("./test.env")

	LogWarning("Some message", nil)
}

func TestLoggingErrorMessage(t *testing.T) {
	env.TryLoad("./test.env")

	LogError(errors.New("some message"), nil)
}

func TestLoggingInfoMessageWithScope(t *testing.T) {
	env.TryLoad("./test.env")

	LogInfo("Some message", getScope())
}

func TestLoggingWarningMessageWithScope(t *testing.T) {
	env.TryLoad("./test.env")

	LogWarning("Some message", getScope())
}

func TestLoggingErrorMessageWithScope(t *testing.T) {
	env.TryLoad("./test.env")

	LogError(errors.New("some message"), getScope())
}

func getScope() *Scope {
	return &Scope{
		device: &Device{
			ID: "777",
		},
		tags: []Tag{
			{
				Title: "Type",
				Value: "Subscriptions",
			},
			{
				Title: "Response status",
				Value: "409",
				Href:  "https://httpstatuses.com/{value}",
			},
		},
		request: nil,
	}
}
