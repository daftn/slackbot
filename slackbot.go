// Package slackbot allows for quick and easy creation of a slackbot.
//
// A bot can interact with users in multiple ways. DirectListeners listen for a command
// specified by regex and execute a corresponding handler. Exchanges allow for an
// interactive conversation between the user and a bot. They listen for a command
// to start the exchange at which point back and forth communication will happen in
// a thread. Scheduled tasks accept a cron expression and a handler function to run
// at the interval specified.
//
// Creating a bot is simple:
// 	func main() {
//		exampleListener := slackbot.Listener{
//	  		Usage: "this tells the user how to use this command",
//			Regex: regexp.MustCompile(`^(?i)(hello|hi|hey|howdy|hola)`),
//			Handler: func(bot *slackbot.Bot, ev *slack.MessageEvent) {
//				bot.Reply(ev.Channel, "Hi! I'm a rad slackbot")
//			},
//		}
//
//		bot := slackbot.Bot{
//			Token: "your_bots_api_token",
//			DirectListeners: []slackbot.Listener{exampleListener},
//		}
//
//		if err := bot.Start(); err != nil {
//			panic(err)
//		}
//	}
//
// For more examples see the /examples directory.
package slackbot

import (
	"bytes"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"github.com/ulule/deepcopier"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	defaultFallback           = "That is not a valid command..."
	circuitBreakerMessage     = "*CIRCUIT BREAKER TRIPPED*\nMore than %d messages were sent in under %d seconds\n\nSelf destruct sequence initiated. Goodbye."
	slackConnectionRetry      = 10
	slackConnectionRetrySleep = 500 * time.Millisecond
)

type (
	Bot struct {

		// Slack bot api token, see https://api.slack.com/bot-users
		Token string

		// Slack api client, through which all slack api interactions will happen.
		// Having the client available on the bot also allows all of the slack api
		// functions to be access by the bot in DirectListeners, Exchanges, and ScheduledTasks.
		API *slackClient

		// If a user chats the bot and the message does not match a regex for any DirectListeners
		// or Exchanges, the Fallback message will be sent as a reply. If FallbackMessage
		// is not set, the constant defaultFallback will be sent.
		FallbackMessage string

		// If the debug channel is set, any string passed to the bot.LogDebug(string) function will
		// be sent to the DebugChannel before being logged to std out.
		DebugChannel      string
		CircuitBreaker    *CircuitBreaker
		DirectListeners   []Listener
		IndirectListeners []Listener
		Exchanges         []Exchange
		ScheduledTasks    []ScheduledTask

		activeExchanges map[string]*Exchange
		userDetails     *slack.UserDetails
		once            sync.Once
	}

	// CircuitBreaker can prevent a bot from sending messages out of control. When a circuit
	// breaker is set on a bot, if more than MaxMessages are sent in the TimeInterval the bot
	// will stop sending messages and self destruct.
	CircuitBreaker struct {
		MaxMessages   int
		TimeInterval  time.Duration
		intervalStart time.Time
		count         int
	}

	// Listeners will listen for an incoming message that matches the Regex. When a match is
	// found the Handler function will be called. There are two types of listeners, direct and indirect.
	// Indirect listeners listed for all messages in channels that the bot is a member of. Messages can
	// match the regex and run the handler even if the message is not directed at the bot. Direct
	// listeners only match the regex and call the handler if the message was sent directly to the bot
	// either through a DM or by @-ing the bot in a channel.
	Listener struct {
		// A string to be presented to users describing how to use the listener.
		Usage   string
		Regex   *regexp.Regexp
		Handler func(bot *Bot, ev *slack.MessageEvent)
	}
)

func (bot *Bot) init() {
	if bot.API == nil {
		bot.API = newSlackClient(bot.Token)
	}
	if bot.FallbackMessage == "" {
		bot.FallbackMessage = defaultFallback
	}
	if bot.DebugChannel != "" {
		var ID string
		if c, err := bot.API.GetChannel(bot.DebugChannel); err != nil {
			if u, err := bot.API.GetUser(bot.DebugChannel); err == nil {
				ID = u.ID
			}
		} else {
			ID = c.ID
		}
		bot.DebugChannel = ID
	}
	bot.activeExchanges = make(map[string]*Exchange)
}

