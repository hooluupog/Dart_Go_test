package controller

import "Test_Pig/model"

func 我是(x string, n int) {
	if n == 0 {
		return
	} else {
		pig.Print(x, n)
		我是(x, n-1)
	}
}

func Run(n int) {
	三 := n
	//猪头 := "猪头"
	猪头 := "pig.png"
	我是(猪头, 三)
}
