package main

import (
	"fmt"
	"sync"
	"testing"
)

const (
	N   = 100
	LEN = 30000000
)

type intSlice []int

func (is intSlice) Filter(f func(int) bool) intSlice {
	vs := make(intSlice, 0)
	for _, v := range is {
		if f(v) {
			vs = append(vs, v)
		}
	}
	return vs
}

func (is intSlice) Map(f func(int) int) intSlice {
	vs := make(intSlice, len(is))
	for i, v := range is {
		vs[i] = f(v)
	}
	return vs
}

func (is intSlice) Reduce(f func(int, int) int) int {
	var res int
	for _, v := range is {
		res = f(v, res)
	}
	return res
}

func loopSum(list []int) {
	res := 0
	for _, v := range list {
		if v%2 == 0 {
			res += v
		}
	}
}

func Sum(list []int) {
	is := intSlice(list)
	is.Filter(func(i int) bool { return i%2 == 0 }).
		Reduce(func(a, b int) int { return a + b })
}

func evenSum(wg *sync.WaitGroup, in <-chan []int, out chan<- int) {
	defer wg.Done()
	list := <-in
	res := 0
	for _, v := range list {
		if v%2 == 0 {
			res += v
		}
	}
	out <- res
}

func totalSum(sum <-chan int) {
	total := 0
	for i := 0; i < N; i++ {
		total += <-sum
	}
}

func main() {
	list := make([]int, LEN)
	in := make(chan []int, N)
	sum := make(chan int, N)
	defer close(sum)
	var wg sync.WaitGroup
	for i := range list {
		list[i] = i
	}

	res1 := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			loopSum(list)
		}
	})
	res2 := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Sum(list)
		}
	})

	res3 := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {

			for i, offset := 0, LEN/N; i < N; i++ {
				// Go does not support transferable data so it directly
				// shares the ownership of accessing list between goroutines
				in <- list[i*offset : i*offset+offset]
				wg.Add(1)
				go evenSum(&wg, in, sum)
			}
			wg.Wait()
			totalSum(sum)
		}
	})
	close(in)
	fmt.Println("Sum using loop     ", res1)
	fmt.Println("Sum using reduce   ", res2)
	fmt.Println("Parallel sum       ", res3)
}
