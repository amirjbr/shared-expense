package dto

type UserRegisterRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	Username    string `json:"username"`
}
type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserLoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
