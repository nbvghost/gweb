package gweb

import (
	"reflect"
	"testing"
)

func TestAttributes_Put(t *testing.T) {
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name string
		att  *Attributes
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.att.Put(tt.args.key, tt.args.value)
		})
	}
}

func TestAttributes_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		att  *Attributes
		args args
		want interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.att.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Attributes.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributes_Delete(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		att  *Attributes
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.att.Delete(tt.args.key)
		})
	}
}
