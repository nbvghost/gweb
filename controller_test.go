package gweb

import (
	"net/http"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "TestServeHTTP"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_fixPath1(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Test_fixPath1", args: args{path: "/sd/f/dsf/ds/f////sd/f/sd/fds//fsd"}, want: "/sd/f/dsf/ds/f/sd/f/sd/fds/fsd"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fixPath(tt.args.path); got != tt.want {
				t.Errorf("fixPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
