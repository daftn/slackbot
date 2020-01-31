package slackbot

import (
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"regexp"
	"sync"
	"testing"
	"time"
)

type mockAPI struct {
	*slack.RTM
	postMessage      func(string, ...slack.MsgOption) (string, string, error)
	getInfo          func() *slack.Info
	manageConnection func()
}

func (m *mockAPI) PostMessage(ch string, opts ...slack.MsgOption) (string, string, error) {
	return m.postMessage(ch, opts...)
}

func (m *mockAPI) GetChannel(identifier string) (slack.Channel, error) {
	return slack.Channel{}, errors.New("unable to find channel with identifier")
}

func (m *mockAPI) GetUser(identifier string) (slack.User, error) {
	return slack.User{}, errors.New("unable to find user with identifier")
}

func (m *mockAPI) GetIncomingEvents() chan slack.RTMEvent {
	return nil
}

func (m *mockAPI) GetInfo() *slack.Info {
	return m.getInfo()
}

func (m *mockAPI) ManageConnection() {
	m.manageConnection()
}

func TestBot_LogDebug(t *testing.T) {
	messageSent := false
	type fields struct {
		API          MessagingClient
		DebugChannel string
	}
	type args struct {
		msg string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		sent   bool
	}{
		{
			name: "should send message to debug channel",
			fields: fields{
				API: &mockAPI{
					postMessage: func(s string, opts ...slack.MsgOption) (string, string, error) {
						messageSent = true
						return "", "", nil
					},
				},
				DebugChannel: "test_channel",
			},
			args: args{
				msg: "the message",
			},
			sent: true,
		},
		{
			name: "should log and not send message to debug channel",
			fields: fields{
				API: &mockAPI{
					postMessage: func(s string, opts ...slack.MsgOption) (string, string, error) {
						messageSent = true
						return "", "", nil
					},
				},
				DebugChannel: "",
			},
			args: args{
				msg: "the message",
			},
			sent: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messageSent = false
			bot := &Bot{
				API:          tt.fields.API,
				DebugChannel: tt.fields.DebugChannel,
			}
			bot.LogDebug(tt.args.msg)
			if tt.sent != messageSent {
				t.Errorf("message sent status incorrect, got = %v, want %v", messageSent, tt.sent)
			}
		})
	}
}

func TestBot_ReplyWithOptions(t *testing.T) {
	type fields struct {
		API MessagingClient
	}
	type args struct {
		channel string
		options []slack.MsgOption
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantRespChannel string
		wantTimestamp   string
		wantErr         bool
	}{
		{
			name: "should return successfully",
			fields: fields{
				API: &mockAPI{
					postMessage: func(s string, opts ...slack.MsgOption) (string, string, error) {
						return "foo", "bar", nil
					},
				},
			},
			args: args{
				channel: "channel",
			},
			wantRespChannel: "foo",
			wantTimestamp:   "bar",
			wantErr:         false,
		},
		{
			name: "should return an error if post message fails",
			fields: fields{
				API: &mockAPI{
					postMessage: func(s string, opts ...slack.MsgOption) (string, string, error) {
						return "foo", "bar", errors.New("here is an error")
					},
				},
			},
			args: args{
				channel: "channel",
			},
			wantRespChannel: "foo",
			wantTimestamp:   "bar",
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bot := &Bot{
				API: tt.fields.API,
			}
			gotRespChannel, gotTimestamp, err := bot.ReplyWithOptions(tt.args.channel, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplyWithOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRespChannel != tt.wantRespChannel {
				t.Errorf("ReplyWithOptions() gotRespChannel = %v, want %v", gotRespChannel, tt.wantRespChannel)
			}
			if gotTimestamp != tt.wantTimestamp {
				t.Errorf("ReplyWithOptions() gotTimestamp = %v, want %v", gotTimestamp, tt.wantTimestamp)
			}
		})
	}
}

