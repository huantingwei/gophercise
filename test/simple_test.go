package test

import (
	"testing"
)

// go test -v
func TestFirst(t *testing.T) {

	t.Parallel()

	t.Logf("Test simple's function `first`\n")

	arr := first(3)
	if len(arr) != 3 {
		t.Fail()
	}
}
func TestSecond(t *testing.T) {

	t.Parallel()

	t.Logf("Test simple's function `second`\n")

	_, err := second("https://global.udn.com/global_vision/index")
	if err != nil {
		t.FailNow()
	}
}

// go test -bench ^(Benchmark*)
// (with memory)
// go test -benchmem -run=^$ -bench ^(Benchmark*)
// 評測預設的運行時間是一秒，如果在這個時間內，無法達到 b.N 的目標值，可以增加這個時間
// go test -benchmem -run=^$ -bench ^(Benchmark*) -benchtime=3s
func BenchmarkFirst(b *testing.B) {
	b.Skip()
	for i := 0; i < b.N; i++ {
		first(10)
	}
}

func BenchmarkSecond(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := second("https://global.udn.com/global_vision/index")
		if err != nil {
			b.FailNow()
		}
	}
}
