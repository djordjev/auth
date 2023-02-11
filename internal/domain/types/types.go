package types

type User struct {
	ID       uint
	Email    string
	Username string
	Password string
	Role     string
	Payload  map[string]any
}
