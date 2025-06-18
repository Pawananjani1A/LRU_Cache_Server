/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

//using: https://pkg.go.dev/github.com/go-playground/validator/v10

package reqmodels

type TestModel struct {
	FirstName string `json:"firstName" binding:"required,eq=pawananjani"`
	LastName  string `json:"lastName" binding:"required,eq=kumar"`
}

type GetModel struct {
	Key string `json:"key"`
}

type SetModel struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	Expiration string `json:"expiration"`
}
