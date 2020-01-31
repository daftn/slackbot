package slackbot

import (
	"testing"
)

func TestSimpleStore_Delete(t *testing.T) {
	tests := []struct {
		name       string
		s          SimpleStore
		key        string
		finalCount int
		wantErr    bool
	}{
		{
			name:       "should error if trying to delete from empty store",
			s:          SimpleStore{},
			key:        "not_there",
			finalCount: 0,
			wantErr:    true,
		},
		{
			name:       "should successfully delete an item from the store",
			s:          SimpleStore{"test_entry": []byte("here it is"), "test_entry2": []byte("here it is again")},
			key:        "test_entry",
			finalCount: 1,
			wantErr:    false,
		},
		{
			name:       "should error if key does not exist",
			s:          SimpleStore{"test_entry": []byte("here it is")},
			key:        "incorrect",
			finalCount: 1,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.s.Delete(tt.key)
			count := len(tt.s)
			if (err != nil) != tt.wantErr || tt.finalCount != count {
				t.Errorf("Delete() error = %v, wantErr %v, count = %d, finalCount = %d", err, tt.wantErr, count, tt.finalCount)
			}
		})
	}
}

func TestSimpleStore_Get_and_Put(t *testing.T) {
	type set struct {
		key string
		val interface{}
	}
	type want struct {
		key    string
		val    interface{}
		putErr bool
		getErr bool
	}
	tests := []struct {
		name string
		set  set
		want want
	}{
		{
			name: "should put and get the correct string value",
			set: set{
				key: "the_key",
				val: "a string",
			},
			want: want{
				key:    "the_key",
				val:    "a string",
				putErr: false,
			},
		},
		{
			name: "should error on nil value",
			set: set{
				key: "the_key",
				val: nil,
			},
			want: want{
				putErr: true,
			},
		},
		{
			name: "should error on invalid key",
			set: set{
				key: "the_key",
				val: "a string",
			},
			want: want{
				key:    "wrong_key",
				putErr: false,
				getErr: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SimpleStore{}
			err := s.Put(tt.set.key, tt.set.val)
			if (err != nil) != tt.want.putErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.want.putErr)
			}

			if err == nil {
				var testString string
				err = s.Get(tt.want.key, &testString)
				if (err != nil) != tt.want.getErr {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.want.getErr)
				}
			}
		})
	}
}
