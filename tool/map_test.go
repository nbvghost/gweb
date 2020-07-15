package tool

import (
	"reflect"
	"testing"
)

/*
性能测试：
1		mapKey =	705840 op			1578 ns/op
2		mapKey =	444418 op			2480 ns/op
3		mapKey =	359131 op			3297 ns/op
10		mapKey =	108100 op			11055 ns/op
*/
func BenchmarkDeepCopyMap(b *testing.B) {

	for i := 0; i < b.N; i++ {

		DeepCopyMap(map[string]interface{}{
			"dsfsd": 555,
			"dsfs1": 555,
			"dsfs2": 555,
			"dsfs3": 555,
			"dsfs4": 555,
			"dsfs5": 555,
			"dsfs6": 555,
			"dsfs7": 555,
			"dsfs8": 555,
			"dsfs9": 555,
		})

	}

}
func TestDeepCopyMap(t *testing.T) {
	type args struct {
		source interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantTarget interface{}
	}{
		{name: "TestDeepCopyMap", args: args{source: map[string]interface{}{"dsfsd": [2]int{55, 55}}}, wantTarget: map[string]interface{}{"dsfsd": [2]int{55, 55}}},
		{name: "TestDeepCopyMap", args: args{source: map[string]interface{}{"dsfsd": []int{55, 55}}}, wantTarget: map[string]interface{}{"dsfsd": []int{55, 55}}},
		{name: "TestDeepCopyMap", args: args{source: map[string]interface{}{"dsfsd": "[]int{55,55}"}}, wantTarget: map[string]interface{}{"dsfsd": "[]int{55,55}"}},
		{name: "TestDeepCopyMap", args: args{source: map[string]interface{}{}}, wantTarget: map[string]interface{}{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTarget := DeepCopyMap(tt.args.source); !reflect.DeepEqual(gotTarget, tt.wantTarget) {
				t.Errorf("DeepCopyMap() = %v, want %v", gotTarget, tt.wantTarget)
			} else {
				tt.args.source.(map[string]interface{})["dsfsd"] = 777
				t.Log(gotTarget.(map[string]interface{}))
			}
		})
	}
}
