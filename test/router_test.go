package test_test

// 选择为数据路由的方式：存Map或存Slice（采用二分法查找）

import "testing"

func BenchmarkSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var j = i % (length - 1)
		if j != getSlice(j) {
			panic(j)
		}
	}
}

func BenchmarkMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var j = i % (length - 1)
		if j != getMap(j) {
			panic(j)
		}
	}
}

var length int = 80

var list []int = make([]int, length)

var m map[int]int = make(map[int]int)

func init() {
	for i := 0; i < length; i++ {
		list[i], m[i] = i, i
	}
}

func getSlice(T int) int {

	var x, y int = 0, len(list) - 1
	var s int
	for {
		s = x + y

		if y-x > 1 && list[s&0b1+s>>1] < T {
			x = s&0b1 + s/2

		} else if y-x > 1 && T < list[s>>1] {
			y = (x + y) / 2

		} else {

			if list[s&0b1+s>>1] == T {
				return s&0b1 + s>>1
			} else if list[s>>1] == T {
				return s >> 1
			} else {
				return -1
			}
		}
	}

}

func getMap(i int) int {
	return m[i]
}
