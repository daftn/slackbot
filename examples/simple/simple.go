package main

import (
	"github.com/nlopes/slack"
	"gitlab-app.eng.qops.net/derekn/slackbot"
	"regexp"
)

func main() {

	// This will set up a simple bot with defaults and a direct listener.
	// It is a working example. You can run it by entering your bot's api token
	// below and go run-ing the file. Once the bot comes up say hi to the bot
	// and it will respond with "Hi there, nice to meet you".

	apiToken := "put_your_token_here"

	bot := slackbot.Bot{
		Token: apiToken,
		DirectListeners: []slackbot.Listener{
			{
				Usage: "say hi and I'll respond",
				Regex: regexp.MustCompile(`^(?i)(hello|hi|hey|howdy|hola)`),
				Handler: func(bot *slackbot.Bot, ev *slack.MessageEvent) {
					bot.Reply(ev.Channel, "Hi there, nice to meet you")
				},
			},
		},
	}

	bot.Start()
}
