package postgres

import "time"

type UserData struct {
	ID       int64     `db:"user_id"`
	Name     string    `db:"user_name"`
	Surname  string    `db:"surname"`
	Email    string    `db:"email"`
	Phone    string    `db:"hashed_phone"`
	Salt     int64     `db:"salt"`
	Domain   int64     `db:"domain_number"`
	CreateAt time.Time `db:"created_at"`
}

type HashedData struct {
	UserID       int64     `json:"user_id"`
	UserName     string    `json:"user_name"`
	Surname      string    `json:"surname"`
	Email        string    `json:"email"`
	HashedPhone  string    `json:"hashed_phone"`
	Salt         int64     `json:"salt"`
	DomainNumber int64     `json:"domain_number"`
	CreatedAt    time.Time `json:"created_at"`
}
type UserHash struct {
	Hash   Hash
	Domain int64 `db:"domain_number"`
}

type UnhashRequest struct {
	HashSalt []Hash `json:"hash"`
	Domain   int64  `json:"domain"`
}
type UnhashResponse struct {
	HashSalt []HashFromUnhashResponse `json:"hash"`
}

type HashFromUnhashResponse struct {
	UserID      int64  `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
}

type Hash struct {
	UserID      int64  `json:"user_id" db:"user_id"`
	PhoneNumber string `json:"phone_number" db:"hashed_phone"`
	Salt        int64  `json:"salt" db:"salt"`
}

type Unhashdata struct {
	HashSalt []Hash `json:"hash"`
	Domain   int64  `json:"domain"`
}
