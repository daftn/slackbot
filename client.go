package slackbot

import (
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"strings"
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
	for _, c := range channels {
		if c.Name == strings.TrimPrefix(identifier, "#") || c.ID == identifier {
			return c, nil
		}
	}
	return slack.Channel{}, errors.Errorf("unable to find channel id for %s", identifier)
}

func (s *slackClient) GetUser(identifier string) (slack.User, error) {
	users, err := s.getUsers()
	if err != nil {
		return slack.User{}, err
	}
	for _, u := range users {
		if u.Name == strings.TrimPrefix(identifier, "@") || u.ID == identifier || u.RealName == identifier {
			return u, nil
		}
	}
	return slack.User{}, nil
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
