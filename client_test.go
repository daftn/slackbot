package slackbot

import (
	"errors"
	"reflect"
	"testing"

	"github.com/slack-go/slack"
)

func Test_slackClient_GetChannel(t *testing.T) {
	type fields struct {
		RTM         *slack.RTM
		getChannels func(bool, ...slack.GetChannelsOption) ([]slack.Channel, error)
		getUsers    func() ([]slack.User, error)
	}
	type args struct {
		identifier string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    slack.Channel
		wantErr bool
	}{
		{
			name: "should return a channel",
			fields: fields{
				getChannels: func(b bool, option ...slack.GetChannelsOption) ([]slack.Channel, error) {
					return []slack.Channel{
						{
							GroupConversation: slack.GroupConversation{
								Name: "channel_name",
							},
						},
					}, nil
				},
			},
			args: args{
				identifier: "channel_name",
			},
			want: slack.Channel{
				GroupConversation: slack.GroupConversation{
					Name: "channel_name",
				},
			},
			wantErr: false,
		},
		{
			name: "should return an error if no channel matches",
			fields: fields{
				getChannels: func(b bool, option ...slack.GetChannelsOption) ([]slack.Channel, error) {
					return []slack.Channel{
						{
							GroupConversation: slack.GroupConversation{
								Name: "blah",
							},
						},
					}, nil
				},
			},
			args: args{
				identifier: "channel_name",
			},
			wantErr: true,
		},
		{
			name: "should return an error if getUsers errors",
			fields: fields{
				getChannels: func(b bool, option ...slack.GetChannelsOption) ([]slack.Channel, error) {
					return nil, errors.New("error")
				},
			},
			args: args{
				identifier: "channel_name",
			},
			wantErr: true,
		},
		{
			name: "should return a channel",
			fields: fields{
				getChannels: func(b bool, option ...slack.GetChannelsOption) ([]slack.Channel, error) {
					return []slack.Channel{
						{
							GroupConversation: slack.GroupConversation{
								Name: "blah",
							},
						},
					}, nil
				},
			},
			args: args{
				identifier: "#blah",
			},
			want: slack.Channel{
				GroupConversation: slack.GroupConversation{
					Name: "blah",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &slackClient{
				RTM:         tt.fields.RTM,
				getChannels: tt.fields.getChannels,
				getUsers:    tt.fields.getUsers,
			}
			got, err := s.GetChannel(tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChannel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChannel() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_slackClient_GetUser(t *testing.T) {
	type fields struct {
		RTM         *slack.RTM
		getChannels func(bool, ...slack.GetChannelsOption) ([]slack.Channel, error)
		getUsers    func() ([]slack.User, error)
	}
	type args struct {
		identifier string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    slack.User
		wantErr bool
	}{
		{
			name: "should get a user on real name match",
			fields: fields{
				getUsers: func() ([]slack.User, error) {
					return []slack.User{
						{
							ID:       "not_match",
							Name:     "not_match",
							RealName: "Should Match",
						},
					}, nil
				},
			},
			args: args{
				identifier: "Should Match",
			},
			want: slack.User{
				RealName: "Should Match",
				Name:     "not_match",
				ID:       "not_match",
			},
			wantErr: false,
		},
		{
			name: "should get a user on ID match",
			fields: fields{
				getUsers: func() ([]slack.User, error) {
					return []slack.User{
						{
							ID:       "should_match",
							Name:     "not_match",
							RealName: "Not Match",
						},
					}, nil
				},
			},
			args: args{
				identifier: "should_match",
			},
			want: slack.User{
				RealName: "Not Match",
				Name:     "not_match",
				ID:       "should_match",
			},
			wantErr: false,
		},
		{
			name: "should get a user on name match",
			fields: fields{
				getUsers: func() ([]slack.User, error) {
					return []slack.User{
						{
							ID:       "not_match",
							Name:     "should_match",
							RealName: "Not Match",
						},
					}, nil
				},
			},
			args: args{
				identifier: "should_match",
			},
			want: slack.User{
				RealName: "Not Match",
				Name:     "should_match",
				ID:       "not_match",
			},
			wantErr: false,
		},
		{
			name: "should get a user with @ in identifier",
			fields: fields{
				getUsers: func() ([]slack.User, error) {
					return []slack.User{
						{
							ID:       "not_match",
							Name:     "should_match",
							RealName: "Not Match",
						},
					}, nil
				},
			},
			args: args{
				identifier: "@should_match",
			},
			want: slack.User{
				RealName: "Not Match",
				Name:     "should_match",
				ID:       "not_match",
			},
			wantErr: false,
		},
		{
			name: "should return error if no user found",
			fields: fields{
				getUsers: func() ([]slack.User, error) {
					return []slack.User{
						{
							ID:       "not_match",
							Name:     "should_match",
							RealName: "Not Match",
						},
					}, nil
				},
			},
			args: args{
				identifier: "@should_not_find",
			},
			wantErr: true,
		},
		{
			name: "should return error if getChannels returns error",
			fields: fields{
				getUsers: func() ([]slack.User, error) {
					return nil, errors.New("error")
				},
			},
			args: args{
				identifier: "should_match",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &slackClient{
				RTM:         tt.fields.RTM,
				getChannels: tt.fields.getChannels,
				getUsers:    tt.fields.getUsers,
			}
			got, err := s.GetUser(tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
