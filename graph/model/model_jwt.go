package model

type JwtToken struct {
	AccessToken string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
