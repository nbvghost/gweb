package gweb

import (
	"fmt"
	"testing"
	"net/http"
)

func TestAll(t *testing.T) {
	fmt.Println(t)
}

func TestStartServer(t *testing.T) {
	type args struct {
		HTTP  *http.Server
		HTTPS *http.Server
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartServer(http.DefaultServeMux,tt.args.HTTP, tt.args.HTTPS)
		})
	}
}
