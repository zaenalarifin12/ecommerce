package userDto

type UserRequest struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Username string `json:"username"`
}
