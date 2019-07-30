# Notification sender

Send notification to stdout and/or Slack and/or Sentry

## Sentry
[Repo](https://github.com/getsentry/sentry-go)

### Environment
- `SENTRY_DSN` - The DSN to use. If the DSN is not set, the client is effectively disabled
- `SENTRY_RELEASE` - The release to be sent with events
- `SENTRY_ENVIRONMENT` - The environment to be sent with events

## Slack
[Repo](https://github.com/nlopes/slack)

### Environment
- `NOTIFY_SLACK_TOKEN` - Webhook URL
- `NOTIFY_SLACK_DEVICE_URL` - Device admin URL
- `NOTIFY_SLACK_USERNAME` - Notification user name
- `NOTIFY_SLACK_ICON` - Notification user icon (for each project)
- `NOTIFY_SLACK_PRETEXT` - Notification pretext (for each project service, description). For example: app, itunes service, android service, :android: Android RTN, A/B tests, Google Ads
- `NOTIFY_SLACK_LOG_IP` - Log client IP
