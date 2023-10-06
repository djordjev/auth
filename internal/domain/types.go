package domain

type User struct {
	ID       uint64
	Email    string
	Username string
	Password string
	Role     string
	Verified bool
	Payload  map[string]any
}

type VerifyAccount struct {
	ID     uint64
	Token  string
	UserID uint64
}

type ForgetPassword struct {
	ID     uint64
	Token  string
	UserID uint64
}

type Session struct {
	ID   string
	User User
}
