package sum

func genList(n int) []int {
	list := make([]int, n)
	for i := range list {
		list[i] = i
	}
	return list
}

func RemSum(n int) int {
	list := genList(n)
	res := 0
	for _, v := range list {
		if v%2 == 0 {
			res += v
		}
	}
	return res
}

func BitAndSum(n int) int {
	list := genList(n)
	res := 0
	for _, v := range list {
		if v&0x1 == 0 {
			res += v
		}
	}
	return res
}
