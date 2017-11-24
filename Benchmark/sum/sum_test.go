package sum

import (
	"testing"
)

func expect(t *testing.T, result int, expected int) {
	if result != expected {
		t.Errorf("sum = [%d] , expected [%d]", result, expected)
	}
}

func TestRemSum(t *testing.T) {
	expect(t, 0, RemSum(0))
	expect(t, 0, RemSum(1))
	expect(t, 249500, RemSum(1000))
	expect(t, 24995000, RemSum(10000))
	expect(t, 224999985000000, RemSum(30000000))
}

func TestBitAndSum(t *testing.T) {
	expect(t, 0, BitAndSum(0))
	expect(t, 0, BitAndSum(1))
	expect(t, 249500, BitAndSum(1000))
	expect(t, 24995000, BitAndSum(10000))
	expect(t, 224999985000000, BitAndSum(30000000))
}

var result int

func benchmarkRemSum(i int, b *testing.B) {
	var res int
	for i := 0; i < b.N; i++ {
		res = RemSum(i)
	}
	result = res
}

func benchmarkBitAndSum(i int, b *testing.B) {
	var res int
	for i := 0; i < b.N; i++ {
		res = BitAndSum(i)
	}
	result = res
}

func BenchmarkRemSum0(b *testing.B)        { benchmarkRemSum(0, b) }
func BenchmarkRemSum1(b *testing.B)        { benchmarkRemSum(1, b) }
func BenchmarkRemSum1000(b *testing.B)     { benchmarkRemSum(1000, b) }
func BenchmarkRemSum10000(b *testing.B)    { benchmarkRemSum(10000, b) }
func BenchmarkRemSum30000000(b *testing.B) { benchmarkRemSum(30000000, b) }

func BenchmarkBitAndSum0(b *testing.B)        { benchmarkBitAndSum(0, b) }
func BenchmarkBitAndSum1(b *testing.B)        { benchmarkBitAndSum(1, b) }
func BenchmarkBitAndSum1000(b *testing.B)     { benchmarkBitAndSum(1000, b) }
func BenchmarkBitAndSum10000(b *testing.B)    { benchmarkBitAndSum(10000, b) }
func BenchmarkBitAndSum30000000(b *testing.B) { benchmarkBitAndSum(30000000, b) }
