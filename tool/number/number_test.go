package number

import (
    "math"
    "testing"
)

func TestParseInt(t *testing.T) {

    //fmt.Println(math.MaxInt64)

    type args struct {
        value interface{}
    }
    tests := []struct {
        name string
        args args
        want int
    }{
        {name: "TestParseInt0", args: args{value: int(4545)}, want: 4545},
        {name: "TestParseInt1", args: args{value: int8(21)}, want: 21},
        {name: "TestParseInt2", args: args{value: "4545"}, want: 4545},
        {name: "TestParseInt3", args: args{value: float64(41541)}, want: 41541},
        {name: "TestParseInt4", args: args{value: nil}, want: 0},
        {name: "TestParseInt5", args: args{value: "45d45"}, want: 0},
        {name: "TestParseInt6", args: args{value: uint64(math.MaxUint64)}, want: 9223372036854775807},
        {name: "TestParseInt7", args: args{value: math.MaxFloat64}, want: 9223372036854775807},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := ParseInt(tt.args.value); got != tt.want {
                t.Errorf("ParseInt() = %v, want %v", got, tt.want)
            }
        })
    }
}

func BenchmarkParseInt(b *testing.B) {

    for i := 0; i < b.N; i++ {
        ParseFloat("45d45")

    }

}
func TestParseFloat(t *testing.T) {
    type args struct {
        value interface{}
    }

    tests := []struct {
        name string
        args args
        want float64
    }{
        {name: "TestParseFloat0", args: args{value: int(4545)}, want: 4545},
        {name: "TestParseFloat1", args: args{value: int8(21)}, want: 21},
        {name: "TestParseFloat2", args: args{value: "4545"}, want: 4545},
        {name: "TestParseFloat3", args: args{value: float64(41541)}, want: 41541},
        {name: "TestParseFloat4", args: args{value: nil}, want: 0},
        {name: "TestParseFloat5", args: args{value: "45d45"}, want: 0},
        {name: "TestParseFloat6", args: args{value: uint64(math.MaxUint64)}, want: float64(uint64(math.MaxUint64))},
        {name: "TestParseFloat7", args: args{value: math.MaxFloat64}, want: math.MaxFloat64},
        {name: "TestParseFloat8", args: args{value: 454545}, want: 454545},
        {name: "TestParseFloat9", args: args{value: uint64(math.MaxUint64)}, want: float64(uint64(math.MaxUint64))},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := ParseFloat(tt.args.value); got != tt.want {
                t.Errorf("ParseFloat() = %v, want %v", got, tt.want)
            }
        })
    }
}
