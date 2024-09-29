package user

func NewUser(loginUserId string) *User {
	return &User{
		loginUserID: loginUserId,
	}
}

type User struct {
	loginUserID string
}
