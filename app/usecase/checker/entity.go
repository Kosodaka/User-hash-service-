package checker

import "mainHashService/app/entity"

type CheckerUcDto struct {
	Hash   []entity.Hash
	Domain int64
}

type HashUc struct {
	PhoneNumber string
	Salt        int64
}