// Start will schedule any Scheduled Tasks on the bot, start managing connections and
// start listening for listener and exchange matches.
func (bot *Bot) Start() error {

	// TODO  - add validation for listeners, exchanges, scheduled tasks before the bot starts

	bot.once.Do(bot.init)
	if err := bot.scheduleTasks(); err != nil {
		return err
	}

	go bot.API.ManageConnection()

	retry := slackConnectionRetry
	for retry > 0 {
		if info := bot.API.GetInfo(); info != nil {
			bot.userDetails = info.User
			break
		}
		time.Sleep(slackConnectionRetrySleep)
		retry--
	}
	if retry == 0 {
		return errors.New("unable to make slack rtm connection")
	}

	bot.LogDebug(bot.buildStartingMessage())
	if err := bot.listen(); err != nil {
		return err
	}
	return nil
}

func (bot *Bot) scheduleTasks() error {
	s := scheduler{cron.New()}
	if err := s.scheduleTasks(bot, bot.ScheduledTasks); err != nil {
		return err
	}
	return nil
}

func (bot *Bot) buildStartingMessage() string {
	var msg strings.Builder
	msg.WriteString("```Starting bot with:\n")
	msg.WriteString(fmt.Sprintf("- %d Direct Listeners\n", len(bot.DirectListeners)))
	msg.WriteString(fmt.Sprintf("- %d Indirect Listeners\n", len(bot.IndirectListeners)))
	msg.WriteString(fmt.Sprintf("- %d Exchanges\n", len(bot.Exchanges)))
	msg.WriteString(fmt.Sprintf("- %d Scheduled Tasks\n", len(bot.ScheduledTasks)))
	if bot.DebugChannel != "" {
		msg.WriteString(fmt.Sprintf("- Debug Channel: %s\n", bot.DebugChannel))
	}
	if bot.FallbackMessage != "" {
		msg.WriteString(fmt.Sprintf("- Fallback Message: \"%s\"\n", bot.FallbackMessage))
	}
	if bot.CircuitBreaker != nil {
		msg.WriteString("- Circuit Breaker Enabled with:\n")
		msg.WriteString(fmt.Sprintf("	- max messages: %d\n", bot.CircuitBreaker.MaxMessages))
		msg.WriteString(fmt.Sprintf("	- interval: %s\n", bot.CircuitBreaker.TimeInterval))
	}
	msg.WriteString("```")
	return msg.String()
}

func (bot *Bot) listen() error {

	// TODO - accept a context in Start, add switch case for <- ctx.Done()

	for {
		select {
		case msg := <-bot.API.IncomingEvents:
			switch ev := msg.Data.(type) {

			case *slack.ConnectedEvent:
				log.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				go bot.processMessage(ev)

			case *slack.RTMError:
				log.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				log.Println("Invalid credentials")
				return errors.New("invalid slack credentials")
			}
		}
	}
}

func (bot *Bot) processMessage(ev *slack.MessageEvent) {
	for _, l := range bot.IndirectListeners {
		if l.Regex.MatchString(ev.Text) {
			if l.Handler != nil {
				l.Handler(bot, ev)
			}
		}
	}

	userPrefix := fmt.Sprintf("<@%s> ", bot.userDetails.ID)
	exchange, activeThread := bot.activeExchanges[ev.ThreadTimestamp]
	if ev.User != "" && ev.User != bot.userDetails.ID && ev.Text != "" &&
		(strings.HasPrefix(ev.Msg.Channel, "D") || strings.HasPrefix(ev.Text, userPrefix) || activeThread) {

		ev.Text = strings.TrimSpace(strings.TrimPrefix(ev.Text, userPrefix))

		if activeThread {
			exchange.continueExecution(ev)
			return
		}

		for _, e := range bot.Exchanges {
			if e.Regex.MatchString(ev.Text) {
				bot.startExchange(ev, &e)
				return
			}
		}
		for _, l := range bot.DirectListeners {
			if l.Regex.MatchString(ev.Text) {
				if l.Handler != nil {
					l.Handler(bot, ev)
				}
				return
			}
		}

		// If there are no exchanges or listeners that match the message, reply with the fallback message.
		if ev.ThreadTimestamp == "" {
			bot.Reply(ev.Channel, bot.FallbackMessage)
		}
	}
}

