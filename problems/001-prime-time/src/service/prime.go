package service

import (
	"math"

	"github.com/hellinvoid/001-prime-time/dto"
)

func IsPrime(req *dto.Request) dto.Response {

	// log.Printf("%v,%v", req, *req.Number)

	res := dto.Response{
		Method: req.Method,
		Prime:  false,
	}

	// Return false if not integer or less than 2
	if math.Trunc(*req.Number) != *req.Number || *req.Number < 2 {
		return res
	}

	var num int64 = int64(*req.Number)
	n := math.Sqrt(float64(num))

	if num%2 == 0 {
		res.Prime = num == 2
		return res
	}

	var i int64
	for i = 3; i <= int64(n); i += 2 {
		if num%i == 0 {
			return res
		}
	}

	res.Prime = true
	return res
}
