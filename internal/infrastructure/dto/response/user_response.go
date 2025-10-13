package response

import "github.com/google/uuid"

type GetUserInfo struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	IconUrl  *string   `json:"iconUrl"`
}
