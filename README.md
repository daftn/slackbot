
# Slackbot

A simple way to create an interactive slackbot. 

One purpose of this package is to allow interactive messages without having to expose 
an http server. Slack allows for [interactive messages](https://api.slack.com/interactive-messages) 
but requires that you to expose an http endpoint that slack can send requests to. This can be a hassle and 
potential security risk because the endpoint has to publicly accessible. This package 
allows for interactive messages through exchanges, which are explained in full detail below and do not require 
an http endpoint. The package also has other easily configurable methods of interaction. 

## Installation 
```
$ go get github.com/daftn/slackbot
```

## GoDocs
https://godoc.org/github.com/daftn/slackbot

## Usage

### Generating a Slack Token
Before you spin up your bot you will need an api token from Slack.  
For more info on setting up a bot in slack see: https://api.slack.com/bot-users

### Configuring the Bot
Create the bot with `bot := slackbot.Bot{}`
```golang
Bot struct {
    Token           string
    API             *slackClient
    FallbackMessage string
    DebugChannel    string
    CircuitBreaker  *CircuitBreaker

    DirectListeners   []Listener
    IndirectListeners []Listener
    Exchanges         []Exchange
    ScheduledTasks    []ScheduledTask
}
```
- **Token** - Slack bot api token, see https://api.slack.com/bot-users
- **API** - optional, this will be set automatically on the bot. 
Slack api client, through which all slack api interactions will happen. 
Having the client available on the bot also allows all of the slack api functions 
to be access by the bot in DirectListeners, Exchanges, and ScheduledTasks.
- **FallbackMessage** - optional, default is "That is not a valid command..." 
If a user directly chats the bot and the message does not match a regex for any DirectListeners 
or Exchanges, the Fallback message will be sent as a reply. If FallbackMessage is not 
set, the constant defaultFallback will be sent.
- **DebugChannel** - optional, if the debug channel is set, any string passed to the `bot.LogDebug(string)` 
function will be sent to the DebugChannel before being logged to std out.
- **CircuitBreaker** - optional, CircuitBreaker can prevent a bot from sending messages out of control. 
When a circuit breaker is set on a bot, if more than MaxMessages are sent in the TimeInterval the bot 
will stop sending messages and self destruct.

Bot also accepts interaction method lists for direct listeners, indirect listeners, 
exchanges, and scheduled tasks. See the Bot Interactions section below for descriptions of 
each of these interaction methods.  


#### Simple Example

```golang
bot := slackbot.Bot{
    Token:           "slack-token-here",
    FallbackMessage: "I couldn't find that command, try again",
    DebugChannel:    "general",
    CircuitBreaker: &slackbot.CircuitBreaker{
        MaxMessages:  30,
        TimeInterval: 10 * time.Second,
    },
    DirectListeners: []slackbot.Listener{
        {
            Usage:   "say hi and I'll respond",
            Regex:   regexp.MustCompile(`^(?i)(hello|hi|hey|howdy|hola)`),
            Handler: func(bot *slackbot.Bot, ev *slack.MessageEvent) {
                bot.Reply(ev.Channel, "Hi there, nice to meet you")
            },
        },
    },
}

bot.Start()
```


#### Full Examples
There are two fully working examples in the /examples dir. 
Comments in the files provide instructions for running the example bots. 

## Bot Interactions
There are 4 ways for the bot and slack users to interact:
- Direct Listeners
- Indirect Listeners 
- Exchanges
- Scheduled Tasks. 

### Listeners
Both direct listeners and indirect listeners implement the same interface.
```golang
type Listener struct {
       Usage   string
       Regex   *regexp.Regexp
       Handler func(bot *Bot, ev *slack.MessageEvent) 
}
```
**Usage** is a description for slack users detailing how this listener is used. 
**Regex** is the regex to look for that will trigger the listener. When an incoming 
message matches the regex, the **Handler** function will be called, passing in the bot and 
the message event that triggered the listener.   

#### Direct Listener
The listener's Handler will only be called if the user's message is 
sent directly to the bot, either through a direct message or by `@`-ing the bot in a channel of which 
the bot is a member and the message matches the `Regex` defined on the listener. 
They are added as a Listener list to the bot as `DirectListeners`.
```golang 
bot := slackbot.Bot{
    Token: apiToken,
    DirectListeners: []slackbot.Listener{
        {
            Usage:   "say hi and I'll respond",
            Regex:   regexp.MustCompile(`^(?i)(hello|hi|hey|howdy|hola)`),
            Handler: func(bot *slackbot.Bot, ev *slack.MessageEvent) {
                bot.Reply(ev.Channel, "Hi there, nice to meet you")
            },
        },
    },
}
```

#### Indirect Listener
The listener's Handler will be called if the Regex matches any message 
that is sent in a channel of which the bot is a member. The message does not have to be sent to the 
bot directly. They are added as a Listener list to the bot as `IndirectListeners`. 
```golang
bot := slackbot.Bot{
    Token: apiToken,
    IndirectListeners: []slackbot.Listener{
        {
            Usage:   "if anyone in the channel says 'trigger indirect listener' I'll respond",
            Regex:   regexp.MustCompile(`trigger indirect listener`),
            Handler: func(bot *slackbot.Bot, ev *slack.MessageEvent) {
                bot.Reply(ev.Channel, "You have triggered an indirect listener")
            },
        },
    },
}
```

### Exchange
Exchanges are a way to have a back and forth conversation between a slack user and a slack bot. 
When a user sends a message that matches the Regex specified in the exchange, the exchange with 
the bot will be initiated in a thread on the original message.  
```golang
Exchange struct {
    Regex *regexp.Regexp
    Usage string
    Steps map[int]*Step
    Store Store
    // there are other fields here also
}

Step struct {
    Name       string
    Message    string
    Handler    func(exchange *Exchange) error
    MsgHandler func(exchange *Exchange, event *slack.MessageEvent) (retry bool, err error)
} 
```
Exchanges contain a list of Steps. Steps have three possible handler types: Message, 
Handler, or MsgHandler. When a step is being executed, if a `Message` is set the message will 
be sent and the exchange will move to the next step. If no message is set for the step and 
the `Handler` function is not nil the Handler function will be called. If the message and handler are not set, 
the MsgHandler will be called. As the exchange moves to the next step if MsgHandler is the 
interaction method, the MsgHandler will not be called until an incoming message event happens 
on the exchange's thread.

See [Exchanges](https://godoc.org/github.com/daftn/slackbot#Exchange) in the godocs for 
functions available on the exchange that will be passed to the Handlers and MsgHandlers.

**Example**:  
```golang
slackbot.Exchange{
    Regex: regexp.MustCompile(`^(?i)start exchange`),
    Usage: "say start exchange and I'll ask you some questions",
    Steps: map[int]*slackbot.Step{
        1: {
            Name:    "send first question",
            Message: "What is your favorite color?",
        },
        2: {
            Name:       "receive favorite color",
            MsgHandler: func(ex *slackbot.Exchange, ev *slack.MessageEvent) (retry bool, err error) {
                if err := ex.Store.Put("color", ev.Text); err != nil {
                    return false, ex.SendDefaultErrorMessage(err)
                }
                return false, nil
            },
        },
        3: {
            Name:    "do something with the data collected",
            Handler: func(ex *slackbot.Exchange) error {
                var name, color string
                if err := ex.Store.Get("name", &name); err != nil {
                    return ex.SendDefaultErrorMessage(err)
                }
                ex.Reply(fmt.Sprintf("Guess what? %s is my favorite color too!", color))
                return nil
            },
        },
    },
}
```
In the example above the steps will be executed in this manner:
1. The `Message` in the first step will be sent to the exchange's thread, and move to the next step. 
2. The second step has a `MsgHandler` so it will pause and wait for the user to post a message in 
the thread before running the `MsgHandler`. After a message is posted in the thread, the users message will 
be passed to the `MsgHandler`. It will then continue to the next step. 
3. The third step has a `Handler` so it will not wait for user input and will be run immediately, and the 
exchange will be complete. 

### Scheduled Task
Scheduled tasks will run a Task function on a cron schedule.
```golang
ScheduledTask struct {
    Schedule string
    Task     func(*Bot)
}
```

**Schedule** takes a cron string defining the times that the **Task** function should run. When the task 
function is executed the bot will be passed to the function. The scheduled tasks will be scheduled 
when the bot is started with `bot.Start()`.

**Example**:
```golang 
slackbot.ScheduledTask{
    Schedule: "0 8 * * *",
    Task:     func(bot *slackbot.Bot) {
        bot.Reply("general", "Hey, its 8am on Monday just in case you were wondering.")
    },
}
```