func TestBot_SendHelp(t *testing.T) {
	type fields struct {
		API MessagingClient
	}
	type args struct {
		channel string
		thread  string
		msg     string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantRespChannel string
		wantTimestamp   string
		wantErr         bool
	}{
		{
			name: "should send help",
			fields: fields{
				API: &mockAPI{
					postMessage: func(s string, opts ...slack.MsgOption) (string, string, error) {
						return "chan", "ts", nil
					},
				},
			},
			args: args{
				channel: "chan",
				thread:  "thre",
				msg:     "mess",
			},
			wantRespChannel: "chan",
			wantTimestamp:   "ts",
			wantErr:         false,
		},
		{
			name: "should return error",
			fields: fields{
				API: &mockAPI{
					postMessage: func(s string, opts ...slack.MsgOption) (string, string, error) {
						return "chan", "ts", errors.New("error")
					},
				},
			},
			args: args{
				channel: "chan",
				thread:  "thre",
				msg:     "mess",
			},
			wantRespChannel: "chan",
			wantTimestamp:   "ts",
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bot := &Bot{
				API: tt.fields.API,
			}
			gotRespChannel, gotTimestamp, err := bot.SendHelp(tt.args.channel, tt.args.thread, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendHelp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRespChannel != tt.wantRespChannel {
				t.Errorf("SendHelp() gotRespChannel = %v, want %v", gotRespChannel, tt.wantRespChannel)
			}
			if gotTimestamp != tt.wantTimestamp {
				t.Errorf("SendHelp() gotTimestamp = %v, want %v", gotTimestamp, tt.wantTimestamp)
			}
		})
	}
}

