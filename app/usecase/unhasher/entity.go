package unhasher

import "mainHashService/app/repo/postgres"

type UserData = postgres.UserData
type UserHash = postgres.UserHash
type Hash = postgres.Hash
type Unhashdata = postgres.Unhashdata
type HashedData = postgres.HashedData

type ResultStruct struct {
	UserID     int64  `json:"user_id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Email      string `json:"email"`
	ClearPhone string `json:"clear_phone"`
}
