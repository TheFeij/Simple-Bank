package requests

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
