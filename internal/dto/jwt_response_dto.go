package dto

type JwtResponse struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
}
