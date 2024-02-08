package responses

import (
	"time"
)

type UserInformationResponse struct {
	Username  string    `json:"username"`
	FullName  string    `json:"fullname"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
