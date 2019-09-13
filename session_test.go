package gweb

import (
	"reflect"
	"testing"
)

func Test_sessionMap_DelectSession(t *testing.T) {
	type args struct {
		k string
	}
	tests := []struct {
		name string
		s    *SessionSafeMap
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.DeleteSession(tt.args.k)
		})
	}
}

func Test_sessionMap_addSession(t *testing.T) {
	type args struct {
		GLSESSIONID string
		session     *Session
	}
	tests := []struct {
		name string
		s    *SessionSafeMap
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.AddSession(tt.args.GLSESSIONID, tt.args.session)
		})
	}
}

func Test_sessionMap_GetSession(t *testing.T) {
	type args struct {
		GLSESSIONID string
	}
	tests := []struct {
		name string
		s    *SessionSafeMap
		args args
		want *Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetSession(tt.args.GLSESSIONID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sessionMap.GetSession() = %v, want %v", got, tt.want)
			}
		})
	}
}
