package gweb

import (
	"testing"
)

func TestInterceptors_Add(t *testing.T) {
	type args struct {
		value Interceptor
	}
	tests := []struct {
		name  string
		inter *Interceptors
		args  args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.inter.AddInterceptor(tt.args.value)
		})
	}
}
