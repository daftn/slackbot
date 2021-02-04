package slackbot

import (
	"errors"
	"reflect"
	"regexp"
	"sync"
	"testing"

	"github.com/slack-go/slack"
)

func TestExchange_GetCurrentStep(t *testing.T) {
	type fields struct {
		Regex       *regexp.Regexp
		Usage       string
		Steps       map[int]*Step
		Store       Store
		Bot         *Bot
		Thread      string
		Channel     string
		User        string
		currentStep int
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Step
		wantErr bool
	}{
		{
			name: "should return step",
			fields: fields{
				currentStep: 1,
				Steps: map[int]*Step{
					1: {
						Name:    "test_name",
						Message: "test message",
					},
				},
			},
			want: &Step{
				Name:    "test_name",
				Message: "test message",
			},
			wantErr: false,
		},
		{
			name: "should error if step does not exist",
			fields: fields{
				currentStep: 2,
				Steps: map[int]*Step{
					1: {
						Name:    "test_name",
						Message: "test message",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ex := &Exchange{
				Regex:       tt.fields.Regex,
				Usage:       tt.fields.Usage,
				Steps:       tt.fields.Steps,
				Store:       tt.fields.Store,
				Bot:         tt.fields.Bot,
				Thread:      tt.fields.Thread,
				Channel:     tt.fields.Channel,
				User:        tt.fields.User,
				currentStep: tt.fields.currentStep,
			}
			got, err := ex.GetCurrentStep()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrentStep() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCurrentStep() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExchange_ReplyWithOptions(t *testing.T) {
	messageSent := false
	type fields struct {
		Regex       *regexp.Regexp
		Usage       string
		Steps       map[int]*Step
		Store       Store
		Bot         *Bot
		Thread      string
		Channel     string
		User        string
		currentStep int
	}
	type args struct {
		options []slack.MsgOption
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		activeCount int
		shouldSend  bool
	}{
		{
			name: "should error if reply fails",
			fields: fields{
				Bot: &Bot{
					API: &mockAPI{
						postMessage: func(s string, opts ...slack.MsgOption) (string, string, error) {
							return "", "", errors.New("error")
						},
					},
					activeExchanges: map[string]*Exchange{
						"test_thread": {
							Steps: map[int]*Step{
								1: {
									Name:    "s1",
									Message: "m1",
								},
								2: {
									Name:    "s2",
									Message: "m2",
								},
							},
						},
					},
				},
				Steps: map[int]*Step{
					1: {
						Name:    "s1",
						Message: "m1",
					},
					2: {
						Name:    "s2",
						Message: "m2",
					},
				},
				currentStep: 1,
				Thread:      "test_thread",
			},
			activeCount: 0,
			shouldSend:  false,
		},
		{
			name: "should send message",
			fields: fields{
				Bot: &Bot{
					API: &mockAPI{
						postMessage: func(s string, opts ...slack.MsgOption) (string, string, error) {
							messageSent = true
							return "", "", nil
						},
					},
					activeExchanges: map[string]*Exchange{
						"test_thread": {
							Steps: map[int]*Step{
								1: {
									Name:    "s1",
									Message: "m1",
								},
								2: {
									Name:    "s2",
									Message: "m2",
								},
							},
						},
					},
				},
				Steps: map[int]*Step{
					1: {
						Name:    "s1",
						Message: "m1",
					},
					2: {
						Name:    "s2",
						Message: "m2",
					},
				},
				currentStep: 1,
				Thread:      "test_thread",
			},
			activeCount: 1,
			shouldSend:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ex := &Exchange{
				Regex:       tt.fields.Regex,
				Usage:       tt.fields.Usage,
				Steps:       tt.fields.Steps,
				Store:       tt.fields.Store,
				Bot:         tt.fields.Bot,
				Thread:      tt.fields.Thread,
				Channel:     tt.fields.Channel,
				User:        tt.fields.User,
				currentStep: tt.fields.currentStep,
			}
			messageSent = false
			ex.ReplyWithOptions()
			if tt.activeCount != len(ex.Bot.activeExchanges) {
				t.Errorf("active exchange count wrong, got = %v, want %v", len(ex.Bot.activeExchanges), tt.activeCount)
			}
			if tt.shouldSend != messageSent {
				t.Errorf("incorrect message sent status, got = %v, want %v", messageSent, tt.shouldSend)
			}
		})
	}
}

func TestExchange_SkipToStep(t *testing.T) {
	type fields struct {
		Regex       *regexp.Regexp
		Usage       string
		Steps       map[int]*Step
		Store       Store
		Bot         *Bot
		Thread      string
		Channel     string
		User        string
		currentStep int
	}
	type args struct {
		i int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "should change step",
			fields: fields{
				Steps: map[int]*Step{
					1: {
						Name: "test_name",
						Handler: func(ex *Exchange) error {
							return nil
						},
					},
					2: {
						Name:    "test_name",
						Message: "test_message",
					},
					3: {
						Name:    "test_name",
						Message: "test_message",
					},
				},
				Bot:         nil,
				currentStep: 1,
			},
			args: args{
				i: 3,
			},
			wantErr: false,
		},
		{
			name: "should error if step doesnt exist",
			fields: fields{
				Steps: map[int]*Step{
					1: {
						Name: "test_name",
						Handler: func(ex *Exchange) error {
							return nil
						},
					},
					2: {
						Name:    "test_name",
						Message: "test_message",
					},
				},
				Bot:         nil,
				currentStep: 1,
			},
			args: args{
				i: 3,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ex := &Exchange{
				Regex:       tt.fields.Regex,
				Usage:       tt.fields.Usage,
				Steps:       tt.fields.Steps,
				Store:       tt.fields.Store,
				Bot:         tt.fields.Bot,
				Thread:      tt.fields.Thread,
				Channel:     tt.fields.Channel,
				User:        tt.fields.User,
				currentStep: tt.fields.currentStep,
			}
			if err := ex.SkipToStep(tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("SkipToStep() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExchange_Terminate(t *testing.T) {
	type fields struct {
		Regex       *regexp.Regexp
		Usage       string
		Steps       map[int]*Step
		Store       Store
		Bot         *Bot
		Thread      string
		Channel     string
		User        string
		currentStep int
	}
	tests := []struct {
		name        string
		fields      fields
		activeCount int
	}{
		{
			name: "should remove exchange",
			fields: fields{
				Bot: &Bot{
					activeExchanges: map[string]*Exchange{
						"test_thread": {
							Regex:       nil,
							Usage:       "",
							Steps:       nil,
							Store:       nil,
							Bot:         nil,
							Thread:      "",
							Channel:     "",
							User:        "",
							currentStep: 0,
						},
					},
					userDetails: nil,
					once:        sync.Once{},
				},
				Thread:      "test_thread",
				currentStep: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ex := &Exchange{
				Regex:       tt.fields.Regex,
				Usage:       tt.fields.Usage,
				Steps:       tt.fields.Steps,
				Store:       tt.fields.Store,
				Bot:         tt.fields.Bot,
				Thread:      tt.fields.Thread,
				Channel:     tt.fields.Channel,
				User:        tt.fields.User,
				currentStep: tt.fields.currentStep,
			}
			ex.Terminate()
			if tt.activeCount != len(ex.Bot.activeExchanges) {
				t.Errorf("active exchange count wrong, got = %v, want %v", len(ex.Bot.activeExchanges), tt.activeCount)
			}
		})
	}
}

func TestExchange_continueExecution(t *testing.T) {
	type fields struct {
		Steps       map[int]*Step
		Bot         *Bot
		Thread      string
		Channel     string
		User        string
		currentStep int
	}
	type args struct {
		ev *slack.MessageEvent
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		activeCount int
	}{
		{
			name: "should return error if getCurrentStep fails",
			fields: fields{
				Steps: map[int]*Step{
					1: {
						Name:    "test_name",
						Message: "test_message",
					},
				},
				Thread:      "test_thread",
				currentStep: 2,
				Bot: &Bot{
					activeExchanges: map[string]*Exchange{
						"test_thread": {
							Thread:  "test_thread",
							Channel: "test_channel",
							User:    "test_user",
						},
					},
				},
			},
			activeCount: 0,
		},
		{
			name: "should return error if handler errors",
			fields: fields{
				Steps: map[int]*Step{
					1: {
						Name: "test_name",
						Handler: func(ex *Exchange) error {
							return errors.New("error")
						},
					},
					2: {
						Name:    "test_name",
						Message: "test_message",
					},
				},
				Thread:      "test_thread",
				currentStep: 1,
				Bot: &Bot{
					activeExchanges: map[string]*Exchange{
						"test_thread": {
							Thread:  "test_thread",
							Channel: "test_channel",
							User:    "test_user",
						},
					},
				},
			},
			activeCount: 0,
		},
		{
			name: "should return error if message handler errors",
			fields: fields{
				Steps: map[int]*Step{
					1: {
						Name: "test_name",
						MsgHandler: func(ex *Exchange, ev *slack.MessageEvent) (bool, error) {
							return false, errors.New("error")
						},
					},
					2: {
						Name:    "test_name",
						Message: "test_message",
					},
				},
				Thread:      "test_thread",
				currentStep: 1,
				Bot: &Bot{
					activeExchanges: map[string]*Exchange{
						"test_thread": {
							Thread:  "test_thread",
							Channel: "test_channel",
							User:    "test_user",
						},
					},
				},
			},
			args: args{
				ev: &slack.MessageEvent{
					Msg: slack.Msg{
						Text: "here is the text",
					},
				}},
			activeCount: 0,
		},
		{
			name: "should finish exchange if on last step",
			fields: fields{
				Steps: map[int]*Step{
					1: {
						Name:    "test_name",
						Message: "message",
					},
					2: {
						Name: "test_name",
						Handler: func(ex *Exchange) error {
							return nil
						},
					},
				},
				Thread:      "test_thread",
				currentStep: 2,
				Bot: &Bot{
					activeExchanges: map[string]*Exchange{
						"test_thread": {
							Thread:  "test_thread",
							Channel: "test_channel",
							User:    "test_user",
						},
					},
				},
			},
			activeCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ex := &Exchange{
				Steps:       tt.fields.Steps,
				Bot:         tt.fields.Bot,
				Thread:      tt.fields.Thread,
				Channel:     tt.fields.Channel,
				User:        tt.fields.User,
				currentStep: tt.fields.currentStep,
			}
			ex.continueExecution(tt.args.ev)
			if tt.activeCount != len(ex.Bot.activeExchanges) {
				t.Errorf("active exchange count wrong, got = %v, want %v", len(ex.Bot.activeExchanges), tt.activeCount)
			}
		})
	}
}

func TestExchange_handleError(t *testing.T) {
	type fields struct {
		Regex       *regexp.Regexp
		Usage       string
		Steps       map[int]*Step
		Store       Store
		Bot         *Bot
		Thread      string
		Channel     string
		User        string
		currentStep int
	}
	type args struct {
		step *Step
		err  error
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		activeCount int
	}{
		{
			name: "should handle error and delete exchange",
			fields: fields{
				Bot: &Bot{
					DebugChannel: "",
					activeExchanges: map[string]*Exchange{
						"test_thread": {
							Thread:      "test_thread",
							Channel:     "test_channel",
							User:        "test_user",
							currentStep: 1,
						},
					},
				},
				Thread:  "test_thread",
				Channel: "test_channel",
				User:    "test_user",
			},
			args: args{
				step: &Step{
					Name:    "test_name",
					Message: "test_message",
				},
				err: errors.New("test error"),
			},
			activeCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ex := &Exchange{
				Regex:       tt.fields.Regex,
				Usage:       tt.fields.Usage,
				Steps:       tt.fields.Steps,
				Store:       tt.fields.Store,
				Bot:         tt.fields.Bot,
				Thread:      tt.fields.Thread,
				Channel:     tt.fields.Channel,
				User:        tt.fields.User,
				currentStep: tt.fields.currentStep,
			}
			ex.handleError(tt.args.step, tt.args.err)
			if tt.activeCount != ex.currentStep {
				t.Errorf("currentStep count is incorrect got %v want %v", ex.currentStep, tt.activeCount)
			}
		})
	}
}

func TestExchange_incrementCurrentStep(t *testing.T) {
	type fields struct {
		Steps       map[int]*Step
		currentStep int
	}
	tests := []struct {
		name        string
		fields      fields
		want        bool
		wantCurrent int
	}{
		{
			name: "should increment current step",
			fields: fields{
				Steps: map[int]*Step{
					1: {
						Name: "one",
					},
					2: {
						Name: "two",
					},
				},
				currentStep: 1,
			},
			want:        true,
			wantCurrent: 2,
		},
		{
			name: "should return false on last step",
			fields: fields{
				Steps: map[int]*Step{
					1: {
						Name: "one",
					},
				},
				currentStep: 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ex := &Exchange{
				Steps:       tt.fields.Steps,
				currentStep: tt.fields.currentStep,
			}
			if got := ex.incrementCurrentStep(); got != tt.want {
				t.Errorf("incrementCurrentStep() = %v, want %v", got, tt.want)
			}
			if tt.want && ex.currentStep != tt.wantCurrent {
				t.Errorf("currentStep = %v, want %v", ex.currentStep, tt.wantCurrent)
			}
		})
	}
}
