package object

import (
	"fmt"
	"math"
	"testing"
)

func TestParseInt(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		//{name: "TestParseInt", args: args{value: int(4545)}, want: 4545},
		//{name: "TestParseInt", args: args{value: int8(21)}, want: 21},
		//{name: "TestParseInt", args: args{value: "4545"}, want: 4545},
		//{name: "TestParseInt", args: args{value: float64(41541)}, want: 41541},
		{name: "TestParseInt", args: args{value: nil}, want: 0},
		//{name: "TestParseInt", args: args{value: "45d45"}, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseInt(tt.args.value); got != tt.want {
				t.Errorf("ParseInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseString(t *testing.T) {

	let := math.MaxFloat32 - float64(1)
	fmt.Println(ParseString(let))

	type args struct {
		value interface{}
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "TestParseString0", args: args{value: 4545}, want: "4545"},
		{name: "TestParseString1", args: args{value: 4545.54545}, want: "4545.54545"},
		{name: "TestParseString2", args: args{value: math.MaxInt64}, want: "9223372036854775807"},
		{name: "TestParseString3", args: args{value: uint64(math.MaxUint64)}, want: "9223372036854775807"},
		{name: "TestParseString4", args: args{value: math.MaxFloat32}, want: "340282346638528860000000000000000000000"},
		{name: "TestParseString5", args: args{value: math.MaxFloat64}, want: "179769313486231570000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseString(tt.args.value); got != tt.want {
				t.Errorf("ParseString() = %v, want %v", got, tt.want)
			} else {
				t.Log(got)
			}
		})
	}
}
