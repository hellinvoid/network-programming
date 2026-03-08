package dto

type Request struct {
	Method string  `json:"method"`
	Number *float64 `json:"number"`
}
