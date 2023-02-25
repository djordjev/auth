package api

import (
	"encoding/json"
	"github.com/djordjev/auth/internal/domain/types"
)

var signUpRequest = `
	{
		"username": "djvukovic",
		"password": "testee",
		"email": "djvukovic@gmail.com",
		"role": "admin",
		"payload": { "foo": "bar", "count": 25 }
	}
`

func userFromRequest(id int) *types.User {
	var user types.User

	_ = json.Unmarshal([]byte(signUpRequest), &user)

	if id != 0 {
		user.ID = uint(id)
	}

	return &user
}
