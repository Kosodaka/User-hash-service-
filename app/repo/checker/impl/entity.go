package impl

// DTO для отправки запроса на сервис расхеширования
type RequestToUnhash struct {
	Hash   []HashDTO `json:"hash"`
	Domain string    `json:"domain"`
}

// вспомогательный DTO для расхеширования
type HashDTO struct {
	PhoneNumber string `json:"phone_number"`
	Salt        int64  `json:"salt"`
}

type VerifyHash struct {
	Hash []UnhsahedData `json:"hash"`
}

type UnhsahedData struct {
	UserID      int64  `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
}
