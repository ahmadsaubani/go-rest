package entities

import (
	"gin/src/entities/auth"
	"gin/src/entities/users"
)

var RegisteredEntities = []any{
	users.User{},
	auth.AccessToken{},
	auth.RefreshToken{},
}
