package domain

type User struct {
	ID       uint
	Email    string
	Username string
	Password string
	Role     string
	Verified bool
}

type VerifyAccount struct {
	ID     uint
	Token  string
	UserID uint
}
