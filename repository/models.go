package repository

import "time"

type User struct {
	Id     int64  `json:"id"`
	Pseudo string `json:"pseudo"`
}

type OperationType string

const (
	TypeAdd       OperationType = "add"
	TypeSubstract OperationType = "substract"
	TypeMultiply  OperationType = "multiply"
	TypeDivide    OperationType = "divide"
	TypeSum       OperationType = "sum"
)

type Operations struct {
	Id        int64         `json:"id"`
	Inputs    []float64     `json:"inputs"`
	Type      OperationType `json:"type"`
	Result    float64       `json:"results"`
	UserId    float64       `json:"user_id"`
	CreatedAt time.Time     `json:"created_at"`
}
