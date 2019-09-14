package slackbot

import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"regexp"
)

const firstStepIndex = 1

type Store interface {
	Put(key string, value interface{}) error
	Get(key string, value interface{}) error
	Delete(key string) error
}

type (
	// Exchanges are a way to have a back and forth conversation between a slack user and a slack bot.
	// When a user sends a message that matches the Regex specified in the exchange, the exchange with
	// the bot will be initiated in a thread on the original message.
	Exchange struct {

		// The Regex to match input from the user if the exchange is initiated through a message.
		Regex *regexp.Regexp

		// Usage describes how to use the exchange. It will be returned with GetHelp().
		Usage string

		// Map of steps in sequential order numbered from 1 -> n, with the step number as the key.
		// They must start with 1 and increase by one for each step.
		Steps map[int]*Step

		// A data store to allow data to be passed between steps.
		Store Store

		// A pointer to the bot that owns the exchange.
		Bot *Bot

		// Thread where the exchange is taking place.
		Thread string

		// Channel where the exchange is taking place.
		Channel string

		// User that initiated the exchange.
		User        string
		currentStep int
	}

	// Exchanges contain a list of Steps. Steps have three potential interaction methods: Message,
	// Handler, or MsgHandler. When a step is being executed, if a Message is set the message will
	// be sent and the exchange will move to the next step. If no message is set the Handler will
	// be checked, if it is set the Handler will be called. If the message and handler are not set,
	// the MsgHandler will be called. As the exchange moves to the next step if MsgHandler is the
	// interaction method, the MsgHandler will not be called until an incoming message event happens
	// on the exchange's thread.
	Step struct {

		// Name of the step, used for readability and in log messages.
		Name string

		// Message to be sent to exchange.Channel in exchange.Thread
		Message string

		// Handler function will be called if Message is not set on the step. If an error is returned
		// when the Handler is called the exchange will be terminated.
		Handler func(exchange *Exchange) error

		// MsgHandler function will be called if Message and Handler are not set on the step and
		// if there is an incoming message event on the exchange thread. If an error is returned
		// the exchange will be terminated. If retry is returned as true, the current step will
		// not increment, the exchange will wait for another incoming message event and the
		// MsgHandler will be retried.
		MsgHandler func(exchange *Exchange, event *slack.MessageEvent) (retry bool, err error)
	}
)

func (ex *Exchange) incrementCurrentStep() bool {
	next := ex.currentStep + 1
	if _, ok := ex.Steps[next]; ok {
		ex.currentStep = next
		return true
	}
	return false
}

func (ex *Exchange) continueExecution(ev *slack.MessageEvent) {
	step, err := ex.GetCurrentStep()
	initialStep := ex.currentStep
	if err != nil {
		ex.handleError(step, err)
		return
	}

	if step.Message != "" {
		ex.Reply(step.Message)
	} else if step.Handler != nil {
		if err := step.Handler(ex); err != nil {
			ex.handleError(step, err)
			return
		}
	} else if step.MsgHandler != nil && ev != nil {
		retry, err := step.MsgHandler(ex, ev)
		if retry {
			ex.continueExecution(nil)
			return
		}
		if err != nil {
			ex.handleError(step, err)
			return
		}
	} else {
		return
	}

	if initialStep == ex.currentStep && !ex.incrementCurrentStep() {
		delete(ex.Bot.activeExchanges, ex.Thread)
		return
	}
	ex.continueExecution(nil)
}

func (ex *Exchange) handleError(step *Step, err error) {
	msg := fmt.Sprintf("An error has occurred in exchange %s-%s, step %d %s: %s", ex.Channel, ex.Thread, ex.currentStep, step.Name, err)
	ex.Bot.LogDebug(msg)
	delete(ex.Bot.activeExchanges, ex.Thread)
}

// GetCurrentStep will get the current step. If there is no step in the exchange with the
// index of e.currentStep an error will be returned.
func (ex *Exchange) GetCurrentStep() (*Step, error) {
	if step, ok := ex.Steps[ex.currentStep]; ok {
		return step, nil
	}
	return nil, errors.New(fmt.Sprintf("exchange step with index %d not found", ex.currentStep))
}

// SkipToStep will change the exchanges current step to the number passed in. If the step
// does not exist an error will be returned.
func (ex *Exchange) SkipToStep(i int) error {
	if _, ok := ex.Steps[i]; ok {
		ex.currentStep = i
		return nil
	}
	return errors.New(fmt.Sprintf("exchange step with index %d not found", ex.currentStep))
}

// Terminate will remove the exchange from the bot's active exchanges list so the next steps will not be executed.
func (ex *Exchange) Terminate() {

	// TODO - figure out if there is a way to kill the currently executing step

	ex.Bot.LogDebug(fmt.Sprintf("killing exchange %s", ex.Thread))
	delete(ex.Bot.activeExchanges, ex.Thread)
}

// Reply will send a message to the exchange's channel and thread.
func (ex *Exchange) Reply(msg string) {
	ex.ReplyWithOptions(slack.MsgOptionText(msg, false))
}

// ReplyWithOptions will send a message to the exchange's channel and thread with the options specified.
// See Bot.ReplyWithOptions method for more information on sending messages with message options.
func (ex *Exchange) ReplyWithOptions(options ...slack.MsgOption) {
	options = append(options, slack.MsgOptionTS(ex.Thread))
	if _, _, err := ex.Bot.ReplyWithOptions(ex.Channel, options...); err != nil {
		if s, _ := ex.GetCurrentStep(); s != nil {
			ex.handleError(s, err)
		}
	}
}

// SendDefaultErrorMessage will send an error message to the exchanges channel/thread and return the error that was passed in.
func (ex *Exchange) SendDefaultErrorMessage(err error) error {
	ex.Reply(fmt.Sprintf("An unrecoverable error has occured. This exchange will be terminated.\nError: %s", err))
	return err
}
