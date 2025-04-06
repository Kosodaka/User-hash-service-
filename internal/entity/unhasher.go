package entity

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
