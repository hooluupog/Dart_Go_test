package main

func swap(a []rune, i int, j int) {
	a[i], a[j] = a[j], a[i]
}

func all(s []rune, n int) {
	if n == 1 {
		//fmt.Println(string(s))
	} else {
		for i := 0; i < n; i++ {
			swap(s, i, n-1)
			all(s, n-1)
			swap(s, i, n-1)
		}
	}
}

func main() {
	a := "♠行到水穷处，坐看云起时"
	all([]rune(a), len([]rune(a)))
}
