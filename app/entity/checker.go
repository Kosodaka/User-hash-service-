package entity

type Checker struct {
	Hash   []Hash `json:"hash"`
	Domain int64  `json:"domain"`
}

type Hash struct {
	UserID      int64  `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	Salt        int64  `json:"salt"`
}
