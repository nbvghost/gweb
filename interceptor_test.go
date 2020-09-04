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
			tt.inter.Add(tt.args.value)
		})
	}
}

func TestInterceptors_Contains(t *testing.T) {
	type args struct {
		interceptor Interceptor
	}
	tests := []struct {
		name  string
		inter *Interceptors
		args  args
		want  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.inter.Contains(tt.args.interceptor); got != tt.want {
				t.Errorf("Interceptors.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
