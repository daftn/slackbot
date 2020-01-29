package slackbot

import (
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"strings"
)

const (
	channelPrefix = "#"
	userPrefix    = "@"
)

type slackClient struct {
	*slack.RTM
	getChannels func(bool, ...slack.GetChannelsOption) ([]slack.Channel, error)
	getUsers    func() ([]slack.User, error)
}

func (s *slackClient) GetChannel(identifier string) (slack.Channel, error) {
	channels, err := s.getChannels(true)
	if err != nil {
		return slack.Channel{}, err
	}
	i := strings.TrimPrefix(identifier, channelPrefix)
	for _, c := range channels {
		if c.Name == i || c.ID == i {
			return c, nil
		}
	}
	return slack.Channel{}, errors.Errorf("unable to find channel with identifier %s", identifier)
}

func (s *slackClient) GetUser(identifier string) (slack.User, error) {
	users, err := s.getUsers()
	if err != nil {
		return slack.User{}, err
	}
	i := strings.TrimPrefix(identifier, userPrefix)
	for _, u := range users {
		if u.Name == i || u.ID == i || u.RealName == i {
			return u, nil
		}
	}
	return slack.User{}, errors.Errorf("unable to find user with identifier %s", identifier)
}

func (s *slackClient) GetIncomingEvents() chan slack.RTMEvent {
	return s.RTM.IncomingEvents
}

func newSlackClient(token string) *slackClient {
	api := slack.New(token)
	c := &slackClient{
		api.NewRTM(),
		nil,
		nil,
	}
	c.getChannels = c.GetChannels
	c.getUsers = c.GetUsers
	return c
}