func (bot *Bot) checkCircuitBreaker(channel string) {
	if bot.CircuitBreaker != nil {
		bot.CircuitBreaker.count += 1
		if bot.CircuitBreaker.intervalStart.Before(time.Now().Add(-bot.CircuitBreaker.TimeInterval)) {
			bot.CircuitBreaker.intervalStart = time.Now()
			bot.CircuitBreaker.count = 1
		} else if bot.CircuitBreaker.count > bot.CircuitBreaker.MaxMessages {
			msg := fmt.Sprintf(circuitBreakerMessage, bot.CircuitBreaker.MaxMessages, bot.CircuitBreaker.TimeInterval/time.Second)
			bot.API.PostMessage(channel, slack.MsgOptionText(msg, false), slack.MsgOptionAsUser(true))
			os.Exit(-1)
		}
	}
}

func (bot *Bot) startExchange(ev *slack.MessageEvent, template *Exchange) {
	ex := &Exchange{}
	if err := deepcopier.Copy(template).To(ex); err != nil {
		bot.LogDebug(fmt.Sprintf("error starting exchange - %s", err))
		return
	}
	for i, step := range template.Steps {
		s := &Step{}
		if err := deepcopier.Copy(step).To(s); err != nil {
			bot.LogDebug(fmt.Sprintf("error starting exchange - %s", err))
			return
		}
		ex.Steps[i] = s
	}

	ex.Bot = bot
	ex.Thread = ev.Timestamp
	ex.Channel = ev.Channel
	ex.User = ev.User
	ex.currentStep = firstStepIndex
	ex.Store = SimpleStore{}
	bot.activeExchanges[ev.Timestamp] = ex
	ex.continueExecution(nil)
}

// LogDebug will send the log message to the bots DebugChannel if set and log the message to the console.
func (bot *Bot) LogDebug(msg string) {
	if bot.DebugChannel != "" {
		bot.checkCircuitBreaker(bot.DebugChannel)
		if _, _, err := bot.API.PostMessage(bot.DebugChannel, slack.MsgOptionText(msg, false), slack.MsgOptionAsUser(true)); err != nil {
			log.Printf("Error sending message to debug channel %s - %s", bot.DebugChannel, err)
		}
	}
	log.Println(msg)
}

// SendHelp will send a message containing all of the Listener and Exchange Usage strings. If msg is passed
// in it will be prepended to the usage help strings
func (bot *Bot) SendHelp(channel string, thread string, msg string) (respChannel string, timestamp string, err error) {
	var buffer bytes.Buffer
	if msg != "" {
		buffer.WriteString(msg + "\n")
	}
	for _, l := range bot.DirectListeners {
		if l.Usage != "" {
			buffer.WriteString(l.Usage + "\n")
		}
	}
	for _, e := range bot.Exchanges {
		if e.Usage != "" {
			buffer.WriteString(e.Usage + "\n")
		}
	}
	return bot.ReplyInThread(channel, thread, buffer.String())
}

// Reply will send a message to the channel specified.
func (bot *Bot) Reply(channel string, text string) (respChannel string, timestamp string, err error) {
	return bot.ReplyWithOptions(channel, slack.MsgOptionText(text, false))
}

// ReplyInThread will send a message to the channel and thread specified.
func (bot *Bot) ReplyInThread(channel string, thread string, text string) (respChannel string, timestamp string, err error) {
	return bot.ReplyWithOptions(channel, slack.MsgOptionText(text, false), slack.MsgOptionTS(thread))
}

// ReplyWithOptions will reply to the channel specified with the message options passed in.
// This is how you would send Attachments or other customizations on messages.
// These options are passed through to the /nlopes/slack package's PostMessage function. To
// see the available MsgOption functions see https://godoc.org/github.com/nlopes/slack#MsgOption
//
// Example:
// 	attachment := slack.Attachment{
//		Pretext: "some pretext",
//		Text:    "some text",
//		Fields: []slack.AttachmentField{
//			slack.AttachmentField{
//				Title: "a",
//				Value: "no",
//			},
//		},
// 	}
//
// 	bot.ReplyWithOptions("example_channel", slack.MsgOptionAttachments(attachment))
func (bot *Bot) ReplyWithOptions(channel string, options ...slack.MsgOption) (respChannel string, timestamp string, err error) {
	bot.checkCircuitBreaker(channel)
	options = append(options, slack.MsgOptionAsUser(true))
	c, t, e := bot.API.PostMessage(channel, options...)
	if e != nil {
		bot.LogDebug(fmt.Sprintf("failure sending message to %s with - %s", channel, e))
	}
	return c, t, e
}
