package follow

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	FollowerId uuid.UUID `json:"followerId"`
	FollowToId uuid.UUID `json:"followToId"`
	FollowAt   time.Time `json:"followAt"`
}
