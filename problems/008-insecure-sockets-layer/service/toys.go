package service

import (
	"strconv"
	"strings"
)

// Get max toy request
func GetToy(req string) string {
	strs := strings.Split(req, ",")
	max := 0
	num := make([]int, 0)
	for i, str := range strs {
		before, _, _ := strings.Cut(str, "x")
		val, err := strconv.Atoi(before)
		if err != nil {
			continue
		}
		num = append(num, val)
		if val > num[max] {
			max = i
		}
	}
	return strs[max]
}
