package data

// Minimal User struct - just what we need from user-service
type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Activated bool   `json:"activated"`
}

// AnonymousUser represents an unauthenticated user
var AnonymousUser = &User{}

// IsAnonymous checks if a User instance is the AnonymousUser
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}
