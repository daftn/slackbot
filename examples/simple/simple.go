package main

import (
	"github.com/daftn/slackbot"
	"github.com/nlopes/slack"
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
					_, _, _ = bot.Reply(ev.Channel, "Hi there, nice to meet you")
				},
			},
		},
	}

	err := bot.Start()
	if err != nil {
		panic("error starting bot")
	}
}
