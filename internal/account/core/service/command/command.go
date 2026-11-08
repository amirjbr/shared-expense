package command

type RegisterUserCommand struct {
	FirstName   string
	LastName    string
	Email       string
	Password    string
	PhoneNumber string
	Username    string
}

type LoginUserCommand struct {
	Username string
	Password string
}
