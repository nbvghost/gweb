package gweb

import (
	"net/http"
	_ "net/http/pprof"
	"testing"
)

func Test_fileNetLoad(t *testing.T) {
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			(Static{}).fileNetLoad(tt.args.writer, tt.args.request)
		})
	}
}

func Test_fileLoad(t *testing.T) {
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			(Static{}).fileLoad(tt.args.writer, tt.args.request)
		})
	}
}

func Test_fileTempLoad(t *testing.T) {
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			(Static{}).fileTempLoad(tt.args.writer, tt.args.request)
		})
	}
}
