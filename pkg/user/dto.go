package user

type UserDTO struct {
	Username string `json:"username"`
	Private  bool   `json:"private"`
}