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

type LoginResponse struct {
	AccessToken           string                  `json:"access_token"`
	AccessTokenExpiresAt  time.Time               `json:"access_token_expires_at"`
	RefreshToken          string                  `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time               `json:"refresh_token_expires_at"`
	UserInformation       UserInformationResponse `json:"user_information"`
}
