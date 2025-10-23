package response

import "github.com/google/uuid"

type GetUserInfo struct {
	Id                    uuid.UUID `json:"id"`
	Username              string    `json:"username"`
	Followers             int       `json:"followers"`
	Following             int       `json:"following"`
	PostCount             int       `json:"postCount"`
	IfCurrentUserFollowed bool      `json:"ifCurrentUserFollowed"`
	Description           *string   `json:"description" binding:"max=100"`
	IconUrl               *string   `json:"iconUrl"`
}
