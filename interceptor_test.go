package gweb

import (
	"reflect"
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

func TestInterceptors_ExecuteAll(t *testing.T) {
	type args struct {
		c *BaseController
	}
	tests := []struct {
		name  string
		inter *Interceptors
		args  args
		want  bool
		want1 Result
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.inter.ExecuteAll(tt.args.c)
			if got != tt.want {
				t.Errorf("Interceptors.ExecuteAll() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Interceptors.ExecuteAll() got1 = %v, want %v", got1, tt.want1)
			}
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
