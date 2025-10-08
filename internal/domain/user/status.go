package user

type UserStatus string

const (
	UserStatusActive UserStatus = "ACTIVE"
	UserStatusBloced UserStatus = "BLOCKED"
)
