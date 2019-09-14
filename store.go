package slackbot

import (
	"bytes"
	"encoding/gob"
	"github.com/pkg/errors"
)

// SimpleStore is an optional store that can be used for the Store on an Exchange.
type SimpleStore map[string][]byte

func (s SimpleStore) Put(key string, value interface{}) error {
	if value == nil {
		return errors.New("error")
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return err
	}
	s[key] = buf.Bytes()
	return nil
}

func (s SimpleStore) Get(key string, value interface{}) error {
	v, ok := s[key]
	if !ok {
		return errors.Errorf("key %s not found", key)
	}
	d := gob.NewDecoder(bytes.NewReader(v))
	return d.Decode(value)
}

func (s SimpleStore) Delete(key string) error {
	_, ok := s[key]
	if !ok {
		return errors.Errorf("key %s not found", key)
	}
	delete(s, key)
	return nil
}
