// 给出一个二维数组，请将这个二维数组按第i列(i从1开始)排序，如果第i列相同，
// 则对相同的行按第i+1列的元素排序，
// 如果第i+1列的元素也相同，则继续比较第i+2列，以此类推，直到最后一列。
// 如果第i列到最后一列都相同，则按原序排列。
// 样例输入：
// 1,4,0,0,0
// 2,3,5,0,0
// 2,3,4,6,0
// 1,3,4,5,2
// 1,3,4,5,1
// 按第2列排序，输出：
// 1,3,4,5,1
// 1,3,4,5,2
// 2,3,4,6,0
// 2,3,5,0,0
// 1,4,0,0,0

package main

import (
	"fmt"
	"sort"
)

func main() {
	a := [][]int{
		{1, 4, 0, 0, 0},
		{2, 3, 5, 0, 0},
		{2, 3, 4, 6, 0},
		{1, 3, 4, 5, 2},
		{1, 3, 4, 5, 1},
	}
	col := 1
	sort.SliceStable(a, func(i, j int) bool { return a[i][col] < a[j][col] })
	for j := col + 1; j < len(a[0]); j++ {
		sort.SliceStable(a, func(i, k int) bool {
			return a[i][j] < a[k][j] && a[i][j-1] == a[k][j-1]
		})
	}
	printSlice(a)
}

func printSlice(a [][]int) {
	for _, v := range a {
		fmt.Println(v)
	}
	fmt.Println()
}
