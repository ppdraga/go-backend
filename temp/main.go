package main

import "fmt"

func main() {
	fmt.Println("Hi!")
	x := []int{}
	x = append(x, 0)
	x = append(x, 1)
	x = append(x, 2)
	x = append(x, 5)
	y := append(x, 3)
	//y = append(y, 3)
	z := append(x, 4)
	//z = append(z, 5)
	fmt.Println(x, y, z)

	//nums := []int{}
	//nums := []int{5}
	//nums := make([]int, 3, 3)
	//fmt.Println(nums, len(nums), cap(nums))

	//changeNum(nums)

	//fmt.Println(nums, len(nums), cap(nums))
}

//func changeNum(ss []int) {
//    ss[0] = 1
//	ss[1] = 2
//	fmt.Println(ss)
//	ss = append(ss, 2)
//	ss[2] = 3
//	fmt.Println(ss)
//	fmt.Println(&ss)
//}