func TestBot_Start(t *testing.T) {
	type fields struct {
		Token             string
		API               MessagingClient
		FallbackMessage   string
		DebugChannel      string
		Store             Store
		CircuitBreaker    *CircuitBreaker
		DirectListeners   []Listener
		IndirectListeners []Listener
		Exchanges         []Exchange
		ScheduledTasks    []ScheduledTask
		activeExchanges   map[string]*Exchange
		userDetails       *slack.UserDetails
		once              sync.Once
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "should error on failed task scheduling",
			fields: fields{
				ScheduledTasks: []ScheduledTask{
					{
						Schedule: "alkmafsdlkmasdf",
						Task:     nil,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "should error on failed connection",
			fields: fields{
				API: &mockAPI{
					getInfo: func() *slack.Info {
						return nil
					},
					manageConnection: func() {},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bot := &Bot{
				Token:             tt.fields.Token,
				API:               tt.fields.API,
				FallbackMessage:   tt.fields.FallbackMessage,
				DebugChannel:      tt.fields.DebugChannel,
				Store:             tt.fields.Store,
				CircuitBreaker:    tt.fields.CircuitBreaker,
				DirectListeners:   tt.fields.DirectListeners,
				IndirectListeners: tt.fields.IndirectListeners,
				Exchanges:         tt.fields.Exchanges,
				ScheduledTasks:    tt.fields.ScheduledTasks,
				activeExchanges:   tt.fields.activeExchanges,
				userDetails:       tt.fields.userDetails,
				once:              tt.fields.once,
			}
			slackConnectionRetry = 1
			if err := bot.Start(); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBot_buildStartingMessage(t *testing.T) {
	type fields struct {
		Token             string
		API               MessagingClient
		FallbackMessage   string
		DebugChannel      string
		Store             Store
		CircuitBreaker    *CircuitBreaker
		DirectListeners   []Listener
		IndirectListeners []Listener
		Exchanges         []Exchange
		ScheduledTasks    []ScheduledTask
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "should return message",
			want: "```Starting bot with:\n- 0 Direct Listeners\n- 0 Indirect Listeners\n- 0 Exchanges\n- 0 Scheduled Tasks\n```",
		},
		{
			name: "should return message with comm counts",
			fields: fields{
				DirectListeners: []Listener{
					{
						Usage: "sdfgdfs",
					},
				},
				IndirectListeners: []Listener{
					{
						Usage: "test",
					},
				},
				Exchanges: []Exchange{
					{
						Usage: "blah",
					},
				},
				ScheduledTasks: []ScheduledTask{
					{
						Schedule: "",
						Task:     nil,
					},
				},
			},
			want: "```Starting bot with:\n- 1 Direct Listeners\n- 1 Indirect Listeners\n- 1 Exchanges\n- 1 Scheduled Tasks\n```",
		},
		{
			name: "should add debug channel to message",
			fields: fields{
				DebugChannel: "channel",
			},
			want: "```Starting bot with:\n- 0 Direct Listeners\n- 0 Indirect Listeners\n- 0 Exchanges\n- 0 Scheduled Tasks\n- Debug Channel: channel\n```",
		},
		{
			name: "should add fallback to message",
			fields: fields{
				DebugChannel:    "channel",
				FallbackMessage: "fallback",
			},
			want: "```Starting bot with:\n- 0 Direct Listeners\n- 0 Indirect Listeners\n- 0 Exchanges\n- 0 Scheduled Tasks\n- Debug Channel: channel\n- Fallback Message: \"fallback\"\n```",
		},
		{
			name: "should add circuit breaker details to message",
			fields: fields{
				CircuitBreaker: &CircuitBreaker{
					MaxMessages:  5,
					TimeInterval: 10 * time.Second,
				},
			},
			want: "```Starting bot with:\n- 0 Direct Listeners\n- 0 Indirect Listeners\n- 0 Exchanges\n- 0 Scheduled Tasks\n- Circuit Breaker Enabled with:\n	- max messages: 5\n	- interval: 10s\n```",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bot := &Bot{
				Token:             tt.fields.Token,
				API:               tt.fields.API,
				FallbackMessage:   tt.fields.FallbackMessage,
				DebugChannel:      tt.fields.DebugChannel,
				Store:             tt.fields.Store,
				CircuitBreaker:    tt.fields.CircuitBreaker,
				DirectListeners:   tt.fields.DirectListeners,
				IndirectListeners: tt.fields.IndirectListeners,
				Exchanges:         tt.fields.Exchanges,
				ScheduledTasks:    tt.fields.ScheduledTasks,
			}
			if got := bot.buildStartingMessage(); got != tt.want {
				t.Errorf("buildStartingMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBot_checkCircuitBreaker(t *testing.T) {
	terminateCalled := false
	type fields struct {
		Token             string
		API               MessagingClient
		FallbackMessage   string
		DebugChannel      string
		Store             Store
		CircuitBreaker    *CircuitBreaker
		DirectListeners   []Listener
		IndirectListeners []Listener
		Exchanges         []Exchange
		ScheduledTasks    []ScheduledTask
		activeExchanges   map[string]*Exchange
		terminate         func(int)
	}
	type args struct {
		channel string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		terminated bool
	}{
		{
			name: "should trip the breaker",
			fields: fields{
				CircuitBreaker: &CircuitBreaker{
					MaxMessages:   1,
					TimeInterval:  10,
					count:         10,
					intervalStart: time.Now().Add(time.Second * 100),
				},
				API: &mockAPI{
					postMessage: func(s string, opts ...slack.MsgOption) (string, string, error) {
						return "foo", "bar", nil
					},
				},
				terminate: func(i int) {
					terminateCalled = true
				},
			},
			args: args{
				channel: "ch",
			},
			terminated: true,
		},
		{
			name: "should skip the breaker",
			fields: fields{
				API: &mockAPI{
					postMessage: func(s string, opts ...slack.MsgOption) (string, string, error) {
						return "foo", "bar", nil
					},
				},
				terminate: func(i int) {
					terminateCalled = true
				},
			},
			args: args{
				channel: "ch",
			},
			terminated: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bot := &Bot{
				Token:             tt.fields.Token,
				API:               tt.fields.API,
				FallbackMessage:   tt.fields.FallbackMessage,
				DebugChannel:      tt.fields.DebugChannel,
				Store:             tt.fields.Store,
				CircuitBreaker:    tt.fields.CircuitBreaker,
				DirectListeners:   tt.fields.DirectListeners,
				IndirectListeners: tt.fields.IndirectListeners,
				Exchanges:         tt.fields.Exchanges,
				ScheduledTasks:    tt.fields.ScheduledTasks,
				activeExchanges:   tt.fields.activeExchanges,
				terminate:         tt.fields.terminate,
			}
			terminateCalled = false
			bot.checkCircuitBreaker(tt.args.channel)
			if terminateCalled != tt.terminated {
				t.Errorf("terminate called wrong, got = %v, want %v", terminateCalled, tt.terminated)
			}
		})
	}
}

func TestBot_processMessage(t *testing.T) {
	handlerCalled := false
	postMessageCalled := false
	type fields struct {
		Token             string
		API               MessagingClient
		FallbackMessage   string
		DebugChannel      string
		Store             Store
		CircuitBreaker    *CircuitBreaker
		DirectListeners   []Listener
		IndirectListeners []Listener
		Exchanges         []Exchange
		ScheduledTasks    []ScheduledTask
		activeExchanges   map[string]*Exchange
		userDetails       *slack.UserDetails
		once              sync.Once
	}
	type args struct {
		ev *slack.MessageEvent
	}
	type want struct {
		handlerCalled     bool
		postMessageCalled bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "should call indirect listener handler",
			fields: fields{
				IndirectListeners: []Listener{
					{
						Usage: "test listener",
						Handler: func(bot *Bot, ev *slack.MessageEvent) {
							handlerCalled = true
						},
						Regex: regexp.MustCompile(`here is the text`),
					},
				},
				userDetails: &slack.UserDetails{
					ID:             "myID",
					Name:           "",
					Created:        0,
					ManualPresence: "",
					Prefs:          slack.UserPrefs{},
				},
			},
			args: args{
				ev: &slack.MessageEvent{
					Msg: slack.Msg{
						Text: "here is the text",
					},
				},
			},
			want: want{
				handlerCalled: true,
			},
		},
		{
			name: "should call direct listener handler",
			fields: fields{
				DirectListeners: []Listener{
					{
						Usage: "test listener",
						Handler: func(bot *Bot, ev *slack.MessageEvent) {
							handlerCalled = true
						},
						Regex: regexp.MustCompile(`here is the text`),
					},
				},
				userDetails: &slack.UserDetails{
					ID: "myID",
				},
			},
			args: args{
				ev: &slack.MessageEvent{
					Msg: slack.Msg{
						Text: "<@myID> here is the text",
						User: "fff",
					},
				},
			},
			want: want{
				handlerCalled: true,
			},
		},
		{
			name: "should call direct listener handler",
			fields: fields{
				DirectListeners: []Listener{
					{
						Usage: "test listener",
						Handler: func(bot *Bot, ev *slack.MessageEvent) {
							handlerCalled = true
						},
						Regex: regexp.MustCompile(`here is the text`),
					},
				},
				userDetails: &slack.UserDetails{
					ID: "myID",
				},
			},
			args: args{
				ev: &slack.MessageEvent{
					Msg: slack.Msg{
						Text: "<@myID> here is the text",
						User: "fff",
					},
				},
			},
			want: want{
				handlerCalled: true,
			},
		},
		{
			name: "should reply with the default message",
			fields: fields{
				userDetails: &slack.UserDetails{
					ID: "myID",
				},
				API: &mockAPI{
					postMessage: func(s string, opts ...slack.MsgOption) (string, string, error) {
						postMessageCalled = true
						return "", "", nil
					},
				},
			},
			args: args{
				ev: &slack.MessageEvent{
					Msg: slack.Msg{
						Text: "<@myID> here is the text",
						User: "fff",
					},
				},
			},
			want: want{
				postMessageCalled: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bot := &Bot{
				Token:             tt.fields.Token,
				API:               tt.fields.API,
				FallbackMessage:   tt.fields.FallbackMessage,
				DebugChannel:      tt.fields.DebugChannel,
				Store:             tt.fields.Store,
				CircuitBreaker:    tt.fields.CircuitBreaker,
				DirectListeners:   tt.fields.DirectListeners,
				IndirectListeners: tt.fields.IndirectListeners,
				Exchanges:         tt.fields.Exchanges,
				ScheduledTasks:    tt.fields.ScheduledTasks,
				activeExchanges:   tt.fields.activeExchanges,
				userDetails:       tt.fields.userDetails,
				once:              tt.fields.once,
			}
			handlerCalled = false
			postMessageCalled = false
			bot.processMessage(tt.args.ev)
			if handlerCalled != tt.want.handlerCalled {
				t.Errorf("handler called wrong, got = %v, want %v", handlerCalled, tt.want.handlerCalled)
			}
			if postMessageCalled != tt.want.postMessageCalled {
				t.Errorf("post message called wrong, got = %v, want %v", postMessageCalled, tt.want.postMessageCalled)
			}
		})
	}
}

func TestBot_startExchange(t *testing.T) {
	type fields struct {
		Token             string
		API               MessagingClient
		FallbackMessage   string
		DebugChannel      string
		Store             Store
		CircuitBreaker    *CircuitBreaker
		DirectListeners   []Listener
		IndirectListeners []Listener
		Exchanges         []Exchange
		ScheduledTasks    []ScheduledTask
		activeExchanges   map[string]*Exchange
		userDetails       *slack.UserDetails
		once              sync.Once
	}
	type args struct {
		ev       *slack.MessageEvent
		template *Exchange
	}
	type want struct {
		key string
		ex  *Exchange
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "should start the exchange",
			fields: fields{
				activeExchanges: make(map[string]*Exchange),
			},
			args: args{
				ev: &slack.MessageEvent{
					Msg: slack.Msg{
						Channel:   "test_chan",
						User:      "test_user",
						Text:      "test_text",
						Timestamp: "here_is_the_timestamp",
					},
				},
				template: &Exchange{
					Regex: regexp.MustCompile(`test_text`),
					Usage: "here is the usage",
					Steps: map[int]*Step{
						1: {
							Name: "step 1",
							MsgHandler: func(ex *Exchange, ev *slack.MessageEvent) (bool, error) {
								return false, nil
							},
						},
						2: {
							Name: "step 2",
							MsgHandler: func(ex *Exchange, ev *slack.MessageEvent) (bool, error) {
								return false, nil
							},
						},
					},
				},
			},
			want: want{
				key: "here_is_the_timestamp",
				ex: &Exchange{
					Regex: regexp.MustCompile(`test_text`),
					Usage: "here is the usage",
					Steps: map[int]*Step{
						1: {
							Name: "step 1",
							MsgHandler: func(ex *Exchange, ev *slack.MessageEvent) (bool, error) {
								return false, nil
							},
						},
						2: {
							Name: "step 2",
							MsgHandler: func(ex *Exchange, ev *slack.MessageEvent) (bool, error) {
								return false, nil
							},
						},
					},
					Store:       SimpleStore{},
					Thread:      "here_is_the_timestamp",
					Channel:     "test_chan",
					User:        "test_user",
					currentStep: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bot := &Bot{
				Token:             tt.fields.Token,
				API:               tt.fields.API,
				FallbackMessage:   tt.fields.FallbackMessage,
				DebugChannel:      tt.fields.DebugChannel,
				Store:             tt.fields.Store,
				CircuitBreaker:    tt.fields.CircuitBreaker,
				DirectListeners:   tt.fields.DirectListeners,
				IndirectListeners: tt.fields.IndirectListeners,
				Exchanges:         tt.fields.Exchanges,
				ScheduledTasks:    tt.fields.ScheduledTasks,
				activeExchanges:   tt.fields.activeExchanges,
				userDetails:       tt.fields.userDetails,
				once:              tt.fields.once,
			}
			bot.startExchange(tt.args.ev, tt.args.template)
			ex, ok := bot.activeExchanges[tt.want.key]
			if !ok && tt.want.key != "" {
				t.Errorf("exchange not added to list of active exchanges")
			}
			if tt.want.ex.User != ex.User ||
				tt.want.ex.Channel != ex.Channel ||
				tt.want.ex.Thread != ex.Thread ||
				tt.want.ex.currentStep != ex.currentStep {
				t.Errorf("active exchange incorrect got = %v, want = %v", ex, tt.want.ex)
			}
		})
	}
}
