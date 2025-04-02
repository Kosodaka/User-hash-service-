package entity

type UserData struct {
	ID       int64
	Name     string
	Surname  string
	Email    string
	UserHash UserDataToUnhash
}

type UserDataToUnhash struct {
	Hash   []UserHash
	Domain int64
}

type UserHash struct {
	PhoneNumber string
	Salt        int64
}

type UnhashedData struct {
	HashSalt []UnhashedNumber `json:"hash"`
}

type UnhashedNumber struct {
	UserID      int64  `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	Hash        string `json:"phone_hash"`
}

type QueryStmt struct {
	Clause string
	Value  interface{}
}
