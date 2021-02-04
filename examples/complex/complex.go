package main

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/daftn/slackbot"
	"github.com/slack-go/slack"
)

// A client to be used to demonstrate using a closure for more functionality.
// For example, this could be a jira or github client that you want to access in the handler of a listener.
type closureExampleClient struct {
	APIToken   string
	HTTPClient *http.Client
}

func main() {

	// This will set up a fully configured bot with an example of each type of interaction.
	// It is a working example. You can run it by entering your bot's api token and
	// and debug channel then go run it.

	apiToken := "add_your_slack_token"
	debugChannel := "add_your_debug_channel"

	exampleClient := &closureExampleClient{
		APIToken:   "random_token",
		HTTPClient: http.DefaultClient,
	}

	bot := slackbot.Bot{
		Token:           apiToken,
		FallbackMessage: "I couldn't find that command, try again",
		DebugChannel:    debugChannel,
		CircuitBreaker: &slackbot.CircuitBreaker{
			MaxMessages:  30,
			TimeInterval: 10 * time.Second,
		},

		DirectListeners:   buildDirectListeners(exampleClient),
		IndirectListeners: buildIndirectListeners(),
		Exchanges:         buildExchanges(),
		ScheduledTasks:    buildScheduledTasks(),
	}

	if err := bot.Start(); err != nil {
		panic(err)
	}
}

// DIRECT LISTENERS
// The handler for these listeners will only be called if the message was
// directly sent to the bot, either in a DM or by @-ing the bot in a channel.
// Also notice that we are passing in the ClosureExample client from main
// so that we can access it in our listener.

func buildDirectListeners(c *closureExampleClient) []slackbot.Listener {
	return []slackbot.Listener{
		{
			Usage: "say hi",
			Regex: regexp.MustCompile(`^(?i)(hello|hi|hey|howdy|hola)`),
			Handler: func(bot *slackbot.Bot, ev *slack.MessageEvent) {
				bot.Reply(ev.Channel, "hey there")
			},
		},
		{
			Usage: "ask for help",
			Regex: regexp.MustCompile(`^(?i)help`),
			Handler: func(bot *slackbot.Bot, ev *slack.MessageEvent) {
				bot.SendHelp(ev.Channel, ev.ThreadTimestamp, "Here are my available commands:")
			},
		},
		{
			Usage:   "say something and Ill use a closure to build the handler",
			Regex:   regexp.MustCompile(`^(?i)something`),
			Handler: getWrappedHandler(c),
		},
	}
}

func getWrappedHandler(c *closureExampleClient) func(bot *slackbot.Bot, ev *slack.MessageEvent) {

	return func(bot *slackbot.Bot, ev *slack.MessageEvent) {

		// in here I can use the client or whatever else I pass to the function
		// so it will be used when the handler is called
		resp, _ := c.HTTPClient.Get("https://jsonplaceholder.typicode.com/todos/1")
		bot.Reply(ev.Channel, fmt.Sprintf("Here is the response from the client - %v", resp))
	}

}

// INDIRECT LISTENERS
// The handler for these listeners will be called if the regex matches any
// message sent in a channel of which the bot is a member.

func buildIndirectListeners() []slackbot.Listener {
	return []slackbot.Listener{
		{
			Usage: "if you start any message in the channel with 'indirect', ill respond",
			Regex: regexp.MustCompile(`^(?i)indirect`),
			Handler: func(bot *slackbot.Bot, ev *slack.MessageEvent) {
				bot.Reply(ev.Channel, "You triggered an indirect listener")
			},
		},
	}
}

// EXCHANGES
// An exchange is started when a the exchange's regex is matched in a message sent to
// the bot either in a DM or by @-ing the bot in a channel of which it is a member.
// The bot will start the exchange in a thread and listen for messages from the user in that thread.

func buildExchanges() []slackbot.Exchange {
	return []slackbot.Exchange{
		{
			Regex: regexp.MustCompile(`^(?i)start exchange`),
			Usage: "say start exchange and I'll ask you some questions",
			Steps: map[int]*slackbot.Step{
				1: {
					Name:    "send first question",
					Message: "What is your favorite color?",
				},
				2: {
					Name: "receive favorite color",
					// Steps with a MsgHandler, will pause and wait for a message in the exchange
					// thread before calling the MsgHandler function. The incoming message will
					// be passed to the msgHandler.
					MsgHandler: func(ex *slackbot.Exchange, ev *slack.MessageEvent) (retry bool, err error) {
						if err := ex.Store.Put("color", ev.Text); err != nil {
							return false, ex.SendDefaultErrorMessage(err)
						}
						ex.Reply("ok, got it, what is your name?")
						return false, nil
					},
				},
				3: {
					Name: "receive name",
					MsgHandler: func(ex *slackbot.Exchange, ev *slack.MessageEvent) (retry bool, err error) {
						if err := ex.Store.Put("name", ev.Text); err != nil {
							return false, ex.SendDefaultErrorMessage(err)
						}
						ex.Reply("ok, that's a rad name")
						return false, nil
					},
				},
				4: {
					Name: "do something with the data collected",
					Handler: func(ex *slackbot.Exchange) error {
						// get data collected from step 2 and 3 from the store
						var name, color string
						if err := ex.Store.Get("name", &name); err != nil {
							return ex.SendDefaultErrorMessage(err)
						}
						if err := ex.Store.Get("color", &color); err != nil {
							return ex.SendDefaultErrorMessage(err)
						}

						ex.Reply(fmt.Sprintf("Guess what %s? %s is my favorite color too!", name, color))
						return nil
					},
				},
			},
		},
	}
}

// SCHEDULED TASKS
// Scheduled tasks will be scheduled when bot.Start() is called. The Task function for
// a scheduled task will be called at the interval specified by the Schedule in cron format.

func buildScheduledTasks() []slackbot.ScheduledTask {
	return []slackbot.ScheduledTask{
		{
			Schedule: "0 8 * * *",
			Task: func(bot *slackbot.Bot) {
				bot.Reply("general", "Hey, its 8am on Monday just in case you were wondering.")
			},
		},
	}
}
