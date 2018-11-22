package gweb

import (
	"fmt"
	"testing"
)

func TestAll(t *testing.T) {
	fmt.Println(t)
}

func TestStartServer(t *testing.T) {
	type args struct {
		HTTP  bool
		HTTPS bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartServer(tt.args.HTTP, tt.args.HTTPS)
		})
	}
}
