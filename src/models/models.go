package models

import (
	"errors"
	"time"
)

type Transaction struct {
	Value       int64     `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	Date        time.Time `json:"realizada_em"`
}

type Info struct {
	Balance int64     `json:"total"`
	Limit   int64     `json:"limite"`
	Date    time.Time `json:"data_extrato"`
}

type ExtractResponse struct {
	UserInfo     Info          `json:"saldo"`
	Transactions []Transaction `json:"ultimas_transacoes"`
}

type TransactionRequest struct {
	Value       int64  `json:"valor" binding:"required"`
	Type        string `json:"tipo" binding:"required"`
	Description string `json:"descricao" binding:"required"`
}

type TransactionResponse struct {
	Balance int64 `json:"saldo"`
	Limit   int64 `json:"limite"`
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInsufficientLimit = errors.New("insufficient limit")
)
