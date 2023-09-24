package domain

type User struct {
	ID       uint
	Email    string
	Username string
	Password string
	Role     string
	Verified bool
	Payload  map[string]any
}

type VerifyAccount struct {
	ID     uint
	Token  string
	UserID uint
}

type ForgetPassword struct {
	ID     uint
	Token  string
	UserID uint
}

type Session struct {
	ID   string
	User User
}
